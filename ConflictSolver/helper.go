package ConflictSolver

import (
	"fmt"
	"strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
	ud "github.com/PlayerR9/MyGoLib/Units/Debugging"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

type HelperElem[T gr.TokenTyper] interface {
	// SetLookahead sets the lookahead of the action.
	//
	// Parameters:
	//   - lookahead: The lookahead to set.
	SetLookahead(lookahead *T)

	// AppendRhs appends a symbol to the right-hand side of the action.
	//
	// Parameters:
	//   - symbol: The symbol to append.
	AppendRhs(symbol T)

	Actioner[T]

	fmt.Stringer
	uc.Copier
}

// Helper represents a helper in a decision table.
type Helper[T gr.TokenTyper] struct {
	// Item is the item of the helper.
	*Item[T]

	// Action is the action of the helper.
	// This can never be nil.
	Action HelperElem[T]
}

// AddChild implements Tree.Noder.
func (h *Helper[T]) AddChild(child tr.Noder) {
	panic("unimplemented")
}

// Cleanup implements Tree.Noder.
func (h *Helper[T]) Cleanup() {
	panic("unimplemented")
}

// DeleteChild implements Tree.Noder.
func (h *Helper[T]) DeleteChild(target tr.Noder) []tr.Noder {
	panic("unimplemented")
}

// FString implements Tree.Noder.
func (h *Helper[T]) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	panic("unimplemented")
}

// GetAncestors implements Tree.Noder.
func (h *Helper[T]) GetAncestors() []tr.Noder {
	panic("unimplemented")
}

// GetFirstChild implements Tree.Noder.
func (h *Helper[T]) GetFirstChild() tr.Noder {
	panic("unimplemented")
}

// GetLeaves implements Tree.Noder.
func (h *Helper[T]) GetLeaves() []tr.Noder {
	panic("unimplemented")
}

// GetParent implements Tree.Noder.
func (h *Helper[T]) GetParent() tr.Noder {
	panic("unimplemented")
}

// IsLeaf implements Tree.Noder.
func (h *Helper[T]) IsLeaf() bool {
	panic("unimplemented")
}

// IsSingleton implements Tree.Noder.
func (h *Helper[T]) IsSingleton() bool {
	panic("unimplemented")
}

// Iterator implements Tree.Noder.
func (h *Helper[T]) Iterator() uc.Iterater[tr.Noder] {
	panic("unimplemented")
}

// LinkChildren implements Tree.Noder.
func (h *Helper[T]) LinkChildren(children []tr.Noder) {
	panic("unimplemented")
}

// RemoveNode implements Tree.Noder.
func (h *Helper[T]) RemoveNode() []tr.Noder {
	panic("unimplemented")
}

// SetParent implements Tree.Noder.
func (h *Helper[T]) SetParent(parent tr.Noder) bool {
	panic("unimplemented")
}

// TreeOf implements Tree.Noder.
func (h *Helper[T]) TreeOf() *tr.Tree {
	panic("unimplemented")
}

// String implements the fmt.Stringer interface.
func (h *Helper[T]) String() string {
	var item_str string

	if h.Item != nil {
		item_str = h.Item.String()
	} else {
		item_str = "no item"
	}

	var builder strings.Builder

	builder.WriteString(item_str)
	builder.WriteString(" (")
	builder.WriteString(h.Action.String())
	builder.WriteRune(')')

	str := builder.String()
	return str
}

// Copy implements the common.Copier interface.
func (h *Helper[T]) Copy() uc.Copier {
	item_copy := h.Item.Copy().(*Item[T])
	act_copy := h.Action.Copy().(HelperElem[T])

	h_copy := &Helper[T]{
		Item:   item_copy,
		Action: act_copy,
	}

	return h_copy
}

// NewHelper is a constructor of Helper.
//
// Parameters:
//   - item: The item of the helper.
//   - action: The action of the helper.
//
// Returns:
//   - *Helper: The pointer to the new Helper. Nil if item or action is nil.
func NewHelper[T gr.TokenTyper](item *Item[T], action HelperElem[T]) *Helper[T] {
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

	ok := lookahead.IsTerminal()
	if !ok {
		return nil
	}

	h.Action.SetLookahead(&lookahead)

	return nil
}

// GetLookahead returns the lookahead of the action.
//
// Returns:
//   - T: The lookahead token ID.
//   - bool: True if the lookahead is set, false otherwise.
func (h *Helper[T]) GetLookahead() (T, bool) {
	lookahead, ok := h.Action.GetLookahead()
	return lookahead, ok
}

// AppendRhs appends a symbol to the right-hand side of the action.
//
// Parameters:
//   - symbol: The symbol to append.
func (h *Helper[T]) AppendRhs(symbol T) {
	h.Action.AppendRhs(symbol)
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
	item_copy := h.Item.ReplaceRhsAt(index, rhs)
	act_copy := h.Action.Copy().(HelperElem[T])

	h_copy := &Helper[T]{
		Item:   item_copy,
		Action: act_copy,
	}
	return h_copy
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
func (h *Helper[T]) SubstituteRhsAt(index int, other_helper *Helper[T]) *Helper[T] {
	if other_helper == nil {
		h_copy := h.Copy().(*Helper[T])
		return h_copy
	}

	item_copy := h.Item.SubstituteRhsAt(index, other_helper.Item)
	act_copy := h.Action.Copy().(HelperElem[T])

	h_copy := &Helper[T]{
		Item:   item_copy,
		Action: act_copy,
	}

	return h_copy
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
	act, ok := h.Action.(Actioner[T])
	uc.Assert(ok, "In Helper.Match: h.Action is not an Actioner")

	err := MatchAction(act, top, stack)

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
func (h *Helper[T]) Size() int {
	size := h.Action.Size()
	return size
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
