package Helpers

import "errors"

type ErrIgnorable struct {
	Err error
}

func (e ErrIgnorable) Error() string {
	return e.Err.Error()
}

func NewErrIgnorable(err error) ErrIgnorable {
	return ErrIgnorable{
		Err: err,
	}
}

func IsErrIgnorable(err error) bool {
	var ignorable *ErrIgnorable

	return errors.As(err, &ignorable)
}
