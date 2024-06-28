package ConflictSolver

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

var (
	// FilterTerminalLeaf filters out terminal leaf nodes.
	//
	// Parameters:
	//   - tn: The tree node to filter.
	//
	// Returns:
	//   - bool: True if the tree node is a terminal leaf, false otherwise.
	FilterTerminalLeaf us.PredicateFilter[*tr.TreeNode[*Helper]]

	// FilterNonTerminalLeaf filters out non-terminal leaf nodes.
	//
	// Parameters:
	//   - tn: The tree node to filter.
	//
	// Returns:
	//   - bool: False if the tree node is a terminal leaf, true otherwise.
	FilterNonTerminalLeaf us.PredicateFilter[*tr.TreeNode[*Helper]]
)

func init() {
	FilterTerminalLeaf = func(tn *tr.TreeNode[*Helper]) bool {
		rhs, err := tn.Data.GetRhsAt(0)
		if err != nil {
			return false
		}

		ok := gr.IsTerminal(rhs)
		return ok
	}

	FilterNonTerminalLeaf = func(tn *tr.TreeNode[*Helper]) bool {
		rhs, err := tn.Data.GetRhsAt(0)
		if err != nil {
			return true
		}

		ok := gr.IsTerminal(rhs)
		return !ok
	}
}
