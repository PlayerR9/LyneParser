package Lexer

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	ers "github.com/PlayerR9/MyGoLib/Units/errors"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// generateEvalTrees adds the matches to a root tree as leaves.
//
// Parameters:
//   - root: The root of the tree to add the leaves to.
//   - matches: The matches to add to the tree evaluator.
func generateEvalTrees(matches []*gr.MatchedResult[*gr.LeafToken]) []*tr.Tree[*CurrentEval] {
	// Get the longest match.
	matches = selectBestMatches(matches)

	children := make([]*tr.Tree[*CurrentEval], 0, len(matches))

	for _, match := range matches {
		currMatch := match.GetMatch()
		ht := newCurrentEval(currMatch)

		tree := tr.NewTree(ht)

		children = append(children, tree)
	}

	return children
}

// filterLeaves processes the leaves in the tree evaluator.
//
// Returns:
//   - bool: True if all leaves are complete, false otherwise.
//   - error: An error of type *ErrAllMatchesFailed if all matches failed.
func filterLeaves(source *cds.Stream[byte], productions []*gr.RegProduction) uc.EvalManyFunc[*CurrentEval, *CurrentEval] {
	filterFunc := func(data *CurrentEval) ([]*CurrentEval, error) {
		nextAt := data.elem.GetPos() + len(data.elem.Data)

		ok := source.IsDone(nextAt, 1)
		if ok {
			data.changeStatus(EvalComplete)
			return nil, nil
		}

		matches, err := matchFrom(source, nextAt, productions)
		if err != nil {
			data.changeStatus(EvalError)
			return nil, nil
		}

		// Get the longest match.
		matches = selectBestMatches(matches)

		children := make([]*CurrentEval, 0, len(matches))

		for _, match := range matches {
			curr := match.GetMatch()
			ht := newCurrentEval(curr)

			children = append(children, ht)
		}

		data.changeStatus(EvalComplete)

		return children, nil
	}

	return filterFunc
}

// REMOVE THIS ONCE MyGoLib is updated
// pruneTree prunes the tree evaluator.
//
// Parameters:
//   - filter: The filter to use to prune the tree.
//
// Returns:
//   - bool: True if no nodes were pruned, false otherwise.
func pruneTree(root *tr.Tree[*CurrentEval], filter us.PredicateFilter[*CurrentEval]) bool {
	for root.Size() != 0 {
		target := root.SearchNodes(filter)
		if target == nil {
			return true
		}

		root.DeleteBranchContaining(target)
	}

	return false
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
func executeLexing(source *cds.Stream[byte], productions []*gr.RegProduction) (*tr.Tree[*CurrentEval], error) {
	rootNode := gr.NewRootToken()
	ce := newCurrentEval(rootNode)
	tree := tr.NewTree(ce)

	matches, err := matchFrom(source, 0, productions)
	if err != nil {
		return tree, ers.NewErrAt(0, "position", err)
	}

	children := generateEvalTrees(matches)

	tree.SetChildren(children)

	tree.Root().Data.changeStatus(EvalComplete)

	shouldContinue := true

	for shouldContinue {
		p := filterLeaves(source, productions)

		err := tree.ProcessLeaves(p)
		if err != nil {
			return tree, err
		}

		ok := pruneTree(tree, filterErrorLeaves)
		if !ok {
			return tree, NewErrAllMatchesFailed()
		}

		leaves := tree.GetLeaves()

		shouldContinue = false
		for _, leaf := range leaves {
			if leaf.Data.status == EvalIncomplete {
				shouldContinue = true
				break
			}
		}
	}

	ok := pruneTree(tree, filterIncompleteLeaves)
	if !ok {
		return tree, NewErrAllMatchesFailed()
	}

	return tree, nil
}
