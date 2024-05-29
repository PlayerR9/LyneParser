package Common

import (
	"fmt"
	"strings"

	trt "github.com/PlayerR9/MyGoLib/TreeLike/Traversor"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
	intf "github.com/PlayerR9/MyGoLib/Units/Common"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"

	gr "github.com/PlayerR9/LyneParser/Grammar"
)

// TTInfo is the information about the token tree.
type TTInfo struct {
	// depth is the depth of each token.
	depth map[gr.Tokener]int
}

// Copy creates a copy of the TTInfo.
//
// Returns:
//   - intf.Copier: A copy of the TTInfo.
func (tti *TTInfo) Copy() intf.Copier {
	ttiCopy := &TTInfo{
		depth: make(map[gr.Tokener]int),
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
//   - error: An error of type *ers.ErrInvalidParameter if the root is nil.
//
// Behaviors:
//   - The depth of the root is set to 0.
func NewTTInfo(root gr.Tokener) (*TTInfo, error) {
	if root == nil {
		return nil, ers.NewErrNilParameter("root")
	}

	info := &TTInfo{
		depth: make(map[gr.Tokener]int),
	}

	info.depth[root] = 0

	return info, nil
}

// SetDepth sets the depth of the tokener.
//
// Parameters:
//   - tokener: The tokener to set the depth of.
//   - depth: The depth to set.
//
// Returns:
//   - bool: True if the depth was set. False if the tokener already has a depth.
func (tti *TTInfo) SetDepth(tokener gr.Tokener, depth int) bool {
	_, ok := tti.depth[tokener]
	if ok {
		return false
	}

	tti.depth[tokener] = depth

	return true
}

// GetDepth gets the depth of the tokener.
//
// Parameters:
//   - tokener: The tokener to get the depth of.
//
// Returns:
//   - int: The depth of the tokener.
//   - bool: True if the depth was found. False if the tokener does not have a depth.
func (tti *TTInfo) GetDepth(tokener gr.Tokener) (int, bool) {
	depth, ok := tti.depth[tokener]
	if !ok {
		return 0, false
	}

	return depth, true
}

// TokenTree is a tree of tokens.
type TokenTree struct {
	// tree is the tree of tokens.
	tree *tr.Tree[gr.Tokener]

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
//   - *ers.ErrInvalidParameter: The root is nil.
//   - *gr.ErrUnknowToken: The root is not a known token.
func NewTokenTree(root gr.Tokener) (*TokenTree, error) {
	treeInfo, err := NewTTInfo(root)
	if err != nil {
		return nil, err
	}

	nextsFunc := func(elem gr.Tokener, h intf.Copier) ([]gr.Tokener, error) {
		hInfo, ok := h.(*TTInfo)
		if !ok {
			return nil, fmt.Errorf("invalid type: %T", h)
		}

		switch elem := elem.(type) {
		case *gr.LeafToken:
			return nil, nil
		case *gr.NonLeafToken:
			for _, child := range elem.Data {
				ok := hInfo.SetDepth(child, hInfo.depth[elem]+1)
				if !ok {
					return nil, NewErrCycleDetected()
				}
			}

			return elem.Data, nil
		default:
			return nil, gr.NewErrUnknowToken(elem)
		}
	}

	var builder trt.Builder[gr.Tokener]

	builder.SetNextFunc(nextsFunc)
	builder.SetInfo(treeInfo)

	tree, err := builder.Build(root)
	if err != nil {
		return nil, err
	}

	return &TokenTree{
		tree: tree,
		Info: treeInfo,
	}, nil
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
		func(elem gr.Tokener, inf intf.Copier) (bool, error) {
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

			switch root := elem.(type) {
			case *gr.LeafToken:
				builder.WriteString(" -> ")
				builder.WriteString(root.Data)
			case *gr.NonLeafToken:
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
//   - [][]gr.Tokener: All the branches of the token tree.
func (tt *TokenTree) GetAllBranches() [][]gr.Tokener {
	return tt.tree.SnakeTraversal()
}
