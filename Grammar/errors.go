package Grammar

import "fmt"

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
	return &ErrMissingArrow{}
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
	return &ErrNoLHSFound{}
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
	return &ErrNoRHSFound{}
}

// ErrNoProductionRulesFound is an error that is returned when no production rules
// are found in a grammar.
type ErrNoProductionRulesFound struct{}

// Error returns the error message: "no production rules found".
//
// Returns:
//   - string: The error message.
func (e *ErrNoProductionRulesFound) Error() string {
	return "no production rules found"
}

// NewErrNoProductionRulesFound creates a new error of type *ErrNoProductionRulesFound.
//
// Returns:
//   - *ErrNoProductionRulesFound: The new error.
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
//   - string: The error message.
func (e *ErrLhsRhsMismatch) Error() string {
	return fmt.Sprintf("lhs of production rule (%s) does not match rhs (%s)", e.Lhs, e.Rhs)
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
	return &ErrLhsRhsMismatch{
		Lhs: lhs,
		Rhs: rhs,
	}
}

// ErrUnknownToken is an error that is returned when an unknown token is found.
type ErrUnknowToken struct {
	// Token is the unknown token.
	Token Tokener
}

// Error returns the error message: "unknown token type: (type)".
//
// Returns:
//   - string: The error message.
//
// Behaviors:
//   - If the token is nil, the error message is "token is nil".
func (e *ErrUnknowToken) Error() string {
	if e.Token == nil {
		return "token is nil"
	} else {
		return fmt.Sprintf("unknown token type: %T", e.Token)
	}
}

// NewErrUnknowToken creates a new error of type *ErrUnknowToken.
//
// Parameters:
//   - token: The unknown token.
//
// Returns:
//   - *ErrUnknowToken: The new error.
func NewErrUnknowToken(token Tokener) *ErrUnknowToken {
	return &ErrUnknowToken{
		Token: token,
	}
}
