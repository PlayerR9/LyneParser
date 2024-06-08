package ConflictSolver

import (
	"fmt"
	"slices"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	intf "github.com/PlayerR9/MyGoLib/Units/Common"
	slext "github.com/PlayerR9/MyGoLib/Units/Slice"
	ers "github.com/PlayerR9/MyGoLib/Units/errors"

	trt "github.com/PlayerR9/MyGoLib/TreeLike/Traversor"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
)

// InfoStruct is the information about the expansion tree.
type InfoStruct struct {
	// seen is a map of helpers that have been seen.
	seen map[*Helper]bool
}

// Copy creates a copy of the InfoStruct.
//
// Returns:
//   - intf.Copier: A copy of the InfoStruct.
func (is *InfoStruct) Copy() intf.Copier {
	isCopy := &InfoStruct{
		seen: make(map[*Helper]bool),
	}

	for k, v := range is.seen {
		isCopy.seen[k] = v
	}

	return isCopy
}

// NewInfoStruct creates a new InfoStruct.
//
// Parameters:
//   - root: The root of the expansion tree.
//
// Returns:
//   - *InfoStruct: The new InfoStruct.
//   - error: An error of type *ers.ErrInvalidParameter if the root is nil.
//
// Behaviors:
//   - The root is set to seen.
func NewInfoStruct(root *Helper) (*InfoStruct, error) {
	if root == nil {
		return nil, ers.NewErrNilParameter("root")
	}

	info := &InfoStruct{
		seen: make(map[*Helper]bool),
	}

	info.seen[root] = true

	return info, nil
}

// IsSNoteen checks if the helper is seen.
//
// Parameters:
//   - h: The helper to check.
//
// Returns:
//   - bool: True if the helper is seen, false otherwise.
func (is *InfoStruct) IsNotSeen(h *Helper) bool {
	return !is.seen[h]
}

// SetSeen sets the helper as seen.
//
// Parameters:
//   - h: The helper to set as seen.
func (is *InfoStruct) SetSeen(h *Helper) {
	is.seen[h] = true
}

// ExpansionTree is a tree of expansion helpers.
type ExpansionTree struct {
	// tree is the tree of expansion helpers.
	tree *tr.Tree[*Helper]

	// info is the information about the expansion tree.
	info *InfoStruct
}

// NewExpansionTree creates a new expansion tree where the root is h and every node is a
// helper whose LHS is the 0th RHS of the parent node. However, the leaves of the tree
// are helpers whose 0th RHS is a terminal symbol.
//
// Parameters:
//   - cs: The conflict solver.
//   - h: The root of the expansion tree.
//
// Returns:
//   - *ExpansionTree: The new expansion tree.
//   - error: An error if the operation failed.
//
// Errors:
//   - *ers.Err0thRhsNotSet: The 0th RHS of the root is not set.
//   - *ers.ErrInvalidParameter: The root is nil.
func NewExpansionTreeRootedAt(cs *ConflictSolver, h *Helper) (*ExpansionTree, error) {
	info, err := NewInfoStruct(h)
	if err != nil {
		return nil, err
	}

	nextsFunc := func(elem *Helper, is intf.Copier) ([]*Helper, error) {
		isInf, ok := is.(*InfoStruct)
		if !ok {
			return nil, fmt.Errorf("invalid type: %T", is)
		}

		rhs, err := elem.GetRhsAt(0)
		if err != nil {
			return nil, NewErr0thRhsNotSet()
		}

		if gr.IsTerminal(rhs) {
			return nil, nil
		}

		result := slext.SliceFilter(cs.GetElemsWithLhs(rhs), isInf.IsNotSeen)

		isInf.SetSeen(elem)

		return result, nil
	}

	var builder trt.Builder[*Helper]

	builder.SetInfo(info)
	builder.SetNextFunc(nextsFunc)

	tree, err := builder.Build(h)
	if err != nil {
		return nil, err
	}

	return &ExpansionTree{
		tree: tree,
		info: info,
	}, nil
}

// PruneNonTerminalLeaves prunes the non-terminal leaves of the expansion tree.
func (et *ExpansionTree) PruneNonTerminalLeaves() {
	leaves := et.tree.GetLeaves()

	todo := slext.SliceFilter(leaves, FilterNonTerminalLeaf)
	if len(todo) == 0 {
		return
	}

	for _, leaf := range todo {
		err := et.tree.DeleteBranchContaining(leaf)
		if err != nil {
			panic(err)
		}
	}
}

// Size returns the size of the expansion tree.
//
// Returns:
//   - int: The size of the expansion tree.
func (et *ExpansionTree) Size() int {
	return et.tree.Size()
}

// Collapse collapses the expansion tree into a slice of strings that
// are the 0th RHS of the terminal leaves.
//
// Returns:
//   - []string: The collapsed expansion tree.
func (et *ExpansionTree) Collapse() []string {
	leaves := et.tree.GetLeaves()

	var result []string

	for _, leaf := range leaves {
		rhs, err := leaf.Data.GetRhsAt(0)
		if err != nil {
			panic(err)
		}

		pos, ok := slices.BinarySearch(result, rhs)
		if !ok {
			result = slices.Insert(result, pos, rhs)
		}
	}

	return result
}

/////////////////////////////////////

/*
func (cs *ConflictSolver) CheckIfLookahead0(index int, h *Helper) ([]*Helper, error) {
	// 1. Take the next symbol of h
	rhs, err := h.GetRhsAt(index + 1)
	if err != nil {
		return nil, NewErrHelper(h, err)
	}

	// 2. Get all the helpers that have the same LHS as rhs
	newHelpers := cs.GetElemsWithLhs(rhs)
	if len(newHelpers) == 0 {
		return nil, nil
	}

	// 3. For each rule, check if the 0th rhs is a terminal symbol
	solutions := make([]*Helper, 0)

	for _, nh := range newHelpers {
		otherRhs, err := nh.GetRhsAt(0)
		if err != nil {
			return solutions, NewErrHelper(nh, err)
		}

		if gr.IsTerminal(otherRhs) {
			solutions = append(solutions, nh)
		} else {

		}
	}

	return solutions, nil
}

*/
