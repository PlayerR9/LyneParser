package Grammar

import (
	"fmt"
	"strings"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	tr "github.com/PlayerR9/tree/tree"
)

// TTInfo is the information about the token tree.
type TTInfo struct {
	// depth is the depth of each token.
	depth map[tr.Noder]int
}

// Copy creates a copy of the TTInfo.
//
// Returns:
//   - uc.Copier: A copy of the TTInfo.
func (tti *TTInfo) Copy() uc.Copier {
	tti_copy := &TTInfo{
		depth: make(map[tr.Noder]int),
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
func NewTTInfo(root tr.Noder) (*TTInfo, error) {
	info := &TTInfo{
		depth: make(map[tr.Noder]int),
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
func (tti *TTInfo) SetDepth(Token tr.Noder, depth int) bool {
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
func (tti *TTInfo) GetDepth(Token tr.Noder) (int, bool) {
	depth, ok := tti.depth[Token]
	if !ok {
		return 0, false
	}

	return depth, true
}

// TokenTree is a tree of tokens.
type TokenTree struct {
	// tree is the tree of tokens.
	tree *tr.Tree

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
func NewTokenTree(root tr.Noder) (*TokenTree, error) {
	tree_info, err := NewTTInfo(root)
	if err != nil {
		return nil, err
	}

	nexts_func := func(elem tr.Noder, h tr.Infoer) ([]tr.Noder, error) {
		h_info, ok := h.(*TTInfo)
		if !ok {
			return nil, fmt.Errorf("invalid type: %T", h)
		}

		ok = elem.IsLeaf()
		if ok {
			return nil, nil
		}

		iter := elem.Iterator()
		if iter == nil {
			return nil, nil
		}

		var children []tr.Noder

		for {
			val, err := iter.Consume()
			ok := uc.IsDone(err)
			if ok {
				break
			} else if err != nil {
				return nil, err
			}

			ok = h_info.SetDepth(val, h_info.depth[elem]+1)
			if !ok {
				return nil, NewErrCycleDetected()
			}

			children = append(children, val)
		}

		return children, nil
	}

	var builder tr.Builder

	builder.SetNextFunc(nexts_func)
	builder.SetInfo(tree_info)

	tree, err := builder.Build(root)
	if err != nil {
		return nil, err
	}

	tt := &TokenTree{
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
func (tt *TokenTree) DebugString(f func(data tr.Noder) string) string {
	var builder strings.Builder

	err := tr.DFS(
		tt.tree,
		tt.Info,
		func(data tr.Noder, info tr.Infoer) (bool, error) {
			h_info, ok := info.(*TTInfo)
			if !ok {
				return false, fmt.Errorf("invalid type: %T", info)
			}

			depth, ok := h_info.GetDepth(data)
			if !ok {
				return false, fmt.Errorf("depth not found for: %v", data)
			}

			builder.WriteString(strings.Repeat("   ", depth))

			str := f(data)
			/*builder.WriteString(data.GetID().String())

			ok = data.IsLeaf()
			if ok {
				builder.WriteString(" -> ")
				builder.WriteString(data.Data.(string))
			} else {
				builder.WriteString(" :")
			}*/

			builder.WriteString(str)
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
func (tt *TokenTree) GetAllBranches() ([][]tr.Noder, error) {
	trav, err := tt.tree.SnakeTraversal()
	if err != nil {
		return nil, err
	}

	return trav, nil
}

// GetRoot returns the root of the token tree.
//
// Returns:
//   - Token: The root of the token tree.
func (tt *TokenTree) GetRoot() tr.Noder {
	root := tt.tree.Root()
	return root
}
