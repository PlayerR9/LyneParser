package ConflictSolver

import (
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	ud "github.com/PlayerR9/MyGoLib/Units/Debugging"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
	lls "github.com/PlayerR9/stack"
)

var (
	debugger *log.Logger = log.New(os.Stdout, "DEBUG: ", log.Lshortfile)
)

var GlobalDebugMode bool = true

// ConflictSolver solves conflicts in a decision table.
type ConflictSolver[T gr.TokenTyper] struct {
	// table is a map of elements in the decision table.
	table map[T][]*HelperNode[T]

	// rt is the rule table.
	rt *RuleTable[T]
}

// FString returns a formatted string representation of the decision table
// multilined with a specific indentation level.
//
// Parameters:
//   - indentLevel: The level of indentation.
//
// Returns:
//   - []string: A formatted string representation of the decision table.
func (cs *ConflictSolver[T]) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	if trav == nil {
		return nil
	}

	counter := 0

	helpers := cs.getHelpers()

	for i, h := range helpers {
		err := trav.AppendString(strconv.Itoa(counter))
		if err != nil {
			return uc.NewErrAt(i, "helper", err)
		}

		err = trav.AppendRune('.')
		if err != nil {
			return uc.NewErrAt(i, "helper", err)
		}

		if h != nil {
			err = trav.AppendRune(' ')
			if err != nil {
				return uc.NewErrAt(i, "helper", err)
			}

			err = trav.AppendString(h.String())
			if err != nil {
				return uc.NewErrAt(i, "helper", err)
			}
		}

		trav.AcceptLine()

		counter++
	}

	return nil
}

// NewConflictSolver is a constructor of ConflictSolver.
//
// Parameters:
//   - symbols: The symbols in the decision table.
//   - rules: The rules in the decision table.
//
// Returns:
//   - *ConflictSolver: The pointer to the new ConflictSolver.
//   - error: An error if the operation failed.
//
// Errors:
//   - *ErrCannotCreateItem: If an item cannot be created.
//   - *uc.ErrInvalidParameter: If the item is nil.
func NewConflictSolver[T gr.TokenTyper](symbols []T, rules []*gr.Production[T]) *ConflictSolver[T] {
	rt := NewRuleTable(symbols, rules)

	cs := &ConflictSolver[T]{
		rt:    rt,
		table: rt.GetBucketsCopy(),
	}
	return cs
}

// getHelpers is a helper function that returns all helpers in the decision table.
//
// Returns:
//   - []*Helper: All helpers in the decision table.
func (cs *ConflictSolver[T]) getHelpers() []*HelperNode[T] {
	var result []*HelperNode[T]

	for _, bucket := range cs.table {
		result = append(result, bucket...)
	}

	return result
}

// GetElemsWithLhs is a method that returns all elements with a specific LHS.
//
// Parameters:
//   - rhs: The RHS to find elements for.
//
// Returns:
//   - []*Helper: The elements with the specified LHS.
func (cs *ConflictSolver[T]) GetElemsWithLhs(rhs T) []*HelperNode[T] {
	filter := func(h *HelperNode[T]) bool {
		ok := h.IsLhsRhs(rhs)
		return ok
	}

	helpers := cs.getHelpers()

	helpers = us.SliceFilter(helpers, filter)

	return helpers
}

// solveSetLookaheadOnShifts is a helper function that evaluates the look-ahead on shifts
// in a simple and coarse way.
func (cs *ConflictSolver[T]) solveSetLookaheadOnShifts() {
	helpers := cs.getHelpers()

	for _, h := range helpers {
		err := h.EvaluateLookahead()
		uc.AssertF(err == nil, "failed to evaluate lookahead: %s", err.Error())
	}
}

// getHelpersWithLookahead is a helper function that returns all helpers with look-ahead
// and groups them by the look-ahead.
//
// Returns:
//   - map[T][]*Helper: The helpers with look-ahead.
func (cs *ConflictSolver[T]) getHelpersWithLookahead() map[T][]*HelperNode[T] {
	groups := make(map[T][]*HelperNode[T])

	todo := cs.getHelpers()

	for _, h := range todo {
		lookahead, ok := h.GetLookahead()
		if ok {
			prev, ok := groups[lookahead]
			if !ok {
				prev = []*HelperNode[T]{h}
			} else {
				prev = append(prev, h)
			}

			groups[lookahead] = prev
		}
	}

	return groups
}

