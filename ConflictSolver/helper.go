package ConflictSolver

import (
	"fmt"
	"strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	ud "github.com/PlayerR9/MyGoLib/Units/Debugging"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

type HelperElem[T uc.Enumer] interface {
	// SetLookahead sets the lookahead of the action.
	//
	// Parameters:
	//   - lookahead: The lookahead to set.
	SetLookahead(lookahead *T)

	fmt.Stringer
	uc.Copier
}

// Helper represents a helper in a decision table.
type Helper[T uc.Enumer] struct {
	// Item is the item of the helper.
	*Item[T]

	// Action is the action of the helper.
	// This can never be nil.
	Action HelperElem[T]
}

// String implements the fmt.Stringer interface.
func (h *Helper[T]) String() string {
	var builder strings.Builder

	if h.Item != nil {
		builder.WriteString(h.Item.String())
	} else {
		builder.WriteString("no item")
	}

	builder.WriteString(" (")
	builder.WriteString(h.Action.String())
	builder.WriteRune(')')

	return builder.String()
}

// Copy implements the common.Copier interface.
func (h *Helper[T]) Copy() uc.Copier {
	return &Helper[T]{
		Item:   h.Item.Copy().(*Item[T]),
		Action: h.Action.Copy().(HelperElem[T]),
	}
}

// NewHelper is a constructor of Helper.
//
// Parameters:
//   - item: The item of the helper.
//   - action: The action of the helper.
//
// Returns:
//   - *Helper: The pointer to the new Helper.
//
// Behaviors:
//   - If the item or action are nil, then nil is returned.
func NewHelper[T uc.Enumer](item *Item[T], action HelperElem[T]) *Helper[T] {
	if item == nil || action == nil {
		return nil
	}

	h := &Helper[T]{
		Item:   item,
		Action: action,
	}

	return h
}

// SetAction sets the action of the helper.
//
// Parameters:
//   - action: The action to set.
//
// Behaviors:
//   - If the action is nil, then the action is not set.
func (h *Helper[T]) SetAction(action HelperElem[T]) {
	if action == nil {
		return
	}

	h.Action = action
}

// EvaluateLookahead evaluates the lookahead of the action.
//
// Returns:
//   - error: An error if the evaluation failed.
func (h *Helper[T]) EvaluateLookahead() error {
	pos := h.Item.GetPos()

	lookahead, err := h.Item.GetRhsAt(pos + 1)
	if err != nil {
		return fmt.Errorf("failed to evaluate lookahead: %w", err)
	}

	ok := gr.IsTerminal(lookahead.String())
	if !ok {
		return nil
	}

	h.Action.SetLookahead(&lookahead)

	return nil
}

// GetLookahead returns the lookahead of the action.
//
// Returns:
//   - *T: The lookahead token ID.
func (h *Helper[T]) GetLookahead() *T {
	var lookahead *T

	switch act := h.Action.(type) {
	case *ActReduce[T]:
		lookahead = act.GetLookahead()
	case *ActShift[T]:
		lookahead = act.GetLookahead()
	}

	return lookahead
}

// AppendRhs appends a symbol to the right-hand side of the action.
//
// Parameters:
//   - symbol: The symbol to append.
//
// Returns:
//   - error: An error of type *ErrNoActionProvided if the action is nil.
func (h *Helper[T]) AppendRhs(symbol T) error {
	switch act := h.Action.(type) {
	case *ActReduce[T]:
		act.AppendRhs(symbol)
	case *ActShift[T]:
		act.AppendRhs(symbol)
	default:
		return uc.NewErrUnexpectedType("action", act)
	}

	return nil
}

// ReplaceRhsAt replaces the right-hand side of the item
// at the specified index with the right-hand side of the other item.
//
// Parameters:
//   - index: The index of the right-hand side to replace.
//   - otherH: The other helper.
//
// Returns:
//   - *Helper: The new helper with the replaced right-hand side.
//   - error: An error if the operation failed.
//
// Errors:
//   - *ers.ErrInvalidParameter: The index is out of bounds, otherH is nil,
//     otherH.Item is nil, or otherH.Item.Rule is nil.
//   - *gr.ErrLhsRhsMismatch: The left-hand side of the item does not match the
//     right-hand side of the other item.
func (h *Helper[T]) ReplaceRhsAt(index int, rhs T) *Helper[T] {
	itemCopy := h.Item.ReplaceRhsAt(index, rhs)

	hCopy := &Helper[T]{
		Item:   itemCopy,
		Action: h.Action.Copy().(HelperElem[T]),
	}
	return hCopy
}

// ReplaceRhsAt replaces the right-hand side of the item
// at the specified index with the right-hand side of the other item.
//
// Parameters:
//   - index: The index of the right-hand side to replace.
//   - otherH: The other helper.
//
// Returns:
//   - *Helper: The new helper with the replaced right-hand side.
//   - error: An error if the operation failed.
//
// Errors:
//   - *ers.ErrInvalidParameter: The index is out of bounds, otherH is nil,
//     otherH.Item is nil, or otherH.Item.Rule is nil.
//   - *gr.ErrLhsRhsMismatch: The left-hand side of the item does not match the
//     right-hand side of the other item.
func (h *Helper[T]) SubstituteRhsAt(index int, otherH *Helper[T]) *Helper[T] {
	if otherH == nil {
		hCopy := h.Copy().(*Helper[T])
		return hCopy
	}

	itemCopy := h.Item.SubstituteRhsAt(index, otherH.Item)

	hCopy := &Helper[T]{
		Item:   itemCopy,
		Action: h.Action.Copy().(HelperElem[T]),
	}

	return hCopy
}

// Match matches the top of the stack with the helper.
//
// Parameters:
//   - top: The top of the stack.
//   - stack: The stack.
//
// Returns:
//   - error: An error if the match failed.
//
// Behaviors:
//   - The stack is refused.
func (h *Helper[T]) Match(top *gr.Token[T], stack *ud.History[lls.Stacker[*gr.Token[T]]]) error {
	var err error

	switch act := h.Action.(type) {
	case *ActReduce[T]:
		err = MatchAction(act.Action, top, stack)
	case *ActShift[T]:
		err = MatchAction(act.Action, top, stack)
	default:
		return uc.NewErrUnexpectedType("action", act)
	}

	// Refuse the stack
	stack.Reject()

	if err != nil {
		return err
	}

	return nil
}

// Size returns the size of the helper.
//
// Returns:
//   - int: The size of the helper.
//
// Behaviors:
//   - If the action is invalid, -1 is returned.
func (h *Helper[T]) Size() int {
	switch act := h.Action.(type) {
	case *ActReduce[T]:
		return act.Size()
	case *ActShift[T]:
		return act.Size()
	default:
		return -1
	}
}

// GetAction returns the action of the helper.
//
// Returns:
//   - Actioner: The action of the helper.
func (h *Helper[T]) GetAction() HelperElem[T] {
	return h.Action
}

// ForceLookahead forces the lookahead of the action.
//
// Parameters:
//   - lookahead: The lookahead to force.
func (h *Helper[T]) ForceLookahead(lookahead T) {
	h.Action.SetLookahead(&lookahead)
}
