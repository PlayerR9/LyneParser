package Lexer

import (
	"fmt"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

type TreeNode[T gr.TokenTyper] struct {
	*tr.StatusInfo[EvalStatus, *gr.Token[T]]
}

func new_tree_node[T gr.TokenTyper](value *gr.Token[T]) *TreeNode[T] {
	si := tr.NewStatusInfo(value, EvalIncomplete)

	tn := &TreeNode[T]{
		StatusInfo: si,
	}

	return tn
}

func convert_branch[T gr.TokenTyper](branch *tr.Branch[*TreeNode[T]]) []*gr.Token[T] {
	slice := branch.Slice()
	slice = slice[1:] // Skip the root.

	result := make([]*gr.Token[T], 0, len(slice))

	for _, tn := range slice {
		token := tn.GetData()

		result = append(result, token)
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
