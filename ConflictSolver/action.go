package ConflictSolver

import (
	"fmt"
	"strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	ud "github.com/PlayerR9/MyGoLib/Units/Debugging"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// Action represents an action in a decision table.
type Action[T uc.Enumer] struct {
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

// Iterator implements the Iterators.Iterater interface.
func (a *Action[T]) Iterator() uc.Iterater[T] {
	iter := uc.NewSimpleIterator(a.rhs)
	return iter
}

// AppendRhs appends a right-hand side token to the action.
//
// Parameters:
//   - rhs: The right-hand side token to append.
func (a *Action[T]) AppendRhs(rhs T) {
	a.rhs = append(a.rhs, rhs)
}

// Size returns the size of the right-hand side tokens.
//
// Returns:
//   - int: The size of the right-hand side tokens.
func (a *Action[T]) Size() int {
	return len(a.rhs)
}

// SetLookahead sets the lookahead token ID for the action.
//
// Parameters:
//   - lookahead: The lookahead token ID to set.
func (a *Action[T]) SetLookahead(lookahead *T) {
	a.lookahead = lookahead
}

// GetLookahead returns the lookahead token ID for the action.
//
// Returns:
//   - *T: The lookahead token ID.
func (a *Action[T]) GetLookahead() *T {
	return a.lookahead
}

// Actioner represents an action that the parser will take.
type Actioner[T uc.Enumer] interface {
	// GetLookahead returns the lookahead token ID for the action.
	//
	// Returns:
	//   - *T: The lookahead token ID.
	GetLookahead() *T

	fmt.Stringer

	uc.Iterable[T]
}

// ActShift represents a shift action.
type ActShift[T uc.Enumer] struct {
	*Action[T]
}

// String implements the Actioner interface.
func (a *ActShift[T]) String() string {
	var builder strings.Builder

	builder.WriteString("shift")
	builder.WriteRune('{')
	builder.WriteString(a.Action.String())
	builder.WriteRune('}')

	return builder.String()
}

// Copy implements the common.Copier interface.
func (a *ActShift[T]) Copy() uc.Copier {
	rhsCopy := make([]T, len(a.Action.rhs))
	copy(rhsCopy, a.Action.rhs)

	return &ActShift[T]{
		Action: &Action[T]{
			lookahead: a.Action.lookahead,
			rhs:       rhsCopy,
		},
	}
}

// NewActShift creates a new shift action.
//
// Returns:
//   - *ActShift: A pointer to the new shift action.
func NewActShift[T uc.Enumer]() *ActShift[T] {
	as := &ActShift[T]{
		Action: &Action[T]{
			lookahead: nil,
			rhs:       make([]T, 0),
		},
	}
	return as
}

// ActReduce represents a reduce action.
type ActReduce[T uc.Enumer] struct {
	*Action[T]

	// rule is the rule to reduce by.
	rule *gr.Production[T]

	// original is the original rule to reduce by.
	// this should never be modified.
	original *gr.Production[T]

	// shouldAccept is true if the reduce action should accept.
	shouldAccept bool
}

// String implements the Actioner interface.
func (a *ActReduce[T]) String() string {
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
func (a *ActReduce[T]) Copy() uc.Copier {
	rhsCopy := make([]T, len(a.Action.rhs))
	copy(rhsCopy, a.Action.rhs)

	ar := &ActReduce[T]{
		Action: &Action[T]{
			lookahead: a.Action.lookahead,
			rhs:       rhsCopy,
		},
		rule:         a.rule.Copy().(*gr.Production[T]),
		original:     a.original,
		shouldAccept: a.shouldAccept,
	}
	return ar
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
func NewActReduce[T uc.Enumer](rule *gr.Production[T], shouldAccept bool) *ActReduce[T] {
	if rule == nil {
		return nil
	}

	ar := &ActReduce[T]{
		Action: &Action[T]{
			lookahead: nil,
			rhs:       make([]T, 0),
		},
		rule:         rule.Copy().(*gr.Production[T]),
		original:     rule,
		shouldAccept: shouldAccept,
	}
	return ar
}

// GetRule returns the rule to reduce by.
//
// Returns:
//   - *gr.Production: The rule to reduce by.
func (a *ActReduce[T]) GetRule() *gr.Production[T] {
	return a.rule
}

// GetOriginal returns the original rule to reduce by.
//
// Returns:
//   - *gr.Production: The original rule to reduce by.
func (a *ActReduce[T]) GetOriginal() *gr.Production[T] {
	return a.original
}

// ShouldAccept returns true if the reduce action should accept.
//
// Returns:
//   - bool: True if the reduce action should accept.
func (a *ActReduce[T]) ShouldAccept() bool {
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
func MatchAction[T uc.Enumer](a Actioner[T], top *gr.Token[T], stack *ud.History[lls.Stacker[*gr.Token[T]]]) error {
	ela := a.GetLookahead()
	tla := top.GetLookahead()

	if ela != nil {
		if tla == nil {
			return uc.NewErrUnexpected("", (*ela).String())
		} else if *ela != tla.GetID() {
			return uc.NewErrUnexpected(top.GoString(), (*ela).String())
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
