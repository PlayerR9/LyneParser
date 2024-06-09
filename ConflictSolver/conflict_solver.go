package ConflictSolver

import (
	"fmt"
	"slices"
	"strconv"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	ds "github.com/PlayerR9/MyGoLib/ListLike/DoubleLL"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

var GlobalDebugMode bool = true

// ConflictSolver solves conflicts in a decision table.
type ConflictSolver struct {
	// table is a map of elements in the decision table.
	table map[string][]*Helper
}

// FString returns a formatted string representation of the decision table
// multilined with a specific indentation level.
//
// Parameters:
//   - indentLevel: The level of indentation.
//
// Returns:
//   - []string: A formatted string representation of the decision table.
func (cs *ConflictSolver) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	if trav == nil {
		return nil
	}

	counter := 0

	helpers := cs.getHelpers()

	for i, h := range helpers {
		err := trav.AppendString(strconv.Itoa(counter))
		if err != nil {
			return ue.NewErrAt(i, "helper", err)
		}

		err = trav.AppendRune('.')
		if err != nil {
			return ue.NewErrAt(i, "helper", err)
		}

		if h != nil {
			err = trav.AppendRune(' ')
			if err != nil {
				return ue.NewErrAt(i, "helper", err)
			}

			err = trav.AppendString(h.String())
			if err != nil {
				return ue.NewErrAt(i, "helper", err)
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
//   - *ers.ErrInvalidParameter: If the item is nil.
func NewConflictSolver(symbols []string, rules []*gr.Production) (*ConflictSolver, error) {
	cs := &ConflictSolver{
		table: make(map[string][]*Helper),
	}

	pairs := groupProdsByRhss(symbols, rules)

	for _, pair := range pairs {
		indices := pair.getIndices()
		lastIndex := pair.Rule.Size() - 1

		for _, i := range indices {
			item, err := NewItem(pair.Rule, i)
			if err != nil {
				return cs, NewErrCannotCreateItem()
			}

			var act Actioner

			if i == lastIndex {
				act = NewActReduce(pair.Rule)
			} else {
				act = NewActShift()
			}

			h, err := NewHelper(item, act)
			if err != nil {
				return cs, err
			}

			cs.table[pair.Symbol] = append(cs.table[pair.Symbol], h)
		}
	}

	return cs, nil
}

// getHelpers is a helper function that returns all helpers in the decision table.
//
// Returns:
//   - []*Helper: All helpers in the decision table.
func (cs *ConflictSolver) getHelpers() []*Helper {
	var result []*Helper

	for _, helpers := range cs.table {
		result = append(result, helpers...)
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
func (cs *ConflictSolver) GetElemsWithLhs(rhs string) []*Helper {
	filter := func(h *Helper) bool {
		return h.IsLhsRhs(rhs)
	}

	helpers := cs.getHelpers()

	helpers = us.SliceFilter(helpers, filter)

	return helpers
}

// DeleteHelper is a method that removes an helper from the decision table.
//
// Parameters:
//   - h: The helper to remove.
func (cs *ConflictSolver) DeleteHelper(h *Helper) {
	if h == nil {
		return
	}

	for symbol, elems := range cs.table {
		index := slices.Index(elems, h)
		if index != -1 {
			cs.table[symbol] = slices.Delete(elems, index, index+1)
		}
	}
}

// solveSetLookaheadOnShifts is a helper function that evaluates the look-ahead on shifts
// in a simple and coarse way.
func (cs *ConflictSolver) solveSetLookaheadOnShifts() {
	helpers := cs.getHelpers()

	for _, h := range helpers {
		h.EvaluateLookahead()
	}
}

// getHelpersWithLookahead is a helper function that returns all helpers with look-ahead
// and groups them by the look-ahead.
//
// Returns:
//   - map[string][]*Helper: The helpers with look-ahead.
func (cs *ConflictSolver) getHelpersWithLookahead() map[string][]*Helper {
	groups := make(map[string][]*Helper)

	todo := cs.getHelpers()

	for _, h := range todo {
		lookahead := h.GetLookahead()

		if lookahead != nil {
			groups[*lookahead] = append(groups[*lookahead], h)
		}
	}

	return groups
}

// AppendHelper is a method that appends a helper to the decision table.
//
// Parameters:
//   - h: The helper to append.
func (cs *ConflictSolver) AppendHelper(h *Helper) {
	if h == nil {
		return
	}

	rhs := h.GetRhs()
	cs.table[rhs] = append(cs.table[rhs], h)
}

// SolveAmbiguousShifts is a method that solves ambiguous shifts in a decision table.
//
// Returns:
//   - error: An error if the operation failed.
//
// Errors:
//   - *ErrHelpersConflictingSize: If the helpers have conflicting sizes.
//   - *ErrHelper: If there is an error appending the right-hand side to the helper.
func (cs *ConflictSolver) SolveAmbiguousShifts() error {
	cs.solveSetLookaheadOnShifts()

	// Now, those shift actions that have the look-ahead are no longer
	// in conflict with their reduce counterparts.
	// However, there still might be conflicts between shift actions
	// with the same look-ahead.

	helpersWithLookahead := cs.getHelpersWithLookahead() // these are potential conflicts

	// If there are at least two helpers with the same look-ahead, then there might be a conflict.

	// To solve this, we have to find the minimal amount of information that is needed to
	// unambiguously determine the next action.

	for _, helpers := range helpersWithLookahead {
		if len(helpers) <= 1 {
			continue
		}

		// Solve conflicts.
		err := solveSubgroup(helpers)
		if err != nil {
			return err
		}
	}

	return nil
}

// CMPerSymbol is a type that represents conflicts per symbol.
type CMPerSymbol map[string]*uc.Pair[[]*Helper, int]

// FindConflicts is a method that finds conflicts for a specific symbol.
//
// Parameters:
//   - symbol: The symbol to find conflicts for.
//
// Returns:
//   - CMPerSymbol: The conflicts per symbol.
func (cs *ConflictSolver) FindConflicts() CMPerSymbol {
	conflictMap := make(CMPerSymbol)

	for symbol, helpers := range cs.table {
		todo := make([]*Helper, len(helpers))
		copy(todo, helpers)

		// 1. Remove every shift action that has a look-ahead.
		todo = us.SliceFilter(todo, func(h *Helper) bool {
			lookahead := h.GetLookahead()

			return lookahead == nil
		})

		conflicts, indexOfConflict := findConflictsPerSymbol(symbol, todo)
		if indexOfConflict != -1 {
			conflictMap[symbol] = uc.NewPair(conflicts, indexOfConflict)
		}
	}

	return conflictMap
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
func (cs *ConflictSolver) MakeExpansionForests(index int, nextRhs map[*Helper]string) (map[*Helper][]string, error) {
	possibleLookaheads := make(map[*Helper][]string)

	for c, rhs := range nextRhs {
		rs := cs.GetElemsWithLhs(rhs)
		if len(rs) == 0 {
			return possibleLookaheads, nil
		}

		var lookaheads []string

		for _, r := range rs {
			tree, err := NewExpansionTreeRootedAt(cs, r)
			if err != nil {
				return possibleLookaheads, NewErrHelper(c, err)
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
			possibleLookaheads[c] = lookaheads
		}
	}

	return possibleLookaheads, nil
}

func (cs *ConflictSolver) SolveAmbiguous(index int, conflicts []*Helper) (bool, error) {
	// 1. Take the next symbol of each conflicting rule
	nextRhs := make(map[*Helper]string)

	for _, c := range conflicts {
		rhs, err := c.GetRhsAt(index + 1)
		if err != nil {
			continue
		}

		nextRhs[c] = rhs
	}

	// DEBUG: Print the next RHS
	fmt.Println("Next RHS:")

	for c, rhs := range nextRhs {
		fmt.Println(c.String(), rhs)
	}
	fmt.Println()

	if len(nextRhs) == 0 {
		return false, nil
	}

	// 2. Make the expansion forests
	possibleLookaheads, err := cs.MakeExpansionForests(index, nextRhs)
	if err != nil {
		return false, err
	} else if len(possibleLookaheads) == 0 {
		return false, nil
	}

	// DEBUG: Print the forests
	fmt.Println("Possible lookaheads:")

	for c, forest := range possibleLookaheads {
		fmt.Println(c.String())
		for _, tree := range forest {
			fmt.Println("\t-", tree)
		}
		fmt.Println()
	}

	for c, forest := range possibleLookaheads {
		cs.DeleteHelper(c)

		for _, lookahead := range forest {
			newR := c.ReplaceRhsAt(index+1, lookahead)

			cs.AppendHelper(newR)
		}
	}

	return true, nil
}

/////////////////////////////////////////////////////////////

// SolveConflicts is a method that solves conflicts in a decision table.
func (cs *ConflictSolver) Solve() error {
	for {
		conflictMap := cs.FindConflicts()
		if len(conflictMap) == 0 {
			// No conflicts found.
			break
		}

		// DEBUG: Print the conflicts
		fmt.Println("Conflicts found:")

		for _, p := range conflictMap {
			conflicts := p.First
			index := p.Second

			for _, c := range conflicts {
				fmt.Println(c.String())
			}

			fmt.Println(index)
			fmt.Println()
		}
		fmt.Println()

		done := false

		for _, p := range conflictMap {
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

		for _, p := range conflictMap {
			conflicts := p.First

			// Solve conflicts.
			err := solveSubgroup(conflicts)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Init is a method that initializes the elements for a specific symbol.
// This can be used in the ExecuteSymbols method.
//
// Parameters:
//   - symbol: The symbol to initialize the elements for.
//
// Returns:
//   - error: An error if the operation failed.
//
// Errors:
//   - *ErrNoElementsFound: If no elements are found for the symbol.
func (cs *ConflictSolver) Init(symbol string) error {
	helpers, ok := cs.table[symbol]
	if !ok {
		return NewErrNoElementsFound(symbol)
	}

	for _, h := range helpers {
		h.Init(symbol)
	}

	return nil
}

// Match is a method that matches the top of the stack with the elements in the decision table.
//
// Parameters:
//   - stack: The stack to match the elements with.
//
// Returns:
//   - []Actioner: The actions to take.
//   - error: An error if the operation failed.
func (cs *ConflictSolver) Match(stack *ds.DoubleStack[gr.Tokener]) ([]Actioner, error) {
	top, err := stack.Peek()
	if err != nil {
		return nil, fmt.Errorf("no top token found")
	}

	id := top.GetID()

	elems, ok := cs.table[id]
	if !ok {
		return nil, fmt.Errorf("no elems found for symbol %s", id)
	}

	if len(elems) == 1 {
		elem := elems[0].Action

		return []Actioner{elem}, nil
	}

	f := func(h *Helper) (*Helper, error) {
		_, err := stack.Pop()
		if err != nil {
			return nil, fmt.Errorf("no top token found")
		}

		err = h.Match(top, stack)
		if err != nil {
			return nil, err
		}

		return h, nil
	}

	successOrFail, ok := us.EvaluateSimpleHelpers(elems, f)
	if !ok {
		// Return the most likely error
		// As of now, we will return the first error
		err := successOrFail[0].GetData().Second

		return nil, err
	}

	success := us.ExtractResults(successOrFail)

	// Get the longest match
	// TO DO: Implement a better way to get the longest match.
	// As of now, every match is considered the longest match.
	firsts := make([]Actioner, 0, len(success))

	for _, final := range success {
		firsts = append(firsts, final.GetAction())
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
