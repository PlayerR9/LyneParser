package Display

import (
	"errors"
)

// ErrESCPressed is an error that indicates that the Ctrl+C key combination was pressed.
type ErrESCPressed struct{}

// Error is a method of errors that returns the error message.
//
// Returns:
//   - string: The error message.
func (e *ErrESCPressed) Error() string {
	return "Ctrl+C was pressed"
}

// NewErrESCPressed creates a new ErrESCPressed error.
//
// Returns:
//   - *ErrESCPressed: A pointer to the new error.
func NewErrESCPressed() *ErrESCPressed {
	return &ErrESCPressed{}
}

// IsErrESCPressed checks if the given error is an ErrESCPressed error.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - bool: True if the error is an ErrESCPressed error, false otherwise.
func IsErrESCPressed(err error) bool {
	var escPressed *ErrESCPressed

	return errors.As(err, &escPressed)
}
