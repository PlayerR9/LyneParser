package Parser

import (
	"fmt"
	"slices"
	"strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	ds "github.com/PlayerR9/MyGoLib/ListLike/DoubleLL"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"

	sfs "github.com/PlayerR9/MyGoLib/Formatting/FString"

	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

type DecisionTable struct {
	table map[string][]*Item

	actions map[string][]Actioner
}

func (dt *DecisionTable) FString(indentLevel int) []string {
	indentation := sfs.NewIndentConfig(sfs.DefaultIndentation, indentLevel, true, false)
	indent := indentation.String()

	result := make([]string, 0)

	counter := 0

	var builder strings.Builder

	for _, items := range dt.table {
		for _, item := range items {
			builder.WriteString(indent)
			builder.WriteString(fmt.Sprintf("%d", counter))
			builder.WriteRune('.')
			builder.WriteRune(' ')
			builder.WriteString(item.String())

			result = append(result, builder.String())
			builder.Reset()

			counter++
		}

		result = append(result, "") // Add a blank line
	}

	result = result[:len(result)-1] // Remove the last blank line

	return result
}

func NewDecisionTable() *DecisionTable {
	return &DecisionTable{
		table:   make(map[string][]*Item),
		actions: make(map[string][]Actioner),
	}
}

func (dt *DecisionTable) GenerateItems(rules []*gr.Production) error {
	if len(rules) == 0 {
		return ers.NewErrInvalidParameter("rules", ers.NewErrEmptySlice())
	}

	symbols := make([]string, 0)

	for _, r := range rules {
		tmp := r.GetSymbols()

		for _, s := range tmp {
			if !slices.Contains(symbols, s) {
				symbols = append(symbols, s)
			}
		}
	}

	for _, s := range symbols {
		items := make([]*Item, 0)

		for j, r := range rules {
			indices := r.IndexOfRhs(s)
			lastIndex := r.Size() - 1

			for _, i := range indices {
				item, err := NewItem(r, i, i == lastIndex, j)
				if err != nil {
					return err
				}

				items = append(items, item)
			}
		}

		dt.table[s] = items
	}

	return nil
}

func (dt *DecisionTable) HasShiftReduceConflict() bool {
	for _, items := range dt.table {
		shifts := make([]*Item, 0)
		reduces := make([]*Item, 0)

		for _, item := range items {
			if item.IsReduce {
				reduces = append(reduces, item)
			} else {
				shifts = append(shifts, item)
			}
		}

		if len(shifts) > 0 && len(reduces) > 0 {
			return true
		}
	}

	return false
}

func (dt *DecisionTable) Match(stack *ds.DoubleStack[gr.Tokener]) Action {
	top := stack.Pop()

	items, ok := dt.table[top.GetID()]
	if !ok {
		return NewErrorAction(fmt.Errorf("no items found for symbol %s", top.GetID()))
	}

	// If there are no reduce items, then we can only shift.
	// Therefore, we do not care about finding the exact shift rule.
	shifts, reduces := SplitShiftReduce(items)

	if len(reduces) == 0 {
		if len(shifts) == 0 {
			return NewErrorAction(fmt.Errorf("no actions found for symbol %s", top.GetID()))
		}

		// We can only shift.
		return NewShiftAction()
	}

	type Helper struct {
		Elem   *Item
		Reason error
	}

	results := make([]*Helper, 0, len(reduces))

	for _, item := range reduces {
		if item.Pos == 0 {
			results = append(results, &Helper{Elem: item, Reason: nil})

			continue
		}

		var reason error = nil

		for i := item.Pos - 1; i >= 0; i-- {
			rhs, err := item.Rule.GetRhsAt(i)
			if err != nil {
				reason = fmt.Errorf("could not get RHS at index %d", i)
				break
			} else if stack.IsEmpty() {
				reason = ers.NewErrUnexpected(nil, rhs)
				break
			}

			if top := stack.Pop(); top.GetID() != rhs {
				reason = ers.NewErrUnexpected(top, rhs)
				break
			}
		}

		results = append(results, &Helper{Elem: item, Reason: reason})

		stack.Refuse()
	}

	success := make([]*Helper, 0)
	fail := make([]*Helper, 0)

	for _, r := range results {
		if r.Reason == nil {
			success = append(success, r)
		} else {
			fail = append(fail, r)
		}
	}

	if len(success) == 0 {
		if len(shifts) == 0 {
			// Return the most likely error
			// As of now, we will return the first error
			return NewErrorAction(fail[0].Reason)
		}

		// We can only shift
		return NewShiftAction()
	}

	// Find the actual reduce item

	weights := slext.ApplyWeightFunc(success, func(h *Helper) (float64, bool) {
		return float64(h.Elem.Pos), true
	})

	final := slext.FilterByPositiveWeight(weights)

	if len(final) == 1 {
		return NewReduceAction(final[0].Elem.ruleIndex)
	}

	// AMBIGUOUS GRAMMAR

	// SHIFT-REDUCE CONFLICT
}

/*
0. key -> [WORD] (reduce : 1)
1. key -> key [WORD] (reduce : 2)

2. arrayObj -> [OP_SQUARE] mapObj CL_SQUARE (shift : 3)

3. arrayObj -> OP_SQUARE mapObj [CL_SQUARE] (reduce : 3)

4. mapObj -> fieldCls OP_CURLY mapObj1 [CL_CURLY] (reduce : 4)

5. fieldCls1 -> [ATTR] (reduce : 8)
6. fieldCls1 -> [ATTR] SEP fieldCls1 (shift : 9)

7. source -> arrayObj [EOF] (reduce : 0)

8. fieldCls -> key [OP_PAREN] fieldCls1 CL_PAREN (shift : 7)

9. fieldCls -> key OP_PAREN fieldCls1 [CL_PAREN] (reduce : 7)

10. fieldCls1 -> ATTR [SEP] fieldCls1 (shift : 9)

11. source -> [arrayObj] EOF (shift : 0)

12. arrayObj -> OP_SQUARE [mapObj] CL_SQUARE (shift : 3)

13. mapObj -> [fieldCls] OP_CURLY mapObj1 CL_CURLY (shift : 4)
14. mapObj1 -> [fieldCls] (reduce : 5)
15. mapObj1 -> [fieldCls] mapObj1 (shift : 6)

16. mapObj -> fieldCls [OP_CURLY] mapObj1 CL_CURLY (shift : 4)

17. fieldCls -> key OP_PAREN [fieldCls1] CL_PAREN (shift : 7)
18. fieldCls1 -> ATTR SEP [fieldCls1] (reduce : 9)


19. key -> [key] WORD (shift : 2)
20. fieldCls -> [key] OP_PAREN fieldCls1 CL_PAREN (shift : 7)

21. mapObj -> fieldCls OP_CURLY [mapObj1] CL_CURLY (shift : 4)
22. mapObj1 -> fieldCls [mapObj1] (reduce : 6)
*/
