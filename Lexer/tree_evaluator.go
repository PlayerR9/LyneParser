package Lexer

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/StatusTree"
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
func filterLeaves(source *cds.Stream[byte], productions []*gr.RegProduction) uc.EvalManyFunc[*tr.TreeNode[EvalStatus, *gr.LeafToken], *uc.Pair[EvalStatus, *gr.LeafToken]] {
	filterFunc := func(leaf *tr.TreeNode[EvalStatus, *gr.LeafToken]) ([]*uc.Pair[EvalStatus, *gr.LeafToken], error) {
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

		children := make([]*uc.Pair[EvalStatus, *gr.LeafToken], 0, len(matches))

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
	tree := tr.NewTree(EvalIncomplete, rootNode)
	rootNode = gr.NewRootToken()
	treeCopy := tr.NewTree(EvalIncomplete, rootNode)

	matches, err := matchFrom(source, 0, productions)
	if err != nil {
		return tree, ers.NewErrAt(0, "position", err)
	}

	children := generateEvalTrees(matches)

	tree.SetChildren(children)

	children = generateEvalTrees(matches)

	treeCopy.SetChildren(children)

	tree.Root().ChangeStatus(EvalComplete)

	treeCopy.Root().ChangeStatus(EvalComplete)

	shouldContinue := true

	for shouldContinue {
		p := filterLeaves(source, productions)

		err := tree.ProcessLeaves(p)
		if err != nil {
			return treeCopy, err
		} else {
			treeCopy.ProcessLeaves(p)
		}

		ok := tree.Prune(filterErrorLeaves)
		if !ok {
			return treeCopy, NewErrAllMatchesFailed()
		} else {
			treeCopy.Prune(filterErrorLeaves)
		}

		leaves := tree.GetLeaves()

		shouldContinue = false
		for _, leaf := range leaves {
			status := leaf.GetStatus()

			if status == EvalIncomplete {
				shouldContinue = true
				break
			}
		}
	}

	ok := tree.Prune(filterIncompleteLeaves)
	if !ok {
		return treeCopy, NewErrAllMatchesFailed()
	}

	return tree, nil
}
