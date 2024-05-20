package FrontierEvaluation

import (
	hlp "github.com/PlayerR9/MyGoLib/CustomData/Helpers"
	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/Common"
)

// Accepter is an interface that represents an accepter.
type Accepter interface {
	// Accept returns true if the accepter accepts the element.
	//
	// Returns:
	//   - bool: True if the accepter accepts the element, false otherwise.
	Accept() bool
}

// Evaluate is the main function of the tree evaluator.
//
// Parameters:
//   - source: The source to evaluate.
//   - root: The root of the tree evaluator.
//
// Returns:
//   - error: An error if lexing fails.
//
// Errors:
//   - *ErrEmptyInput: The source is empty.
//   - *ers.ErrAt: An error occurred at a specific index.
//   - *ErrAllMatchesFailed: All matches failed.
func FrontierEvaluate[T Accepter](elem T, matcher uc.EvalManyFunc[T, T]) ([]T, error) {
	if matcher == nil {
		return nil, nil
	} else if elem.Accept() {
		return []T{elem}, nil
	}

	solutions := make([]hlp.HResult[T], 0)

	S := Stacker.NewArrayStack(elem)

	for {
		elem, err := S.Pop()
		if err != nil {
			break
		}

		nexts, err := matcher(elem)
		if err != nil {
			solutions = append(solutions, hlp.NewHResult(elem, err))

			continue
		}

		for _, next := range nexts {
			if next.Accept() {
				solutions = append(solutions, hlp.NewHResult(next, nil))

				continue
			} else {
				err := S.Push(next)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	solutions, ok := hlp.FilterSuccessOrFail(solutions)
	if !ok {
		// Determine the most likely error.
		// As of now, we will just return the first error.
		return nil, solutions[0].Second
	}

	// TODO: Fix once the MyGoLib is updated.

	result := make([]T, 0, len(solutions))

	for _, solution := range solutions {
		result = append(result, solution.First)
	}

	return result, nil
}
