package ConflictSolver

import (
	"fmt"
	"slices"
	"strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	intf "github.com/PlayerR9/MyGoLib/Units/Common"
	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"

	tr "github.com/PlayerR9/MyGoLib/CustomData/Tree"
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
func (cs *ConflictSolver) FString(indentLevel int) []string {
	indentation := ffs.NewIndentConfig(ffs.DefaultIndentation, indentLevel, false)
	indent := indentation.String()

	result := make([]string, 0)

	counter := 0

	var builder strings.Builder

	helpers := cs.getHelpers()

	for _, h := range helpers {
		builder.WriteString(indent)
		builder.WriteString(fmt.Sprintf("%d", counter))
		builder.WriteRune('.')

		if h != nil {
			builder.WriteRune(' ')
			builder.WriteString(h.String())
		}

		result = append(result, builder.String())
		builder.Reset()

		counter++
	}

	return result
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
				act, err = NewActReduce(pair.Rule)
				if err != nil {
					return cs, err
				}
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
	result := make([]*Helper, 0)

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

	return slext.SliceFilter(cs.getHelpers(), filter)
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

// FindConflicts is a method that finds conflicts for a specific symbol.
//
// Parameters:
//   - symbol: The symbol to find conflicts for.
//
// Returns:
//   - []*Helper: The conflicting helpers.
//   - int: The index of the position of the conflict.
func (cs *ConflictSolver) FindConflicts() ([]*Helper, int) {
	for symbol, helpers := range cs.table {
		todo := make([]*Helper, len(helpers))
		copy(todo, helpers)

		// 1. Remove every shift action that has a look-ahead.
		todo = slext.SliceFilter(todo, func(h *Helper) bool {
			return h.GetLookahead() == nil
		})

		conflicts, indexOfConflict := findConflictsPerSymbol(symbol, todo)
		if indexOfConflict != -1 {
			return conflicts, indexOfConflict
		}
	}

	return nil, -1
}

type InfoStruct struct {
	seen map[*Helper]bool
}

func (is *InfoStruct) Copy() intf.Copier {
	isCopy := &InfoStruct{
		seen: make(map[*Helper]bool),
	}

	for k, v := range is.seen {
		isCopy.seen[k] = v
	}

	return isCopy
}

func (cs *ConflictSolver) GenerateTreeRootedAt(h *Helper) (*tr.Tree[*Helper], error) {
	tree, err := tr.MakeTree(h, &InfoStruct{
		seen: make(map[*Helper]bool),
	}, func(elem *Helper, is *InfoStruct) ([]*Helper, error) {
		rhs, err := elem.GetRhsAt(0)
		if err != nil {
			return nil, NewErr0thRhsNotSet()
		}

		seenFilter := func(h *Helper) bool {
			return !is.seen[h]
		}

		return slext.SliceFilter(cs.GetElemsWithLhs(rhs), seenFilter), nil
	})
	if err != nil {
		return nil, err
	}

	return tree, nil
}

func (cs *ConflictSolver) CheckIfLookahead0(index int, h *Helper) ([]*Helper, error) {
	// 1. Take the next symbol of h
	rhs, err := h.GetRhsAt(index + 1)
	if err != nil {
		return nil, NewErrHelper(h, err)
	}

	// 2. Get all the helpers that have the same LHS as rhs
	newHelpers := cs.GetElemsWithLhs(rhs)
	if len(newHelpers) == 0 {
		return nil, nil
	}

	// 3. For each rule, check if the 0th rhs is a terminal symbol
	solutions := make([]*Helper, 0)

	for _, nh := range newHelpers {
		otherRhs, err := nh.GetRhsAt(0)
		if err != nil {
			return solutions, NewErrHelper(nh, err)
		}

		if gr.IsTerminal(otherRhs) {
			solutions = append(solutions, nh)
		} else {

		}
	}

	return solutions, nil
}

func (cs *ConflictSolver) SolveAmbiguous(index int, conflicts []*Helper) (bool, error) {
	newHelpers := make(map[*Helper][]*Helper)

	for _, c := range conflicts {
		// 1. Take the next symbol of each conflicting rule
		rhs, err := c.GetRhsAt(index + 1)
		if err != nil {
			continue
		}

		// 2. Replace the current symbol with every rule

		rs := cs.GetElemsWithLhs(rhs)

		/*
			// remove the current rule from the list ???
			index := slices.Index(rs, c)
			if index != -1 {
				rs = slices.Delete(rs, index, index+1)
			}
		*/

		if len(rs) != 0 {
			newHelpers[c] = rs
		}
	}

	if len(newHelpers) == 0 {
		return false, nil
	}

	for c, rs := range newHelpers {
		cs.DeleteHelper(c)

		for _, r := range rs {
			newR, err := c.ReplaceRhsAt(index+1, r)
			if err != nil {
				return false, NewErrHelper(c, err)
			}

			cs.AppendHelper(newR)
		}
	}

	return true, nil
}

/////////////////////////////////////////////////////////////

// SolveConflicts is a method that solves conflicts in a decision table.
func (cs *ConflictSolver) Solve() error {
	for {
		conflicts, limit := cs.FindConflicts()
		if limit == -1 {
			// No conflicts found.
			break
		}

		fmt.Println("Conflicts found:")
		for _, c := range conflicts {
			fmt.Println(c.String())
		}
		fmt.Println()

		ok, err := cs.SolveAmbiguous(limit, conflicts)
		if err != nil {
			return err
		}

		if !ok {
			break
		}

		if GlobalDebugMode {
			return nil
		}

		err = solveSubgroup(conflicts)
		if err != nil {
			return err
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
