package Lexer

import (
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// ErrNoMatches is an error that is returned when there are no
// matches at a position.
type ErrNoMatches struct{}

// Error returns the error message: "no matches".
//
// Returns:
//   - string: The error message.
func (e *ErrNoMatches) Error() string {
	return "no matches"
}

// NewErrNoMatches creates a new error of type *ErrNoMatches.
//
// Returns:
//   - *ErrNoMatches: The new error.
func NewErrNoMatches() *ErrNoMatches {
	return &ErrNoMatches{}
}

// ErrAllMatchesFailed is an error that is returned when all matches
// fail.
type ErrAllMatchesFailed struct{}

// Error returns the error message: "all matches failed".
//
// Returns:
//   - string: The error message.
func (e *ErrAllMatchesFailed) Error() string {
	return "all matches failed"
}

// NewErrAllMatchesFailed creates a new error of type *ErrAllMatchesFailed.
//
// Returns:
//   - *ErrAllMatchesFailed: The new error.
func NewErrAllMatchesFailed() *ErrAllMatchesFailed {
	return &ErrAllMatchesFailed{}
}

// ErrInvalidElement is an error that is returned when an invalid element
// is found.
type ErrInvalidElement struct{}

// Error returns the error message: "invalid element".
//
// Returns:
//   - string: The error message.
func (e *ErrInvalidElement) Error() string {
	return "invalid element"
}

// NewErrInvalidElement creates a new error of type *ErrInvalidElement.
//
// Returns:
//   - *ErrInvalidElement: The new error.
func NewErrInvalidElement() *ErrInvalidElement {
	return &ErrInvalidElement{}
}

// IsDone checks if an error is a completion error or nil.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - bool: True if the error is a completion error or nil.
//     False otherwise.
func IsDone(err error) bool {
	if err == nil {
		return true
	}

	ok := uc.Is[*uc.ErrExhaustedIter](err)
	return ok
}
