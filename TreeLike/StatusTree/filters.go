package StatusTree

import uc "github.com/PlayerR9/MyGoLib/Units/common"

// FilterNilTree is a filter that returns true if the tree is not nil.
//
// Parameters:
//   - tree: The tree to filter.
//
// Returns:
//   - bool: True if the tree is not nil, false otherwise.
func FilterNilTree[S uc.Enumer, T any](tree *Tree[S, T]) bool {
	return tree != nil && tree.root != nil
}
