package ConflictSolver

import (
	"strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	ud "github.com/PlayerR9/MyGoLib/Units/Debugging"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// Actioner represents an action that the parser will take.
type Actioner interface {
	// GetLookahead returns the lookahead token ID for the action.
	//
	// Returns:
	//   - *string: The lookahead token ID.
	GetLookahead() *string

	// Iterator returns an iterator of the right-hand side tokens.
	//
	// Returns:
	//   - uc.Iterater[string]: An iterator of the right-hand side tokens.
	Iterator() uc.Iterater[string]
}

// Action represents an action in a decision table.
type Action struct {
	// lookahead is the lookahead token ID for the shift action.
	lookahead *string

	// Rhs is the right-hand side tokens of the rule.
	rhs []string
}

// String implements the fmt.Stringer interface.
func (a *Action) String() string {
	if a.lookahead == nil {
		return strings.Join(a.rhs, " ")
	}

	var builder strings.Builder

	builder.WriteString(*a.lookahead)
	builder.WriteString(" <- ")
	builder.WriteString(strings.Join(a.rhs, " "))

	return builder.String()
}

// Iterator implements the Iterators.Iterater interface.
func (a *Action) Iterator() uc.Iterater[string] {
	return uc.NewSimpleIterator(a.rhs)
}

// AppendRhs appends a right-hand side token to the action.
//
// Parameters:
//   - rhs: The right-hand side token to append.
func (a *Action) AppendRhs(rhs string) {
	a.rhs = append(a.rhs, rhs)
}

// Size returns the size of the right-hand side tokens.
//
// Returns:
//   - int: The size of the right-hand side tokens.
func (a *Action) Size() int {
	return len(a.rhs)
}

// SetLookahead sets the lookahead token ID for the action.
//
// Parameters:
//   - lookahead: The lookahead token ID to set.
func (a *Action) SetLookahead(lookahead *string) {
	a.lookahead = lookahead
}

// GetLookahead returns the lookahead token ID for the action.
//
// Returns:
//   - *string: The lookahead token ID.
func (a *Action) GetLookahead() *string {
	return a.lookahead
}

// ActShift represents a shift action.
type ActShift struct {
	*Action
}

// String implements the fmt.Stringer interface.
func (a *ActShift) String() string {
	var builder strings.Builder

	builder.WriteString("shift")
	builder.WriteRune('{')
	builder.WriteString(a.Action.String())
	builder.WriteRune('}')

	return builder.String()
}

// Copy implements the common.Copier interface.
func (a *ActShift) Copy() uc.Copier {
	rhsCopy := make([]string, len(a.Action.rhs))
	copy(rhsCopy, a.Action.rhs)

	return &ActShift{
		Action: &Action{
			lookahead: a.Action.lookahead,
			rhs:       rhsCopy,
		},
	}
}

// NewActShift creates a new shift action.
//
// Returns:
//   - *ActShift: A pointer to the new shift action.
func NewActShift() *ActShift {
	return &ActShift{
		Action: &Action{
			lookahead: nil,
			rhs:       make([]string, 0),
		},
	}
}

// ActReduce represents a reduce action.
type ActReduce struct {
	*Action

	// rule is the rule to reduce by.
	rule *gr.Production

	// original is the original rule to reduce by.
	// this should never be modified.
	original *gr.Production

	// shouldAccept is true if the reduce action should accept.
	shouldAccept bool
}

// String implements the fmt.Stringer interface.
func (a *ActReduce) String() string {
	var builder strings.Builder

	if a.shouldAccept {
		builder.WriteString("accept{")
	} else {
		builder.WriteString("reduce{")
	}

	builder.WriteString(a.Action.String())
	builder.WriteRune('}')

	return builder.String()
}

// Copy implements the common.Copier interface.
func (a *ActReduce) Copy() uc.Copier {
	rhsCopy := make([]string, len(a.Action.rhs))
	copy(rhsCopy, a.Action.rhs)

	return &ActReduce{
		Action: &Action{
			lookahead: a.Action.lookahead,
			rhs:       rhsCopy,
		},
		rule:         a.rule.Copy().(*gr.Production),
		original:     a.original,
		shouldAccept: a.shouldAccept,
	}
}

// NewActReduce creates a new reduce action.
//
// Parameters:
//   - rule: The rule to reduce by.
//   - shouldAccept: True if the reduce action should accept.
//
// Returns:
//   - *ActReduce: A pointer to the new reduce action.
//
// Behaviors:
//   - If the rule is nil, nil is returned.
func NewActReduce(rule *gr.Production, shouldAccept bool) *ActReduce {
	if rule == nil {
		return nil
	}

	return &ActReduce{
		Action: &Action{
			lookahead: nil,
			rhs:       make([]string, 0),
		},
		rule:         rule.Copy().(*gr.Production),
		original:     rule,
		shouldAccept: shouldAccept,
	}
}

// GetRule returns the rule to reduce by.
//
// Returns:
//   - *gr.Production: The rule to reduce by.
func (a *ActReduce) GetRule() *gr.Production {
	return a.rule
}

// GetOriginal returns the original rule to reduce by.
//
// Returns:
//   - *gr.Production: The original rule to reduce by.
func (a *ActReduce) GetOriginal() *gr.Production {
	return a.original
}

// ShouldAccept returns true if the reduce action should accept.
//
// Returns:
//   - bool: True if the reduce action should accept.
func (a *ActReduce) ShouldAccept() bool {
	return a.shouldAccept
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
func MatchAction(a Actioner, top gr.Tokener, stack *ud.History[lls.Stacker[gr.Tokener]]) error {
	ela := a.GetLookahead()
	tla := top.GetLookahead()

	if ela != nil {
		if tla == nil {
			return uc.NewErrUnexpected("", *ela)
		} else if *ela != tla.GetID() {
			return uc.NewErrUnexpected(top.GoString(), *ela)
		}
	}

	iter := a.Iterator()

	for {
		rhs, err := iter.Consume()
		if err != nil {
			break
		}

		cmd := lls.NewPop[gr.Tokener]()
		err = stack.ExecuteCommand(cmd)
		if err != nil {
			return uc.NewErrUnexpected("", rhs)
		}
		top := cmd.Value()

		id := top.GetID()
		if id != rhs {
			return uc.NewErrUnexpected(top.GoString(), rhs)
		}
	}

	return nil
}
