package TreeExplorer

import (
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// FilterIncompleteLeaves is a filter that filters out incomplete leaves.
//
// Parameters:
//   - leaf: The leaf to filter.
//
// Returns:
//   - bool: True if the leaf is incomplete, false otherwise.
func FilterIncompleteLeaves[O any](h *uc.Pair[EvalStatus, O]) bool {
	if h == nil {
		return true
	}

	return h.First == EvalIncomplete
}

// FilterCompleteTokens is a filter that filters complete helper tokens.
//
// Parameters:
//   - h: The helper tokens to filter.
//
// Returns:
//   - bool: True if the helper tokens are incomplete, false otherwise.
func FilterCompleteTokens[O any](h []*uc.Pair[EvalStatus, O]) bool {
	if len(h) == 0 {
		return false
	}

	status := h[len(h)-1].First

	return status == EvalComplete
}

// HelperWeightFunc is a weight function that returns the length of the helper tokens.
//
// Parameters:
//   - h: The helper tokens to weigh.
//
// Returns:
//   - float64: The weight of the helper tokens.
//   - bool: True if the weight is valid, false otherwise.
func HelperWeightFunc[O any](h []*uc.Pair[EvalStatus, O]) (float64, bool) {
	return float64(len(h)), true
}

// FilterErrorLeaves is a filter that filters out leaves that are in error.
//
// Parameters:
//   - leaf: The leaf to filter.
//
// Returns:
//   - bool: True if the leaf is in error, false otherwise.
func FilterErrorLeaves[O any](h *uc.Pair[EvalStatus, O]) bool {
	if h == nil {
		return true
	}

	return h.First == EvalError
}
