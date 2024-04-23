package Helpers

import "errors"

// ErrIgnorable is an error that can be ignored.
type ErrIgnorable struct {
	// Err is the error.
	Err error
}

// Error is a method of error interface.
//
// Returns:
//
//   - string: The error message.
func (e ErrIgnorable) Error() string {
	return e.Err.Error()
}

// NewErrIgnorable creates a new ErrIgnorable.
//
// Parameters:
//
//   - err: The error.
//
// Returns:
//
//   - ErrIgnorable: The new error.
func NewErrIgnorable(err error) *ErrIgnorable {
	return &ErrIgnorable{
		Err: err,
	}
}

// IsErrIgnorable checks if an error is ignorable.
//
// Parameters:
//
//   - err: The error to check.
//
// Returns:
//
//   - bool: True if the error is ignorable, false otherwise.
func IsErrIgnorable(err error) bool {
	var ignorable *ErrIgnorable

	return errors.As(err, &ignorable)
}
