package Lexer0

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	ers "github.com/PlayerR9/MyGoLib/Units/errors"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// TreeEvaluator is a tree evaluator that uses a grammar to tokenize a string.
type TreeEvaluator struct {
	// root is the root node of the tree evaluator.
	root *tr.Tree[*CurrentEval]

	// source is the source to match.
	source *cds.Stream[byte]
}

// addMatchLeaves adds the matches to a root tree as leaves.
//
// Parameters:
//   - root: The root of the tree to add the leaves to.
//   - matches: The matches to add to the tree evaluator.
func (te *TreeEvaluator) addMatchLeaves(root *tr.Tree[*CurrentEval], matches []*gr.MatchedResult[*gr.LeafToken]) {
	// Get the longest match.
	weights := us.ApplyWeightFunc(matches, MatchWeightFunc)
	pairs := us.FilterByPositiveWeight(weights)

	matches = us.ExtractResults(pairs)

	children := make([]*tr.Tree[*CurrentEval], 0, len(matches))

	for _, match := range matches {
		ht := NewCurrentEval(match.GetMatch())
		children = append(children, tr.NewTree(ht))
	}

	root.SetChildren(children)
}

// processLeaves processes the leaves in the tree evaluator.
//
// Returns:
//   - bool: True if all leaves are complete, false otherwise.
//   - error: An error of type *ErrAllMatchesFailed if all matches failed.
func (te *TreeEvaluator) processLeaves() uc.EvalManyFunc[*CurrentEval, *CurrentEval] {
	filterFunc := func(data *CurrentEval) ([]*CurrentEval, error) {
		nextAt := data.Elem.GetPos() + len(data.Elem.Data)

		if te.source.IsDone(nextAt, 1) {
			data.SetStatus(EvalComplete)

			return nil, nil
		}

		matches, err := MatchFrom(te.source, nextAt, Productions)
		if err != nil {
			data.SetStatus(EvalError)

			return nil, nil
		}

		// Get the longest match.
		weights := us.ApplyWeightFunc(matches, MatchWeightFunc)
		pairs := us.FilterByPositiveWeight(weights)

		matches = us.ExtractResults(pairs)

		children := make([]*CurrentEval, 0, len(matches))

		for _, match := range matches {
			ht := NewCurrentEval(match.GetMatch())
			children = append(children, ht)
		}

		data.SetStatus(EvalComplete)

		return children, nil
	}

	return filterFunc
}

// canContinue returns true if the tree evaluator can continue.
//
// Returns:
//   - bool: True if the tree evaluator can continue, false otherwise.
func (te *TreeEvaluator) canContinue() bool {
	for _, leaf := range te.root.GetLeaves() {
		if leaf.Data.Status == EvalIncomplete {
			return true
		}
	}

	return false
}

// GetBranches returns the tokens that have been lexed.
//
// Remember to use Lexer.RemoveToSkipTokens() to remove tokens that
// are not needed for the parser (i.e., marked as to skip in the grammar).
//
// Returns:
//   - result: The tokens that have been lexed.
//   - reason: An error if the tree evaluator has not been run yet.
func (te *TreeEvaluator) GetBranches() ([]*cds.Stream[*gr.LeafToken], error) {
	if te.root == nil {
		return nil, ers.NewErrInvalidUsage(
			ers.NewErrNilValue(),
			"must call TreeEvaluator.Evaluate() first",
		)
	}

	tokenBranches := te.root.SnakeTraversal()

	branches, invalidTokIndex := filterInvalidBranches(tokenBranches)
	if invalidTokIndex != -1 {
		var result []*cds.Stream[*gr.LeafToken]
		for _, branch := range branches {
			conv := convertBranchToTokenStream(branch)

			result = append(result, conv)
		}

		return result, ers.NewErrAt(invalidTokIndex, "token", NewErrInvalidElement())
	}

	branches, err := RemoveToSkipTokens(branches)
	if err != nil {
		var result []*cds.Stream[*gr.LeafToken]
		for _, branch := range branches {
			conv := convertBranchToTokenStream(branch)

			result = append(result, conv)
		}

		return result, err
	}

	te.root = nil

	var result []*cds.Stream[*gr.LeafToken]
	for _, branch := range branches {
		conv := convertBranchToTokenStream(branch)

		result = append(result, conv)
	}

	return result, nil
}
