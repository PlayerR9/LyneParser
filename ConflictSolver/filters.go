package ConflictSolver

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
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
