package Lexer

import (
	"fmt"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
	ud "github.com/PlayerR9/MyGoLib/Units/Debugging"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

type leavesResult struct {
	leaves []*tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.LeafToken]]
}

func (lr *leavesResult) Size() int {
	return len(lr.leaves)
}

func (lr *leavesResult) getFirst() *tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.LeafToken]] {
	if len(lr.leaves) == 0 {
		return nil
	}

	first := lr.leaves[0]
	lr.leaves = lr.leaves[1:]

	return first
}

// SourceIterator is an iterator that uses a grammar to tokenize a string.
type SourceIterator struct {
	// source is the source to lex.
	source *cds.Stream[byte]

	// tree is the tree to use.
	tree *ud.History[*tr.Tree[*tr.StatusInfo[EvalStatus, *gr.LeafToken]]]

	// productions are the production rules to use.
	productions []*gr.RegProduction

	// isFirst is a flag that indicates if the lexer is the first.
	isFirst bool

	// canContinue is a flag that indicates if the lexer can continue.
	canContinue bool

	// errBranches are the branches that have errors.
	errBranches []*tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.LeafToken]]
}

// Size implements the Iterater interface.
//
// Size is calculated by the number of leaves in the tree.
// of course, this is just an approximation as, to get the exact size,
// we would need to traverse the entire tree.
func (si *SourceIterator) Size() (count int) {
	f := func(data *tr.Tree[*tr.StatusInfo[EvalStatus, *gr.LeafToken]]) {
		count = data.Size()
	}

	si.tree.ReadData(f)

	return
}

// generateEvalTrees adds the matches to a root tree as leaves.
//
// Parameters:
//   - root: The root of the tree to add the leaves to.
//   - matches: The matches to add to the tree evaluator.
func generateEvalTrees(matches []*gr.MatchedResult[*gr.LeafToken]) []*tr.Tree[*tr.StatusInfo[EvalStatus, *gr.LeafToken]] {
	// DEBUG: Display the matches
	fmt.Println("Matches:")
	for _, match := range matches {
		fmt.Printf("\t%+v\n", match.Matched)
	}

	// Get the longest match.
	matches = selectBestMatches(matches)

	children := make([]*tr.Tree[*tr.StatusInfo[EvalStatus, *gr.LeafToken]], 0, len(matches))

	for _, match := range matches {
		currMatch := match.GetMatch()

		inf := tr.NewStatusInfo(currMatch, EvalIncomplete)

		tree := tr.NewTree(inf)
		children = append(children, tree)
	}

	return children
}

// lexOne lexes one branch of the tree.
//
// Returns:
//   - error: An error if lexing fails.
func (si *SourceIterator) lexOne() error {
	if si.isFirst {
		matches, err := matchFrom(si.source, 0, si.productions)
		if err != nil {
			return uc.NewErrAt(0, "position", err)
		}

		children := generateEvalTrees(matches)

		cmd := tr.NewSetChildrenCmd(children)
		si.tree.ExecuteCommand(cmd)

		si.tree.ReadData(func(data *tr.Tree[*tr.StatusInfo[EvalStatus, *gr.LeafToken]]) {
			root := data.Root()
			ld := root.Data
			ld.ChangeStatus(EvalComplete)
		})

		// DEBUG: Display the resulting tree
		fmt.Println("Resulting Tree:")
		si.tree.ReadData(func(data *tr.Tree[*tr.StatusInfo[EvalStatus, *gr.LeafToken]]) {
			p, trav := ffs.NewStdPrinter(
				ffs.NewFormatter(ffs.NewIndentConfig("   ", 0)),
			)

			err := data.FString(trav)
			if err != nil {
				panic(err)
			}

			pages := ffs.Stringfy(p.GetPages())

			fmt.Println(pages[0])
		})

		si.isFirst = false
	} else {
		p := filterLeaves(si.source, si.productions)

		cmd1 := tr.NewProcessLeavesCmd(p)
		err := si.tree.ExecuteCommand(cmd1)
		if err != nil {
			si.tree.UndoLastCommand()
			return err
		}

		si.tree.Accept()

		var leaves []*tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.LeafToken]]

		si.tree.ReadData(func(data *tr.Tree[*tr.StatusInfo[EvalStatus, *gr.LeafToken]]) {
			leaves = data.GetLeaves()
		})

		var success []*tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.LeafToken]]
		var failed []*tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.LeafToken]]

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
		si.errBranches = append(si.errBranches, failed...)

		/*
			cmd2 := tr.NewPruneTreeCmd(filterErrorLeaves)
			err = si.tree.ExecuteCommand(cmd2)
			if err != nil {
				si.tree.UndoLastCommand()
				return err
			}

			si.tree.Accept()
		*/
	}

	return nil
}

