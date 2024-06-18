package Tree

import (
	"slices"

	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// TreeNode is a generic data structure that represents a node in a tree.
type TreeNode[T any] struct {
	// Data is the value of the node.
	Data T

	// parent is the parent of the node.
	parent *TreeNode[T]

	// children is the children of the node.
	children []*TreeNode[T]
}

// Equals implements common.Objecter.
func (t *TreeNode[T]) Equals(other uc.Objecter) bool {
	panic("unimplemented")
}

// String implements common.Objecter.
func (t *TreeNode[T]) String() string {
	panic("unimplemented")
}

// FString returns a formatted string representation of the node.
//
// Parameters:
//   - indentLevel: The level of indentation to use.
//
// Returns:
//   - []string: A slice of strings that represent the node.
func (t *TreeNode[T]) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	err := trav.AddLine(uc.StringOf(t.Data))
	if err != nil {
		return err
	}

	err = ffs.ApplyFormMany(
		trav.GetConfig(
			ffs.WithModifiedIndent(1),
		),
		trav,
		t.children,
	)
	if err != nil {
		return err
	}

	return nil
}

// Copy returns a deep copy of the node.
//
// Returns:
//   - *Node[T]: A pointer to the deep copy of the node.
//
// Behaviors:
//   - This function is recursive.
//   - The parent is not copied.
//   - The data is shallow copied.
func (n *TreeNode[T]) Copy() uc.Copier {
	node := &TreeNode[T]{
		Data:     n.Data,
		parent:   nil,
		children: make([]*TreeNode[T], 0, len(n.children)),
	}

	for _, child := range n.children {
		node.children = append(node.children, child.Copy().(*TreeNode[T]))
	}

	return node
}

// NewTreeNode creates a new node with the given data.
//
// Parameters:
//   - data: The value of the node.
//
// Returns:
//   - *Node[T]: A pointer to the newly created node.
func NewTreeNode[T any](data T) *TreeNode[T] {
	return &TreeNode[T]{
		Data:     data,
		children: make([]*TreeNode[T], 0),
	}
}

// IsLeaf returns true if the node is a leaf.
//
// Returns:
//   - bool: True if the node is a leaf, false otherwise.
func (n *TreeNode[T]) IsLeaf() bool {
	return len(n.children) == 0
}

// IsRoot returns true if the node does not have a parent.
//
// Returns:
//   - bool: True if the node is the root, false otherwise.
func (n *TreeNode[T]) IsRoot() bool {
	return n.parent == nil
}

// Leaves returns all the leaves of the tree rooted at n.
//
// Returns:
//   - []*Node[T]: A slice of pointers to the leaves of the tree.
//
// Behaviors:
//   - The leaves are returned in the order of a DFS traversal.
func (n *TreeNode[T]) Leaves() []*TreeNode[T] {
	leaves := make([]*TreeNode[T], 0)

	S := Stacker.NewLinkedStack(n)

	for {
		top, ok := S.Pop()
		if !ok {
			break
		}

		ok = top.IsLeaf()
		if ok {
			leaves = append(leaves, top)
		} else {
			for _, child := range top.children {
				S.Push(child)
			}
		}
	}

	return leaves
}

// ToTree converts the node to a tree.
//
// Returns:
//   - *Tree[T]: A pointer to the tree.
func (n *TreeNode[T]) ToTree() *Tree[T] {
	if len(n.children) == 0 {
		return &Tree[T]{
			root:   n,
			leaves: []*TreeNode[T]{n},
			size:   1,
		}
	} else {
		return &Tree[T]{
			root:   n,
			leaves: n.Leaves(),
			size:   n.Size(),
		}
	}
}

// AddChild adds a new child to the node with the given data.
//
// Parameters:
//   - child: The child to add.
func (n *TreeNode[T]) AddChild(child *TreeNode[T]) {
	if child == nil {
		return
	}

	child.parent = n
	n.children = append(n.children, child)
}

// AddChildren adds zero or more children to the node.
//
// Parameters:
//   - children: The children to add.
//
// Behaviors:
//   - This is just a more efficient way to add multiple children.
func (n *TreeNode[T]) AddChildren(children ...*TreeNode[T]) {
	children = us.FilterNilValues(children)
	if len(children) == 0 {
		return
	}

	for _, child := range children {
		child.parent = n
	}

	n.children = append(n.children, children...)
}

