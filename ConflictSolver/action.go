package ConflictSolver

import (
	"fmt"
	"strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	ds "github.com/PlayerR9/MyGoLib/ListLike/DoubleLL"
	ui "github.com/PlayerR9/MyGoLib/Units/Iterators"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	ers "github.com/PlayerR9/MyGoLib/Units/errors"
)

// Actioner represents an action that the parser will take.
type Actioner interface {
	// AppendRhs appends a right-hand side token to the action.
	//
	// Parameters:
	//   - rhs: The right-hand side token to append.
	AppendRhs(rhs string)

	// Size returns the size of the action.
	//
	// Returns:
	//   - int: The size of the action.
	Size() int

	fmt.Stringer

	uc.Copier

	ui.Iterable[string]
}

// ActShift represents a shift action.
type ActShift struct {
	// Lookahead is the lookahead token ID for the shift action.
	Lookahead *string

	// Rhs is the right-hand side tokens of the rule.
	Rhs []string
}

// String implements the Actioner interface.
func (a *ActShift) String() string {
	var builder strings.Builder

	builder.WriteString("shift")
	builder.WriteRune('{')

	if a.Lookahead != nil {
		builder.WriteString(*a.Lookahead)
		builder.WriteString(" -> ")
	}

	builder.WriteString(strings.Join(a.Rhs, " "))
	builder.WriteRune('}')

	return builder.String()
}

// Copy implements the Actioner interface.
func (a *ActShift) Copy() uc.Copier {
	rhsCopy := make([]string, len(a.Rhs))
	copy(rhsCopy, a.Rhs)

	return &ActShift{
		Lookahead: a.Lookahead,
		Rhs:       rhsCopy,
	}
}

// AppendRhs implements the Actioner interface.
func (a *ActShift) AppendRhs(rhs string) {
	a.Rhs = append(a.Rhs, rhs)
}

// Iterator implements the Actioner interface.
func (a *ActShift) Iterator() ui.Iterater[string] {
	return ui.NewSimpleIterator(a.Rhs)
}

// Size implements the Actioner interface.
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
	// Rule is the rule to reduce by.
	Rule *gr.Production

	// Rhs is the right-hand side tokens of the rule.
	Rhs []string
}

// String implements the Actioner interface.
func (a *ActReduce) String() string {
	var builder strings.Builder

	builder.WriteString("reduce{")
	builder.WriteString(strings.Join(a.Rhs, " "))
	builder.WriteRune('}')

	return builder.String()
}

// Copy implements the Actioner interface.
func (a *ActReduce) Copy() uc.Copier {
	rhsCopy := make([]string, len(a.Rhs))
	copy(rhsCopy, a.Rhs)

	return &ActReduce{
		Rule: a.Rule.Copy().(*gr.Production),
		Rhs:  rhsCopy,
	}
}

// AppendRhs implements the Actioner interface.
func (a *ActReduce) AppendRhs(rhs string) {
	a.Rhs = append(a.Rhs, rhs)
}

// Iterator implements the Actioner interface.
func (a *ActReduce) Iterator() ui.Iterater[string] {
	return ui.NewSimpleIterator(a.Rhs)
}

// Size implements the Actioner interface.
func (a *ActReduce) Size() int {
	return len(a.Rhs)
}

// NewActReduce creates a new reduce action.
//
// Parameters:
//   - rule: The rule to reduce by.
//
// Returns:
//   - *ActReduce: A pointer to the new reduce action.
//
// Behaviors:
//   - If the rule is nil, nil is returned.
func NewActReduce(rule *gr.Production) *ActReduce {
	if rule == nil {
		return nil
	}

	return &ActReduce{
		Rule: rule,
	}
}

// GetRule returns the rule to reduce by.
//
// Returns:
//   - *gr.Production: The rule to reduce by.
func (a *ActReduce) GetRule() *gr.Production {
	return a.Rule
}

// ActAccept represents an accept action.
type ActAccept struct {
	// Rule is the rule to reduce by.
	Rule *gr.Production

	// Rhs is the right-hand side tokens of the rule.
	Rhs []string
}

// String implements the Actioner interface.
func (a *ActAccept) String() string {
	var builder strings.Builder

	builder.WriteString("accept{")
	builder.WriteString(strings.Join(a.Rhs, " "))
	builder.WriteRune('}')

	return builder.String()
}

// Copy implements the Actioner interface.
func (a *ActAccept) Copy() uc.Copier {
	rhsCopy := make([]string, len(a.Rhs))
	copy(rhsCopy, a.Rhs)

	return &ActAccept{
		Rule: a.Rule.Copy().(*gr.Production),
		Rhs:  rhsCopy,
	}
}

// AppendRhs implements the Actioner interface.
func (a *ActAccept) AppendRhs(rhs string) {
	a.Rhs = append(a.Rhs, rhs)
}

// Iterator implements the Actioner interface.
func (a *ActAccept) Iterator() ui.Iterater[string] {
	return ui.NewSimpleIterator(a.Rhs)
}

// Size implements the Actioner interface.
func (a *ActAccept) Size() int {
	return len(a.Rhs)
}

// NewAcceptAction creates a new accept action.
//
// Parameters:
//   - rule: The rule to reduce by.
//
// Returns:
//   - *ActAccept: A pointer to the new accept action.
//
// Behaviors:
//   - If the rule is nil, nil is returned.
func NewAcceptAction(rule *gr.Production) *ActAccept {
	if rule == nil {
		return nil
	}

	return &ActAccept{
		Rule: rule,
		Rhs:  make([]string, 0),
	}
}

// GetRule returns the rule to reduce by.
//
// Returns:
//   - *gr.Production: The rule to reduce by.
func (a *ActAccept) GetRule() *gr.Production {
	return a.Rule
}

// Match matches the action with the top of the stack.
//
// Parameters:
//   - a: The action to match.
//   - top: The top of the stack.
//   - stack: The stack.
//
// Returns:
//   - error: An error if the action does not match the top of the stack.
func MatchAction(a Actioner, top gr.Tokener, stack *ds.DoubleStack[gr.Tokener]) error {
	aShift, ok := a.(*ActShift)
	if ok && aShift != nil {
		lookahead := top.GetLookahead()

		if lookahead == nil {
			return ers.NewErrUnexpected(nil, *aShift.Lookahead)
		} else if lookahead.ID != *aShift.Lookahead {
			return ers.NewErrUnexpected(lookahead, *aShift.Lookahead)
		}
	}

	iter := a.Iterator()

	for {
		rhs, err := iter.Consume()
		if err != nil {
			break
		}

		top, err := stack.Pop()
		if err != nil {
			return ers.NewErrUnexpected(nil, rhs)
		}

		id := top.GetID()
		if id != rhs {
			return ers.NewErrUnexpected(top, rhs)
		}
	}

	return nil
}
