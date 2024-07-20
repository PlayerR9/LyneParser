package ConflictSolver

import (
	"fmt"
	"strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	ud "github.com/PlayerR9/MyGoLib/Units/Debugging"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	lls "github.com/PlayerR9/stack"
)

// Actioner represents an action that the parser will take.
type Actioner[T gr.TokenTyper] interface {
	// GetLookahead returns the lookahead token ID for the action.
	//
	// Returns:
	//   - T: The lookahead token ID.
	//   - bool: True if there is a lookahead token ID.
	GetLookahead() (T, bool)

	// ActionSize returns the size of the right-hand side tokens.
	//
	// Returns:
	//   - int: The size of the right-hand side tokens.
	ActionSize() int

	fmt.Stringer

	uc.Iterable[T]
}

// Action represents an action in a decision table.
type Action[T gr.TokenTyper] struct {
	// lookahead is the lookahead token ID for the shift action.
	lookahead *T

	// Rhs is the right-hand side tokens of the rule.
	rhs []T
}

// String implements the Actioner interface.
func (a *Action[T]) String() string {
	if a.lookahead == nil {
		values := make([]string, 0, len(a.rhs))
		for _, symbol := range a.rhs {
			values = append(values, symbol.String())
		}

		str := strings.Join(values, " ")
		return str
	}

	var builder strings.Builder

	builder.WriteString((*a.lookahead).String())
	builder.WriteString(" <- ")

	values := make([]string, 0, len(a.rhs))
	for _, symbol := range a.rhs {
		values = append(values, symbol.String())
	}

	str := strings.Join(values, " ")

	builder.WriteString(str)

	return builder.String()
}

// Copy implements the common.Copier interface.
func (a *Action[T]) Copy() uc.Copier {
	rhs_copy := make([]T, len(a.rhs))
	copy(rhs_copy, a.rhs)

	a_copy := &Action[T]{
		lookahead: a.lookahead,
		rhs:       rhs_copy,
	}

	return a_copy
}

// Iterator implements the Iterators.Iterater interface.
func (a *Action[T]) Iterator() uc.Iterater[T] {
	iter := uc.NewSimpleIterator(a.rhs)
	return iter
}

// newAction creates a new action.
//
// Parameters:
//   - lookahead: The lookahead token ID for the action.
//   - rhs: The right-hand side tokens of the rule.
//
// Returns:
//   - *Action: A pointer to the new action.
func newAction[T gr.TokenTyper](lookahead *T, rhs []T) *Action[T] {
	act := &Action[T]{
		lookahead: lookahead,
		rhs:       rhs,
	}
	return act
}

// AppendRhs appends a right-hand side token to the action.
//
// Parameters:
//   - rhs: The right-hand side token to append.
func (a *Action[T]) AppendRhs(rhs T) {
	a.rhs = append(a.rhs, rhs)
}

// ActionSize returns the size of the right-hand side tokens.
//
// Returns:
//   - int: The size of the right-hand side tokens.
func (a *Action[T]) ActionSize() int {
	return len(a.rhs)
}

// SetLookahead sets the lookahead token ID for the action.
//
// Parameters:
//   - lookahead: The lookahead token ID to set.
func (a *Action[T]) SetLookahead(lookahead *T) {
	a.lookahead = lookahead
}

// GetLookahead implements the Actioner interface.
func (a *Action[T]) GetLookahead() (T, bool) {
	if a.lookahead == nil {
		return *new(T), false
	}

	return *a.lookahead, true
}

// ActShift represents a shift action.
type ActShift[T gr.TokenTyper] struct {
	*Action[T]
}

// String implements the Actioner interface.
func (a *ActShift[T]) String() string {
	actStr := a.Action.String()

	var builder strings.Builder

	builder.WriteString("shift{")
	builder.WriteString(actStr)
	builder.WriteRune('}')

	str := builder.String()

	return str
}

// Copy implements the common.Copier interface.
func (a *ActShift[T]) Copy() uc.Copier {
	act_copy := a.Action.Copy().(*Action[T])

	a_copy := &ActShift[T]{
		Action: act_copy,
	}

	return a_copy
}

// NewActShift creates a new shift action.
//
// Returns:
//   - *ActShift: A pointer to the new shift action.
func NewActShift[T gr.TokenTyper]() *ActShift[T] {
	act := newAction(nil, make([]T, 0))

	as := &ActShift[T]{
		Action: act,
	}
	return as
}