// Consume implements the Iterater interface.
func (si *SourceIterator) Consume() (*leavesResult, error) {
	var result *leavesResult

	for {
		if !si.canContinue {
			return nil, uc.NewErrExhaustedIter()
		}

		var leaves []*tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.LeafToken]]

		err := si.lexOne()
		if err != nil {
			si.canContinue = false

			f := func(data *tr.Tree[*tr.StatusInfo[EvalStatus, *gr.LeafToken]]) {
				leaves = data.GetLeaves()
			}

			si.tree.ReadData(f)
		} else {
			leaves, si.canContinue = si.getCompletedBranch()
		}

		// Ignore error leaves.
		leaves = us.SliceFilter(leaves, func(leaf *tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.LeafToken]]) bool {
			ld := leaf.Data
			status := ld.GetStatus()

			return status != EvalError
		})

		if len(leaves) > 0 {
			result = &leavesResult{
				leaves: leaves,
			}

			break
		}
	}

	return result, nil
}

// Restart implements the Iterater interface.
func (si *SourceIterator) Restart() {
	si.isFirst = true
	si.canContinue = true

	rootNode := gr.NewRootToken()

	p := tr.NewStatusInfo(rootNode, EvalIncomplete)

	tree := tr.NewTreeWithHistory(p)
	si.tree = tree
}

// newSourceIterator creates a new source iterator.
//
// Parameters:
//   - source: The source to use.
//   - productions: The production rules to use.
//
// Returns:
//   - *SourceIterator: The new source iterator.
func newSourceIterator(source *cds.Stream[byte], productions []*gr.RegProduction) *SourceIterator {
	rootNode := gr.NewRootToken()

	p := tr.NewStatusInfo(rootNode, EvalIncomplete)

	tree := tr.NewTreeWithHistory(p)

	si := &SourceIterator{
		source:      source,
		productions: productions,
		tree:        tree,
		isFirst:     true,
		canContinue: true,
	}

	return si
}

// getCompletedBranch gets the completed branch of the tree.
//
// Returns:
//   - []*tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.LeafToken]]: The completed branch.
//   - bool: True if the branch can continue, false otherwise.
func (si *SourceIterator) getCompletedBranch() ([]*tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.LeafToken]], bool) {
	var leaves []*tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.LeafToken]]

	f := func(data *tr.Tree[*tr.StatusInfo[EvalStatus, *gr.LeafToken]]) {
		leaves = data.GetLeaves()
	}

	si.tree.ReadData(f)

	canContinue := false

	var completedLeaves []*tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.LeafToken]]

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
func (si *SourceIterator) deleteBranch(leaf *tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.LeafToken]]) {
	cmd := tr.NewDeleteBranchContainingCmd(leaf)
	si.tree.ExecuteCommand(cmd)
	si.tree.Accept()
}

// LexerIterator is an iterator that uses a grammar to tokenize a string.
type LexerIterator struct {
	// toSkip are the tokens to skip.
	toSkip []string

	// completedLeaves are the leaves that have been completed.
	completedLeaves *leavesResult

	// sourceIter is the source iterator.
	sourceIter *SourceIterator
}

// Size implements the Iterater interface.
//
// Size is calculated by the number of leaves in the source iterator
// and the number of completed leaves.
func (li *LexerIterator) Size() (count int) {
	count = li.sourceIter.Size()

	count += li.completedLeaves.Size()

	return
}

// Consume implements the Iterater interface.
func (li *LexerIterator) Consume() (*cds.Stream[*gr.LeafToken], error) {
	var branch *cds.Stream[*gr.LeafToken]

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
func (li *LexerIterator) Restart() {
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
func executeLexing(source *cds.Stream[byte], productions []*gr.RegProduction) (*tr.Tree[*tr.StatusInfo[EvalStatus, *gr.LeafToken]], error) {
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

	tree.ReadData(func(data *tr.Tree[*tr.StatusInfo[EvalStatus, *gr.LeafToken]]) {
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

		var completedLeaves []*tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.LeafToken]]

		tree.ReadData(func(data *tr.Tree[*tr.StatusInfo[EvalStatus, *gr.LeafToken]]) {
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
