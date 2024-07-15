package ConflictSolver

import (
	"slices"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
	tr "github.com/PlayerR9/tree/Tree"
	tn "github.com/PlayerR9/treenode"
)

// InfoStruct is the information about the expansion tree.
type InfoStruct[T gr.TokenTyper] struct {
	// seen is a map of helpers that have been seen.
	seen map[*HelperNode[T]]bool
}

// Copy implements the Copier interface.
func (is *InfoStruct[T]) Copy() uc.Copier {
	is_copy := &InfoStruct[T]{
		seen: make(map[*HelperNode[T]]bool),
	}

	for k, v := range is.seen {
		is_copy.seen[k] = v
	}

	return is_copy
}

// NewInfoStruct creates a new InfoStruct.
//
// Parameters:
//   - root: The root of the expansion tree.
//
// Returns:
//   - *InfoStruct: The new InfoStruct.
//
// Behaviors:
//   - The root is set to seen.
//   - If the root is nil, then nil is returned.
func NewInfoStruct[T gr.TokenTyper](root *HelperNode[T]) *InfoStruct[T] {
	if root == nil {
		return nil
	}

	info := &InfoStruct[T]{
		seen: make(map[*HelperNode[T]]bool),
	}

	info.seen[root] = true

	return info
}

// IsSNoteen checks if the helper is seen.
//
// Parameters:
//   - h: The helper to check.
//
// Returns:
//   - bool: True if the helper is seen, false otherwise.
func (is *InfoStruct[T]) IsNotSeen(h *HelperNode[T]) bool {
	return !is.seen[h]
}

// SetSeen sets the helper as seen.
//
// Parameters:
//   - h: The helper to set as seen.
func (is *InfoStruct[T]) SetSeen(h *HelperNode[T]) {
	is.seen[h] = true
}

// ExpansionTree is a tree of expansion helpers.
type ExpansionTree[T gr.TokenTyper] struct {
	// tree is the tree of expansion helpers.
	tree *tr.Tree

	// info is the information about the expansion tree.
	info *InfoStruct[T]
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
func NewExpansionTreeRootedAt[T gr.TokenTyper](cs *ConflictSolver[T], h *HelperNode[T]) (*ExpansionTree[T], error) {
	info := NewInfoStruct(h)

	nexts_func := func(data tn.Noder, is tr.Infoer) ([]tn.Noder, error) {
		is_inf, ok := is.(*InfoStruct[T])
		if !ok {
			return nil, uc.NewErrUnexpectedType("is", is)
		}

		hn, ok := data.(*HelperNode[T])
		if !ok {
			return nil, uc.NewErrUnexpectedType("HelperNode[T] node", data)
		}

		rhs, err := hn.GetRhsAt(0)
		if err != nil {
			return nil, NewErr0thRhsNotSet()
		}

		ok = rhs.IsTerminal()
		if ok {
			return nil, nil
		}

		result := cs.GetElemsWithLhs(rhs)

		result = us.SliceFilter(result, is_inf.IsNotSeen)

		is_inf.SetSeen(hn)

		var children []tn.Noder

		for _, r := range result {
			children = append(children, r)
		}

		return children, nil
	}

	var builder tr.Builder

	builder.SetInfo(info)
	builder.SetNextFunc(nexts_func)

	tree, err := builder.Build(h)
	if err != nil {
		return nil, err
	}

	ext := &ExpansionTree[T]{
		tree: tree,
		info: info,
	}

	return ext, nil
}

// PruneNonTerminalLeaves prunes the non-terminal leaves of the expansion tree.
func (et *ExpansionTree[T]) PruneNonTerminalLeaves() {
	leaves := et.tree.GetLeaves()

	todo := us.SliceFilter(leaves, FilterNonTerminalLeaf[T])
	if len(todo) == 0 {
		return
	}

	for _, leaf := range todo {
		err := et.tree.DeleteBranchContaining(leaf)
		uc.AssertF(err == nil, "unexpected error: %s", err.Error())
	}
}

// Size returns the size of the expansion tree.
//
// Returns:
//   - int: The size of the expansion tree.
func (et *ExpansionTree[T]) Size() int {
	size := et.tree.Size()
	return size
}

// Collapse collapses the expansion tree into a slice of strings that
// are the 0th RHS of the terminal leaves.
//
// Returns:
//   - []T: The collapsed expansion tree.
func (et *ExpansionTree[T]) Collapse() []T {
	leaves := et.tree.GetLeaves()

	var result []T

	for _, leaf := range leaves {
		tn, ok := leaf.(*HelperNode[T])
		uc.Assert(ok, "Must be a *HelperNode[T]")

		rhs, err := tn.GetRhsAt(0)
		uc.AssertF(err == nil, "unexpected error: %s", err.Error())

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
