package FrontierEvaluation

import (
	hlp "github.com/PlayerR9/MyGoLib/CustomData/Helpers"
	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

type EvalManyFunc[E, R any] func(E) ([]R, error)

func FrontierEvaluation[T any](todo []T, accept slext.PredicateFilter[T], f EvalManyFunc[T, T]) ([]hlp.HResult[T], bool) {
	if len(todo) == 0 {
		// Nothing to do
		return nil, true
	} else if accept == nil {
		// Nothing to accept
		results := make([]hlp.HResult[T], 0, len(todo))

		for _, t := range todo {
			results = append(results, hlp.NewHResult(t, NewErrNoAcceptance()))
		}

		return results, false
	}

	if f == nil {
		solutions, ok := slext.SFSeparateEarly(todo, accept)

		results := make([]hlp.HResult[T], 0, len(solutions))

		if ok {
			for _, s := range solutions {
				results = append(results, hlp.NewHResult(s, nil))
			}
		} else {
			for _, s := range solutions {
				results = append(results, hlp.NewHResult(s, NewErrNoAcceptance()))
			}
		}

		return results, ok
	}

	results := make([]hlp.HResult[T], 0)

	success, fail := slext.SFSeparate(todo, accept)
	if len(success) > 0 {
		for _, s := range success {
			results = append(results, hlp.NewHResult(s, nil))
		}
	}

	todo = fail

	newTodo := make([]T, 0)

	for _, t := range todo {
		tmp, err := f(t)
		if err != nil {
			results = append(results, hlp.NewHResult(t, err))
		} else {
			newTodo = append(newTodo, tmp...)
		}
	}

	for _, t := range todo {
		results, err := f(t)
	}

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
