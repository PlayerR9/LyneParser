package ConflictSolver

import (
	"fmt"
	"slices"
	"strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// Item represents an item in a decision table.
type Item[T gr.TokenTyper] struct {
	// Rule is the production rule that the item represents.
	Rule *gr.Production[T]

	// Pos is the position of the item in the production rule.
	Pos int

	// rule_index is the index of the rule in the decision table.
	rule_index int
}

// String implements the fmt.Stringer interface.
func (i *Item[T]) String() string {
	var builder strings.Builder

	builder.WriteString(i.Rule.GetLhs().String())
	builder.WriteRune(' ')
	builder.WriteString(gr.LeftToRight)

	iter := i.Rule.Iterator()

	for pos := 0; ; pos++ {
		rhs, err := iter.Consume()
		if err != nil {
			break
		}

		builder.WriteRune(' ')

		if pos == i.Pos {
			builder.WriteRune('[')
			builder.WriteString(rhs.String())
			builder.WriteRune(']')
		} else {
			builder.WriteString(rhs.String())
		}
	}

	builder.WriteRune(' ')
	builder.WriteRune('(')

	builder.WriteString(fmt.Sprintf("%d", i.rule_index))

	builder.WriteRune(')')

	return builder.String()
}

// Copy implements the Copier interface.
func (i *Item[T]) Copy() uc.Copier {
	rule_copy := i.Rule.Copy().(*gr.Production[T])

	i_copy := &Item[T]{
		Rule:       rule_copy,
		Pos:        i.Pos,
		rule_index: i.rule_index,
	}

	return i_copy
}

// NewItem is a constructor of Item.
//
// Parameters:
//   - rule: The production rule that the item represents.
//   - pos: The position of the item in the production rule.
//   - ruleIndex: The index of the rule in the decision table.
//
// Returns:
//   - *Item: The pointer to the new Item.
//   - error: An error of type *uc.ErrInvalidParameter if the rule is nil or
//     the pos is out of bounds.
func NewItem[T gr.TokenTyper](rule *gr.Production[T], pos int, rule_index int) (*Item[T], error) {
	if rule == nil {
		return nil, uc.NewErrNilParameter("rule")
	}

	size := rule.Size()

	if pos < 0 || pos >= size {
		return nil, uc.NewErrInvalidParameter(
			"pos",
			uc.NewErrOutOfBounds(pos, 0, size),
		)
	}

	item := &Item[T]{
		Rule:       rule,
		Pos:        pos,
		rule_index: rule_index,
	}
	return item, nil
}

// GetPos returns the position of the item in the production rule.
//
// Returns:
//   - int: The position of the item.
func (item *Item[T]) GetPos() int {
	return item.Pos
}

// GetRhsAt returns the right-hand side of the production rule at the specified index.
//
// Parameters:
//   - index: The index of the right-hand side to get.
//
// Returns:
//   - T: The right-hand side of the production rule.
//   - error: An error if it is unable to get the right-hand side.
//
// Errors:
//   - *uc.ErrInvalidParameter: If the index is out of bounds or the item's rule
//     is nil.
func (item *Item[T]) GetRhsAt(index int) (T, error) {
	rhs, err := item.Rule.GetRhsAt(index)
	return rhs, err
}

// GetRhs returns the right-hand side of the production rule at the current position.
//
// Returns:
//   - T: The right-hand side of the production rule.
func (item *Item[T]) GetRhs() T {
	rhs, err := item.Rule.GetRhsAt(item.Pos)
	uc.AssertF(err == nil, "GetRhs: %s", err.Error())

	return rhs
}

// GetSymbolsUpToPos returns the symbols of the production rule up to the current position.
//
// Returns:
//   - []T: The symbols of the production rule up to the current position.
//
// Behaviors:
//   - The symbols are reversed. Thus, the symbol at index 0 is the current symbol
//     of the item.
func (item *Item[T]) GetSymbolsUpToPos() []T {
	symbols := item.Rule.GetSymbols()

	symbols = symbols[:item.Pos+1]

	slices.Reverse(symbols)

	return symbols
}

// IsReduce returns true if the item is a reduce item.
//
// Returns:
//   - bool: True if the item is a reduce item. Otherwise, false.
//
// Behaviors:
//   - If the item's rule is nil, it returns false.
func (item *Item[T]) IsReduce() bool {
	size := item.Rule.Size()
	return item.Pos == size
}

// ReplaceRhsAt replaces the right-hand side of the production rule at the given
// index with the right-hand side of the other item.
//
// Parameters:
//   - index: The index of the right-hand side to replace.
//   - otherI: The other item to replace the right-hand side with.
//
// Returns:
//   - *Item: The new item with the replaced right-hand side.
//   - error: An error if it is unable to replace the right-hand side.
//
// Errors:
//   - *uc.ErrInvalidParameter: If the other item is nil, otherI.Rule is nil,
//     or the index is out of bounds.
//   - *gr.ErrLhsRhsMismatch: If the left-hand side of the production rule does
//     not match the right-hand side.
func (item *Item[T]) ReplaceRhsAt(index int, rhs T) *Item[T] {
	rule_copy := item.Rule.ReplaceRhsAt(index, rhs)

	item_copy := &Item[T]{
		Rule:       rule_copy,
		Pos:        item.Pos,
		rule_index: item.rule_index,
	}

	return item_copy
}

// ReplaceRhsAt replaces the right-hand side of the production rule at the given
// index with the right-hand side of the other item.
//
// Parameters:
//   - index: The index of the right-hand side to replace.
//   - otherI: The other item to replace the right-hand side with.
//
// Returns:
//   - *Item: The new item with the replaced right-hand side.
//   - error: An error if it is unable to replace the right-hand side.
//
// Errors:
//   - *uc.ErrInvalidParameter: If the other item is nil, otherI.Rule is nil,
//     or the index is out of bounds.
//   - *gr.ErrLhsRhsMismatch: If the left-hand side of the production rule does
//     not match the right-hand side.
func (item *Item[T]) SubstituteRhsAt(index int, other_item *Item[T]) *Item[T] {
	if other_item == nil {
		item_copy := item.Copy().(*Item[T])
		return item_copy
	}

	rule_copy := item.Rule.SubstituteRhsAt(index, other_item.Rule)

	item_copy := &Item[T]{
		Rule:       rule_copy,
		Pos:        item.Pos,
		rule_index: item.rule_index,
	}
	return item_copy
}

// GetRule returns the production rule that the item represents.
//
// Returns:
//   - *gr.Production: The production rule that the item represents.
func (item *Item[T]) GetRule() *gr.Production[T] {
	return item.Rule
}

// IsLhsRhs returns true if the left-hand side of the production rule matches the
// right-hand side.
//
// Parameters:
//   - rhs: The right-hand side to compare with the left-hand side.
//
// Returns:
//   - bool: True if the left-hand side matches the right-hand side. Otherwise, false.
func (item *Item[T]) IsLhsRhs(rhs T) bool {
	lhs := item.Rule.GetLhs()
	return lhs == rhs
}

// IndicesOfRhs returns the indices of the right-hand side of the item
// that match the specified right-hand side.
//
// Parameters:
//   - rhs: The right-hand side to search for.
//
// Returns:
//   - []int: The indices of the right-hand side.
func (item *Item[T]) IndicesOfRhs(rhs T) []int {
	indices := item.Rule.IndicesOfRhs(rhs)
	return indices
}
