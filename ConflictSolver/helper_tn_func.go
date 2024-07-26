package ConflictSolver

import (
	"fmt"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	ud "github.com/PlayerR9/MyGoLib/Units/Debugging"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	lls "github.com/PlayerR9/stack/stack"
)

// SetAction sets the action of the helper.
//
// Parameters:
//   - action: The action to set.
//
// Behaviors:
//   - If the action is nil, then the action is not set.
func (h *HelperNode[T]) SetAction(action HelperElem[T]) {
	if action == nil {
		return
	}

	h.Action = action
}

// EvaluateLookahead evaluates the lookahead of the action.
//
// Returns:
//   - error: An error if the evaluation failed.
func (h *HelperNode[T]) EvaluateLookahead() error {
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
func (h *HelperNode[T]) GetLookahead() (T, bool) {
	lookahead, ok := h.Action.GetLookahead()
	return lookahead, ok
}

// AppendRhs appends a symbol to the right-hand side of the action.
//
// Parameters:
//   - symbol: The symbol to append.
func (h *HelperNode[T]) AppendRhs(symbol T) {
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
func (h *HelperNode[T]) ReplaceRhsAt(index int, rhs T) *HelperNode[T] {
	item_copy := h.Item.ReplaceRhsAt(index, rhs)
	act_copy := h.Action.Copy().(HelperElem[T])

	h_copy := &HelperNode[T]{
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
func (h *HelperNode[T]) SubstituteRhsAt(index int, other_helper *HelperNode[T]) *HelperNode[T] {
	if other_helper == nil {
		h_copy := h.Copy().(*HelperNode[T])
		return h_copy
	}

	item_copy := h.Item.SubstituteRhsAt(index, other_helper.Item)
	act_copy := h.Action.Copy().(HelperElem[T])

	h_copy := &HelperNode[T]{
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
func (h *HelperNode[T]) Match(top *gr.Token[T], stack *ud.History[lls.Stacker[*gr.Token[T]]]) error {
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

// ActionSize returns the size of the helper.
//
// Returns:
//   - int: The size of the helper.
func (h *HelperNode[T]) ActionSize() int {
	size := h.Action.ActionSize()
	return size
}

// GetAction returns the action of the helper.
//
// Returns:
//   - Actioner: The action of the helper.
func (h *HelperNode[T]) GetAction() HelperElem[T] {
	return h.Action
}

// ForceLookahead forces the lookahead of the action.
//
// Parameters:
//   - lookahead: The lookahead to force.
func (h *HelperNode[T]) ForceLookahead(lookahead T) {
	h.Action.SetLookahead(&lookahead)
}
