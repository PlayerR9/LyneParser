package ConflictSolver

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// FilterTerminalLeaf filters out terminal leaf nodes.
//
// Parameters:
//   - tn: The tree node to filter.
//
// Returns:
//   - bool: True if the tree node is a terminal leaf, false otherwise.
func FilterTerminalLeaf[T uc.Enumer](tn *tr.TreeNode[*Helper[T]]) bool {
	rhs, err := tn.Data.GetRhsAt(0)
	if err != nil {
		return false
	}

	ok := gr.IsTerminal(rhs.String())
	return ok
}

// FilterNonTerminalLeaf filters out non-terminal leaf nodes.
//
// Parameters:
//   - tn: The tree node to filter.
//
// Returns:
//   - bool: False if the tree node is a terminal leaf, true otherwise.
func FilterNonTerminalLeaf[T uc.Enumer](tn *tr.TreeNode[*Helper[T]]) bool {
	rhs, err := tn.Data.GetRhsAt(0)
	if err != nil {
		return true
	}

	ok := gr.IsTerminal(rhs.String())
	return !ok
}
