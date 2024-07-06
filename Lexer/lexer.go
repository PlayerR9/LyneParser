package Lexer

import (
	"fmt"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

type leavesResult[T uc.Enumer] struct {
	leaves []*tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.Token[T]]]
}

func (lr *leavesResult[T]) Size() int {
	return len(lr.leaves)
}

func (lr *leavesResult[T]) getFirst() *tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.Token[T]]] {
	if len(lr.leaves) == 0 {
		return nil
	}

	first := lr.leaves[0]
	lr.leaves = lr.leaves[1:]

	return first
}

// SourceIterator is an iterator that uses a grammar to tokenize a string.
type SourceIterator[T uc.Enumer] struct {
	// source is the source to lex.
	source *cds.Stream[byte]

	// tree is the tree to use.
	tree *tr.Tree[*tr.StatusInfo[EvalStatus, *gr.Token[T]]]

	// productions are the production rules to use.
	productions []*gr.RegProduction[T]

	// canContinue is a flag that indicates if the lexer can continue.
	canContinue bool

	// errBranches are the branches that have errors.
	errBranches []*tr.Branch[*tr.StatusInfo[EvalStatus, *gr.Token[T]]]

	// logger is a flag that indicates if the lexer should be verbose.
	logger *Verbose
}

// Size implements the Iterater interface.
//
// Size is calculated by the number of leaves in the tree.
// of course, this is just an approximation as, to get the exact size,
// we would need to traverse the entire tree.
func (si *SourceIterator[T]) Size() (count int) {
	count = si.tree.Size()

	return
}

// generateEvalTrees adds the matches to a root tree as leaves.
//
// Parameters:
//   - matches: The matches to add to the tree evaluator.
//   - logger: A verbose logger.
func generateEvalTrees[T uc.Enumer](matches []*gr.MatchedResult[T], logger *Verbose) []*tr.StatusInfo[EvalStatus, *gr.Token[T]] {
	logger.DoIf(func(p *Printer) {
		// DEBUG: Display the matches
		p.Print("Matches:")

		for _, match := range matches {
			p.Printf("\t%+v", match.Matched)
		}
	})

	// Get the longest match.
	matches = selectBestMatches(matches, logger)

	children := make([]*tr.StatusInfo[EvalStatus, *gr.Token[T]], 0, len(matches))

	for _, match := range matches {
		inf := tr.NewStatusInfo(match.Matched, EvalIncomplete)

		children = append(children, inf)
	}

	return children
}

// lexOne lexes one branch of the tree.
//
// Parameters:
//   - logger: A verbose logger.
//
// Returns:
//   - error: An error if lexing fails.
func (si *SourceIterator[T]) lexOne(logger *Verbose) error {
	p := filterLeaves(si.source, si.productions, logger)

	err := si.tree.ProcessLeaves(p)
	if err != nil {
		return fmt.Errorf("failed to process leaves: %w", err)
	}

	logger.DoIf(func(p *Printer) {
		// DEBUG: Display the resulting tree
		p.Print("Resulting Tree:")

		printer, trav := ffs.NewStdPrinter(
			ffs.NewFormatter(ffs.NewIndentConfig("   ", 0)),
		)

		err = si.tree.FString(trav)
		if err != nil {
			panic(err)
		}

		pages := ffs.Stringfy(printer.GetPages(), 1)

		p.Print(pages[0])
	})

	leaves := si.tree.GetLeaves()

	var success []*tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.Token[T]]]
	var failed []*tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.Token[T]]]

	for _, leaf := range leaves {
		ld := leaf.Data

		status := ld.GetStatus()
		if status == EvalError {
			failed = append(failed, leaf)
		} else {
			success = append(success, leaf)
		}
	}

	// Add the failed branches to the error branches.
	for _, leaf := range failed {
		branch := si.tree.ExtractBranch(leaf, true)
		if branch != nil {
			si.errBranches = append(si.errBranches, branch)
		}
	}

	/*
		cmd2 := tr.NewPruneTreeCmd(filterErrorLeaves)
		err = si.tree.ExecuteCommand(cmd2)
		if err != nil {
			si.tree.UndoLastCommand()
			return err
		}

		si.tree.Accept()
	*/

	return nil
}

// Consume implements the Iterater interface.
func (si *SourceIterator[T]) Consume() (*leavesResult[T], error) {
	var result *leavesResult[T]

	for {
		if !si.canContinue {
			return nil, uc.NewErrExhaustedIter()
		}

		var leaves []*tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.Token[T]]]

		err := si.lexOne(si.logger)
		if err != nil {
			si.canContinue = false

			leaves = si.tree.GetLeaves()
		} else {
			leaves, si.canContinue = si.getCompletedBranch()
		}

		// Ignore error leaves.
		f := func(leaf *tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.Token[T]]]) bool {
			ld := leaf.Data
			status := ld.GetStatus()

			return status != EvalError
		}

		leaves = us.SliceFilter(leaves, f)

		if len(leaves) > 0 {
			result = &leavesResult[T]{
				leaves: leaves,
			}

			break
		}
	}

	return result, nil
}

// Restart implements the Iterater interface.
func (si *SourceIterator[T]) Restart() {
	si.canContinue = true

	//rootNode := gr.RootToken()

	// p := tr.NewStatusInfo(rootNode, EvalIncomplete)

	// si.tree = tr.NewTree(p)

	panic("SourceIterator.Restart() not implemented")
}