// SolveAmbiguousShifts is a method that solves ambiguous shifts in a decision table.
//
// Returns:
//   - error: An error if the operation failed.
//
// Errors:
//   - *ErrHelpersConflictingSize: If the helpers have conflicting sizes.
//   - *ErrHelper: If there is an error appending the right-hand side to the helper.
func (cs *ConflictSolver[T]) SolveAmbiguousShifts() error {
	cs.solveSetLookaheadOnShifts()

	// Now, those shift actions that have the look-ahead are no longer
	// in conflict with their reduce counterparts.
	// However, there still might be conflicts between shift actions
	// with the same look-ahead.

	helpers_with_lookahead := cs.getHelpersWithLookahead() // these are potential conflicts

	// If there are at least two helpers with the same look-ahead, then there might be a conflict.

	// To solve this, we have to find the minimal amount of information that is needed to
	// unambiguously determine the next action.

	for _, bucket := range helpers_with_lookahead {
		if len(bucket) < 2 {
			continue
		}

		// Solve conflicts.
		err := solve_subgroup(bucket)
		if err != nil {
			return err
		}
	}

	return nil
}

// CMPerSymbol is a type that represents conflicts per symbol.
type CMPerSymbol[T gr.TokenTyper] map[T]uc.Pair[[]*HelperNode[T], int]

// FindConflicts is a method that finds conflicts for a specific symbol.
//
// Parameters:
//   - symbol: The symbol to find conflicts for.
//
// Returns:
//   - CMPerSymbol: The conflicts per symbol.
func (cs *ConflictSolver[T]) FindConflicts() CMPerSymbol[T] {
	conflict_map := make(CMPerSymbol[T])

	for symbol, bucket := range cs.table {
		todo := make([]*HelperNode[T], len(bucket))
		copy(todo, bucket)

		// 1. Remove every shift action that has a look-ahead.
		todo = us.SliceFilter(todo, func(h *HelperNode[T]) bool {
			_, ok := h.GetLookahead()
			return !ok
		})

		conflicts, index_of_conflict := find_conflicts_per_symbol(symbol, todo)
		if index_of_conflict != -1 {
			conflict_map[symbol] = uc.NewPair(conflicts, index_of_conflict)
		}
	}

	return conflict_map
}

// MakeExpansionForests creates a forest of expansion trees rooted at the next symbol of the
// conflicting rules.
//
// Parameters:
//   - index: The index of the conflicting rules.
//   - nextRhs: The next symbol of the conflicting rules.
//
// Returns:
//   - map[*Helper][]*ExpansionTree: The forest of expansion trees.
//   - error: An error of type *ErrHelper if the operation failed.
func (cs *ConflictSolver[T]) MakeExpansionForests(index int, next_rhs map[*HelperNode[T]]T) (map[*HelperNode[T]][]T, error) {
	possible_lookaheads := make(map[*HelperNode[T]][]T)

	for c, rhs := range next_rhs {
		rs := cs.GetElemsWithLhs(rhs)
		if len(rs) == 0 {
			return possible_lookaheads, nil
		}

		var lookaheads []T

		for _, r := range rs {
			tree, err := NewExpansionTreeRootedAt(cs, r)
			if err != nil {
				return possible_lookaheads, NewErrHelper(c, err)
			}

			tree.PruneNonTerminalLeaves()

			collapsed := tree.Collapse()

			if len(collapsed) == 0 {
				continue
			}

			for _, c := range collapsed {
				pos, ok := slices.BinarySearch(lookaheads, c)
				if !ok {
					lookaheads = slices.Insert(lookaheads, pos, c)
				}
			}
		}

		if len(lookaheads) != 0 {
			possible_lookaheads[c] = lookaheads
		}
	}

	return possible_lookaheads, nil
}

