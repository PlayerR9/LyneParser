package ConflictSolver

/////////////////////////////////////////////////////////////

// FilterNonShiftHelper filters out non-shift helpers.
//
// Parameters:
//   - h: The helper to filter.
//
// Returns:
//   - bool: True if the helper is a shift helper, false otherwise.
func FilterNonShiftHelper(h *Helper) bool {
	return h != nil && h.IsShift()
}
