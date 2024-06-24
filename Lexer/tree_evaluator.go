package Lexer

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/StatusTree"
	ud "github.com/PlayerR9/MyGoLib/Units/Debugging"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	ers "github.com/PlayerR9/MyGoLib/Units/errors"
)

// generateEvalTrees adds the matches to a root tree as leaves.
//
// Parameters:
//   - root: The root of the tree to add the leaves to.
//   - matches: The matches to add to the tree evaluator.
func generateEvalTrees(matches []*gr.MatchedResult[*gr.LeafToken]) []*tr.Tree[EvalStatus, *gr.LeafToken] {
	// Get the longest match.
	matches = selectBestMatches(matches)

	children := make([]*tr.Tree[EvalStatus, *gr.LeafToken], 0, len(matches))

	for _, match := range matches {
		currMatch := match.GetMatch()

		tree := tr.NewTree(EvalIncomplete, currMatch)
		children = append(children, tree)
	}

	return children
}

// filterLeaves processes the leaves in the tree evaluator.
//
// Returns:
//   - bool: True if all leaves are complete, false otherwise.
//   - error: An error of type *ErrAllMatchesFailed if all matches failed.
func filterLeaves(source *cds.Stream[byte], productions []*gr.RegProduction) uc.EvalManyFunc[*tr.TreeNode[EvalStatus, *gr.LeafToken], uc.Pair[EvalStatus, *gr.LeafToken]] {
	filterFunc := func(leaf *tr.TreeNode[EvalStatus, *gr.LeafToken]) ([]uc.Pair[EvalStatus, *gr.LeafToken], error) {
		nextAt := leaf.Data.GetPos() + len(leaf.Data.Data)

		if nextAt >= source.Size() {
			leaf.ChangeStatus(EvalComplete)
			return nil, nil
		}

		matches, err := matchFrom(source, nextAt, productions)
		if err != nil {
			leaf.ChangeStatus(EvalError)
			return nil, nil
		}

		// Get the longest match.
		matches = selectBestMatches(matches)

		children := make([]uc.Pair[EvalStatus, *gr.LeafToken], 0, len(matches))

		for _, match := range matches {
			curr := match.GetMatch()
			p := uc.NewPair(EvalIncomplete, curr)

			children = append(children, p)
		}

		leaf.ChangeStatus(EvalComplete)

		return children, nil
	}

	return filterFunc
}

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
//   - *ers.ErrAt: An error occurred at a specific index.
//   - *ErrAllMatchesFailed: All matches failed.
func executeLexing(source *cds.Stream[byte], productions []*gr.RegProduction) (*tr.Tree[EvalStatus, *gr.LeafToken], error) {
	rootNode := gr.NewRootToken()

	tree := tr.NewTreeWithHistory(EvalIncomplete, rootNode)

	matches, err := matchFrom(source, 0, productions)
	if err != nil {
		data := extractData(tree)

		return data, ers.NewErrAt(0, "position", err)
	}

	children := generateEvalTrees(matches)

	cmd := tr.NewSetChildrenCmd(children)
	tree.ExecuteCommand(cmd)

	tree.ReadData(func(data *tr.Tree[EvalStatus, *gr.LeafToken]) {
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

		var completedLeaves []*tr.TreeNode[EvalStatus, *gr.LeafToken]

		tree.ReadData(func(data *tr.Tree[EvalStatus, *gr.LeafToken]) {
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

// FIXME: Remove this function once MyGoLib is updated.
func extractData(tree *ud.History[*tr.Tree[EvalStatus, *gr.LeafToken]]) *tr.Tree[EvalStatus, *gr.LeafToken] {
	var tmp *tr.Tree[EvalStatus, *gr.LeafToken]

	tree.ReadData(func(data *tr.Tree[EvalStatus, *gr.LeafToken]) {
		tmp = data
	})

	return tmp
}

type ActiveLexer struct {
	source      *cds.Stream[byte]
	productions []*gr.RegProduction
	skip        []string

	tree *ud.History[*tr.Tree[EvalStatus, *gr.LeafToken]]

	completedLeaves []*tr.TreeNode[EvalStatus, *gr.LeafToken]

	canContinue bool
	isFirst     bool
}

func newLexState(source *cds.Stream[byte], productions []*gr.RegProduction, toSkip []string) *ActiveLexer {
	rootNode := gr.NewRootToken()

	tree := tr.NewTreeWithHistory(EvalIncomplete, rootNode)

	ls := &ActiveLexer{
		source:      source,
		productions: productions,
		skip:        toSkip,
		tree:        tree,
		isFirst:     true,
		canContinue: true,
	}

	return ls
}

func (ls *ActiveLexer) lexOne() error {
	if ls.isFirst {
		matches, err := matchFrom(ls.source, 0, ls.productions)
		if err != nil {
			return ers.NewErrAt(0, "position", err)
		}

		children := generateEvalTrees(matches)

		cmd := tr.NewSetChildrenCmd(children)
		ls.tree.ExecuteCommand(cmd)

		ls.tree.ReadData(func(data *tr.Tree[EvalStatus, *gr.LeafToken]) {
			root := data.Root()
			root.ChangeStatus(EvalComplete)
		})

		ls.isFirst = false
	} else {
		p := filterLeaves(ls.source, ls.productions)

		cmd1 := tr.NewProcessLeavesCmd(p)
		err := ls.tree.ExecuteCommand(cmd1)
		if err != nil {
			ls.tree.UndoLastCommand()
			return err
		}

		ls.tree.Accept()

		cmd2 := tr.NewPruneTreeCmd(filterErrorLeaves)
		err = ls.tree.ExecuteCommand(cmd2)
		if err != nil {
			ls.tree.UndoLastCommand()
			return err
		}

		ls.tree.Accept()
	}

	return nil
}

func (ls *ActiveLexer) getCompletedBranch() ([]*tr.TreeNode[EvalStatus, *gr.LeafToken], bool) {
	var leaves []*tr.TreeNode[EvalStatus, *gr.LeafToken]

	ls.tree.ReadData(func(data *tr.Tree[EvalStatus, *gr.LeafToken]) {
		leaves = data.GetLeaves()
	})

	canContinue := false

	var completedLeaves []*tr.TreeNode[EvalStatus, *gr.LeafToken]

	for _, leaf := range leaves {
		status := leaf.GetStatus()

		if status == EvalComplete {
			completedLeaves = append(completedLeaves, leaf)
		} else if status == EvalIncomplete {
			canContinue = true
		}
	}

	return completedLeaves, canContinue
}

func (ls *ActiveLexer) GetBranch() (*cds.Stream[*gr.LeafToken], error) {
	if len(ls.completedLeaves) == 0 {
		if !ls.canContinue {
			return nil, nil
		}

		err := ls.lexOne()
		if err != nil {
			ls.canContinue = false

			var leaves []*tr.TreeNode[EvalStatus, *gr.LeafToken]

			ls.tree.ReadData(func(data *tr.Tree[EvalStatus, *gr.LeafToken]) {
				leaves = data.GetLeaves()
			})

			ls.completedLeaves = leaves
		} else {
			leaves, canContinue := ls.getCompletedBranch()

			ls.canContinue = canContinue
			ls.completedLeaves = leaves
		}
	}

	// Extract the entire branch
	leaf := ls.completedLeaves[0]

	anch := leaf.GetAncestors()
	anch = append(anch, leaf)

	ls.completedLeaves = ls.completedLeaves[1:]

	// Remove the branch from the tree
	cmd := tr.NewDeleteBranchContainingCmd(leaf)
	ls.tree.ExecuteCommand(cmd)
	ls.tree.Accept()

	branch := convertBranchToTokenStream(anch, ls.skip)

	return branch, nil
}
