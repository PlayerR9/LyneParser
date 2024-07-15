package ConflictSolver

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	tn "github.com/PlayerR9/treenode"
)

// FilterTerminalLeaf filters out terminal leaf nodes.
//
// Parameters:
//   - tn: The tree node to filter.
//
// Returns:
//   - bool: True if the tree node is a terminal leaf, false otherwise.
func FilterTerminalLeaf[T gr.TokenTyper](tn *HelperNode[T]) bool {
	rhs, err := tn.GetRhsAt(0)
	if err != nil {
		return false
	}

	ok := rhs.IsTerminal()
	return ok
}

// FilterNonTerminalLeaf filters out non-terminal leaf nodes.
//
// Parameters:
//   - tn: The tree node to filter.
//
// Returns:
//   - bool: False if the tree node is a terminal leaf, true otherwise.
func FilterNonTerminalLeaf[T gr.TokenTyper](n tn.Noder) bool {
	tn, ok := n.(*HelperNode[T])
	if !ok {
		return false
	}

	rhs, err := tn.GetRhsAt(0)
	if err != nil {
		return true
	}

	ok = rhs.IsTerminal()
	return !ok
}
