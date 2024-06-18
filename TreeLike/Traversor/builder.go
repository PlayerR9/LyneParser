package Traversor

import (
	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// Builder is a struct that builds a tree.
type Builder[T any] struct {
	// info is the info of the builder.
	info uc.Copier

	// f is the next function.
	f NextsFunc[T]
}

// SetInfo sets the info of the builder.
//
// Parameters:
//   - info: The info to set.
func (b *Builder[T]) SetInfo(info uc.Copier) {
	b.info = info
}

// NextsFunc is a function that returns the next elements.
//
// Parameters:
//   - elem: The element to get the next elements from.
//   - info: The info of the element.
//
// Returns:
//   - []T: A slice of the next elements.
//   - error: An error if the function fails.
type NextsFunc[T any] func(elem T, info uc.Copier) ([]T, error)

// SetNextFunc sets the next function of the builder.
//
// Parameters:
//   - f: The function to set.
func (b *Builder[T]) SetNextFunc(f NextsFunc[T]) {
	b.f = f
}

// MakeTree creates a tree from the given element.
//
// Parameters:
//   - elem: The element to start the tree from.
//   - info: The info of the element.
//   - f: The function that, given an element and info, returns the next elements.
//     (i.e., the children of the element).
//
// Returns:
//   - *Tree[T]: The tree created from the element.
//   - error: An error if the tree cannot be created.
//
// Behaviors:
//   - The 'info' parameter is copied for each node and it specifies the initial info
//     before traversing the tree.
func (b *Builder[T]) Build(elem T) (*tr.Tree[T], error) {
	if b.f == nil {
		return nil, nil
	}

	// 1. Handle the root node
	nexts, err := b.f(elem, b.info)
	if err != nil {
		return nil, err
	}

	tree := tr.NewTree(elem)

	if len(nexts) == 0 {
		return tree, nil
	}

	S := Stacker.NewLinkedStack[*stackElement[T]]()

	for _, next := range nexts {
		se := newStackElement(tree.Root(), next, b.info)

		S.Push(se)
	}

	for {
		top, ok := S.Pop()
		if !ok {
			break
		}

		topData, ok := top.getData()
		if !ok {
			panic("Missing data")
		}
		topInfo := top.getInfo()

		nexts, err := b.f(topData, topInfo)
		if err != nil {
			return nil, err
		}

		ok = top.linkToPrev()
		if !ok {
			panic("Cannot link to previous node")
		}

		if len(nexts) == 0 {
			continue
		}

		topElem := top.getElem()

		for _, next := range nexts {
			se := newStackElement(topElem, next, topInfo)

			S.Push(se)
		}
	}

	b.Reset()

	tree.RegenerateLeaves()

	return tree, nil
}

// Reset resets the builder.
func (b *Builder[T]) Reset() {
	b.info = nil
	b.f = nil
}
