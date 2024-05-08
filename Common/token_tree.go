package Common

import (
	"strings"

	tr "github.com/PlayerR9/MyGoLib/CustomData/Tree"
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

	nextsFunc := func(elem gr.Tokener, h *TTInfo) ([]gr.Tokener, error) {
		switch elem := elem.(type) {
		case *gr.LeafToken:
			return nil, nil
		case *gr.NonLeafToken:
			for _, child := range elem.Data {
				_, ok := h.depth[child]
				if ok {
					return nil, NewErrCycleDetected()
				}

				h.depth[child] = h.depth[elem] + 1
			}

			return elem.Data, nil
		default:
			return nil, gr.NewErrUnknowToken(elem)
		}
	}

	tree, err := tr.MakeTree(root, treeInfo, nextsFunc)

	return &TokenTree{
		tree: tree,
		Info: treeInfo,
	}, err
}

// DebugString returns a string representation of the token tree.
//
// Returns:
//   - string: The string representation of the token tree.
//
// Information: This is a debug function.
func (tt *TokenTree) DebugString() string {
	var builder strings.Builder

	trav := tr.Traverse(
		tt.tree,
		tt.Info,
		func(elem gr.Tokener, inf *TTInfo) error {
			builder.WriteString(strings.Repeat("   ", inf.depth[elem]))
			builder.WriteString(elem.GetID())

			switch root := elem.(type) {
			case *gr.LeafToken:
				builder.WriteString(" -> ")
				builder.WriteString(root.Data)
			case *gr.NonLeafToken:
				builder.WriteString(" :")
			}

			builder.WriteRune('\n')

			return nil
		},
	)

	err := trav.DFS()
	if err != nil {
		panic(err)
	}

	return builder.String()
}
