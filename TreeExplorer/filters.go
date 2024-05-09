package TreeExplorer

// FilterIncompleteLeaves is a filter that filters out incomplete leaves.
//
// Parameters:
//   - leaf: The leaf to filter.
//
// Returns:
//   - bool: True if the leaf is incomplete, false otherwise.
func FilterIncompleteLeaves[O any](h *CurrentEval[O]) bool {
	return h == nil || h.Status == EvalIncomplete
}

// FilterIncompleteTokens is a filter that filters out incomplete tokens.
//
// Parameters:
//   - h: The helper tokens to filter.
//
// Returns:
//   - bool: True if the helper tokens are incomplete, false otherwise.
func FilterIncompleteTokens[O any](h []*CurrentEval[O]) bool {
	return len(h) != 0 && h[len(h)-1].Status == EvalComplete
}

// HelperWeightFunc is a weight function that returns the length of the helper tokens.
//
// Parameters:
//   - h: The helper tokens to weigh.
//
// Returns:
//   - float64: The weight of the helper tokens.
//   - bool: True if the weight is valid, false otherwise.
func HelperWeightFunc[O any](h []*CurrentEval[O]) (float64, bool) {
	return float64(len(h)), true
}
