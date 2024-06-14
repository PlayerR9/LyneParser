package Lexer0

import (
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

var (
	// FilterIncompleteLeaves is a filter that filters out incomplete leaves.
	//
	// Parameters:
	//   - leaf: The leaf to filter.
	//
	// Returns:
	//   - bool: True if the leaf is incomplete, false otherwise.
	FilterIncompleteLeaves us.PredicateFilter[*CurrentEval]

	// FilterIncompleteTokens is a filter that filters out incomplete tokens.
	//
	// Parameters:
	//   - h: The helper tokens to filter.
	//
	// Returns:
	//   - bool: True if the helper tokens are incomplete, false otherwise.
	FilterIncompleteTokens us.PredicateFilter[[]*CurrentEval]

	// HelperWeightFunc is a weight function that returns the length of the helper tokens.
	//
	// Parameters:
	//   - h: The helper tokens to weigh.
	//
	// Returns:
	//   - float64: The weight of the helper tokens.
	//   - bool: True if the weight is valid, false otherwise.
	HelperWeightFunc us.WeightFunc[[]*CurrentEval]
)

func init() {
	FilterIncompleteLeaves = func(h *CurrentEval) bool {
		return h == nil || h.Status == EvalIncomplete
	}

	FilterIncompleteTokens = func(h []*CurrentEval) bool {
		return len(h) != 0 && h[len(h)-1].Status == EvalComplete
	}

	HelperWeightFunc = func(h []*CurrentEval) (float64, bool) {
		return float64(len(h)), true
	}
}
