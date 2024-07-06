package Lexer

import (
	"fmt"
	"strconv"
	"strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

type ErrLexerError[T uc.Enumer] struct {
	At   int
	Prev []*gr.Token[T]
}

func (e *ErrLexerError[T]) Error() string {
	var builder strings.Builder

	builder.WriteString("no matches found at ")
	builder.WriteString(strconv.Itoa(e.At))

	return builder.String()
}

func NewErrLexerError[T uc.Enumer](at int, prev []*gr.Token[T]) *ErrLexerError[T] {
	e := &ErrLexerError[T]{
		At:   at,
		Prev: prev,
	}
	return e
}

type TreeNode[T uc.Enumer] struct {
	*tr.StatusInfo[EvalStatus, *gr.Token[T]]
}

func newTreeNode[T uc.Enumer](value *gr.Token[T]) *TreeNode[T] {
	si := tr.NewStatusInfo(value, EvalIncomplete)

	tn := &TreeNode[T]{
		StatusInfo: si,
	}

	return tn
}

func convBranch[T uc.Enumer](branch *tr.Branch[*TreeNode[T]]) []*gr.Token[T] {
	slice := branch.Slice()
	slice = slice[1:] // Skip the root.

	result := make([]*gr.Token[T], 0, len(slice))

	for _, tn := range slice {
		token := tn.GetData()

		result = append(result, token)
	}

	return result
}

func lastOfBranch[T uc.Enumer](branch []*gr.Token[T]) int {
	len := len(branch)

	if len == 0 {
		return -1
	}

	last := branch[len-1]

	return last.At
}

func Lex[T uc.Enumer](s *cds.Stream[byte], productions []*gr.RegProduction[T], v *Verbose) error {
	f := func(tn *TreeNode[T]) ([]*TreeNode[T], error) {
		data := tn.GetData()

		var nextAt int

		if data.ID.String() == gr.RootTokenID {
			nextAt = 0
		} else {
			nextAt = data.At + len(data.Data.(string))
		}

		results, err := matchFrom(s, nextAt, productions)
		if err != nil {
			return nil, err
		}

		// Get the longest match.
		results = selectBestMatches(results, v)

		children := make([]*TreeNode[T], 0, len(results))

		for _, result := range results {
			tn := newTreeNode(result.Matched)

			children = append(children, tn)
		}

		return children, nil
	}

	iter := newCoreIter(f)

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
