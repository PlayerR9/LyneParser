package TreeExplorer

import (
	"fmt"

	uc "github.com/PlayerR9/MyGoLib/Units/Common"
)

// EvalStatus represents the status of an evaluation.
type EvalStatus int8

const (
	// EvalComplete represents a completed evaluation.
	EvalComplete EvalStatus = iota

	// EvalIncomplete represents an incomplete evaluation.
	EvalIncomplete

	// EvalError represents an evaluation that has an error.
	EvalError
)

// String is a method of fmt.Stringer that returns the string representation of the EvalStatus.
//
// Returns:
//   - string: The string representation of the EvalStatus.
func (s EvalStatus) String() string {
	return [...]string{
		"complete",
		"incomplete",
		"error",
	}[s]
}

// CurrentEval is a struct that holds the current evaluation of the TreeExplorer.
type CurrentEval[T any] struct {
	// Status is the status of the current evaluation.
	Status EvalStatus

	// Elem is the element of the current evaluation.
	Elem T
}

// String is a method of fmt.Stringer that returns the string representation of the CurrentEval.
//
// Returns:
//   - string: The string representation of the CurrentEval.
func (ce *CurrentEval[T]) String() string {
	return fmt.Sprintf("%s [%s]", uc.StringOf(ce.Elem), ce.Status.String())
}

// NewCurrentEval creates a new CurrentEval with the given element.
//
// Parameters:
//   - elem: The element of the CurrentEval.
//
// Returns:
//   - *CurrentEval: The new CurrentEval.
func NewCurrentEval[T any](elem T) *CurrentEval[T] {
	return &CurrentEval[T]{
		Status: EvalIncomplete,
		Elem:   elem,
	}
}

// SetStatus sets the status of the CurrentEval.
//
// Parameters:
//   - status: The status to set.
func (ce *CurrentEval[T]) SetStatus(status EvalStatus) {
	ce.Status = status
}

// GetStatus returns the status of the CurrentEval.
//
// Returns:
//   - EvalStatus: The status of the CurrentEval.
func (ce *CurrentEval[T]) GetStatus() EvalStatus {
	return ce.Status
}

// GetElem returns the element of the CurrentEval.
//
// Returns:
//   - T: The element of the CurrentEval.
func (ce *CurrentEval[T]) GetElem() T {
	return ce.Elem
}
