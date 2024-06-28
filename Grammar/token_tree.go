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
type TTInfo struct {
	// depth is the depth of each token.
	depth map[Token]int
}

// Copy creates a copy of the TTInfo.
//
// Returns:
//   - uc.Copier: A copy of the TTInfo.
func (tti *TTInfo) Copy() uc.Copier {
	ttiCopy := &TTInfo{
		depth: make(map[Token]int),
	}

	for k, v := range tti.depth {
		ttiCopy.depth[k] = v
	}

	return ttiCopy
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
func NewTTInfo(root Token) (*TTInfo, error) {
	info := &TTInfo{
		depth: make(map[Token]int),
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
func (tti *TTInfo) SetDepth(Token Token, depth int) bool {
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
func (tti *TTInfo) GetDepth(Token Token) (int, bool) {
	depth, ok := tti.depth[Token]
	if !ok {
		return 0, false
	}

	return depth, true
}

// TokenTree is a tree of tokens.
type TokenTree struct {
	// tree is the tree of tokens.
	tree *tr.Tree[Token]

	// Info is the information about the tree.
	Info *TTInfo
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
func NewTokenTree(root Token) (*TokenTree, error) {
	treeInfo, err := NewTTInfo(root)
	if err != nil {
		return nil, err
	}

	nextsFunc := func(elem Token, h uc.Copier) ([]Token, error) {
		hInfo, ok := h.(*TTInfo)
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

		children := elem.Data.([]Token)

		for _, child := range children {
			ok := hInfo.SetDepth(child, hInfo.depth[elem]+1)
			if !ok {
				return nil, NewErrCycleDetected()
			}
		}

		return children, nil
	}

	var builder trt.Builder[Token]

	builder.SetNextFunc(nextsFunc)
	builder.SetInfo(treeInfo)

	tree, err := builder.Build(root)
	if err != nil {
		return nil, err
	}

	tt := &TokenTree{
		tree: tree,
		Info: treeInfo,
	}

	return tt, nil
}

// DebugString returns a string representation of the token tree.
//
// Returns:
//   - string: The string representation of the token tree.
//
// Information: This is a debug function.
func (tt *TokenTree) DebugString() string {
	var builder strings.Builder

	err := trt.DFS(
		tt.tree,
		tt.Info,
		func(elem Token, inf uc.Copier) (bool, error) {
			hInfo, ok := inf.(*TTInfo)
			if !ok {
				return false, fmt.Errorf("invalid type: %T", inf)
			}

			depth, ok := hInfo.GetDepth(elem)
			if !ok {
				return false, fmt.Errorf("depth not found for: %v", elem)
			}

			builder.WriteString(strings.Repeat("   ", depth))
			builder.WriteString(elem.GetID())

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
	if err != nil {
		panic(err)
	}

	return builder.String()
}

// GetAllBranches returns all the branches of the token tree.
//
// Returns:
//   - [][]Token: All the branches of the token tree.
func (tt *TokenTree) GetAllBranches() [][]Token {
	trav := tt.tree.SnakeTraversal()
	return trav
}

// GetRoot returns the root of the token tree.
//
// Returns:
//   - Token: The root of the token tree.
func (tt *TokenTree) GetRoot() Token {
	root := tt.tree.Root()
	data := root.Data

	return data
}
