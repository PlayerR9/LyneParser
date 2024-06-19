package StatusTree

import (
	fsp "github.com/PlayerR9/MyGoLib/Formatting/FString"
	"github.com/PlayerR9/MyGoLib/ListLike/Queuer"
	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// Tree is a generic data structure that represents a tree.
type Tree[S uc.Enumer, T any] struct {
	// root is the root of the tree.
	root *TreeNode[S, T]

	// leaves is the leaves of the tree.
	leaves []*TreeNode[S, T]

	// size is the number of nodes in the tree.
	size int
}

// FString returns a formatted string representation of the tree.
//
// Parameters:
//   - indentLevel: The level of indentation to use.
//
// Returns:
//   - []string: A slice of strings that represent the tree.
func (t *Tree[S, T]) FString(trav *fsp.Traversor) error {
	if t.root == nil {
		return nil
	}

	form := fsp.NewFormatter(
		fsp.NewIndentConfig("| ", 0),
	)

	err := fsp.ApplyForm(form, trav, t.root)
	if err != nil {
		return err
	}

	return nil
}

// NewTree creates a new tree with the given root.
//
// Parameters:
//   - data: The value of the root.
//
// Returns:
//   - *Tree[S, T]: A pointer to the newly created tree.
func NewTree[S uc.Enumer, T any](status S, data T) *Tree[S, T] {
	root := NewTreeNode(status, data)

	return &Tree[S, T]{
		root:   root,
		leaves: []*TreeNode[S, T]{root},
		size:   1,
	}
}

// SetChildren sets the children of the root of the tree.
//
// Parameters:
//   - children: The children to set.
//
// Returns:
//   - error: An error of type *ErrMissingRoot if the tree does not have a root.
func (t *Tree[S, T]) SetChildren(children []*Tree[S, T]) error {
	children = us.SliceFilter(children, FilterNilTree)
	if len(children) == 0 {
		return nil
	}

	root := t.root

	if root == nil {
		return NewErrMissingRoot()
	}

	t.leaves = make([]*TreeNode[S, T], 0)
	t.size = 1

	root.children = make([]*TreeNode[S, T], 0, len(children))

	for _, child := range children {
		t.leaves = append(t.leaves, child.leaves...)
		t.size += child.size

		child.root.setParent(root)
		root.children = append(root.children, child.root)
	}

	return nil
}

// IsSingleton returns true if the tree has only one node.
//
// Returns:
//   - bool: True if the tree has only one node, false otherwise.
func (t *Tree[S, T]) IsSingleton() bool {
	return t.size == 1
}

// Size returns the number of nodes in the tree.
//
// Returns:
//   - int: The number of nodes in the tree.
func (t *Tree[S, T]) Size() int {
	return t.size
}

// Root returns the root of the tree.
//
// Returns:
//   - *Node[T]: A pointer to the root of the tree.
func (t *Tree[S, T]) Root() *TreeNode[S, T] {
	return t.root
}

// GetChildren returns all the children of the tree in a DFS order.
//
// Returns:
//   - []T: A slice of the values of the nodes in the tree.
//
// Behaviors:
//   - The root is the first element in the slice.
//   - If the tree does not have a root, it returns nil.
func (t *Tree[S, T]) GetChildren() []*uc.Pair[S, T] {
	root := t.root

	if root == nil {
		return nil
	}

	var children []*uc.Pair[S, T]

	St := Stacker.NewLinkedStack(root)

	for {
		node, ok := St.Pop()
		if !ok {
			break
		}

		p := uc.NewPair(node.status, node.Data)

		children = append(children, p)

		for _, child := range node.children {
			St.Push(child)
		}
	}

	return children
}

// Cleanup removes every node in the tree.
//
// Behaviors:
//   - This function is recursive and so, it is expensive.
func (t *Tree[S, T]) Cleanup() {
	root := t.root

	recCleanup(root)

	t.root = nil
}

// GetLeaves returns all the leaves of the tree.
//
// Returns:
//   - []*Node[T]: A slice of pointers to the leaves of the tree.
//
// Behaviors:
//   - Always returns at least one leaf.
//   - It returns the leaves that are stored in the tree. Make sure to call
//     any update function before calling this function if the tree has been modified
//     unexpectedly.
func (t *Tree[S, T]) GetLeaves() []*TreeNode[S, T] {
	return t.leaves
}

// RegenerateLeaves regenerates the leaves of the tree and returns them.
//
// Returns:
//   - []*Node[T]: A slice of pointers to the leaves of the tree.
//
// Behaviors:
//   - The leaves are updated in a DFS order.
//   - Expensive operation; use it only when necessary (i.e., leaves changed unexpectedly.)
//   - This also updates the size of the tree.
func (t *Tree[S, T]) RegenerateLeaves() []*TreeNode[S, T] {
	root := t.root

	if root == nil {
		t.leaves = make([]*TreeNode[S, T], 0)
		t.size = 0

		return t.leaves
	}

	leaves := make([]*TreeNode[S, T], 0)

	St := Stacker.NewLinkedStack(root)

	t.size = 0

	for {
		top, ok := St.Pop()
		if !ok {
			break
		}

		t.size++

		if len(top.children) == 0 {
			leaves = append(leaves, top)
		} else {
			for _, child := range top.children {
				St.Push(child)
			}
		}
	}

	t.leaves = leaves

	return leaves
}

// UpdateLeaves updates the leaves of the tree and returns them.
//
// Returns:
//   - []*Node[T]: A slice of pointers to the leaves of the tree.
//
// Behaviors:
//   - The leaves are updated in a DFS order.
//   - Less expensive than RegenerateLeaves. However, if nodes has been deleted
//     from the tree, this may give unexpected results.
//   - This also updates the size of the tree.
func (t *Tree[S, T]) UpdateLeaves() []*TreeNode[S, T] {
	if len(t.leaves) == 0 {
		t.size = 0

		return t.leaves
	}

	newLeaves := make([]*TreeNode[S, T], 0)

	St := Stacker.NewLinkedStack(t.leaves...)

	t.size -= len(t.leaves)

	for {
		top, ok := St.Pop()
		if !ok {
			break
		}

		t.size++

		if len(top.children) == 0 {
			newLeaves = append(newLeaves, top)
		} else {
			for _, child := range top.children {
				St.Push(child)
			}
		}
	}

	t.leaves = newLeaves

	return newLeaves
}

// SnakeTraversal returns all the paths from the root to the leaves of the tree.
//
// Returns:
//   - [][]T: A slice of slices of the values of the nodes in the paths.
//
// Behaviors:
//   - The paths are returned in the order of a DFS traversal.
//   - If the tree is empty, it returns an empty slice.
func (t *Tree[S, T]) SnakeTraversal() [][]*uc.Pair[S, T] {
	root := t.root

	trav := recSnakeTraversal(root)
	return trav
}

// HasChild returns true if the tree has the given child in any of its nodes
// in a BFS order.
//
// Parameters:
//   - filter: The filter to apply.
//
// Returns:
//   - bool: True if the tree has the child, false otherwise.
func (t *Tree[S, T]) HasChild(filter us.PredicateFilter[*uc.Pair[S, T]]) bool {
	if filter == nil {
		return false
	}

	root := t.root
	if root == nil {
		return false
	}

	Q := Queuer.NewLinkedQueue(root)

	for {
		node, ok := Q.Dequeue()
		if !ok {
			break
		}

		p := uc.NewPair(node.status, node.Data)

		ok = filter(p)
		if ok {
			return true
		}

		for _, child := range node.children {
			Q.Enqueue(child)
		}
	}

	return false
}

// FilterChildren returns all the children of the tree that satisfy the given filter
// in a BFS order.
//
// Parameters:
//   - filter: The filter to apply.
//
// Returns:
//   - []*Node[T]: A slice of pointers to the children of the node.
func (t *Tree[S, T]) FilterChildren(filter us.PredicateFilter[*uc.Pair[S, T]]) []*TreeNode[S, T] {
	if filter == nil {
		return nil
	}

	root := t.root
	if root == nil {
		return nil
	}

	Q := Queuer.NewLinkedQueue(root)

	solutions := make([]*TreeNode[S, T], 0)

	for {
		node, ok := Q.Dequeue()
		if !ok {
			break
		}

		p := uc.NewPair(node.status, node.Data)

		ok = filter(p)
		if ok {
			solutions = append(solutions, node)
		}

		for _, child := range node.children {
			Q.Enqueue(child)
		}
	}

	return solutions
}

// PruneBranches removes all the children of the node that satisfy the given filter.
// The filter is a function that takes the value of a node and returns a boolean.
// If the filter returns true for a child, the child is removed along with its children.
//
// Parameters:
//   - filter: The filter to apply.
//
// Returns:
//   - bool: True if the whole tree can be deleted, false otherwise.
//
// Behaviors:
//   - If the root satisfies the filter, the tree is cleaned up.
//   - It is a recursive function.
func (t *Tree[S, T]) PruneBranches(filter us.PredicateFilter[*uc.Pair[S, T]]) bool {
	if filter == nil {
		return false
	}

	root := t.root

	if root == nil {
		return false
	}

	highest, ok := recPruneFunc(filter, nil, root)
	if ok {
		return true
	}

	t.leaves = highest.Leaves()
	t.size = highest.Size()

	return false
}

// SkipFunc removes all the children of the tree that satisfy the given filter
// without removing any of their children. Useful for removing unwanted nodes from the tree.
//
// Parameters:
//   - filter: The filter to apply.
//
// Returns:
//   - []*Tree[S, T]: A slice of pointers to the trees obtained after removing the nodes.
//
// Behaviors:
//   - If this function returns only one tree, this is the updated tree. But, if
//     it returns more than one tree, then we have deleted the root of the tree and
//     obtained a forest.
func (t *Tree[S, T]) SkipFilter(filter us.PredicateFilter[*uc.Pair[S, T]]) []*Tree[S, T] {
	frontier := make([]*TreeNode[S, T], len(t.leaves))
	copy(frontier, t.leaves)

	seen := make(map[*TreeNode[S, T]]bool)
	newLeaves := make([]*TreeNode[S, T], 0)

	forest := make([]*Tree[S, T], 0)

	for len(frontier) > 0 {
		leaf := frontier[0]
		seen[leaf] = true

		// Remove any node that has been seen from the frontier.
		frontier = us.SliceFilter(frontier, func(n *TreeNode[S, T]) bool {
			return !seen[n]
		})

		p := uc.NewPair(leaf.status, leaf.Data)

		ok := filter(p)
		if !ok {
			if leaf.parent == nil {
				// We reached the root
				frontier = frontier[1:]
			} else {
				if len(leaf.children) == 0 {
					newLeaves = append(newLeaves, leaf)
				}

				if !seen[leaf.parent] {
					frontier[0] = leaf.parent
				} else {
					frontier = frontier[1:]
				}
			}
		} else {
			children := leaf.removeNode()

			if len(children) != 0 {
				// We obtained a forest as we reached the root

				for _, child := range children {
					forest = append(forest, child.ToTree())
				}

				// We reached the root
				frontier = frontier[1:]
			} else {
				if !seen[leaf.parent] {
					frontier[0] = leaf.parent
				} else {
					frontier = frontier[1:]
				}

				t.size--
			}
		}
	}

	if len(forest) == 0 {
		t.leaves = newLeaves

		return []*Tree[S, T]{t}
	} else {
		return forest
	}
}

// replaceLeafWithTree is a helper function that replaces a leaf with a tree.
//
// Parameters:
//   - at: The index of the leaf to replace.
//   - children: The children of the leaf.
//
// Behaviors:
//   - The leaf is replaced with the children.
//   - The size of the tree is updated.
func (t *Tree[S, T]) replaceLeafWithTree(at int, children []*uc.Pair[S, T]) {
	leaf := t.leaves[at]

	// Make the subtree
	leaf.children = make([]*TreeNode[S, T], 0, len(children))
	for _, child := range children {
		node := NewTreeNode(child.First, child.Second)
		node.parent = leaf

		leaf.children = append(leaf.children, node)
	}

	// Update the size of the tree
	t.size += len(leaf.children)

	// Replace the current leaf with the leaf's children
	if at == len(t.leaves)-1 {
		t.leaves = append(t.leaves[:at], leaf.children...)
	} else if at == 0 {
		t.leaves = append(leaf.children, t.leaves[at+1:]...)
	} else {
		t.leaves = append(t.leaves[:at], append(leaf.children, t.leaves[at+1:]...)...)
	}
}

// ProcessLeaves applies the given function to the leaves of the tree and replaces
// the leaves with the children returned by the function.
//
// Parameters:
//   - f: The function to apply to the leaves.
//
// Returns:
//   - error: An error returned by the function.
//
// Behaviors:
//   - The function is applied to the leaves in order.
//   - The function must return a slice of values of type T.
//   - If the function returns an error, the process stops and the error is returned.
//   - The leaves are replaced with the children returned by the function.
func (t *Tree[S, T]) ProcessLeaves(f uc.EvalManyFunc[*TreeNode[S, T], *uc.Pair[S, T]]) error {
	for i, leaf := range t.leaves {
		children, err := f(leaf)
		if err != nil {
			return err
		}

		if len(children) != 0 {
			t.replaceLeafWithTree(i, children)
		}
	}

	return nil
}

// SearchNodes searches for the first node that satisfies the given filter in a BFS order.
//
// Parameters:
//   - f: The filter to apply.
//
// Returns:
//   - *treeNode[T]: A pointer to the node that satisfies the filter.
func (t *Tree[S, T]) SearchNodes(f us.PredicateFilter[*uc.Pair[S, T]]) *TreeNode[S, T] {
	root := t.root

	Q := Queuer.NewLinkedQueue(root)

	for {
		first, ok := Q.Dequeue()
		if !ok {
			break
		}

		p := uc.NewPair(first.status, first.Data)

		ok = f(p)
		if ok {
			return first
		}

		for _, child := range first.children {
			Q.Enqueue(child)
		}
	}

	return nil
}

// FindBranchingPoint returns the first node in the path from n to the root
// such that has more than one sibling.
//
// Returns:
//   - *treeNode[T]: A pointer to the one node before the branching point.
//   - *treeNode[T]: A pointer to the branching point.
//   - bool: True if the node has a branching point, false otherwise.
//
// Behaviors:
//   - If there is no branching point, it returns the root of the tree.
func (tn *TreeNode[S, T]) FindBranchingPoint() (*TreeNode[S, T], *TreeNode[S, T], bool) {
	if tn.parent == nil {
		return nil, tn, false
	}

	node := tn
	parent := tn.parent

	hasBranchingPoint := false

	for parent.parent != nil {
		if len(parent.children) > 1 {
			hasBranchingPoint = true
			break
		}

		node = parent
		parent = parent.parent

		if node.parent != parent {
			panic("Parent is not the parent of the node")
		}
	}

	return node, parent, hasBranchingPoint
}

// DeleteBranchContaining deletes the branch containing the given node.
//
// Parameters:
//   - tn: The node to delete.
//
// Returns:
//   - error: An error if the node is not a part of the tree.
func (t *Tree[S, T]) DeleteBranchContaining(tn *TreeNode[S, T]) error {
	child, parent, hasBranching := tn.FindBranchingPoint()
	if !hasBranching {
		root := t.root

		if parent != root {
			return NewErrNodeNotPartOfTree()
		}

		t.Cleanup()
	}

	children := parent.DeleteChild(child)

	for _, child := range children {
		recCleanup(child)
	}

	t.leaves = t.RegenerateLeaves()

	return nil
}

// GetDirectChildren returns the direct children of the root of the tree.
//
// Returns:
//   - []*Node[T]: A slice of pointers to the direct children of the root.
func (t *Tree[S, T]) GetDirectChildren() []*TreeNode[S, T] {
	root := t.root
	if root == nil {
		return nil
	}

	return root.children
}

// PruneTree prunes the tree using the given filter.
//
// Parameters:
//   - filter: The filter to use to prune the tree.
//
// Returns:
//   - bool: True if no nodes were pruned, false otherwise.
func (t *Tree[S, T]) Prune(filter us.PredicateFilter[*uc.Pair[S, T]]) bool {
	for t.Size() != 0 {
		target := t.SearchNodes(filter)
		if target == nil {
			return true
		}

		err := t.DeleteBranchContaining(target)
		if err != nil {
			panic(err)
		}
	}

	return false
}