// SolveAmbiguous is a method that solves ambiguous conflicts in a decision table.
//
// Parameters:
//   - index: The index of the conflicting rules.
//   - conflicts: The conflicting rules.
//
// Returns:
//   - bool: A boolean value indicating if the operation was successful.
//   - error: An error if the operation failed.
func (cs *ConflictSolver[T]) SolveAmbiguous(index int, conflicts []*HelperNode[T]) (bool, error) {
	// 1. Take the next symbol of each conflicting rule
	next_rhs := make(map[*HelperNode[T]]T)

	for _, c := range conflicts {
		rhs, err := c.GetRhsAt(index + 1)
		if err != nil {
			continue
		}

		next_rhs[c] = rhs
	}

	if len(next_rhs) == 0 {
		return false, nil
	}

	// 2. Make the expansion forests
	possible_lookaheads, err := cs.MakeExpansionForests(index, next_rhs)
	if err != nil {
		return false, err
	} else if len(possible_lookaheads) == 0 {
		return false, nil
	}

	// If there are more than one possible lookaheads,
	// then we have to pick one of them.
	// As of now, we will pick the first one.
	for c, forest := range possible_lookaheads {
		if len(forest) > 1 {
			debugger.Println("Found more than one possible lookaheads. Picking the first one.")
		}

		new_rule := c.ReplaceRhsAt(index+1, forest[0])
		new_rule.ForceLookahead(forest[0])

		for key, bucket := range cs.table {
			slice := make([]*HelperNode[T], 0, len(bucket))

			for _, h := range bucket {
				if h == c {
					slice = append(slice, new_rule)
				} else {
					slice = append(slice, h)
				}
			}

			cs.table[key] = slice
		}
	}

	return true, nil
}

/////////////////////////////////////////////////////////////

// SolveConflicts is a method that solves conflicts in a decision table.
//
// Returns:
//   - error: An error if the operation failed.
func (cs *ConflictSolver[T]) Solve() error {
	for {
		conflict_map := cs.FindConflicts()
		if len(conflict_map) == 0 {
			// No conflicts found.
			break
		}

		done := false

		for _, p := range conflict_map {
			ok, err := cs.SolveAmbiguous(p.Second, p.First)
			if err != nil {
				return err
			}

			if ok {
				done = true
			}
		}

		if !done {
			break
		}

		for _, p := range conflict_map {
			conflicts := p.First

			// Solve conflicts.
			err := solve_subgroup(conflicts)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Match is a method that matches the top of the stack with the elements in the decision table.
//
// Parameters:
//   - stack: The stack to match the elements with.
//
// Returns:
//   - []HelperElem: The elements that match the top of the stack.
//   - error: An error if the operation failed.
func (cs *ConflictSolver[T]) Match(stack *ud.History[lls.Stacker[*gr.Token[T]]]) ([]HelperElem[T], error) {
	var top *gr.Token[T]
	var ok bool

	stack.ReadData(func(data lls.Stacker[*gr.Token[T]]) {
		top, ok = data.Peek()
	})
	if !ok {
		return nil, errors.New("no top token found")
	}

	id := top.GetID()

	elems, ok := cs.table[id]
	if !ok {
		return nil, fmt.Errorf("no elems found for symbol %s", id)
	}

	f := func(h *HelperNode[T]) (*HelperNode[T], error) {
		cmd := lls.NewPop[*gr.Token[T]]()
		err := stack.ExecuteCommand(cmd)
		if err != nil {
			return nil, errors.New("no top token found")
		}
		top := cmd.Value()

		err = h.Match(top, stack)
		if err != nil {
			return nil, err
		}

		return h, nil
	}

	slice := make([]*HelperNode[T], len(elems))
	copy(slice, elems)

	success_or_fail, ok := us.EvaluateSimpleHelpers(slice, f)
	if !ok {
		// Return the most likely error
		// As of now, we will return the first error
		data := success_or_fail[0].GetData()
		err := data.Second

		return nil, err
	}

	success := us.ExtractResults(success_or_fail)

	// Get the longest match
	// TO DO: Implement a better way to get the longest match.
	// As of now, every match is considered the longest match.
	firsts := make([]HelperElem[T], 0, len(success))

	for _, final := range success {
		act := final.GetAction()
		firsts = append(firsts, act)
	}

	if len(success) == 1 {
		return firsts, nil
	}

	return firsts, nil
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
