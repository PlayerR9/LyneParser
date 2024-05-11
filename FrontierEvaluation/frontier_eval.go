package FrontierEvaluation

import (
	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

type EvalManyFunc[E, R any] func(E) ([]R, error)

func FrontierEvaluation[T any](todo []T, accept slext.PredicateFilter[T], f EvalManyFunc[T, T]) ([]T, error) {
	if len(todo) == 0 {
		return nil, nil
	}

	if accept == nil {
		return nil, nil
	}

	if f == nil {
		solution, ok := slext.SFSeparateEarly(todo, accept)
		if ok {
			return solution, nil
		}

		return solution, nil //
	}

	newTodo := make([]T, 0)
	lastErrors := make([]error, 0)

	for _, t := range todo {
		ress, err := f(t)
		if err != nil {
			lastErrors = append(lastErrors, err)
		} else {
			newTodo = append(newTodo, ress...)
		}
	}

	done := make([]T, 0)

	for len(todo) > 0 {
		newElem := make([]T, 0)

		s1, s2 := slext.SFSeparate(todo, accept)
		if len(s1) > 0 {
			done = append(done, s1...)
		}

		if len(s2) == 0 {
			break
		}

		for _, elem := range s2 {
			others, err := f(elem)
			if err != nil {
				continue
			}

			newElem = append(newElem, others...)
		}

		todo = newElem
	}

	return done, nil
}
