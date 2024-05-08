package ConflictSolver

import (
	"fmt"
	"slices"
	"strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"

	intf "github.com/PlayerR9/MyGoLib/Units/Common"
)

// Item represents an item in a decision table.
type Item struct {
	// Rule is the production rule that the item represents.
	Rule *gr.Production

	// Pos is the position of the item in the production rule.
	Pos int

	// ruleIndex is the index of the rule in the decision table.
	ruleIndex int
}

// String returns a string representation of the item.
//
// Returns:
//   - string: The string representation of the item.
func (i *Item) String() string {
	var builder strings.Builder

	builder.WriteString(i.Rule.GetLhs())
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
			builder.WriteString(rhs)
			builder.WriteRune(']')
		} else {
			builder.WriteString(rhs)
		}
	}

	builder.WriteRune(' ')
	builder.WriteRune('(')

	builder.WriteString(fmt.Sprintf("%d", i.ruleIndex))

	builder.WriteRune(')')

	return builder.String()
}

// Copy creates a copy of the item.
//
// Returns:
//   - intf.Copier: The copy of the item.
func (i *Item) Copy() intf.Copier {
	return &Item{
		Rule:      i.Rule.Copy().(*gr.Production),
		Pos:       i.Pos,
		ruleIndex: i.ruleIndex,
	}
}

// NewItem is a constructor of Item.
//
// Parameters:
//   - rule: The production rule that the item represents.
//   - pos: The position of the item in the production rule.
//
// Returns:
//   - *Item: The pointer to the new Item.
//   - error: An error of type *ers.ErrInvalidParameter if the rule is nil or
//     the pos is out of bounds.
func NewItem(rule *gr.Production, pos int) (*Item, error) {
	if rule == nil {
		return nil, ers.NewErrNilParameter("rule")
	}

	size := rule.Size()

	if pos < 0 || pos >= size {
		return nil, ers.NewErrInvalidParameter(
			"pos",
			ers.NewErrOutOfBounds(pos, 0, size),
		)
	}

	return &Item{
		Rule: rule,
		Pos:  pos,
	}, nil
}

// GetPos returns the position of the item in the production rule.
//
// Returns:
//   - int: The position of the item.
func (item *Item) GetPos() int {
	return item.Pos
}

// GetRhsAt returns the right-hand side of the production rule at the specified index.
//
// Parameters:
//   - index: The index of the right-hand side to get.
//
// Returns:
//   - string: The right-hand side of the production rule.
//   - error: An error if it is unable to get the right-hand side.
//
// Errors:
//   - *ers.ErrInvalidParameter: If the index is out of bounds or the item's rule
//     is nil.
func (item *Item) GetRhsAt(index int) (string, error) {
	return item.Rule.GetRhsAt(index)
}

// GetRhs returns the right-hand side of the production rule at the current position.
//
// Returns:
//   - string: The right-hand side of the production rule.
func (item *Item) GetRhs() string {
	rhs, err := item.Rule.GetRhsAt(item.Pos)
	if err != nil {
		return ""
	}

	return rhs
}

// GetSymbolsUpToPos returns the symbols of the production rule up to the current position.
//
// Returns:
//   - []string: The symbols of the production rule up to the current position.
//
// Behaviors:
//   - The symbols are reversed. Thus, the symbol at index 0 is the current symbol
//     of the item.
func (item *Item) GetSymbolsUpToPos() []string {
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
func (item *Item) IsReduce() bool {
	return item.Pos == item.Rule.Size()
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
//   - *ers.ErrInvalidParameter: If the other item is nil, otherI.Rule is nil,
//     or the index is out of bounds.
//   - *gr.ErrLhsRhsMismatch: If the left-hand side of the production rule does
//     not match the right-hand side.
func (item *Item) ReplaceRhsAt(index int, otherI *Item) (*Item, error) {
	if otherI == nil {
		return nil, ers.NewErrNilParameter("otherI")
	}

	newItem := item.Copy().(*Item)

	var err error

	newItem.Rule, err = newItem.Rule.ReplaceRhsAt(index, otherI.Rule)
	if err != nil {
		return nil, err
	}

	return newItem, nil
}

// GetRule returns the production rule that the item represents.
//
// Returns:
//   - *gr.Production: The production rule that the item represents.
func (item *Item) GetRule() *gr.Production {
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
func (item *Item) IsLhsRhs(rhs string) bool {
	return item.Rule.GetLhs() == rhs
}

// IndicesOfRhs returns the indices of the right-hand side of the item
// that match the specified right-hand side.
//
// Parameters:
//   - rhs: The right-hand side to search for.
//
// Returns:
//   - []int: The indices of the right-hand side.
func (item *Item) IndicesOfRhs(rhs string) []int {
	return item.Rule.IndicesOfRhs(rhs)
}

/////////////////////////////////////////////////////////////
