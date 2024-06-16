package Lexer

import (
	"strings"

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

// String implements the fmt.Stringer interface.
func (s EvalStatus) String() string {
	return [...]string{
		"complete",
		"incomplete",
		"error",
	}[s]
}

// CurrentEval is a struct that holds the current evaluation of the TreeExplorer.
type CurrentEval struct {
	// status is the status of the current evaluation.
	status EvalStatus

	// elem is the element of the current evaluation.
	elem *gr.LeafToken
}

// String implements the fmt.Stringer interface.
func (ce *CurrentEval) String() string {
	var builder strings.Builder

	builder.WriteString(ce.elem.GoString())
	builder.WriteString(" [")
	builder.WriteString(ce.status.String())
	builder.WriteRune(']')

	return builder.String()
}

// newCurrentEval creates a new CurrentEval with the given element.
//
// Parameters:
//   - elem: The element of the CurrentEval.
//
// Returns:
//   - *CurrentEval: The new CurrentEval.
func newCurrentEval(elem *gr.LeafToken) *CurrentEval {
	return &CurrentEval{
		status: EvalIncomplete,
		elem:   elem,
	}
}

// changeStatus sets the status of the CurrentEval.
//
// Parameters:
//   - status: The status to set.
func (ce *CurrentEval) changeStatus(status EvalStatus) {
	ce.status = status
}

// getStatus returns the status of the CurrentEval.
//
// Returns:
//   - EvalStatus: The status of the CurrentEval.
func (ce *CurrentEval) getStatus() EvalStatus {
	return ce.status
}

// getElem returns the element of the CurrentEval.
//
// Returns:
//   - T: The element of the CurrentEval.
func (ce *CurrentEval) getElem() *gr.LeafToken {
	return ce.elem
}
