package Lexer

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"

	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
)

// FilterIncompleteLeaves is a filter that filters out incomplete leaves.
//
// Parameters:
//   - leaf: The leaf to filter.
//
// Returns:
//   - bool: True if the leaf is incomplete, false otherwise.
func FilterIncompleteLeaves(h *helperToken) bool {
	return h == nil || h.Status == TkIncomplete
}

// FilterErrorLeaves is a filter that filters out leaves that are in error.
//
// Parameters:
//   - leaf: The leaf to filter.
//
// Returns:
//   - bool: True if the leaf is in error, false otherwise.
func FilterErrorLeaves(h *helperToken) bool {
	return h == nil || h.Status == TkError
}

// FilterEmptyTokenStream is a filter that filters out empty token streams.
//
// Parameters:
//   - branch: The token stream to filter.
//
// Returns:
//   - bool: True if the token stream is empty, false otherwise.
func FilterEmptyTokenStream(branch *cds.Stream[*gr.LeafToken]) bool {
	return branch.IsEmpty()
}

// MatchWeightFunc is a weight function that returns the length of the match.
//
// Parameters:
//   - match: The match to weigh.
//
// Returns:
//   - float64: The weight of the match.
//   - bool: True if the weight is valid, false otherwise.
func MatchWeightFunc(match gr.MatchedResult[*gr.LeafToken]) (float64, bool) {
	return float64(len(match.Matched.Data)), true
}

// FilterIncompleteTokens is a filter that filters out incomplete tokens.
//
// Parameters:
//   - h: The helper tokens to filter.
//
// Returns:
//   - bool: True if the helper tokens are incomplete, false otherwise.
func FilterIncompleteTokens(h []*helperToken) bool {
	return len(h) != 0 && h[len(h)-1].Status == TkComplete
}

// FilterEmptyBranch is a filter that filters out empty branches.
//
// Parameters:
//   - branch: The branch to filter.
//
// Returns:
//   - bool: True if the branch is not empty, false otherwise.
func FilterEmptyBranch(branch []*helperToken) bool {
	return len(branch) != 0
}

// HelperWeightFunc is a weight function that returns the length of the helper tokens.
//
// Parameters:
//   - h: The helper tokens to weigh.
//
// Returns:
//   - float64: The weight of the helper tokens.
//   - bool: True if the weight is valid, false otherwise.
func HelperWeightFunc(h []*helperToken) (float64, bool) {
	return float64(len(h)), true
}
