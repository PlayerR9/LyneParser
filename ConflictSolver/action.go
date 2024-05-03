package ConflictSolver

import (
	"errors"
	"fmt"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	ds "github.com/PlayerR9/MyGoLib/ListLike/DoubleLL"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// Actioner represents an action that the parser will take.
type Actioner interface {
	AppendRhs(rhs string) error
	Match(top gr.Tokener, stack *ds.DoubleStack[gr.Tokener]) error
	Size() int
}

// ActShift represents a shift action.
type ActShift struct {
	// Lookahead is the lookahead token ID for the shift action.
	Lookahead *string

	// Rhs is the right-hand side tokens of the rule.
	Rhs []string
}

// AppendRhs appends a right-hand side token to the shift action.
// It never returns an error.
//
// Parameters:
//   - rhs: The right-hand side token to append.
//
// Returns:
//   - error: An error if the right-hand side token could not be appended.
func (a *ActShift) AppendRhs(rhs string) error {
	a.Rhs = append(a.Rhs, rhs)

	return nil
}

func (a *ActShift) Match(top gr.Tokener, stack *ds.DoubleStack[gr.Tokener]) error {
	if a.Lookahead != nil {
		lookahead := top.GetLookahead()

		if lookahead == nil {
			return ers.NewErrUnexpected(nil, *a.Lookahead)
		} else if lookahead.ID != *a.Lookahead {
			return ers.NewErrUnexpected(top.GetLookahead(), *a.Lookahead)
		}
	}

	for _, rhs := range a.Rhs {
		if stack.IsEmpty() {
			return ers.NewErrUnexpected(nil, rhs)
		}

		top := stack.Pop()

		if top.GetID() != rhs {
			return ers.NewErrUnexpected(top, rhs)
		}
	}

	return nil
}

func (a *ActShift) Size() int {
	return len(a.Rhs)
}

// NewActShift creates a new shift action.
//
// Returns:
//   - *ActShift: A pointer to the new shift action.
func NewActShift() *ActShift {
	return &ActShift{
		Lookahead: nil,
		Rhs:       make([]string, 0),
	}
}

// SetLookahead sets the lookahead token ID for the shift action.
//
// Parameters:
//   - lookahead: The lookahead token ID.
func (a *ActShift) SetLookahead(lookahead *string) {
	a.Lookahead = lookahead
}

// ActReduce represents a reduce action.
type ActReduce struct {
	// RuleIndex is the index of the rule to reduce by.
	RuleIndex int

	// Rhs is the right-hand side tokens of the rule.
	Rhs []string
}

// AppendRhs appends a right-hand side token to the reduce action.
// It never returns an error.
//
// Parameters:
//   - rhs: The right-hand side token to append.
//
// Returns:
//   - error: An error if the right-hand side token could not be appended.
func (a *ActReduce) AppendRhs(rhs string) error {
	a.Rhs = append(a.Rhs, rhs)

	return nil
}

func (a *ActReduce) Match(top gr.Tokener, stack *ds.DoubleStack[gr.Tokener]) error {
	for _, rhs := range a.Rhs {
		if stack.IsEmpty() {
			return ers.NewErrUnexpected(nil, rhs)
		}

		top := stack.Pop()

		if top.GetID() != rhs {
			return ers.NewErrUnexpected(top, rhs)
		}
	}

	return nil
}

func (a *ActReduce) Size() int {
	return len(a.Rhs)
}

// NewActReduce creates a new reduce action.
//
// If the rule index is less than 0, an error action will be returned instead.
//
// Parameters:
//   - ruleIndex: The index of the rule to reduce by.
//
// Returns:
//   - Actioner: The new reduce or error action.
func NewActReduce(ruleIndex int) Actioner {
	if ruleIndex < 0 {
		reason := ers.NewErrInvalidParameter(
			"ruleIndex",
			fmt.Errorf("value (%d) must be greater than or equal to 0", ruleIndex),
		)

		return &ActError{
			Reason: reason,
		}
	}

	return &ActReduce{
		RuleIndex: ruleIndex,
	}
}

// ActError represents an error action.
type ActError struct {
	// Reason is the reason for the error.
	Reason error
}

// AppendRhs appends a right-hand side token to the error action.
// It always returns an error.
//
// Parameters:
//   - rhs: The right-hand side token to append.
//
// Returns:
//   - error: An error if the right-hand side token could not be appended.
func (a *ActError) AppendRhs(rhs string) error {
	return errors.New("cannot append right-hand side token to error action")
}

func (a *ActError) Match(top gr.Tokener, stack *ds.DoubleStack[gr.Tokener]) error {
	return a.Reason
}

func (a *ActError) Size() int {
	return 0
}

// NewErrorAction creates a new error action.
//
// Parameters:
//   - reason: The reason for the error.
//
// Returns:
//   - *ActError: A pointer to the new error action.
func NewErrorAction(reason error) *ActError {
	return &ActError{
		Reason: reason,
	}
}

// ActAccept represents an accept action.
type ActAccept struct {
	// RuleIndex is the index of the rule to reduce by.
	RuleIndex int

	// Rhs is the right-hand side tokens of the rule.
	Rhs []string
}

// AppendRhs appends a right-hand side token to the accept action.
// It never returns an error.
//
// Parameters:
//   - rhs: The right-hand side token to append.
//
// Returns:
//   - error: An error if the right-hand side token could not be appended.
func (a *ActAccept) AppendRhs(rhs string) error {
	a.Rhs = append(a.Rhs, rhs)

	return nil
}

func (a *ActAccept) Match(top gr.Tokener, stack *ds.DoubleStack[gr.Tokener]) error {
	for _, rhs := range a.Rhs {
		if stack.IsEmpty() {
			return ers.NewErrUnexpected(nil, rhs)
		}

		top := stack.Pop()

		if top.GetID() != rhs {
			return ers.NewErrUnexpected(top, rhs)
		}
	}

	return nil
}

func (a *ActAccept) Size() int {
	return len(a.Rhs)
}

// NewAcceptAction creates a new accept action.
//
// Parameters:
//   - ruleIndex: The index of the rule to reduce by.
//
// Returns:
//   - Actioner: The new accept or error action.
func NewAcceptAction(ruleIndex int) Actioner {
	if ruleIndex < 0 {
		reason := ers.NewErrInvalidParameter(
			"ruleIndex",
			fmt.Errorf("value (%d) must be greater than or equal to 0", ruleIndex),
		)

		return &ActError{
			Reason: reason,
		}
	}

	return &ActAccept{
		RuleIndex: ruleIndex,
		Rhs:       make([]string, 0),
	}
}
