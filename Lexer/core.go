package Lexer

import (
	"fmt"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	tr "github.com/PlayerR9/tree/tree"
)

func convert_branch[T gr.TokenTyper](branch *tr.Branch) []*gr.Token[T] {
	slice := branch.Slice()
	slice = slice[1:] // Skip the root.

	result := make([]*gr.Token[T], 0, len(slice))

	for _, tn := range slice {
		n, ok := tn.(*TokenNode[T])
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
	f := func(tn *TokenNode[T]) ([]*TokenNode[T], error) {
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
