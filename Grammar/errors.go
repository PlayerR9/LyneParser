package Grammar

import (
	"strings"
)

// ErrMissingArrow is an error that is returned when an arrow is missing in a rule.
type ErrMissingArrow struct{}

// Error implements the error interface.
//
// Message: "missing arrow in rule".
func (e *ErrMissingArrow) Error() string {
	return "missing arrow in rule"
}

// NewErrMissingArrow creates a new error of type *ErrMissingArrow.
//
// Returns:
//   - *ErrMissingArrow: The new error.
func NewErrMissingArrow() *ErrMissingArrow {
	e := &ErrMissingArrow{}
	return e
}

// ErrNoLHSFound is an error that is returned when no left-hand side is found in a rule.
type ErrNoLHSFound struct{}

// Error implements the error interface.
//
// Message: "no left-hand side in rule".
func (e *ErrNoLHSFound) Error() string {
	return "no left-hand side in rule"
}

// NewErrNoLHSFound creates a new error of type *ErrNoLHSFound.
//
// Returns:
//   - *ErrNoLHSFound: The new error.
func NewErrNoLHSFound() *ErrNoLHSFound {
	e := &ErrNoLHSFound{}
	return e
}

// ErrNoRHSFound is an error that is returned when no right-hand side is found in a rule.
type ErrNoRHSFound struct{}

// Error implements the error interface.
//
// Message: "no right-hand side in rule".
func (e *ErrNoRHSFound) Error() string {
	return "no right-hand side in rule"
}

// NewErrNoRHSFound creates a new error of type *ErrNoRHSFound.
//
// Returns:
//   - *ErrNoRHSFound: The new error.
func NewErrNoRHSFound() *ErrNoRHSFound {
	e := &ErrNoRHSFound{}
	return e
}

// ErrNoProductionRulesFound is an error that is returned when no production rules
// are found in a grammar.
type ErrNoProductionRulesFound struct{}

// Error implements the error interface.
//
// Message: "no production rules found".
func (e *ErrNoProductionRulesFound) Error() string {
	return "no production rules found"
}

// NewErrNoProductionRulesFound creates a new error of type *ErrNoProductionRulesFound.
//
// Returns:
//   - *ErrNoProductionRulesFound: The new error.
func NewErrNoProductionRulesFound() *ErrNoProductionRulesFound {
	e := &ErrNoProductionRulesFound{}
	return e
}

// ErrLhsRhsMismatch is an error that is returned when the lhs of a production rule
// does not match the rhs.
type ErrLhsRhsMismatch struct {
	// Lhs is the left-hand side of the production rule.
	Lhs string

	// Rhs is the right-hand side of the production rule.
	Rhs string
}

// Error implements the error interface.
//
// Message: "lhs of production rule (lhs) does not match rhs (rhs)".
func (e *ErrLhsRhsMismatch) Error() string {
	var builder strings.Builder

	builder.WriteString("lhs of production rule (")
	builder.WriteString(e.Lhs)
	builder.WriteString(") does not match rhs (")
	builder.WriteString(e.Rhs)
	builder.WriteRune(')')

	return builder.String()
}

// NewErrLhsRhsMismatch creates a new error of type *ErrLhsRhsMismatch.
//
// Parameters:
//   - lhs: The left-hand side of the production rule.
//   - rhs: The right-hand side of the production rule.
//
// Returns:
//   - *ErrLhsRhsMismatch: The new error.
func NewErrLhsRhsMismatch(lhs, rhs string) *ErrLhsRhsMismatch {
	e := &ErrLhsRhsMismatch{
		Lhs: lhs,
		Rhs: rhs,
	}
	return e
}

// ErrCycleDetected is an error that is returned when a cycle is detected.
type ErrCycleDetected struct{}

// Error implements the error interface.
//
// Message: "cycle detected".
func (e *ErrCycleDetected) Error() string {
	return "cycle detected"
}

// NewErrCycleDetected creates a new error of type *ErrCycleDetected.
//
// Returns:
//   - *ErrCycleDetected: The new error.
func NewErrCycleDetected() *ErrCycleDetected {
	e := &ErrCycleDetected{}
	return e
}
