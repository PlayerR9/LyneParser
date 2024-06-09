package ConflictSolver

import (
	"fmt"
	"strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	ds "github.com/PlayerR9/MyGoLib/ListLike/DoubleLL"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

type HelperElem interface {
	// SetLookahead sets the lookahead of the action.
	//
	// Parameters:
	//   - lookahead: The lookahead to set.
	SetLookahead(lookahead *string)

	fmt.Stringer
	uc.Copier
}

// Helper represents a helper in a decision table.
type Helper struct {
	// Item is the item of the helper.
	*Item

	// Action is the action of the helper.
	// This can never be nil.
	Action HelperElem
}

// String implements the fmt.Stringer interface.
func (h *Helper) String() string {
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
func (h *Helper) Copy() uc.Copier {
	return &Helper{
		Item:   h.Item.Copy().(*Item),
		Action: h.Action.Copy().(HelperElem),
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
func NewHelper(item *Item, action HelperElem) *Helper {
	if item == nil || action == nil {
		return nil
	}

	return &Helper{
		Item:   item,
		Action: action,
	}
}

// Init initializes the helper with the specified symbol.
//
// Parameters:
//   - symbol: The symbol to initialize the helper with.
func (h *Helper) Init(symbol string) {
	if !h.Item.IsReduce() {
		h.Action = NewActShift()

		return
	}

	r := h.Item.GetRule()

	if symbol == gr.EOFTokenID {
		h.Action = NewActReduce(r, true)
	} else {
		h.Action = NewActReduce(r, false)
	}
}

// SetAction sets the action of the helper.
//
// Parameters:
//   - action: The action to set.
//
// Behaviors:
//   - If the action is nil, then the action is not set.
func (h *Helper) SetAction(action HelperElem) {
	if action == nil {
		return
	}

	h.Action = action
}

// EvaluateLookahead evaluates the lookahead of the action.
//
// Returns:
//   - error: An error if the evaluation failed.
func (h *Helper) EvaluateLookahead() error {
	pos := h.Item.GetPos()

	lookahead, err := h.Item.GetRhsAt(pos + 1)
	if err != nil {
		return fmt.Errorf("failed to evaluate lookahead: %w", err)
	}

	ok := gr.IsTerminal(lookahead)
	if !ok {
		return nil
	}

	h.Action.SetLookahead(&lookahead)

	return nil
}

// GetLookahead returns the lookahead of the action.
//
// Returns:
//   - *string: The lookahead token ID.
func (h *Helper) GetLookahead() *string {
	var lookahead *string

	switch act := h.Action.(type) {
	case *ActReduce:
		lookahead = act.GetLookahead()
	case *ActShift:
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
func (h *Helper) AppendRhs(symbol string) error {
	switch act := h.Action.(type) {
	case *ActReduce:
		act.AppendRhs(symbol)
	case *ActShift:
		act.AppendRhs(symbol)
	default:
		return ue.NewErrUnexpectedType("action", act)
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
func (h *Helper) ReplaceRhsAt(index int, rhs string) *Helper {
	itemCopy := h.Item.ReplaceRhsAt(index, rhs)

	return &Helper{
		Item:   itemCopy,
		Action: h.Action.Copy().(HelperElem),
	}
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
func (h *Helper) SubstituteRhsAt(index int, otherH *Helper) *Helper {
	if otherH == nil {
		return h.Copy().(*Helper)
	}

	itemCopy := h.Item.SubstituteRhsAt(index, otherH.Item)

	return &Helper{
		Item:   itemCopy,
		Action: h.Action.Copy().(HelperElem),
	}
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
func (h *Helper) Match(top gr.Tokener, stack *ds.DoubleStack[gr.Tokener]) error {
	var err error

	switch act := h.Action.(type) {
	case *ActReduce:
		err = MatchAction(act.Action, top, stack)
	case *ActShift:
		err = MatchAction(act.Action, top, stack)
	default:
		return ue.NewErrUnexpectedType("action", act)
	}

	// Refuse the stack
	stack.Refuse()

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
func (h *Helper) Size() int {
	switch act := h.Action.(type) {
	case *ActReduce:
		return act.Size()
	case *ActShift:
		return act.Size()
	default:
		return -1
	}
}

// GetAction returns the action of the helper.
//
// Returns:
//   - Actioner: The action of the helper.
func (h *Helper) GetAction() HelperElem {
	return h.Action
}

// ForceLookahead forces the lookahead of the action.
//
// Parameters:
//   - lookahead: The lookahead to force.
func (h *Helper) ForceLookahead(lookahead string) {
	h.Action.SetLookahead(&lookahead)
}

/////////////////////////////////////////////////////////////
