package Parser

import (
	"errors"

	cs "github.com/PlayerR9/LyneParser/ConflictSolver"
	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// evaluate evaluates the frontier evaluator given an element.
//
// Parameters:
//   - elem: The element to evaluate.
//
// Behaviors:
//   - If the element is accepted, the solutions will be set to the element.
//   - If the element is not accepted, the solutions will be set to the results of the matcher.
//   - If the matcher returns an error, the solutions will be set to the error.
//   - The evaluations assume that, the more the element is elaborated, the more the weight increases.
//     Thus, it is assumed to be the most likely solution as it is the most elaborated. Euristic: Depth.
func evaluate(dt *cs.ConflictSolver, source *cds.Stream[gr.Token], elem *CurrentEval) []*us.WeightedHelper[*CurrentEval] {
	ok := elem.Accept()
	if ok {
		h := us.NewWeightedHelper(elem, nil, 0.0)

		sols := []*us.WeightedHelper[*CurrentEval]{h}

		return sols
	}

	var sols []*us.WeightedHelper[*CurrentEval]

	p := uc.NewPair(elem, 0.0)
	S := Stacker.NewArrayStack(p)

	for {
		p, ok := S.Pop()
		if !ok {
			break
		}

		nexts, err := p.First.Parse(source, dt)
		if err != nil {
			h := us.NewWeightedHelper(p.First, err, p.Second)

			sols = append(sols, h)
			continue
		}

		newPairs := make([]uc.Pair[*CurrentEval, float64], 0, len(nexts))

		for _, next := range nexts {
			p := uc.NewPair(next, p.Second+1.0)

			newPairs = append(newPairs, p)
		}

		for _, pair := range newPairs {
			ok := pair.First.Accept()
			if ok {
				h := us.NewWeightedHelper(pair.First, nil, pair.Second)

				sols = append(sols, h)
			} else {
				S.Push(pair)
			}
		}
	}

	return sols
}

// extractResults gets the results of the frontier evaluator.
//
// Returns:
//   - []T: The results of the frontier evaluator.
//   - error: An error if the frontier evaluator failed.
//
// Behaviors:
//   - If the solutions are empty, the function returns nil.
//   - If the solutions contain errors, the function returns the first error.
//   - Otherwise, the function returns the solutions.
func extractResults(sols []*us.WeightedHelper[*CurrentEval]) ([]*CurrentEval, error) {
	if len(sols) == 0 {
		return nil, nil
	}

	results, ok := us.SuccessOrFail(sols, true)

	extracted := us.ExtractResults(results)

	if !ok {
		// Determine the most likely error.
		// As of now, we will just return the first error.
		data := results[0].GetData()
		return extracted, data.Second
	} else {
		return extracted, nil
	}
}

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
func ParseBranch(parser *Parser, source *cds.Stream[gr.Token]) ([]*gr.TokenTree, error) {
	err := Parse(parser, source)
	if err != nil {
		return nil, err
	}

	roots, err := parser.GetParseTree()
	if err != nil {
		return roots, uc.NewErrIgnorable(err)
	}

	if len(roots) == 0 {
		return nil, uc.NewErrIgnorable(errors.New("no roots found"))
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
func ParseIS(parser *Parser, branches []*cds.Stream[*LeafToken]) ([]gr.NonLeafToken, error) {
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
