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

type leaves_result[T gr.TokenTyper] struct {
	leaves []*TreeNode[T]
}

func new_leaves_result[T gr.TokenTyper](nodes []tr.Noder) *leaves_result[T] {
	var valid_nodes []*TreeNode[T]

	for _, node := range nodes {
		n, ok := node.(*TreeNode[T])
		if !ok {
			continue
		}

		valid_nodes = append(valid_nodes, n)
	}

	lr := &leaves_result[T]{
		leaves: valid_nodes,
	}

	return lr
}

func (lr *leaves_result[T]) Size() int {
	return len(lr.leaves)
}

func (lr *leaves_result[T]) get_first() *TreeNode[T] {
	if len(lr.leaves) == 0 {
		return nil
	}

	first := lr.leaves[0]
	lr.leaves = lr.leaves[1:]

	return first
}

// SourceIterator is an iterator that uses a grammar to tokenize a string.
type SourceIterator[T gr.TokenTyper] struct {
	// source is the source to lex.
	source *cds.Stream[byte]

	// tree is the tree to use.
	tree *tr.Tree

	// productions are the production rules to use.
	productions []*gr.RegProduction[T]

	// can_continue is a flag that indicates if the lexer can continue.
	can_continue bool

	// err_branches are the branches that have errors.
	err_branches []*tr.Branch

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

// generate_eval_trees adds the matches to a root tree as leaves.
//
// Parameters:
//   - matches: The matches to add to the tree evaluator.
//   - logger: A verbose logger.
func generate_eval_trees[T gr.TokenTyper](matches []*gr.MatchedResult[T], logger *Verbose) []*TreeNode[T] {
	logger.DoIf(func(p *Printer) {
		// DEBUG: Display the matches
		p.Print("Matches:")

		for _, match := range matches {
			p.Printf("\t%+v", match.Matched)
		}
	})

	// Get the longest match.
	matches = select_best_matches(matches, logger)

	children := make([]*TreeNode[T], 0, len(matches))

	for _, match := range matches {
		tn := NewTreeNode(match.Matched)

		children = append(children, tn)
	}

	return children
}

// lex_one lexes one branch of the tree.
//
// Parameters:
//   - logger: A verbose logger.
//
// Returns:
//   - error: An error if lexing fails.
func (si *SourceIterator[T]) lex_one(logger *Verbose) error {
	p := filter_leaves(si.source, si.productions, logger)

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
		uc.AssertF(err == nil, "FString failed: %s", err.Error())

		pages := ffs.Stringfy(printer.GetPages(), 1)

		p.Print(pages[0])
	})

	leaves := si.tree.GetLeaves()

	var success []*TreeNode[T]
	var failed []*TreeNode[T]

	for _, leaf := range leaves {
		tn, ok := leaf.(*TreeNode[T])
		uc.Assert(ok, "Must be a *TreeNode[T]")

		if tn.Status == EvalError {
			failed = append(failed, tn)
		} else {
			success = append(success, tn)
		}
	}

	// Add the failed branches to the error branches.
	for i, leaf := range failed {
		branch, err := si.tree.ExtractBranch(leaf, true)
		if err != nil {
			return uc.NewErrWhileAt("extracting", i+1, "branch", err)
		}

		if branch != nil {
			si.err_branches = append(si.err_branches, branch)
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
func (si *SourceIterator[T]) Consume() (*leaves_result[T], error) {
	var result *leaves_result[T]

	for {
		if !si.can_continue {
			return nil, uc.NewErrExhaustedIter()
		}

		var leaves []tr.Noder

		err := si.lex_one(si.logger)
		if err != nil {
			si.can_continue = false

			leaves = si.tree.GetLeaves()
		} else {
			leaves, si.can_continue = si.get_completed_branch()
		}

		// Ignore error leaves.
		f := func(leaf tr.Noder) bool {
			tn, ok := leaf.(*TreeNode[T])
			if !ok {
				return false
			}

			return tn.Status != EvalError
		}

		leaves = us.SliceFilter(leaves, f)

		if len(leaves) > 0 {
			result = new_leaves_result[T](leaves)
			break
		}
	}

	return result, nil
}

// Restart implements the Iterater interface.
func (si *SourceIterator[T]) Restart() {
	si.can_continue = true

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
func newSourceIterator[T gr.TokenTyper](source *cds.Stream[byte], productions []*gr.RegProduction[T], logger *Verbose) *SourceIterator[T] {
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

// get_completed_branch gets the completed branch of the tree.
//
// Returns:
//   - []Tree.Noder: The completed branch.
//   - bool: True if the branch can continue, false otherwise.
func (si *SourceIterator[T]) get_completed_branch() ([]tr.Noder, bool) {
	leaves := si.tree.GetLeaves()

	can_continue := false

	var completed_leaves []tr.Noder

	for _, leaf := range leaves {
		tn, ok := leaf.(*TreeNode[T])
		uc.Assert(ok, "Must be a *TreeNode[T]")

		switch tn.Status {
		case EvalComplete:
			completed_leaves = append(completed_leaves, leaf)
		case EvalIncomplete:
			can_continue = true
		}
	}

	return completed_leaves, can_continue
}

// delete_branch deletes a branch from the tree.
//
// Parameters:
//   - leaf: The leaf to delete.
func (si *SourceIterator[T]) delete_branch(leaf tr.Noder) {
	err := si.tree.DeleteBranchContaining(leaf)
	uc.AssertF(err == nil, "DeleteBranchContaining failed: %s", err.Error())
}

// LexerIterator is an iterator that uses a grammar to tokenize a string.
type LexerIterator[T gr.TokenTyper] struct {
	// to_skip are the tokens to skip.
	to_skip []T

	// completed_leaves are the leaves that have been completed.
	completed_leaves *leaves_result[T]

	// source_iter is the source iterator.
	source_iter *SourceIterator[T]
}

// Size implements the Iterater interface.
//
// Size is calculated by the number of leaves in the source iterator
// and the number of completed leaves.
func (li *LexerIterator[T]) Size() (count int) {
	count = li.source_iter.Size()

	count += li.completed_leaves.Size()

	return
}

// Consume implements the Iterater interface.
func (li *LexerIterator[T]) Consume() (*cds.Stream[*gr.Token[T]], error) {
	var branch *cds.Stream[*gr.Token[T]]

	for {
		if li.completed_leaves.Size() == 0 {
			res, err := li.source_iter.Consume()
			if err != nil {
				return nil, err
			}

			li.completed_leaves = res
		}

		// Extract the entire branch
		leaf := li.completed_leaves.get_first()

		anch := leaf.GetAncestors()
		anch = append(anch, leaf)

		li.source_iter.delete_branch(leaf)

		branch = convert_branch_to_token_stream(anch, li.to_skip)
		if branch.Size() > 0 {
			break
		}
	}

	return branch, nil
}

// Restart implements the Iterater interface.
func (li *LexerIterator[T]) Restart() {
	li.completed_leaves = nil
	li.source_iter.Restart()
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
