package Grammar

import (
	"errors"
	"fmt"
	"strings"

	trt "github.com/PlayerR9/MyGoLib/TreeLike/Traversor"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// TTInfo is the information about the token tree.
type TTInfo[T TokenTyper] struct {
	// depth is the depth of each token.
	depth map[*Token[T]]int
}

// Copy creates a copy of the TTInfo.
//
// Returns:
//   - uc.Copier: A copy of the TTInfo.
func (tti *TTInfo[T]) Copy() uc.Copier {
	tti_copy := &TTInfo[T]{
		depth: make(map[*Token[T]]int),
	}

	for k, v := range tti.depth {
		tti_copy.depth[k] = v
	}

	return tti_copy
}

// NewTTInfo creates a new TTInfo.
//
// Parameters:
//   - root: The root of the token tree.
//
// Returns:
//   - *TTInfo: The new TTInfo.
//   - error: An error of type *uc.ErrInvalidParameter if the root is nil.
//
// Behaviors:
//   - The depth of the root is set to 0.
func NewTTInfo[T TokenTyper](root *Token[T]) (*TTInfo[T], error) {
	info := &TTInfo[T]{
		depth: make(map[*Token[T]]int),
	}

	info.depth[root] = 0

	return info, nil
}

// SetDepth sets the depth of the Token.
//
// Parameters:
//   - Token: The Token to set the depth of.
//   - depth: The depth to set.
//
// Returns:
//   - bool: True if the depth was set. False if the Token already has a depth.
func (tti *TTInfo[T]) SetDepth(Token *Token[T], depth int) bool {
	_, ok := tti.depth[Token]
	if ok {
		return false
	}

	tti.depth[Token] = depth

	return true
}

// GetDepth gets the depth of the Token.
//
// Parameters:
//   - Token: The Token to get the depth of.
//
// Returns:
//   - int: The depth of the Token.
//   - bool: True if the depth was found. False if the Token does not have a depth.
func (tti *TTInfo[T]) GetDepth(Token *Token[T]) (int, bool) {
	depth, ok := tti.depth[Token]
	if !ok {
		return 0, false
	}

	return depth, true
}

// TokenTree is a tree of tokens.
type TokenTree[T TokenTyper] struct {
	// tree is the tree of tokens.
	tree *tr.Tree[*Token[T]]

	// Info is the information about the tree.
	Info *TTInfo[T]
}

// NewTokenTree creates a new token tree.
//
// Parameters:
//   - root: The root of the token tree.
//
// Returns:
//   - *TokenTree: The new token tree.
//   - error: An error if the token tree could not be created.
//
// Errors:
//   - *ErrCycleDetected: A cycle is detected in the token tree.
//   - *uc.ErrInvalidParameter: The root is nil.
//   - *ErrUnknowToken: The root is not a known token.
func NewTokenTree[T TokenTyper](root *Token[T]) (*TokenTree[T], error) {
	tree_info, err := NewTTInfo(root)
	if err != nil {
		return nil, err
	}

	nexts_func := func(elem *Token[T], h uc.Copier) ([]*Token[T], error) {
		h_info, ok := h.(*TTInfo[T])
		if !ok {
			return nil, fmt.Errorf("invalid type: %T", h)
		}

		ok = elem.IsLeaf()
		if ok {
			return nil, nil
		}

		ok = elem.IsNonLeaf()
		if !ok {
			return nil, errors.New("token is not leaf or non-leaf")
		}

		children := elem.Data.([]*Token[T])

		for _, child := range children {
			ok := h_info.SetDepth(child, h_info.depth[elem]+1)
			if !ok {
				return nil, NewErrCycleDetected()
			}
		}

		return children, nil
	}

	var builder trt.Builder[*Token[T]]

	builder.SetNextFunc(nexts_func)
	builder.SetInfo(tree_info)

	tree, err := builder.Build(root)
	if err != nil {
		return nil, err
	}

	tt := &TokenTree[T]{
		tree: tree,
		Info: tree_info,
	}

	return tt, nil
}

// DebugString returns a string representation of the token tree.
//
// Returns:
//   - string: The string representation of the token tree.
//
// Information: This is a debug function.
func (tt *TokenTree[T]) DebugString() string {
	var builder strings.Builder

	err := trt.DFS(
		tt.tree,
		tt.Info,
		func(elem *Token[T], inf uc.Copier) (bool, error) {
			h_info, ok := inf.(*TTInfo[T])
			if !ok {
				return false, fmt.Errorf("invalid type: %T", inf)
			}

			depth, ok := h_info.GetDepth(elem)
			if !ok {
				return false, fmt.Errorf("depth not found for: %v", elem)
			}

			builder.WriteString(strings.Repeat("   ", depth))
			builder.WriteString(elem.GetID().String())

			ok = elem.IsLeaf()
			if ok {
				builder.WriteString(" -> ")
				builder.WriteString(elem.Data.(string))
			} else {
				builder.WriteString(" :")
			}

			builder.WriteRune('\n')

			return true, nil
		},
	)
	uc.AssertF(err == nil, "DFS failed: %s", err.Error())

	return builder.String()
}

// GetAllBranches returns all the branches of the token tree.
//
// Returns:
//   - [][]Token: All the branches of the token tree.
func (tt *TokenTree[T]) GetAllBranches() [][]*Token[T] {
	trav := tt.tree.SnakeTraversal()
	return trav
}

// GetRoot returns the root of the token tree.
//
// Returns:
//   - Token: The root of the token tree.
func (tt *TokenTree[T]) GetRoot() *Token[T] {
	root := tt.tree.Root()
	data := root.Data

	return data
}
