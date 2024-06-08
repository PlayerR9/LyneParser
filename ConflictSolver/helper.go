package ConflictSolver

import (
	"strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	ds "github.com/PlayerR9/MyGoLib/ListLike/DoubleLL"
	intf "github.com/PlayerR9/MyGoLib/Units/Common"
	ers "github.com/PlayerR9/MyGoLib/Units/errors"
)

// Helper represents a helper in a decision table.
type Helper struct {
	// Item is the item of the helper.
	*Item

	// Action is the action of the helper.
	Action Actioner
}

// String returns a string representation of the helper.
//
// Returns:
//   - string: The string representation of the helper.
func (h *Helper) String() string {
	var builder strings.Builder

	if h.Item != nil {
		builder.WriteString(h.Item.String())
	} else {
		builder.WriteString("no item")
	}

	builder.WriteRune(' ')
	builder.WriteRune('(')
	if h.Action != nil {
		builder.WriteString(h.Action.String())
	}
	builder.WriteRune(')')

	return builder.String()
}

// NewHelper is a constructor of Helper.
//
// Parameters:
//   - item: The item of the helper.
//   - action: The action of the helper.
//
// Returns:
//   - *Helper: The pointer to the new Helper.
//   - error: An error of type *ers.ErrInvalidParameter if the item is nil.
func NewHelper(item *Item, action Actioner) (*Helper, error) {
	if item == nil {
		return nil, ers.NewErrNilParameter("item")
	}

	return &Helper{
		Item:   item,
		Action: action,
	}, nil
}

// SetAction sets the action of the helper.
//
// Parameters:
//   - action: The action to set.
func (h *Helper) SetAction(action Actioner) {
	h.Action = action
}

// EvaluateLookahead evaluates the lookahead of the shift action. If the action is not a
// shift action, this method does nothing.
func (h *Helper) EvaluateLookahead() {
	if h.Action == nil {
		return
	}

	pos := h.Item.GetPos()

	lookahead, err := h.Item.GetRhsAt(pos + 1)
	if err != nil || !gr.IsTerminal(lookahead) {
		return
	}

	act, ok := h.Action.(*ActShift)
	if ok {
		act.SetLookahead(&lookahead)
	}
}

// GetLookahead returns the lookahead of the shift action. If the action is not a shift action,
// this method returns nil.
//
// Returns:
//   - *string: The lookahead token ID.
func (h *Helper) GetLookahead() *string {
	if h.Action == nil {
		return nil
	}

	act, ok := h.Action.(*ActShift)
	if !ok {
		return nil
	}

	return act.Lookahead
}

// AppendRhs appends a symbol to the right-hand side of the action.
//
// Parameters:
//   - symbol: The symbol to append.
//
// Returns:
//   - error: An error of type *ErrNoActionProvided if the action is nil.
func (h *Helper) AppendRhs(symbol string) error {
	if h.Action == nil {
		return NewErrNoActionProvided()
	}

	h.Action.AppendRhs(symbol)

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
func (h *Helper) ReplaceRhsAt(index int, otherH *Helper) (*Helper, error) {
	if otherH == nil {
		return nil, ers.NewErrNilParameter("otherH")
	}

	newH := &Helper{
		Item:   h.Item.Copy().(*Item),
		Action: h.Action.Copy().(Actioner),
	}

	var err error

	newH.Item, err = newH.Item.ReplaceRhsAt(index, otherH.Item)
	if err != nil {
		return nil, err
	}

	return newH, nil
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
	err := h.Action.Match(top, stack)

	// Refuse the stack
	stack.Refuse()

	return err
}

// Size returns the size of the helper.
//
// Returns:
//   - int: The size of the helper.
func (h *Helper) Size() int {
	return h.Action.Size()
}

// GetAction returns the action of the helper.
//
// Returns:
//   - Actioner: The action of the helper.
func (h *Helper) GetAction() Actioner {
	return h.Action
}

/////////////////////////////////////////////////////////////

func (h *Helper) Copy() intf.Copier {
	return &Helper{
		Item:   h.Item.Copy().(*Item),
		Action: h.Action.Copy().(Actioner),
	}
}

// Init initializes the helper with the specified symbol.
//
// Parameters:
//   - symbol: The symbol to initialize the helper with.
//
// Returns:
//   - error: An error of type *ers.ErrInvalidParameter if the rule is nil.
func (h *Helper) Init(symbol string) error {
	if !h.Item.IsReduce() {
		h.Action = NewActShift()

		return nil
	}

	r := h.Item.GetRule()

	var err error

	if symbol == gr.EOFTokenID {
		h.Action, err = NewAcceptAction(r)
	} else {
		h.Action, err = NewActReduce(r)
	}

	return err
}

// IsShift returns true if the action is a shift action.
//
// Returns:
//   - bool: True if the action is a shift action. Otherwise, false.
func (h *Helper) IsShift() bool {
	if h.Action == nil {
		return false
	}

	_, ok := h.Action.(*ActShift)

	return ok
}
