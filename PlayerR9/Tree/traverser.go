package Tree

import (
	intf "github.com/PlayerR9/LyneParser/PlayerR9/Common"
	up "github.com/PlayerR9/LyneParser/PlayerR9/Pair"
	"github.com/PlayerR9/MyGoLib/ListLike/Queuer"
	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"
)

// ObserverFunc is a function that observes a node.
//
// Parameters:
//   - data: The data of the node.
//   - info: The info of the node.
//
// Returns:
//   - error: An error if the observation fails.
type ObserverFunc[T any, I intf.Copier] func(data T, info I) error

// helper is a helper struct for traversing the tree.
type helper[T any, I intf.Copier] up.Pair[*treeNode[T], I]

// newHelper creates a new helper.
//
// Parameters:
//   - node: The current node.
//   - info: The info of the current node.
//
// Returns:
//   - *helper[T, I]: A pointer to the helper.
func newHelper[T any, I intf.Copier](node *treeNode[T], info I) *helper[T, I] {
	return &helper[T, I]{
		First:  node,
		Second: info.Copy().(I),
	}
}

// Traversor is a struct that traverses a tree.
type Traversor[T any, I intf.Copier] struct {
	// The helper struct.
	h *helper[T, I]

	// The observer function.
	observe ObserverFunc[T, I]
}

// DFS traverses the tree in depth-first order.
//
// Returns:
//   - error: An error if the traversal fails.
func (t *Traversor[T, I]) DFS() error {
	if t.h.First == nil || t.observe == nil {
		return nil
	}

	S := Stacker.NewLinkedStack(newHelper(t.h.First, t.h.Second))

	for {
		top, err := S.Pop()
		if err != nil {
			break
		}

		err = t.observe(top.First.Data, top.Second)
		if err != nil {
			return err
		}

		for _, child := range top.First.children {
			err := S.Push(newHelper(child, top.Second))
			if err != nil {
				panic(err)
			}
		}
	}

	return nil
}

// BFS traverses the tree in breadth-first order.
//
// Returns:
//   - error: An error if the traversal fails.
func (t *Traversor[T, I]) BFS() error {
	if t.h.First == nil || t.observe == nil {
		return nil
	}

	Q := Queuer.NewLinkedQueue(newHelper(t.h.First, t.h.Second))

	for {
		first, err := Q.Dequeue()
		if err != nil {
			break
		}

		err = t.observe(first.First.Data, first.Second)
		if err != nil {
			return err
		}

		for _, child := range first.First.children {
			err := Q.Enqueue(newHelper(child, first.Second))
			if err != nil {
				panic(err)
			}
		}
	}

	return nil
}

// Traverse creates a new traversor for the tree.
//
// Parameters:
//   - tree: The tree to traverse.
//   - init: The initial info.
//   - f: The observer function.
//
// Returns:
//   - Traversor[T, I]: The traversor.
func Traverse[T any, I intf.Copier](tree *Tree[T], init I, f ObserverFunc[T, I]) *Traversor[T, I] {
	var root *treeNode[T]

	if tree == nil {
		root = nil
	} else {
		root = tree.root
	}

	return &Traversor[T, I]{
		h: &helper[T, I]{
			First:  root,
			Second: init,
		},
		observe: f,
	}
}

type NextsFunc[T any, I intf.Copier] func(elem T, info I) ([]T, error)

// MakeTree creates a tree from the given element.
//
// Parameters:
//   - elem: The element to start the tree from.
//
// Returns:
//   - *Tree[T]: The tree created from the element.
//   - error: An error if the tree cannot be created.
func MakeTree[T any, I intf.Copier](elem T, info I, f NextsFunc[T, I]) (*Tree[T], error) {
	// 1. Handle the first element
	h := newHelper(newTreeNode(elem), info)

	nexts, err := f(h.First.Data, h.Second)
	if err != nil {
		return nil, err
	}

	tree := &Tree[T]{
		root: h.First,
	}

	if len(nexts) == 0 {
		return tree, nil
	}

	// 2. Create a stack and push the first element
	type StackElement struct {
		Prev *treeNode[T]
		Elem *helper[T, I]
	}

	S := Stacker.NewLinkedStack[StackElement]()

	for _, next := range nexts {
		err := S.Push(StackElement{Prev: tree.root, Elem: newHelper(newTreeNode(next), h.Second)})
		if err != nil {
			panic(err)
		}
	}

	for {
		top, err := S.Pop()
		if err != nil {
			break
		}

		nexts, err := f(top.Elem.First.Data, top.Elem.Second)
		if err != nil {
			return nil, err
		}

		top.Prev.AddChildren(top.Elem.First)

		for _, next := range nexts {
			err := S.Push(StackElement{Prev: top.Elem.First, Elem: newHelper(newTreeNode(next), top.Elem.Second)})
			if err != nil {
				panic(err)
			}
		}
	}

	return tree, nil
}
