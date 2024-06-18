package Tree

// FilterNilTree is a filter that returns true if the tree is not nil.
//
// Parameters:
//   - tree: The tree to filter.
//
// Returns:
//   - bool: True if the tree is not nil, false otherwise.
func FilterNilTree[T any](tree *Tree[T]) bool {
	return tree != nil && tree.root != nil
}
