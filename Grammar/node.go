package Grammar

import (
	"errors"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// Noder is an interface for a node.
type Noder interface {
	// SetParent sets the parent.
	//
	// Parameters:
	//   - n: The parent. Never nil.
	SetParent(n Noder)

	// SetFirstChild sets the first child.
	//
	// This should also set the parent of the first child.
	//
	// Parameters:
	//   - n: The first child. Never nil.
	SetFirstChild(n Noder)

	// SetLastChild sets the last child.
	//
	// This should also set the parent of the last child.
	//
	// Parameters:
	//   - n: The last child. Never nil.
	SetLastChild(n Noder)

	// SetNextSibling sets the next sibling.
	//
	// This should also set the previous sibling of the next sibling (if any).
	//
	// Parameters:
	//   - n: The next sibling. Never nil.
	SetNextSibling(n Noder)
}

// LinkSubNodes links the sub nodes.
//
// Parameters:
//   - sub_nodes: The sub nodes.
func LinkSubNodes[T Noder](sub_nodes []T) {
	for i := 0; i < len(sub_nodes)-1; i++ {
		sn := sub_nodes[i]

		sn.SetNextSibling(sub_nodes[i+1])
	}
}

// LinkParent links the parent to the children.
//
// Parameters:
//   - parent: The parent.
//   - children: The children.
func LinkParent[T Noder](parent T, children []T) {
	switch len(children) {
	case 0:
		// Do nothing.
	case 1:
		parent.SetFirstChild(children[0])
		parent.SetLastChild(children[0])
	case 2:
		parent.SetFirstChild(children[0])
		parent.SetLastChild(children[1])
	default:
		parent.SetFirstChild(children[0])

		for i := 1; i < len(children)-1; i++ {
			child := children[i]

			child.SetParent(parent)
		}

		parent.SetLastChild(children[len(children)-1])
	}
}

// AstRecFunc is a function to recursively extract sub nodes.
//
// Parameters:
//   - tok: The token.
//
// Returns:
//   - []N: The sub nodes. Nil if there are no sub nodes.
type AstRecFunc[N Noder, T uc.Enumer] func(tok *Token[T]) []N

// ExtractSubNodes is a helper function to extract sub nodes.
//
// Parameters:
//   - f: The function to recursively extract sub nodes.
//   - tok: The token.
//
// Returns:
//   - []*Node: The sub nodes. Nil if there are no sub nodes.
//   - bool: True if the function succeeds. (i.e., tok.Data is a []*Token[T])
//
// Behaviors:
//   - If the function is nil or tok is nil, the function will return nil and true.
func ExtractSubNodes[N Noder, T uc.Enumer](f AstRecFunc[N, T], tok *Token[T]) ([]N, bool) {
	if f == nil || tok == nil {
		return nil, true
	}

	children, ok := tok.Data.([]*Token[T])
	if !ok {
		return nil, false
	}

	var sub_nodes []N

	for _, child := range children {
		ns := f(child)
		if len(ns) != 0 {
			sub_nodes = append(sub_nodes, ns...)
		}
	}

	if len(sub_nodes) == 0 {
		return nil, true
	}

	return sub_nodes, true
}

// Ast generates the AST.
//
// Parameters:
//   - f: The function to recursively extract sub nodes.
//   - token: The token.
//
// Returns:
//   - N: The root node.
//   - error: An error if the AST generation fails.
func Ast[N Noder, T uc.Enumer](f AstRecFunc[N, T], tok *Token[T]) (N, error) {
	if f == nil {
		return *new(N), uc.NewErrNilParameter("f")
	} else if tok == nil {
		return *new(N), uc.NewErrNilParameter("tok")
	}

	sub_nodes := f(tok)
	if len(sub_nodes) == 0 {
		return *new(N), errors.New("ast cleaned up to nothing")
	}

	if len(sub_nodes) > 1 {
		return *new(N), errors.New("ast generated a forest")
	}

	n := sub_nodes[0]

	return n, nil
}
