package Parser

import (
	"errors"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"

	com "github.com/PlayerR9/LyneParser/Common"
)

/////////////////////////////////////////////////////////////

// ParseBranch is a function that, given a parser and an input stream of tokens,
// returns a slice of non-leaf tokens.
//
// Parameters:
//
//   - parser: The parser to use.
//   - inputStream: The input stream of tokens to parse.
//
// Returns:
//
//   - []gr.NonLeafToken: A slice of non-leaf tokens.
//   - error: An error if the branch cannot be parsed.
func ParseBranch(parser *Parser, source *com.TokenStream) ([]*com.TokenTree, error) {
	err := parser.Parse(source)
	if err != nil {
		return nil, err
	}

	roots, err := parser.GetParseTree()
	if err != nil {
		return roots, ers.NewErrIgnorable(err)
	}

	if len(roots) == 0 {
		return nil, ers.NewErrIgnorable(errors.New("no roots found"))
	}

	return roots, nil
}

/*

// ParseIS is a function that, given a parser and a slice of branches of tokens,
// returns a slice of non-leaf tokens.
//
// Parameters:
//
//   - parser: The parser to use.
//   - branches: The branches of tokens to parse.
//
// Returns:
//
//   - []gr.NonLeafToken: A slice of non-leaf tokens.
//   - error: An error if the branches cannot be parsed.
func ParseIS(parser *Parser, branches []*com.TokenStream) ([]gr.NonLeafToken, error) {
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
		if solutions[i].Second != nil {
			if len(solutions) == 1 {
				return nil, solutions[i].Second
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
		extracted[i] = sol.First
	}

	return extracted, nil
}
*/
