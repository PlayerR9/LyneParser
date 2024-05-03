package ConflictSolver

import (
	"errors"
	"strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
	intf "github.com/PlayerR9/MyGoLib/Units/Interfaces"
)

type Helper struct {
	Item   *Item
	Action Actioner
}

func (h *Helper) String() string {
	if h == nil {
		return ""
	}

	var builder strings.Builder

	builder.WriteString(h.Item.String())
	builder.WriteRune(' ')
	builder.WriteRune('(')

	if h.Action != nil {
		builder.WriteString(h.Action.String())
	}

	builder.WriteRune(')')

	return builder.String()
}

func (h *Helper) Copy() intf.Copier {
	return &Helper{
		Item:   h.Item.Copy().(*Item),
		Action: h.Action.Copy().(Actioner),
	}
}

func NewHelper(item *Item, action Actioner) *Helper {
	return &Helper{
		Item:   item,
		Action: action,
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
func (h *Helper) ReplaceRhsAt(index int, otherH *Helper) (*Helper, error) {
	if otherH == nil {
		return nil, ers.NewErrNilParameter("otherH")
	}

	newH := &Helper{
		Action: h.Action.Copy().(Actioner),
	}

	var err error

	newH.Item, err = h.Item.ReplaceRhsAt(index, otherH.Item)
	if err != nil {
		return nil, err
	}

	return newH, nil
}

// HasRhs returns true if the right-hand side of the item
// matches the specified right-hand side.
//
// Parameters:
//   - rhs: The right-hand side to search for.
//
// Returns:
//   - bool: True if the right-hand side matches.
func (h *Helper) HasRhs(rhs string) bool {
	if h.Item == nil || h.Item.Rule == nil {
		return false
	}

	return h.Item.Rule.HasRhs(rhs)
}

// IndicesOfRhs returns the indices of the right-hand side of the item
// that match the specified right-hand side.
//
// Parameters:
//   - rhs: The right-hand side to search for.
//
// Returns:
//   - []int: The indices of the right-hand side.
func (h *Helper) IndicesOfRhs(rhs string) []int {
	if h.Item == nil || h.Item.Rule == nil {
		return nil
	}

	return h.Item.Rule.IndicesOfRhs(rhs)
}

// GetRhsAt returns the right-hand side of the item at the specified index.
//
// Parameters:
//   - index: The index of the right-hand side to get.
//
// Returns:
//   - string: The right-hand side of the item.
//   - error: An error if the operation failed.
//
// Errors:
//   - *ers.ErrInvalidParameter: The index is out of bounds, the item is nil,
//     or the item's rule is nil.
func (h *Helper) GetRhsAt(index int) (string, error) {
	if h.Item == nil {
		return "", ers.NewErrNilParameter("h.Item")
	}

	return h.Item.GetRhsAt(index)
}

// GetPos returns the position of the item in the production rule.
//
// Returns:
//   - int: The position of the item. If the item is nil, -1 is returned.
func (h *Helper) GetPos() int {
	if h.Item == nil {
		return -1
	}

	return h.Item.GetPos()
}

// IsReduce returns true if the item is a reduce item.
//
// Returns:
//   - bool: True if the item is a reduce item.
//
// Behavior:
//   - If the item is nil or the item's rule is nil, false is returned.
func (h *Helper) IsReduce() bool {
	if h.Item == nil {
		return false
	}

	return h.Item.IsReduce()
}

// GetRuleIndex returns the index of the item's rule.
//
// Returns:
//   - int: The index of the item's rule. If the item is nil, -1 is returned.
func (h *Helper) GetRuleIndex() int {
	if h.Item == nil {
		return -1
	}

	return h.Item.GetRuleIndex()
}

// SetAction sets the action of the helper.
//
// Parameters:
//   - action: The action to set.
func (h *Helper) SetAction(action Actioner) {
	h.Action = action
}

// Init initializes the helper with the specified symbol.
//
// Parameters:
//   - symbol: The symbol to initialize the helper with.
//
// Returns:
//   - error: An error of type *ErrItemIsNil if the item is nil.
func (h *Helper) Init(symbol string) error {
	if h.Item == nil {
		return NewErrItemIsNil()
	}

	if !h.Item.IsReduce() {
		h.Action = NewActShift()

		return nil
	}

	ri := h.Item.GetRuleIndex()

	if symbol == gr.EOFTokenID {
		h.Action = NewAcceptAction(ri)
	} else {
		h.Action = NewActReduce(ri)
	}

	return nil
}

// GetSymbolsUpToPos returns the symbols of the production rule up to the current position.
//
// Returns:
//   - []string: The symbols of the production rule up to the current position.
//
// Behavior:
//   - Symbols are reversed. Thus, the symbol at index 0 is the current symbol.
func (h *Helper) GetSymbolsUpToPos() []string {
	if h.Item == nil {
		return nil
	}

	return h.Item.GetSymbolsUpToPos()
}

// GetRule returns the rule of the item.
//
// Returns:
//   - *gr.Production: The rule of the item.
func (h *Helper) GetRule() *gr.Production {
	if h.Item == nil {
		return nil
	}

	return h.Item.GetRule()
}

// IsLhsRhs returns true if the left-hand side of the item
// matches the specified right-hand side.
//
// Parameters:
//   - rhs: The right-hand side to compare.
//
// Returns:
//   - bool: True if the left-hand side matches the right-hand side. False otherwise.
func (h *Helper) IsLhsRhs(rhs string) bool {
	if h.Item == nil {
		return false
	}

	return h.Item.IsLhsRhs(rhs)
}

// GetRhs returns the top right-hand side of the item.
//
// Returns:
//   - string: The right-hand side of the item.
//   - error: An error of type *ErrItemIsNil if the item is nil.
func (h *Helper) GetRhs() (string, error) {
	if h.Item == nil {
		return "", NewErrItemIsNil()
	}

	return h.Item.GetRhs(), nil
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

// SetLookahead sets the lookahead of the shift action. If the action is not a shift action,
// this method does nothing.
//
// Returns:
//   - error: An error of type *ErrItemIsNil if the item is nil.
func (h *Helper) SetLookahead() error {
	if h.Item == nil {
		return NewErrItemIsNil()
	} else if h.Action == nil {
		return nil
	}

	pos := h.Item.GetPos()
	if pos == -1 {
		return errors.New("pos is -1")
	}

	lookahead, err := h.Item.GetRhsAt(pos + 1)
	if err != nil || !gr.IsTerminal(lookahead) {
		return nil
	}

	act, ok := h.Action.(*ActShift)
	if ok {
		act.Lookahead = &lookahead
	}

	return nil
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