// newSourceIterator creates a new source iterator.
//
// Parameters:
//   - source: The source to use.
//   - productions: The production rules to use.
//   - logger: A verbose logger.
//
// Returns:
//   - *SourceIterator: The new source iterator.
func newSourceIterator[T uc.Enumer](source *cds.Stream[byte], productions []*gr.RegProduction[T], logger *Verbose) *SourceIterator[T] {
	// rootNode := gr.RootToken()

	// p := tr.NewStatusInfo(rootNode, EvalIncomplete)

	// tree := tr.NewTree(p)

	// si := &SourceIterator[T]{
	// 	source:      source,
	// 	productions: productions,
	// 	tree:        tree,
	// 	canContinue: true,
	// 	logger:      logger,
	// }

	// return si

	panic("newSourceIterator() not implemented")
}

// getCompletedBranch gets the completed branch of the tree.
//
// Returns:
//   - []*tr.TreeNode[*tr.StatusInfo[EvalStatus, gr.Token]]: The completed branch.
//   - bool: True if the branch can continue, false otherwise.
func (si *SourceIterator[T]) getCompletedBranch() ([]*tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.Token[T]]], bool) {
	leaves := si.tree.GetLeaves()

	canContinue := false

	var completedLeaves []*tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.Token[T]]]

	for _, leaf := range leaves {
		ld := leaf.Data
		status := ld.GetStatus()

		if status == EvalComplete {
			completedLeaves = append(completedLeaves, leaf)
		} else if status == EvalIncomplete {
			canContinue = true
		}
	}

	return completedLeaves, canContinue
}

// deleteBranch deletes a branch from the tree.
//
// Parameters:
//   - leaf: The leaf to delete.
func (si *SourceIterator[T]) deleteBranch(leaf *tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.Token[T]]]) {
	si.tree.DeleteBranchContaining(leaf)
}

// LexerIterator is an iterator that uses a grammar to tokenize a string.
type LexerIterator[T uc.Enumer] struct {
	// toSkip are the tokens to skip.
	toSkip []T

	// completedLeaves are the leaves that have been completed.
	completedLeaves *leavesResult[T]

	// sourceIter is the source iterator.
	sourceIter *SourceIterator[T]
}

// Size implements the Iterater interface.
//
// Size is calculated by the number of leaves in the source iterator
// and the number of completed leaves.
func (li *LexerIterator[T]) Size() (count int) {
	count = li.sourceIter.Size()

	count += li.completedLeaves.Size()

	return
}

// Consume implements the Iterater interface.
func (li *LexerIterator[T]) Consume() (*cds.Stream[*gr.Token[T]], error) {
	var branch *cds.Stream[*gr.Token[T]]

	for {
		if li.completedLeaves.Size() == 0 {
			res, err := li.sourceIter.Consume()
			if err != nil {
				return nil, err
			}

			li.completedLeaves = res
		}

		// Extract the entire branch
		leaf := li.completedLeaves.getFirst()

		anch := leaf.GetAncestors()
		anch = append(anch, leaf)

		li.sourceIter.deleteBranch(leaf)

		branch = convertBranchToTokenStream(anch, li.toSkip)
		if branch.Size() > 0 {
			break
		}
	}

	return branch, nil
}

// Restart implements the Iterater interface.
func (li *LexerIterator[T]) Restart() {
	li.completedLeaves = nil
	li.sourceIter.Restart()
}

/*

// executeLexing is the main function of the tree evaluator.
//
// Parameters:
//   - source: The source to executeLexing.
//   - root: The root of the tree evaluator.
//
// Returns:
//   - error: An error if lexing fails.
//
// Errors:
//   - *ErrEmptyInput: The source is empty.
//   - *uc.ErrAt: An error occurred at a specific index.
//   - *ErrAllMatchesFailed: All matches failed.
func executeLexing(source *cds.Stream[byte], productions []*gr.RegProduction) (*tr.Tree[*tr.StatusInfo[EvalStatus, gr.Token]], error) {
	rootNode := gr.NewRootToken()

	tree := tr.NewTreeWithHistory(EvalIncomplete, rootNode)

	matches, err := matchFrom(source, 0, productions)
	if err != nil {
		data := extractData(tree)

		return data, uc.NewErrAt(0, "position", err)
	}

	children := generateEvalTrees(matches)

	cmd := tr.NewSetChildrenCmd(children)
	tree.ExecuteCommand(cmd)

	tree.ReadData(func(data *tr.Tree[*tr.StatusInfo[EvalStatus, gr.Token]]) {
		root := data.Root()
		root.ChangeStatus(EvalComplete)
	})

	shouldContinue := true

	for shouldContinue {
		tree.Accept() // Accept the changes made by the last command.

		p := filterLeaves(source, productions)

		cmd1 := tr.NewProcessLeavesCmd(p)
		err := tree.ExecuteCommand(cmd1)
		if err != nil {
			tree.UndoLastCommand()

			data := extractData(tree)
			return data, err
		}

		tree.Accept()

		cmd2 := tr.NewPruneTreeCmd(filterErrorLeaves)
		err = tree.ExecuteCommand(cmd2)
		if err != nil {
			tree.UndoLastCommand()

			data := extractData(tree)
			return data, err
		}

		tree.Accept()

		var completedLeaves []*tr.TreeNode[*tr.StatusInfo[EvalStatus, gr.Token]]

		tree.ReadData(func(data *tr.Tree[*tr.StatusInfo[EvalStatus, gr.Token]]) {
			completedLeaves = data.GetLeaves()
		})

		shouldContinue = false
		for _, leaf := range completedLeaves {
			status := leaf.GetStatus()

			if status == EvalIncomplete {
				shouldContinue = true
				break
			}
		}
	}

	cmd3 := tr.NewPruneTreeCmd(filterIncompleteLeaves)
	err = tree.ExecuteCommand(cmd3)
	if err != nil {
		tree.UndoLastCommand()

		data := extractData(tree)
		return data, NewErrAllMatchesFailed()
	}

	tree.Accept()

	data := extractData(tree)

	return data, nil
}
*/
