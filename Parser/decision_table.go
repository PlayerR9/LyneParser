package Parser

import (
	"fmt"
	"slices"
	"strings"

	cs "github.com/PlayerR9/LyneParser/ConflictSolver"
	gr "github.com/PlayerR9/LyneParser/Grammar"
	hlp "github.com/PlayerR9/MyGoLib/CustomData/Helpers"
	ds "github.com/PlayerR9/MyGoLib/ListLike/DoubleLL"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"

	sfs "github.com/PlayerR9/MyGoLib/Formatting/FString"

	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

type DecisionTable struct {
	table map[string][]*cs.Helper
}

func (dt *DecisionTable) FString(indentLevel int) []string {
	indentation := sfs.NewIndentConfig(sfs.DefaultIndentation, indentLevel, true, false)
	indent := indentation.String()

	result := make([]string, 0)

	counter := 0

	var builder strings.Builder

	for _, elems := range dt.table {
		for _, elem := range elems {
			builder.WriteString(indent)
			builder.WriteString(fmt.Sprintf("%d", counter))
			builder.WriteRune('.')
			builder.WriteRune(' ')

			builder.WriteString(elem.String())

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
		table: make(map[string][]*cs.Helper),
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
		elems := make([]*cs.Helper, 0)

		for j, r := range rules {
			indices := r.IndexOfRhs(s)
			lastIndex := r.Size() - 1

			for _, i := range indices {
				elem, err := cs.NewItem(r, i, i == lastIndex, j)
				if err != nil {
					return err
				}

				elems = append(elems, cs.NewHelper(elem, nil))
			}
		}

		dt.table[s] = elems
	}

	return nil
}

func (dt *DecisionTable) Match(stack *ds.DoubleStack[gr.Tokener]) cs.Actioner {
	top := stack.Pop()

	elems, ok := dt.table[top.GetID()]
	if !ok {
		return cs.NewErrorAction(fmt.Errorf("no elems found for symbol %s", top.GetID()))
	}

	if len(elems) == 1 {
		return elems[0].Action
	}

	results := make([]hlp.HResult[cs.Actioner], 0, len(elems))

	for _, elem := range elems {
		err := elem.Action.Match(top, stack)
		results = append(results, hlp.NewHResult(elem.Action, err))

		// Refuse the stack
		stack.Refuse()

		// Pop the top token
		stack.Pop()
	}

	success := make([]hlp.HResult[cs.Actioner], 0)
	fail := make([]hlp.HResult[cs.Actioner], 0)

	for _, r := range results {
		if r.Reason == nil {
			success = append(success, r)
		} else {
			fail = append(fail, r)
		}
	}

	if len(success) == 0 {
		// Return the most likely error
		// As of now, we will return the first error
		return cs.NewErrorAction(fail[0].Reason)
	} else if len(success) == 1 {
		return success[0].Result
	}

	// Get the longest match
	weights := slext.ApplyWeightFunc(success, func(h hlp.HResult[cs.Actioner]) (float64, bool) {
		return float64(h.Result.Size()), true
	})

	finals := slext.FilterByPositiveWeight(weights)

	if len(finals) == 1 {
		return finals[0].Result
	} else {
		return cs.NewErrorAction(fmt.Errorf("ambiguous grammar"))
	}
}

func (dt *DecisionTable) FixConflicts() error {
	items := make([]*cs.Item, 0)

	for _, elems := range dt.table {
		for _, elem := range elems {
			items = append(items, elem.Item)
		}
	}

	solver := cs.NewConflictSolver(items)

	err := solver.SolveConflicts()
	if err != nil {
		return err
	}

	/*
		WHILE TRUE DO:
			 conflict <- findConflict(s)

			 IF len(conflict) = 0 THEN:
				  BREAK

			 ok, err <- solveAmbiguous(s)
			 IF NOT(err = NULL) THEN:
				  ERROR err

			 IF NOT(ok) THEN:
				  // We cannot continue
				  BREAK

			 solve(rules)
	*/

	return nil
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
