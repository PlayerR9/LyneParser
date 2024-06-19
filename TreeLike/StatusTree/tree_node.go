package StatusTree

import (
	"slices"
	"strings"

	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// TreeNode is a generic data structure that represents a node in a tree.
type TreeNode[S uc.Enumer, T any] struct {
	// Data is the value of the node.
	Data T

	// status is the status of the current evaluation.
	status S

	// parent is the parent of the node.
	parent *TreeNode[S, T]

	// children is the children of the node.
	children []*TreeNode[S, T]
}

// String implements fmt.Stringer.
func (t *TreeNode[S, T]) String() string {
	var builder strings.Builder

	builder.WriteString(uc.StringOf(t.Data))
	builder.WriteString(" [")
	builder.WriteString(t.status.String())
	builder.WriteRune(']')

	return builder.String()
}

// FString returns a formatted string representation of the node.
//
// Parameters:
//   - indentLevel: The level of indentation to use.
//
// Returns:
//   - []string: A slice of strings that represent the node.
func (t *TreeNode[S, T]) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	err := trav.AddJoinedLine("", uc.StringOf(t.Data), " [", t.status.String(), "]")
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
func (n *TreeNode[S, T]) Copy() uc.Copier {
	node := &TreeNode[S, T]{
		Data:     n.Data,
		parent:   nil,
		status:   n.status,
		children: make([]*TreeNode[S, T], 0, len(n.children)),
	}

	for _, child := range n.children {
		node.children = append(node.children, child.Copy().(*TreeNode[S, T]))
	}

	return node
}

// NewTreeNode creates a new node with the given data.
//
// Parameters:
//   - status: The status of the node.
//   - data: The value of the node.
//
// Returns:
//   - *Node[T]: A pointer to the newly created node.
func NewTreeNode[S uc.Enumer, T any](status S, data T) *TreeNode[S, T] {
	return &TreeNode[S, T]{
		Data:     data,
		status:   status,
		children: make([]*TreeNode[S, T], 0),
	}
}

// IsLeaf returns true if the node is a leaf.
//
// Returns:
//   - bool: True if the node is a leaf, false otherwise.
func (n *TreeNode[S, T]) IsLeaf() bool {
	return len(n.children) == 0
}

// IsRoot returns true if the node does not have a parent.
//
// Returns:
//   - bool: True if the node is the root, false otherwise.
func (n *TreeNode[S, T]) IsRoot() bool {
	return n.parent == nil
}

// Leaves returns all the leaves of the tree rooted at n.
//
// Returns:
//   - []*Node[T]: A slice of pointers to the leaves of the tree.
//
// Behaviors:
//   - The leaves are returned in the order of a DFS traversal.
func (n *TreeNode[S, T]) Leaves() []*TreeNode[S, T] {
	leaves := make([]*TreeNode[S, T], 0)

	St := Stacker.NewLinkedStack(n)

	for {
		top, ok := St.Pop()
		if !ok {
			break
		}

		ok = top.IsLeaf()
		if ok {
			leaves = append(leaves, top)
		} else {
			for _, child := range top.children {
				St.Push(child)
			}
		}
	}

	return leaves
}

// ToTree converts the node to a tree.
//
// Returns:
//   - *Tree[S, T]: A pointer to the tree.
func (n *TreeNode[S, T]) ToTree() *Tree[S, T] {
	if len(n.children) == 0 {
		return &Tree[S, T]{
			root:   n,
			leaves: []*TreeNode[S, T]{n},
			size:   1,
		}
	} else {
		return &Tree[S, T]{
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
func (n *TreeNode[S, T]) AddChild(child *TreeNode[S, T]) {
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
func (n *TreeNode[S, T]) AddChildren(children ...*TreeNode[S, T]) {
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
func (n *TreeNode[S, T]) GetChildren() []*TreeNode[S, T] {
	return n.children
}

// HasChild returns true if the node has the given child.
//
// Parameters:
//   - target: The child to check for.
//
// Returns:
//   - bool: True if the node has the child, false otherwise.
func (n *TreeNode[S, T]) HasChild(target *TreeNode[S, T]) bool {
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
func (n *TreeNode[S, T]) DeleteChild(target *TreeNode[S, T]) []*TreeNode[S, T] {
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
func (n *TreeNode[S, T]) Parent() *TreeNode[S, T] {
	return n.parent
}

// Size returns the number of nodes in the tree rooted at n.
//
// Returns:
//   - int: The number of nodes in the tree.
//
// Behaviors:
//   - This function is expensive since size is not stored.
func (n *TreeNode[S, T]) Size() int {
	size := 0

	St := Stacker.NewLinkedStack(n)

	for {
		top, ok := St.Pop()
		if !ok {
			break
		}

		size++

		for _, child := range top.children {
			St.Push(child)
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
func (n *TreeNode[S, T]) GetAncestors() []*TreeNode[S, T] {
	ancestors := make([]*TreeNode[S, T], 0)

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
func (n *TreeNode[S, T]) IsChildOf(target *TreeNode[S, T]) bool {
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
func (n *TreeNode[S, T]) removeNode() []*TreeNode[S, T] {
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
func (n *TreeNode[S, T]) GetData() T {
	return n.Data
}

// ChangeStatus sets the status of the tree node.
//
// Parameters:
//   - status: The status to set.
func (tn *TreeNode[S, T]) ChangeStatus(status S) {
	tn.status = status
}

// GetStatus returns the status of the tree node.
//
// Returns:
//   - S: The status of the tree node.
func (tn *TreeNode[S, T]) GetStatus() S {
	return tn.status
}

// setParent sets the parent of the tree node.
//
// Parameters:
//   - parent: The parent to set.
func (tn *TreeNode[S, T]) setParent(parent *TreeNode[S, T]) {
	tn.parent = parent
}
