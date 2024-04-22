package Parser

import (
	"errors"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	hp "github.com/PlayerR9/LyneParser/Helpers"
)

func ParseBranch(parser *Parser, branch []gr.LeafToken) ([]gr.NonLeafToken, error) {
	err := parser.SetInputStream(branch)
	if err != nil {
		return nil, hp.NewErrIgnorable(err)
	}

	err = parser.Parse()
	if err != nil {
		return nil, err
	}

	roots, err := parser.GetParseTree()
	if err != nil {
		return roots, hp.NewErrIgnorable(err)
	}

	if len(roots) == 0 {
		return nil, hp.NewErrIgnorable(errors.New("no roots found"))
	}

	return roots, nil
}

func ParseIS(parser *Parser, branches [][]gr.LeafToken) ([]gr.NonLeafToken, error) {
	solutions := make([]hp.HResult[gr.NonLeafToken], 0)

	for _, branch := range branches {
		results := hp.EvaluateMany(func() ([]gr.NonLeafToken, error) {
			return ParseBranch(parser, branch)
		})

		solutions = append(solutions, results...)
	}

	// Filter out solutions with errors
	// FIXME: Finish this
	for i := 0; i < len(solutions); {
		if solutions[i].Reason != nil {
			if len(solutions) == 1 {
				return nil, solutions[i].Reason
			}

			solutions = append(solutions[:i], solutions[i+1:]...)
		} else {
			i++
		}
	}

	if len(solutions) == 0 {
		return nil, errors.New("no solutions found")
	}

	// Extract the results
	extracted := make([]gr.NonLeafToken, len(solutions))

	for i, sol := range solutions {
		extracted[i] = sol.Result
	}

	return extracted, nil
}
