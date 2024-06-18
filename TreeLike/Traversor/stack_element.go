package Traversor

import (
	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// stackElement is a stack element.
type stackElement[T any] struct {
	// prev is the previous node.
	prev *tr.TreeNode[T]

	// elem is the current node.
	elem *tr.TreeNode[T]

	// info is the info of the current node.
	info uc.Copier
}

// newStackElement creates a new stack element.
//
// Parameters:
//   - prev: The previous node.
//   - elem: The helper.
//
// Returns:
//   - *stackElement[T, E]: A pointer to the stack element.
func newStackElement[T any](prev *tr.TreeNode[T], data T, info uc.Copier) *stackElement[T] {
	return &stackElement[T]{
		prev: prev,
		elem: tr.NewTreeNode(data),
		info: info,
	}
}

// getData returns the data of the stack element.
//
// Returns:
//   - T: The data of the stack element.
//   - bool: True if the data is valid, otherwise false.
func (se *stackElement[T]) getData() (T, bool) {
	if se.elem == nil {
		return *new(T), false
	}

	return se.elem.Data, true
}

// getInfo returns the info of the stack element.
//
// Returns:
//   - uc.Objecter: The info of the stack element.
func (se *stackElement[T]) getInfo() uc.Copier {
	return se.info
}

// linkToPrev links the current node to the previous node.
func (se *stackElement[T]) linkToPrev() bool {
	if se.prev == nil {
		return false
	}

	se.prev.AddChild(se.elem)

	return true
}

// getElem returns the current node.
//
// Returns:
//   - *tr.TreeNode[T]: The current node.
func (se *stackElement[T]) getElem() *tr.TreeNode[T] {
	return se.elem
}
