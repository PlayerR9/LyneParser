package Parser

import "fmt"

// ErrNoAccept is an error that is returned when the parser reaches the end of the
// input stream without accepting the input stream.
type ErrNoAccept struct{}

// Error is a method of the error interface.
//
// Returns:
//
//   - string: The error message.
func (e *ErrNoAccept) Error() string {
	return "reached end of input stream without accepting"
}

// NewErrNoAccept creates a new ErrNoAccept error.
//
// Returns:
//
//   - *ErrNoAccept: A pointer to the new ErrNoAccept error.
func NewErrNoAccept() *ErrNoAccept {
	return &ErrNoAccept{}
}

// ErrAfter is an error that is returned when the parser encounters an error after
// a certain point in the input stream.
type ErrAfter struct {
	// After is the position in the input stream where the error occurred.
	After string

	// Reason is the reason for the error.
	Reason error
}

// Error is a method of the error interface.
//
// Returns:
//  - string: The error message.
func (e *ErrAfter) Error() string {
	if e.Reason == nil {
		return fmt.Sprintf("something went wrong after %s", e.After)
	} else {
		return fmt.Sprintf("after %s: %s", e.After, e.Reason.Error())
	}
}

// NewErrAfter creates a new ErrAfter error.
//
// Parameters:
//   - after: The position in the input stream where the error occurred.
//   - reason: The reason for the error.
//
// Returns:
//   - *ErrAfter: A pointer to the new ErrAfter error.
func NewErrAfter(after string, reason error) *ErrAfter {
	return &ErrAfter{
		After:  after,
		Reason: reason,
	}
}

// ErrUnknownAction is an error that is returned when the parser encounters an unknown
// action.
type ErrUnknownAction struct {
	// Action is the action that was attempted.
	Action any
}

// Error is a method of the error interface.
//
// Returns:
//  - string: The error message.
func (e *ErrUnknownAction) Error() string {
	return fmt.Sprintf("unknown action: %T", e.Action)
}

// NewErrUnknownAction creates a new ErrUnknownAction error.
//
// Parameters:
//   - action: The action that was attempted.
//
// Returns:
//   - *ErrUnknownAction: A pointer to the new ErrUnknownAction error.
func NewErrUnknownAction(action any) *ErrUnknownAction {
	return &ErrUnknownAction{
		Action: action,
	}
}
