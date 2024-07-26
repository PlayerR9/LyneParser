package Lexer

import (
	"fmt"
	"slices"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
	tr "github.com/PlayerR9/tree/tree"
)

type CoreIter[T gr.TokenTyper] struct {
	do_func uc.EvalManyFunc[tr.Noder, tr.Noder]
	tree    *tr.Tree

	// data is the data.
	data []rune

	// Tokens is the list of tokens.
	Tokens []*gr.Token[T]

	// idx is the index.
	idx int

	// SyntaxError is the syntax error.
	SyntaxError SyntaxErrorer
}

func (it *CoreIter[T]) Size() (count int) {
	count = it.tree.Size()
	return
}

func (it *CoreIter[T]) can_continue() bool {
	size := it.tree.Size()

	if size != 1 {
		return true
	}

	root := it.tree.Root()

	n, ok := root.(*TokenNode[T])
	uc.Assert(ok, "Must be a *TokenNode[T]")

	return n.Status != EvalComplete
}

func (it *CoreIter[T]) Consume() ([][]*gr.Token[T], error) {
	var errs uc.ErrOrSol[any]

	for {
		ok := it.can_continue()
		if !ok {
			break
		}

		err := it.tree.ProcessLeaves(it.do_func)
		if err != nil {
			return nil, fmt.Errorf("could not process leaves: %w", err)
		}

		leaves := it.tree.GetLeaves()

		f := func(n tr.Noder) bool {
			tn, ok := n.(*TokenNode[T])
			if !ok {
				return false
			}

			ok = tr.Status != EvalIncomplete
			return ok
		}

		leaves_done := us.SliceFilter(leaves, f)

		var results [][]*gr.Token[T]

		for i, leaf := range leaves_done {
			// Extract the branch.
			branch, err := it.tree.ExtractBranch(leaf, true)
			if err != nil {
				return nil, uc.NewErrWhileAt("extracting", i+1, "branch", err)
			}

			tn, ok := leaf.(*TokenNode[T])
			uc.Assert(ok, "Must be a *TokenNode[T]")

			converted := convert_branch[T](branch)
			level := last_of_branch(converted)

			if tr.Status == EvalError {
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

func (it *CoreIter[T]) Restart() {
	// root := gr.RootToken()

	// tn := newTreeNode(root)

	// tree := tr.NewTree(tn)
	// it.tree = tree

	panic("Restart not implemented yet")
}

func new_core_iter[T gr.TokenTyper](doFunc uc.EvalManyFunc[*TokenNode[T], *TokenNode[T]]) *CoreIter[T] {
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

type ActiveLexer[T gr.TokenTyper] struct {
	ci   *CoreIter[T]
	iter uc.Iterater[[]*gr.Token[T]]
}

func NewActiveLexer[T gr.TokenTyper]() *ActiveLexer[T] {
	al := &ActiveLexer[T]{}
	return al
}

// GetTokens gets the tokens.
//
// Returns:
//   - []*gr.Token: The tokens. Nil if lexer has finished.
//   - SyntaxErrorer: The syntax error.
func (al *ActiveLexer[T]) GetTokens() ([]*gr.Token[T], SyntaxErrorer) {
	for {
		if al.iter == nil {
			return nil, nil
		}

		tokens, err := al.iter.Consume()
		if err == nil {
			return tokens, nil
		}

		ok := uc.Is[*uc.ErrExhaustedIter](err)
		if !ok {
			se := NewGenericSyntaxError(0, err.Error(), "")
			return nil, se
		}

		branches, err := al.ci.Consume()
		if err == nil {
			al.iter = uc.NewSimpleIterator(branches)
		} else {
			ok := uc.Is[*uc.ErrExhaustedIter](err)
			if ok {
				al.iter = nil
			} else {
				se := NewGenericSyntaxError(0, err.Error(), "")
				return nil, se
			}
		}
	}
}

type SolutionIterator[T any] struct {
	iter   uc.Iterater[T]
	source uc.Iterater[[]T]
}

func (si *SolutionIterator[T]) Consume() (T, error) {
	var value T
	var values []T
	var err error

	for {
		if si.iter != nil {
			value, err = si.iter.Consume()
			if err == nil {
				break
			}

			ok := uc.Is[*uc.ErrExhaustedIter](err)
			if !ok {
				break
			}
		}

		values, err = si.source.Consume()
		if err != nil {
			break
		}

		si.iter = uc.NewSimpleIterator(values)
	}

	if err != nil {
		return *new(T), err
	} else {
		return value, nil
	}
}

func NewSolutionIterator[T any](source uc.Iterater[[]T]) *SolutionIterator[T] {
	si := &SolutionIterator[T]{
		source: source,
	}

	return si
}