// GetChildren returns all the children of the node.
// If the node has no children, it returns nil.
//
// Returns:
//   - []*Node[T]: A slice of pointers to the children of the node.
func (n *TreeNode[T]) GetChildren() []*TreeNode[T] {
	return n.children
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
func (tn *TreeNode[T]) FindBranchingPoint() (*TreeNode[T], *TreeNode[T], bool) {
	node := tn
	var prev *TreeNode[T] = nil

	for node.parent != nil {
		if len(node.parent.children) > 1 {
			return prev, node.parent, true
		}

		prev = node
		node = node.parent
	}

	return prev, node, false
}

// HasChild returns true if the node has the given child.
//
// Parameters:
//   - target: The child to check for.
//
// Returns:
//   - bool: True if the node has the child, false otherwise.
func (n *TreeNode[T]) HasChild(target *TreeNode[T]) bool {
	if target == nil {
		return false
	}

	for _, c := range n.children {
		if c == target {
			return true
		}
	}

	return false
}

// DeleteChild removes the given child from the children of the node.
//
// Parameters:
//   - target: The child to remove.
//
// Returns:
//   - []*treeNode[T]: A slice of pointers to the children of the node.
//
// Behaviors:
//   - If the node has no children, it returns nil.
func (n *TreeNode[T]) DeleteChild(target *TreeNode[T]) []*TreeNode[T] {
	if target == nil {
		// No target to delete
		return nil
	}

	switch len(n.children) {
	case 0:
		// No children to delete
	case 1:
		if n.children[0] == target {
			n.children = nil
			target.parent = nil

			return target.children
		}
	default:
		index := slices.Index(n.children, target)
		if index == -1 {
			// target is not a child of n
			return nil
		}

		n.children = slices.Delete(n.children, index, index+1)
		target.parent = nil

		return target.children
	}

	return nil
}

// Parent is a getter for the parent of the node.
//
// Returns:
//   - *Node[T]: A pointer to the parent of the node. Nil if the node has no parent.
func (n *TreeNode[T]) Parent() *TreeNode[T] {
	return n.parent
}

// Size returns the number of nodes in the tree rooted at n.
//
// Returns:
//   - int: The number of nodes in the tree.
//
// Behaviors:
//   - This function is expensive since size is not stored.
func (n *TreeNode[T]) Size() int {
	size := 0

	S := Stacker.NewLinkedStack(n)

	for {
		top, ok := S.Pop()
		if !ok {
			break
		}

		size++

		for _, child := range top.children {
			S.Push(child)
		}
	}

	return size
}

// GetAncestors returns all the ancestors of the node.
//
// Returns:
//   - []*Node[T]: A slice of pointers to the ancestors of the node.
//
// Behaviors:
//   - The ancestors are returned in the opposite order of a DFS traversal.
//     Therefore, the first element is the parent of the node.
func (n *TreeNode[T]) GetAncestors() []*TreeNode[T] {
	ancestors := make([]*TreeNode[T], 0)

	for node := n; node.parent != nil; node = node.parent {
		ancestors = append(ancestors, node.parent)
	}

	slices.Reverse(ancestors)

	return ancestors
}

// IsChildOf returns true if the node is a child of the parent.
//
// Parameters:
//   - target: The target parent to check for.
//
// Returns:
//   - bool: True if the node is a child of the parent, false otherwise.
func (n *TreeNode[T]) IsChildOf(target *TreeNode[T]) bool {
	if target == nil {
		return false
	}

	parents := target.GetAncestors()

	for node := n; node.parent != nil; node = node.parent {
		if slices.Contains(parents, node.parent) {
			return true
		}
	}

	return false
}

// removeNode removes the node from the tree.
//
// Returns:
//   - []*treeNode[T]: A slice of pointers to the children of the node if
//     the node is the root. Nil otherwise.
func (n *TreeNode[T]) removeNode() []*TreeNode[T] {
	if n.parent == nil {
		// The node is the root
		return n.children
	}

	parent := n.parent

	// Remove the node from the parent's children
	index := slices.Index(parent.children, n)
	parent.children = slices.Delete(parent.children, index, index+1)

	if len(n.children) != 0 {
		for i := 0; i < len(n.children); i++ {
			n.children[i].parent = parent
		}

		parent.children = slices.Insert(parent.children, index, n.children...)
	}

	return nil
}

// GetData is a getter for the data of the node.
//
// Returns:
//   - T: The data of the node.
func (n *TreeNode[T]) GetData() T {
	return n.Data
}
