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
	if err != nil {
		return false
	}

	ok := gr.IsTerminal(rhs)
	return ok
}

// FilterNonTerminalLeaf filters non-terminal leaf nodes.
//
// Parameters:
//   - tn: The tree node to filter.
//
// Returns:
//   - bool: False if the tree node is a terminal leaf, true otherwise.
func FilterNonTerminalLeaf(tn *tr.TreeNode[*Helper]) bool {
	rhs, err := tn.Data.GetRhsAt(0)
	if err != nil {
		return true
	}

	ok := gr.IsTerminal(rhs)
	return !ok
}
