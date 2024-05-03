package Grammar

import "fmt"

// ErrNoProductionRulesFound is an error that is returned when no production rules
// are found in a grammar.
type ErrNoProductionRulesFound struct{}

// Error returns the error message: "no production rules found".
//
// Returns:
//  	- string: The error message.
func (e *ErrNoProductionRulesFound) Error() string {
	return "no production rules found"
}

// NewErrNoProductionRulesFound creates a new error of type *ErrNoProductionRulesFound.
//
// Returns:
//  	- *ErrNoProductionRulesFound: The new error.
func NewErrNoProductionRulesFound() *ErrNoProductionRulesFound {
	return &ErrNoProductionRulesFound{}
}

// ErrLhsRhsMismatch is an error that is returned when the lhs of a production rule
// does not match the rhs.
type ErrLhsRhsMismatch struct {
	// Lhs is the left-hand side of the production rule.
	Lhs string

	// Rhs is the right-hand side of the production rule.
	Rhs string
}

// Error returns the error message: "lhs of production rule (lhs) does not match rhs (rhs)".
//
// Returns:
//  	- string: The error message.
func (e *ErrLhsRhsMismatch) Error() string {
	return fmt.Sprintf("lhs of production rule (%s) does not match rhs (%s)", e.Lhs, e.Rhs)
}

// NewErrLhsRhsMismatch creates a new error of type *ErrLhsRhsMismatch.
//
// Parameters:
//  	- lhs: The left-hand side of the production rule.
//  	- rhs: The right-hand side of the production rule.
//
// Returns:
//  	- *ErrLhsRhsMismatch: The new error.
func NewErrLhsRhsMismatch(lhs, rhs string) *ErrLhsRhsMismatch {
	return &ErrLhsRhsMismatch{
		Lhs: lhs,
		Rhs: rhs,
	}
}
