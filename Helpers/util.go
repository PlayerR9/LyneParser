package Helpers

// Remove HResult when MyGoLib is updated
type HResult[T any] struct {
	Result T
	Reason error
}

func NewHResult[T any](result T, reason error) HResult[T] {
	return HResult[T]{
		Result: result,
		Reason: reason,
	}
}

func EvaluateFunc[T any](f func() (T, error)) HResult[T] {
	res, err := f()
	return HResult[T]{
		Result: res,
		Reason: err,
	}
}

func EvaluateMany[T any](f func() ([]T, error)) []HResult[T] {
	res, err := f()

	if len(res) == 0 {
		return []HResult[T]{NewHResult[T](*new(T), err)}
	}

	results := make([]HResult[T], len(res))

	for i, r := range res {
		results[i] = NewHResult(r, err)
	}

	return results
}
