package Parser

import (
	"fmt"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// ActionType represents the type of action that the parser will take.
type ActionType int8

const (
	// ActShift represents the action of shifting a token onto the stack.
	ActShift ActionType = iota

	// ActReduce represents the action of reducing the stack by a rule.
	ActReduce

	// ActError represents the action of encountering an error.
	ActError

	// ActAccept represents the action of accepting the input.
	// (It must be a reduce action and not a shift action)
	ActAccept
)

// Action represents an action that the parser will take.
type Action struct {
	// Type represents the type of action that the parser will take.
	Type ActionType

	// Data represents the data associated with the action.
	Data any
}

// NewShiftAction creates a new shift action.
//
// Returns:
//   - Action: The new shift action.
func NewShiftAction() Action {
	return Action{Type: ActShift}
}

// NewReduceAction creates a new reduce action.
//
// If the rule index is less than 0, an error action will be returned instead.
//
// Parameters:
//
//   - ruleIndex: The index of the rule to reduce by.
//
// Returns:
//
//   - Action: The new reduce action.
func NewReduceAction(ruleIndex int) Action {
	if ruleIndex < 0 {
		reason := ers.NewErrInvalidParameter(
			"ruleIndex",
			fmt.Errorf("value (%d) must be greater than or equal to 0", ruleIndex),
		)

		return Action{
			Type: ActError,
			Data: reason,
		}
	}

	return Action{Type: ActReduce, Data: ruleIndex}
}

// NewErrorAction creates a new error action.
//
// Parameters:
//
//   - reason: The reason for the error.
//
// Returns:
//
//   - Action: The new error action.
func NewErrorAction(reason error) Action {
	return Action{Type: ActError, Data: reason}
}

// NewAcceptAction creates a new accept action.
//
// Parameters:
//
//   - ruleIndex: The index of the rule to reduce by.
//
// Returns:
//
//   - Action: The new accept action.
func NewAcceptAction(ruleIndex int) Action {
	if ruleIndex < 0 {
		reason := ers.NewErrInvalidParameter(
			"ruleIndex",
			fmt.Errorf("value (%d) must be greater than or equal to 0", ruleIndex),
		)

		return Action{
			Type: ActError,
			Data: reason,
		}
	}

	return Action{Type: ActAccept, Data: ruleIndex}
}
