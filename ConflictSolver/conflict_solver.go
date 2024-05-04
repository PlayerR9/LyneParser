package ConflictSolver

import (
	"fmt"
	"slices"
	"strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

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
	indentation := ffs.NewIndentConfig(ffs.DefaultIndentation, indentLevel, true, false)
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

// solveSetLookaheadOnShifts is a helper function that evaluates the look-ahead on shifts
// in a simple and coarse way.
func (cs *ConflictSolver) solveSetLookaheadOnShifts() {
	helpers := cs.getHelpers()

	for _, h := range helpers {
		h.EvaluateLookahead()
	}
}

/////////////////////////////////////////////////////////////

// SolveConflicts is a method that solves conflicts in a decision table.
func (cs *ConflictSolver) Solve() error {
	cs.solveSetLookaheadOnShifts()

	// Now, those shift actions that have the look-ahead are no longer
	// in conflict with their reduce counterparts.
	// However, there still might be conflicts between shift actions
	// with the same look-ahead.

	laConflicts := cs.getLookAheadConflicts()

	for _, elems := range laConflicts {
		if len(elems) != 1 {
			// Solve conflicts.
			err := solveSubgroup(elems)
			if err != nil {
				return err
			}
		}
	}

	// FIXME:

	// AMBIGUOUS GRAMMAR

	// SHIFT-REDUCE CONFLICT

	for {
		conflicts, limit := cs.FindConflicts()
		if limit == -1 {
			// No conflicts found.
			break
		}

		ok, err := cs.SolveAmbiguous(limit, conflicts)
		if err != nil {
			return err
		}

		if !ok {
			break
		}

		err = solveSubgroup(conflicts)
		if err != nil {
			return err
		}
	}

	return nil
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

// RemoveElement is a method that removes an element from the decision table.
//
// Parameters:
//   - elem: The element to remove.
func (cs *ConflictSolver) RemoveElement(elem *Helper) {
	if elem == nil {
		return
	}

	for symbol, elems := range cs.table {
		index := slices.Index(elems, elem)
		if index != -1 {
			cs.table[symbol] = slices.Delete(elems, index, index+1)
		}
	}
}

// FindConflicts is a method that finds conflicts for a specific symbol.
//
// Parameters:
//   - symbol: The symbol to find conflicts for.
//
// Returns:
//   - []*Helper: The conflicting helpers.
//   - error: An error of type *ErrNoElementsFound if no elements are found.
func (cs *ConflictSolver) FindConflicts() ([]*Helper, int) {
	for symbol, helpers := range cs.table {
		conflicts, limit := findConflictsPerSymbol(symbol, helpers)
		if limit == -1 {
			return conflicts, limit
		}
	}

	return nil, -1
}

// GetElemsWithLhs is a method that returns all elements with a specific LHS.
//
// Parameters:
//   - rhs: The RHS to find elements for.
//
// Returns:
//   - []*Helper: The elements with the specified LHS.
func (cs *ConflictSolver) GetElemsWithLhs(rhs string) []*Helper {
	return slext.SliceFilter(cs.getHelpers(), func(h *Helper) bool {
		return h.IsLhsRhs(rhs)
	})
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

		// remove the current rule from the list
		index := slices.Index(rs, c)
		if index != -1 {
			rs = slices.Delete(rs, index, index+1)
		}

		if len(rs) != 0 {
			newHelpers[c] = rs
		}
	}

	if len(newHelpers) == 0 {
		return false, nil
	}

	for c, rs := range newHelpers {
		cs.RemoveElement(c)

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

// getLookAheadConflicts is a helper function that returns all look-ahead conflicts.
//
// Returns:
//   - map[string][]*Helper: A map of look-ahead conflicts.
func (cs *ConflictSolver) getLookAheadConflicts() map[string][]*Helper {
	laConflicts := make(map[string][]*Helper)

	todo := cs.getHelpers()

	for _, h := range todo {
		if lookahead := h.GetLookahead(); lookahead != nil {
			laConflicts[*lookahead] = append(laConflicts[*lookahead], h)
		}
	}

	return laConflicts
}
