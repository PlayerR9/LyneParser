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

type ErrLexerError struct {
	At   int
	Prev []*gr.LeafToken
}

func (e *ErrLexerError) Error() string {
	var builder strings.Builder

	builder.WriteString("no matches found at ")
	builder.WriteString(strconv.Itoa(e.At))

	return builder.String()
}

func NewErrLexerError(at int, prev []*gr.LeafToken) *ErrLexerError {
	e := &ErrLexerError{
		At:   at,
		Prev: prev,
	}
	return e
}

type TreeNode struct {
	*tr.StatusInfo[EvalStatus, *gr.LeafToken]
}

func newTreeNode(value *gr.LeafToken) *TreeNode {
	si := tr.NewStatusInfo(value, EvalIncomplete)

	tn := &TreeNode{
		StatusInfo: si,
	}

	return tn
}

func convBranch(branch *tr.Branch[*TreeNode]) []*gr.LeafToken {
	slice := branch.Slice()
	slice = slice[1:] // Skip the root.

	result := make([]*gr.LeafToken, 0, len(slice))

	for _, tn := range slice {
		token := tn.GetData()

		result = append(result, token)
	}

	return result
}

func lastOfBranch(branch []*gr.LeafToken) int {
	len := len(branch)

	if len == 0 {
		return -1
	}

	last := branch[len-1]

	return last.At
}

func Lex(s *cds.Stream[byte], productions []*gr.RegProduction, v *Verbose) error {
	f := func(tn *TreeNode) ([]*TreeNode, error) {
		data := tn.GetData()

		var nextAt int

		if data.ID == gr.RootTokenID {
			nextAt = 0
		} else {
			nextAt = data.At + len(data.Data)
		}

		results, err := matchFrom(s, nextAt, productions)
		if err != nil {
			return nil, err
		}

		// Get the longest match.
		results = selectBestMatches(results, v)

		children := make([]*TreeNode, 0, len(results))

		for _, result := range results {
			curr := result.GetMatch()
			tn := newTreeNode(curr)

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
