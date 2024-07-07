package ConflictSolver

import (
	"fmt"

	gr "github.com/PlayerR9/LyneParser/Grammar"
)

// ErrHelpersConflictingSize is an error that is returned when helpers have conflicting sizes.
type ErrHelpersConflictingSize struct{}

// Error implements the error interface.
//
// Message: "helpers have conflicting sizes".
func (e *ErrHelpersConflictingSize) Error() string {
	return "helpers have conflicting sizes"
}

// NewErrHelpersConflictingSize creates a new error of type *ErrHelpersConflictingSize.
//
// Returns:
//   - *ErrHelpersConflictingSize: A pointer to the new error.
func NewErrHelpersConflictingSize() *ErrHelpersConflictingSize {
	e := &ErrHelpersConflictingSize{}
	return e
}

// Err0thRhsNotSet is an error that is returned when the 0th right-hand side is not set.
type Err0thRhsNotSet struct{}

// Error implements the error interface.
//
// Message: "0th RHS not set".
func (e *Err0thRhsNotSet) Error() string {
	return "0th RHS not set"
}

// NewErr0thRhsNotSet creates a new error of type *Err0thRhsNotSet.
//
// Returns:
//   - *Err0thRhsNotSet: A pointer to the new error.
func NewErr0thRhsNotSet() *Err0thRhsNotSet {
	e := &Err0thRhsNotSet{}
	return e
}

/////////////////////////////////////////////////////////////

// ErrHelper is an error that is returned when something goes wrong
// with a helper.
type ErrHelper[T gr.TokenTyper] struct {
	// Elem is the helper that caused the error.
	Elem *Helper[T]

	// Reason is the reason for the error.
	Reason error
}

// Error implements the error interface.
//
// Messages:
//   - "something went wrong with helper (no helper)" if Elem is nil.
//   - "helper (Elem) error: Reason" if Reason is not nil.
func (e *ErrHelper[T]) Error() string {
	var elem string

	if e.Elem == nil {
		elem = "no helper"
	} else {
		elem = e.Elem.String()
	}

	if e.Reason == nil {
		return fmt.Sprintf("something went wrong with helper (%s)", elem)
	} else {
		return fmt.Sprintf("helper (%s) error: %s", elem, e.Reason.Error())
	}
}

// NewErrHelper creates a new error of type *ErrHelper.
//
// Parameters:
//   - elem: The helper that caused the error.
//   - reason: The reason for the error.
//
// Returns:
//   - *ErrHelper: A pointer to the new error.
func NewErrHelper[T gr.TokenTyper](elem *Helper[T], reason error) *ErrHelper[T] {
	e := &ErrHelper[T]{
		Elem:   elem,
		Reason: reason,
	}

	return e
}
