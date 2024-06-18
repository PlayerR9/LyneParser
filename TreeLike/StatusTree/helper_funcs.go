package StatusTree

import (
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// recCleanup is a helper function that removes every node in the tree rooted at n.
//
// Behaviors:
//   - This function is recursive.
func recCleanup[S uc.Enumer, T any](n *TreeNode[S, T]) {
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
func recSnakeTraversal[S uc.Enumer, T any](n *TreeNode[S, T]) [][]*uc.Pair[S, T] {
	if len(n.children) == 0 {
		p := uc.NewPair(n.status, n.Data)

		return [][]*uc.Pair[S, T]{
			{p},
		}
	}

	var result [][]*uc.Pair[S, T]

	for _, child := range n.children {
		subResults := recSnakeTraversal(child)

		p := uc.NewPair(n.status, n.Data)

		for _, tmp := range subResults {
			result = append(result, append([]*uc.Pair[S, T]{p}, tmp...))
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
//   - *TreeNode[S, T]: A pointer to the highest ancestor of the pruned node.
//   - bool: True if the node satisfies the filter, false otherwise.
//
// Behaviors:
//   - This function is recursive.
func recPruneFunc[S uc.Enumer, T any](filter us.PredicateFilter[*uc.Pair[S, T]], highest *TreeNode[S, T], n *TreeNode[S, T]) (*TreeNode[S, T], bool) {
	p := uc.NewPair(n.status, n.Data)

	ok := filter(p)
	if ok {
		// Delete all children
		recCleanup(n)

		nodes := FindCommonAncestor(highest, n)

		return nodes, true
	}

	var top int

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
