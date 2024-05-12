package Lexer

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	teval "github.com/PlayerR9/MyGoLib/Evaluations/TreeExplorer"
)

// MatchWeightFunc is a weight function that returns the length of the match.
//
// Parameters:
//   - match: The match to weigh.
//
// Returns:
//   - float64: The weight of the match.
//   - bool: True if the weight is valid, false otherwise.
func MatchWeightFunc(match *gr.MatchedResult[*gr.LeafToken]) (float64, bool) {
	return float64(len(match.Matched.Data)), true
}

// FilterEmptyBranch is a filter that filters out empty branches.
//
// Parameters:
//   - branch: The branch to filter.
//
// Returns:
//   - bool: True if the branch is not empty, false otherwise.
func FilterEmptyBranch(branch []*teval.CurrentEval[*gr.LeafToken]) bool {
	return len(branch) != 0
}
