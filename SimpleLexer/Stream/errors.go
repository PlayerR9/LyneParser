package Stream

import "errors"

type ErrStreamExhausted struct{}

func (e *ErrStreamExhausted) Error() string {
	return "stream is exhausted; no more data"
}

func NewErrStreamExhausted() *ErrStreamExhausted {
	e := &ErrStreamExhausted{}
	return e
}

func IsExhausted(err error) bool {
	var exhaustedErr *ErrStreamExhausted

	ok := errors.As(err, &exhaustedErr)
	return ok
}
