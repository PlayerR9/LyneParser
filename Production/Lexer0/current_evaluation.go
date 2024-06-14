package Lexer0

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
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
type CurrentEval struct {
	// Status is the status of the current evaluation.
	Status EvalStatus

	// Elem is the element of the current evaluation.
	Elem *gr.LeafToken
}

// NewCurrentEval creates a new CurrentEval with the given element.
//
// Parameters:
//   - elem: The element of the CurrentEval.
//
// Returns:
//   - *CurrentEval: The new CurrentEval.
func NewCurrentEval(elem *gr.LeafToken) *CurrentEval {
	return &CurrentEval{
		Status: EvalIncomplete,
		Elem:   elem,
	}
}

// SetStatus sets the status of the CurrentEval.
//
// Parameters:
//   - status: The status to set.
func (ce *CurrentEval) SetStatus(status EvalStatus) {
	ce.Status = status
}

// GetStatus returns the status of the CurrentEval.
//
// Returns:
//   - EvalStatus: The status of the CurrentEval.
func (ce *CurrentEval) GetStatus() EvalStatus {
	return ce.Status
}

// GetElem returns the element of the CurrentEval.
//
// Returns:
//   - *gr.Token: The element of the CurrentEval.
func (ce *CurrentEval) GetElem() *gr.LeafToken {
	return ce.Elem
}
