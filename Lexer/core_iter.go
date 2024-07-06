package Lexer

import (
	"slices"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

type CoreIter[T uc.Enumer] struct {
	doFunc uc.EvalManyFunc[*TreeNode[T], *TreeNode[T]]
	tree   *tr.Tree[*TreeNode[T]]
}

func (it *CoreIter[T]) Size() (count int) {
	count = it.tree.Size()
	return
}

func (it *CoreIter[T]) canContinue() bool {
	size := it.tree.Size()

	if size != 1 {
		return true
	}

	root := it.tree.Root()
	status := root.Data.GetStatus()

	return status != EvalComplete
}

func (it *CoreIter[T]) Consume() ([][]*gr.Token[T], error) {
	var errs uc.ErrOrSol[any]

	for {
		ok := it.canContinue()
		if !ok {
			ok := errs.HasError()
			if ok {
				errs := errs.GetErrors()

				return nil, uc.NewErrPossibleError(
					NewErrAllMatchesFailed(),
					errs[0],
				)
			}

			return nil, uc.NewErrExhaustedIter()
		}

		err := it.tree.ProcessLeaves(it.doFunc)
		if err != nil {
			return nil, err
		}

		leaves := it.tree.GetLeaves()

		f := func(tn *tr.TreeNode[*TreeNode[T]]) bool {
			data := tn.Data
			status := data.GetStatus()

			ok := status != EvalIncomplete
			return ok
		}

		leavesDone := us.SliceFilter(leaves, f)

		var results [][]*gr.Token[T]

		for _, leaf := range leavesDone {
			// Extract the branch.
			branch := it.tree.ExtractBranch(leaf, true)
			status := leaf.Data.GetStatus()

			converted := convBranch(branch)
			level := lastOfBranch(converted)

			if status == EvalError {
				err := NewErrLexerError(level, converted)

				errs.AddErr(err, level)
			} else {
				results = append(results, converted)
			}
		}

		if len(results) > 0 {
			f := func(a, b []*gr.Token[T]) int {
				return len(a) - len(b)
			}

			slices.SortStableFunc(results, f)

			return results, nil
		}
	}
}

func (it *CoreIter[T]) Restart() {
	// root := gr.RootToken()

	// tn := newTreeNode(root)

	// tree := tr.NewTree(tn)
	// it.tree = tree

	panic("Restart not implemented yet")
}

func newCoreIter[T uc.Enumer](doFunc uc.EvalManyFunc[*TreeNode[T], *TreeNode[T]]) *CoreIter[T] {
	// root := gr.RootToken()

	// tn := newTreeNode(root)

	// tree := tr.NewTree(tn)

	// it := &CoreIter[T]{
	// 	tree:   tree,
	// 	doFunc: doFunc,
	// }

	// return it

	panic("newCoreIter not implemented yet")
}
