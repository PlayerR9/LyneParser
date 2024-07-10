package Lexer

import (
	"fmt"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// TreeNode is a node in the tree.
type TreeNode[T gr.TokenTyper] struct {
	Parent   *TreeNode[T]
	Children []*TreeNode[T]

	Token  *gr.Token[T]
	Status EvalStatus
}

// SetParent implements the Noder interface.
func (tn *TreeNode[T]) SetParent(parent tr.Noder) bool {
	n, ok := parent.(*TreeNode[T])
	if !ok {
		return false
	}

	tn.Parent = n

	return true
}

// GetParent implements the Noder interface.
func (tn *TreeNode[T]) GetParent() tr.Noder {
	return tn.Parent
}

// LinkWithParent implements the Noder interface.
func (tn *TreeNode[T]) LinkChildren(children []tr.Noder) {
	var valid_children []*TreeNode[T]

	for _, child := range children {
		n, ok := child.(*TreeNode[T])
		if !ok {
			continue
		}

		valid_children = append(valid_children, n)
	}

	for _, child := range valid_children {
		child.Parent = tn
	}

	tn.Children = valid_children
}

// IsLeaf implements the Noder interface.
func (tn *TreeNode[T]) IsLeaf() bool {
	return len(tn.Children) == 0
}

// IsSingleton implements the Noder interface.
func (tn *TreeNode[T]) IsSingleton() bool {
	return len(tn.Children) == 1
}

// GetLeaves implements the Noder interface.
func (tn *TreeNode[T]) GetLeaves() []tr.Noder {
	var leaves []tr.Noder

	for _, child := range tn.Children {
		ok := child.IsLeaf()
		if ok {
			leaves = append(leaves, child)
			continue
		}

		sub_leaves := child.GetLeaves()
		leaves = append(leaves, sub_leaves...)
	}

	return leaves
}

// GetAncestors implements the Noder interface.
func (tn *TreeNode[T]) GetAncestors() []tr.Noder {
	panic("Not implemented")
}

// GetFirstChild implements the Noder interface.
func (tn *TreeNode[T]) GetFirstChild() tr.Noder {
	panic("Not implemented")
}

// DeleteChild implements the Noder interface.
func (tn *TreeNode[T]) DeleteChild(target tr.Noder) []tr.Noder {
	panic("Not implemented")
}

// Size implements the Noder interface.
func (tn *TreeNode[T]) Size() int {
	panic("Not implemented")
}

// AddChild implements the Noder interface.
func (tn *TreeNode[T]) AddChild(child tr.Noder) {
	panic("Not implemented")
}

// RemoveNode implements the Noder interface.
func (tn *TreeNode[T]) RemoveNode() []tr.Noder {
	panic("Not implemented")
}

// TreeOf implements the Noder interface.
func (tn *TreeNode[T]) TreeOf() *tr.Tree {
	panic("Not implemented")
}

// Iterator implements the Noder interface.
func (tn *TreeNode[T]) Iterator() uc.Iterater[tr.Noder] {
	panic("Not implemented")
}

// Copy implements the Noder interface.
func (tn *TreeNode[T]) Copy() uc.Copier {
	panic("Not implemented")
}

// FString implements the Noder interface.
func (tn *TreeNode[T]) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	panic("Not implemented")
}

// Cleanup implements the Noder interface.
func (tn *TreeNode[T]) Cleanup() {
	panic("Not implemented")
}

func NewTreeNode[T gr.TokenTyper](value *gr.Token[T]) *TreeNode[T] {
	tn := &TreeNode[T]{
		Token:  value,
		Status: EvalIncomplete,
	}

	return tn
}

func convert_branch[T gr.TokenTyper](branch *tr.Branch) []*gr.Token[T] {
	slice := branch.Slice()
	slice = slice[1:] // Skip the root.

	result := make([]*gr.Token[T], 0, len(slice))

	for _, tn := range slice {
		n, ok := tn.(*TreeNode[T])
		uc.Assert(ok, "Must be a *TreeNode[T]")

		result = append(result, n.Token)
	}

	return result
}

func last_of_branch[T gr.TokenTyper](branch []*gr.Token[T]) int {
	len := len(branch)

	if len == 0 {
		return -1
	}

	last := branch[len-1]

	return last.At
}

func Lex[T gr.TokenTyper](s *cds.Stream[byte], productions []*gr.RegProduction[T], v *Verbose) error {
	f := func(tn *TreeNode[T]) ([]*TreeNode[T], error) {
		// data := tn.GetData()

		// var nextAt int

		// if data.ID.String() == gr.RootTokenID {
		// 	nextAt = 0
		// } else {
		// 	nextAt = data.At + len(data.Data.(string))
		// }

		// results, err := matchFrom(s, nextAt, productions)
		// if err != nil {
		// 	return nil, err
		// }

		// Get the longest match.
		// results = selectBestMatches(results, v)

		// children := make([]*TreeNode[T], 0, len(results))

		// for _, result := range results {
		// 	tn := newTreeNode(result.Matched)

		// 	children = append(children, tn)
		// }

		// return children, nil

		panic("Lex: not implemented")
	}

	iter := new_core_iter(f)

	for {
		branches, err := iter.Consume()
		if err != nil {
			ok := uc.Is[*uc.ErrExhaustedIter](err)
			if !ok {
				return err
			}

			break
		}

		// DEBUG: Print solutions
		fmt.Println("Solutions:")

		for i, solution := range branches {
			fmt.Printf("Solution %d:\n", i)

			for _, token := range solution {
				fmt.Printf("%+v\n", token)
			}

			fmt.Println()
		}
	}

	return nil
}
