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

// FilterNiltreeNode is a filter that returns true if the node is not nil.
//
// Parameters:
//   - node: The node to filter.
//
// Returns:
//   - bool: True if the node is not nil, false otherwise.
func FilterNiltreeNode[T any](node *treeNode[T]) bool {
	return node != nil
}
