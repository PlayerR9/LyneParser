package ConflictSolver

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	hlp "github.com/PlayerR9/MyGoLib/CustomData/Helpers"
	tr "github.com/PlayerR9/MyGoLib/CustomData/Tree"
)

// FilterTerminalLeaf filters out terminal leaf nodes.
//
// Parameters:
//   - tn: The tree node to filter.
//
// Returns:
//   - bool: True if the tree node is a terminal leaf, false otherwise.
func FilterTerminalLeaf(tn *tr.TreeNode[*Helper]) bool {
	rhs, err := tn.Data.GetRhsAt(0)
	return err != nil && !gr.IsTerminal(rhs)
}

// HResultWeightFunc is a function that, given a HResult, returns the weight of the result.
//
// Parameters:
//   - h: The HResult to calculate the weight of.
//
// Returns:
//   - float64: The weight of the result.
//   - bool: True if the weight is valid, false otherwise.
func HResultWeightFunc(h hlp.HResult[Actioner]) (float64, bool) {
	return float64(h.First.Size()), true
}

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
