package Tree

import slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"

// recCleanup is a helper function that removes every node in the tree rooted at n.
//
// Behaviors:
//   - This function is recursive.
func recCleanup[T any](n *treeNode[T]) {
	n.parent = nil

	for _, child := range n.children {
		recCleanup(child)
	}

	for i := range n.children {
		n.children[i] = nil
	}

	n.children = nil
}

// recSnakeTraversal is an helper function that returns all the paths
// from n to the leaves of the tree rooted at n.
//
// Returns:
//   - [][]T: A slice of slices of the values of the nodes in the paths.
//
// Behaviors:
//   - The paths are returned in the order of a BFS traversal.
//   - It is a recursive function.
func recSnakeTraversal[T any](n *treeNode[T]) [][]T {
	if len(n.children) == 0 {
		return [][]T{
			{n.Data},
		}
	}

	result := make([][]T, 0)

	for _, child := range n.children {
		subResults := recSnakeTraversal(child)

		for _, tmp := range subResults {
			result = append(result, append([]T{n.Data}, tmp...))
		}
	}

	return result
}

// recPruneFunc is an helper function that removes all the children of the
// node that satisfy the given filter including all of their children.
//
// Parameters:
//   - filter: The filter to apply.
//   - n: The node to prune.
//
// Returns:
//   - *treeNode[T]: A pointer to the highest ancestor of the pruned node.
//   - bool: True if the node satisfies the filter, false otherwise.
//
// Behaviors:
//   - This function is recursive.
func recPruneFunc[T any](filter slext.PredicateFilter[T], highest *treeNode[T], n *treeNode[T]) (*treeNode[T], bool) {
	if filter(n.Data) {
		// Delete all children
		recCleanup(n)

		return FindCommonAncestor(highest, n), true
	}

	top := 0

	for i := 0; i < len(n.children); i++ {
		high, ok := recPruneFunc(filter, highest, n.children[i])
		if ok {
			n.children[i] = nil

			highest = FindCommonAncestor(highest, high)
		} else {
			n.children[top] = n.children[i]
			top++
		}
	}

	n.children = n.children[:top]

	return highest, false
}

// recSkipFunc is an helper function that removes all the children of the
// node that satisfy the given filter without removing their children.
//
// Parameters:
//   - filter: The filter to apply.
//
// Returns:
//   - []*treeNode[T]: A slice of pointers to the children of the node.
//   - int: The total number of nodes removed.
//   - bool: True if the node satisfies the filter, false otherwise.
func recSkipFunc[T any](filter slext.PredicateFilter[T], n *treeNode[T]) ([]*treeNode[T], int, bool) {
	// 1. Check if the children satisfy the filter
	newChildren := make([]*treeNode[T], 0)

	total := 0

	for i := 0; i < len(n.children); i++ {
		sub, amount, ok := recSkipFunc(filter, n.children[i])
		if ok {
			n.children[i] = nil
			newChildren = append(newChildren, sub...)
		} else {
			newChildren = append(newChildren, n.children[i])
		}

		total += amount

		newChildren = append(newChildren, sub...)
	}

	n.children = newChildren

	// 2. Check if the node satisfies the filter
	if filter(n.Data) {
		n.parent = nil

		return n.children, total + 1, true
	}

	// 3. Update the parent of the children
	for i := 0; i < len(n.children); i++ {
		n.children[i].parent = n
	}

	return nil, total, false
}
