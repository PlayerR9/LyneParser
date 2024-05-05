package ConflictSolver

import (
	"fmt"
	"strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	ds "github.com/PlayerR9/MyGoLib/ListLike/DoubleLL"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"

	intf "github.com/PlayerR9/MyGoLib/Units/Interfaces"
)

// Actioner represents an action that the parser will take.
type Actioner interface {
	// AppendRhs appends a right-hand side token to the action.
	//
	// Parameters:
	//   - rhs: The right-hand side token to append.
	AppendRhs(rhs string)

	// Match matches the action with the top of the stack.
	//
	// Parameters:
	//   - top: The top of the stack.
	//   - stack: The stack.
	//
	// Returns:
	//   - error: An error if the action does not match the top of the stack.
	Match(top gr.Tokener, stack *ds.DoubleStack[gr.Tokener]) error

	// Size returns the size of the action.
	//
	// Returns:
	//   - int: The size of the action.
	Size() int

	fmt.Stringer

	intf.Copier
}

// ActShift represents a shift action.
type ActShift struct {
	// Lookahead is the lookahead token ID for the shift action.
	Lookahead *string

	// Rhs is the right-hand side tokens of the rule.
	Rhs []string
}

// String returns a string representation of the shift action.
//
// Returns:
//   - string: The string representation of the shift action.
func (a *ActShift) String() string {
	var builder strings.Builder

	builder.WriteString("shift")
	builder.WriteRune('{')

	if a.Lookahead != nil {
		builder.WriteString(*a.Lookahead)
		builder.WriteRune(' ')
		builder.WriteRune('|')
		builder.WriteRune(' ')
	}

	builder.WriteString(strings.Join(a.Rhs, " "))
	builder.WriteRune('}')

	return builder.String()
}

// Copy creates a copy of the shift action.
//
// Returns:
//   - intf.Copier: The copy of the shift action.
func (a *ActShift) Copy() intf.Copier {
	rhsCopy := make([]string, len(a.Rhs))
	copy(rhsCopy, a.Rhs)

	return &ActShift{
		Lookahead: a.Lookahead,
		Rhs:       rhsCopy,
	}
}

// AppendRhs appends a right-hand side token to the shift action.
// It never returns an error.
//
// Parameters:
//   - rhs: The right-hand side token to append.
func (a *ActShift) AppendRhs(rhs string) {
	a.Rhs = append(a.Rhs, rhs)
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
	// Rule is the rule to reduce by.
	Rule *gr.Production

	// Rhs is the right-hand side tokens of the rule.
	Rhs []string
}

// String returns a string representation of the reduce action.
//
// Returns:
//   - string: The string representation of the reduce action.
func (a *ActReduce) String() string {
	return fmt.Sprintf(
		"reduce{%s}",
		strings.Join(a.Rhs, " "),
	)
}

// Copy creates a copy of the reduce action.
//
// Returns:
//   - intf.Copier: The copy of the reduce action.
func (a *ActReduce) Copy() intf.Copier {
	rhsCopy := make([]string, len(a.Rhs))
	copy(rhsCopy, a.Rhs)

	return &ActReduce{
		Rule: a.Rule.Copy().(*gr.Production),
		Rhs:  rhsCopy,
	}
}

// AppendRhs appends a right-hand side token to the reduce action.
// It never returns an error.
//
// Parameters:
//   - rhs: The right-hand side token to append.
func (a *ActReduce) AppendRhs(rhs string) {
	a.Rhs = append(a.Rhs, rhs)
}

// NewActReduce creates a new reduce action.
//
// Parameters:
//   - rule: The rule to reduce by.
//
// Returns:
//   - *ActReduce: A pointer to the new reduce action.
//   - error: An error of type *ers.ErrInvalidParameter if the rule is nil.
func NewActReduce(rule *gr.Production) (*ActReduce, error) {
	if rule == nil {
		return nil, ers.NewErrNilParameter("rule")
	}

	return &ActReduce{
		Rule: rule,
	}, nil
}

// ActAccept represents an accept action.
type ActAccept struct {
	// Rule is the rule to reduce by.
	Rule *gr.Production

	// Rhs is the right-hand side tokens of the rule.
	Rhs []string
}

// String returns a string representation of the accept action.
//
// Returns:
//   - string: The string representation of the accept action.
func (a *ActAccept) String() string {
	return fmt.Sprintf(
		"accept{%s}",
		strings.Join(a.Rhs, " "),
	)
}

// Copy creates a copy of the accept action.
//
// Returns:
//   - intf.Copier: The copy of the accept action.
func (a *ActAccept) Copy() intf.Copier {
	rhsCopy := make([]string, len(a.Rhs))
	copy(rhsCopy, a.Rhs)

	return &ActAccept{
		Rule: a.Rule.Copy().(*gr.Production),
		Rhs:  rhsCopy,
	}
}

// AppendRhs appends a right-hand side token to the accept action.
// It never returns an error.
//
// Parameters:
//   - rhs: The right-hand side token to append.
func (a *ActAccept) AppendRhs(rhs string) {
	a.Rhs = append(a.Rhs, rhs)
}

// NewAcceptAction creates a new accept action.
//
// Parameters:
//   - rule: The rule to reduce by.
//
// Returns:
//   - *ActAccept: A pointer to the new accept action.
//   - error: An error of type *ers.ErrInvalidParameter if the rule is nil.
func NewAcceptAction(rule *gr.Production) (*ActAccept, error) {
	if rule == nil {
		return nil, ers.NewErrNilParameter("rule")
	}

	return &ActAccept{
		Rule: rule,
		Rhs:  make([]string, 0),
	}, nil
}

/////////////////////////////////////////////////////////////

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
		top, err := stack.Pop()
		if err != nil {
			return ers.NewErrUnexpected(nil, rhs)
		}

		if top.GetID() != rhs {
			return ers.NewErrUnexpected(top, rhs)
		}
	}

	return nil
}

func (a *ActShift) Size() int {
	return len(a.Rhs)
}

func (a *ActReduce) Match(top gr.Tokener, stack *ds.DoubleStack[gr.Tokener]) error {
	for _, rhs := range a.Rhs {
		top, err := stack.Pop()
		if err != nil {
			return ers.NewErrUnexpected(nil, rhs)
		}

		if top.GetID() != rhs {
			return ers.NewErrUnexpected(top, rhs)
		}
	}

	return nil
}

func (a *ActReduce) Size() int {
	return len(a.Rhs)
}

// GetRule returns the rule to reduce by.
//
// Returns:
//   - *gr.Production: The rule to reduce by.
func (a *ActReduce) GetRule() *gr.Production {
	return a.Rule
}

func (a *ActAccept) Match(top gr.Tokener, stack *ds.DoubleStack[gr.Tokener]) error {
	for _, rhs := range a.Rhs {
		top, err := stack.Pop()
		if err != nil {
			return ers.NewErrUnexpected(nil, rhs)
		}

		if top.GetID() != rhs {
			return ers.NewErrUnexpected(top, rhs)
		}
	}

	return nil
}

func (a *ActAccept) Size() int {
	return len(a.Rhs)
}

// GetRule returns the rule to reduce by.
//
// Returns:
//   - *gr.Production: The rule to reduce by.
func (a *ActAccept) GetRule() *gr.Production {
	return a.Rule
}
