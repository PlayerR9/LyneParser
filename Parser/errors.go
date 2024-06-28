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

// ErrUnknownAction is an error that is returned when the parser encounters an unknown
// action.
type ErrUnknownAction struct {
	// TODO: Remove this once the MyGoLib/Units/Errors package is updated.

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
