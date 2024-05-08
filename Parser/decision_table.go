package Parser

import (
	"fmt"

	cs "github.com/PlayerR9/LyneParser/ConflictSolver"
	gr "github.com/PlayerR9/LyneParser/Grammar"
	hlp "github.com/PlayerR9/MyGoLib/CustomData/Helpers"
	ds "github.com/PlayerR9/MyGoLib/ListLike/DoubleLL"

	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

/////////////////////////////////////////////////////////////

type DecisionTable struct {
	table map[string][]*cs.Helper
}

func (dt *DecisionTable) Match(stack *ds.DoubleStack[gr.Tokener]) (cs.Actioner, error) {
	top, err := stack.Pop()
	if err != nil {
		return nil, fmt.Errorf("no top token found")
	}

	elems, ok := dt.table[top.GetID()]
	if !ok {
		return nil, fmt.Errorf("no elems found for symbol %s", top.GetID())
	}

	if len(elems) == 1 {
		return elems[0].Action, nil
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
		if r.First == nil {
			success = append(success, r)
		} else {
			fail = append(fail, r)
		}
	}

	if len(success) == 0 {
		// Return the most likely error
		// As of now, we will return the first error
		return nil, fail[0].Second
	} else if len(success) == 1 {
		return success[0].First, nil
	}

	// Get the longest match
	weights := slext.ApplyWeightFunc(success, HResultWeightFunc)

	finals := slext.FilterByPositiveWeight(weights)

	if len(finals) == 1 {
		return finals[0].First, nil
	} else {
		return nil, fmt.Errorf("ambiguous grammar")
	}
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
