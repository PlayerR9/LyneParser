package ConflictSolver

import (
	"fmt"
	"strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// Item represents an item in a decision table.
type Item struct {
	// Rule is the production rule that the item represents.
	Rule *gr.Production

	// Pos is the position of the item in the production rule.
	Pos int

	// IsReduce is a flag that indicates if the item is a reduce item.
	IsReduce bool

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

	if i.IsReduce {
		builder.WriteString("reduce")
	} else {
		builder.WriteString("shift")
	}

	builder.WriteRune(' ')

	builder.WriteRune(':')

	builder.WriteRune(' ')

	builder.WriteString(fmt.Sprintf("%d", i.ruleIndex))

	builder.WriteRune(')')

	return builder.String()
}

// NewItem is a constructor of Item.
//
// Parameters:
//   - rule: The production rule that the item represents.
//   - pos: The position of the item in the production rule.
//   - isReduce: A flag that indicates if the item is a reduce item.
//
// Returns:
//   - *Item: The pointer to the new Item.
//   - error: An error of type *ers.ErrInvalidParameter if the rule is nil or
//     the pos is out of bounds.
//
// Behaviors:
//   - By default the IsReduce flag is set to false.
func NewItem(rule *gr.Production, pos int, isReduce bool, ruleIndex int) (*Item, error) {
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
		Rule:      rule,
		Pos:       pos,
		IsReduce:  isReduce,
		ruleIndex: ruleIndex,
	}, nil
}

// Size returns the size of the production rule that the item represents.
//
// Returns:
//   - int: The size of the production rule.
func (i *Item) Size() int {
	return i.Rule.Size()
}

// SetAction sets the action of the item.
//
// Parameters:
//   - isReduce: A flag that indicates if the item is a reduce item.
func (i *Item) SetAction(isReduce bool) {
	i.IsReduce = isReduce
}

// SplitShiftReduce splits the items into two slices: one for the shifts and
// one for the reduces.
//
// Parameters:
//   - items: The items to split.
//
// Returns:
//   - []*Item: The shifts.
//   - []*Item: The reduces.
func SplitShiftReduce(items []*Item) ([]*Item, []*Item) {
	shifts := make([]*Item, 0)
	reduces := make([]*Item, 0)

	for _, item := range items {
		if item.IsReduce {
			reduces = append(reduces, item)
		} else {
			shifts = append(shifts, item)
		}
	}

	return shifts, reduces
}
