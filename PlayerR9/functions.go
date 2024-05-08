package PlayerR9

import (
	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

// EvalManyFunc is a function that evaluates many elements.
//
// Parameters:
//   - T: The type of elements to evaluate.
//
// Returns:
//   - []T: The elements that were evaluated.
//   - error: An error if the evaluation failed.
type EvalManyFunc[T any] func(T) ([]T, error)

// DoWhile performs a do-while loop on a slice of elements.
//
// Parameters:
//   - todo: The elements to perform the do-while loop on.
//   - accept: The predicate filter to accept elements.
//   - f: The evaluation function to perform on the elements.
//
// Returns:
//   - []T: The elements that were accepted.
func DoWhile[T any](todo []T, accept slext.PredicateFilter[T], f EvalManyFunc[T]) []T {
	if len(todo) == 0 {
		return nil
	} else if accept == nil {
		return nil
	} else if f == nil {
		s1, _ := slext.SFSeparate(todo, accept)

		return s1
	}

	done := make([]T, 0)

	for len(todo) > 0 {
		s1, s2 := slext.SFSeparate(todo, accept)
		if len(s1) > 0 {
			done = append(done, s1...)
		}

		if len(s2) == 0 {
			break
		}

		newElem := make([]T, 0)

		for _, elem := range s2 {
			others, err := f(elem)
			if err != nil || len(others) == 0 {
				continue
			}

			newElem = append(newElem, others...)
		}

		if len(newElem) == 0 {
			break
		}

		todo = newElem
	}

	return done
}