// ActReduce represents a reduce action.
type ActReduce[T gr.TokenTyper] struct {
	*Action[T]

	// Rule is the Rule to reduce by.
	Rule *gr.Production[T]

	// Original is the Original rule to reduce by.
	// this should never be modified.
	Original *gr.Production[T]
}

// String implements the Actioner interface.
func (a *ActReduce[T]) String() string {
	act_str := a.Action.String()

	var builder strings.Builder

	builder.WriteString("reduce{")
	builder.WriteString(act_str)
	builder.WriteRune('}')

	str := builder.String()

	return str
}

// Copy implements the common.Copier interface.
func (a *ActReduce[T]) Copy() uc.Copier {
	act_copy := a.Action.Copy().(*Action[T])
	rule_copy := a.Rule.Copy().(*gr.Production[T])

	a_copy := &ActReduce[T]{
		Action:   act_copy,
		Rule:     rule_copy,
		Original: a.Original,
	}
	return a_copy
}

// NewActReduce creates a new reduce action.
//
// Parameters:
//   - rule: The rule to reduce by.
//
// Returns:
//   - *ActReduce: A pointer to the new reduce action. Nil if the rule is nil.
func NewActReduce[T gr.TokenTyper](rule *gr.Production[T]) *ActReduce[T] {
	if rule == nil {
		return nil
	}

	act := newAction(nil, make([]T, 0))

	rule_copy := rule.Copy().(*gr.Production[T])

	ar := &ActReduce[T]{
		Action:   act,
		Rule:     rule_copy,
		Original: rule,
	}
	return ar
}

// ActAccept represents an accept action.
type ActAccept[T gr.TokenTyper] struct {
	*Action[T]

	// Rule is the Rule to reduce by.
	Rule *gr.Production[T]

	// Original is the Original rule to reduce by.
	// this should never be modified.
	Original *gr.Production[T]
}

// String implements the Actioner interface.
func (a *ActAccept[T]) String() string {
	act_str := a.Action.String()

	var builder strings.Builder

	builder.WriteString("accept{")
	builder.WriteString(act_str)
	builder.WriteRune('}')

	str := builder.String()

	return str
}

// Copy implements the common.Copier interface.
func (a *ActAccept[T]) Copy() uc.Copier {
	act_copy := a.Action.Copy().(*Action[T])

	rule_copy := a.Rule.Copy().(*gr.Production[T])

	aCopy := &ActAccept[T]{
		Action:   act_copy,
		Rule:     rule_copy,
		Original: a.Original,
	}
	return aCopy
}

// NewActAccept creates a new accept action.
//
// Parameters:
//   - rule: The rule to reduce by.
//
// Returns:
//   - *ActAccept: A pointer to the new reduce action. Nil if the rule is nil.
func NewActAccept[T gr.TokenTyper](rule *gr.Production[T]) *ActAccept[T] {
	if rule == nil {
		return nil
	}

	act := newAction(nil, make([]T, 0))

	rule_copy := rule.Copy().(*gr.Production[T])

	ar := &ActAccept[T]{
		Action:   act,
		Rule:     rule_copy,
		Original: rule,
	}
	return ar
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
func MatchAction[T gr.TokenTyper](a Actioner[T], top *gr.Token[T], stack *ud.History[lls.Stacker[*gr.Token[T]]]) error {
	if a == nil {
		return uc.NewErrNilParameter("a")
	} else if top == nil {
		return uc.NewErrNilParameter("top")
	} else if stack == nil {
		return uc.NewErrNilParameter("stack")
	}

	tla := top.GetLookahead()

	ela, ok := a.GetLookahead()
	if ok {
		if tla == nil {
			return uc.NewErrUnexpected("", ela.String())
		} else if ela != tla.GetID() {
			return uc.NewErrUnexpected(top.GoString(), ela.String())
		}
	}

	iter := a.Iterator()

	for {
		rhs, err := iter.Consume()
		if err != nil {
			break
		}

		cmd := lls.NewPop[*gr.Token[T]]()
		err = stack.ExecuteCommand(cmd)
		if err != nil {
			return uc.NewErrUnexpected("", rhs.String())
		}
		top := cmd.Value()

		id := top.GetID()
		if id != rhs {
			return uc.NewErrUnexpected(top.GoString(), rhs.String())
		}
	}

	return nil
}
