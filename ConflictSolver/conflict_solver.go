package ConflictSolver

import (
	"errors"
	"slices"

	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

// ConflictSolver solves conflicts in a decision table.
type ConflictSolver struct {
	// elems is a map of elements in the decision table.
	elems map[string][]*Helper
}

// NewConflictSolver creates a new conflict solver.
//
// Parameters:
//   - items: The items in the decision table.
//
// Returns:
//   - *ConflictSolver: A pointer to the new conflict solver.
func NewConflictSolver(elems map[string][]*Helper) *ConflictSolver {
	for _, helpers := range elems {
		for _, h := range helpers {
			h.SetAction(nil)
		}
	}

	return &ConflictSolver{
		elems: elems,
	}
}

// AppendHelper is a method that appends a helper to the decision table.
//
// Parameters:
//   - h: The helper to append.
//
// Returns:
//   - error: An error of type *ErrItemIsNil if the helper is nil.
func (cs *ConflictSolver) AppendHelper(h *Helper) error {
	if h == nil {
		return nil
	}

	rhs, err := h.GetRhs()
	if err != nil {
		return err
	}

	cs.elems[rhs] = append(cs.elems[rhs], h)

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
//   - *ErrHelper: If something goes wrong with a helper.
func (cs *ConflictSolver) Init(symbol string) error {
	helpers, ok := cs.elems[symbol]
	if !ok {
		return NewErrNoElementsFound(symbol)
	}

	for _, h := range helpers {
		err := h.Init(symbol)
		if err != nil {
			return NewErrHelper(h, err)
		}
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

	for symbol, elems := range cs.elems {
		index := slices.Index(elems, elem)
		if index != -1 {
			cs.elems[symbol] = slices.Delete(elems, index, index+1)
		}
	}
}

// findConflict is a helper function that finds conflicts between a subgroup of helpers.
//
// Parameters:
//   - limit: The limit of the subgroup.
//   - helpers: The helpers to find conflicts for.
//
// Returns:
//   - []*Helper: The conflicting helpers.
func findConflict(limit int, helpers []*Helper) []*Helper {
	// TO DO: Check if this is correct.

	if len(helpers) <= 1 {
		return nil
	}

	// 1. Fill the matrix of symbols.
	matrixOfSymbols := make([][]string, limit)

	for i, h := range helpers {
		symbols := h.GetSymbolsUpToPos()

		for l, symbol := range symbols {
			matrixOfSymbols[l][i] = symbol
		}
	}

	// 2. Create a conflict matrix.
	conflictMatrix := make([][]int, len(helpers))

	for i := 0; i < len(helpers); i++ {
		conflictMatrix[i] = make([]int, len(helpers))
	}

	// 3. Evaluate the matrix of conflicts.
	for _, row := range matrixOfSymbols {
		for i := 0; i < len(helpers); i++ {
			for j := 0; j < len(helpers); j++ {
				if row[i] == row[j] {
					conflictMatrix[i][j]++
				}
			}
		}
	}

	// 4. Find conflicts.
	countMap := make(map[int][]*Helper)

	for j := 0; j < len(helpers); j++ {
		count := 0

		for i := 0; i < len(helpers); i++ {
			count += conflictMatrix[i][j]
		}

		countMap[count] = append(countMap[count], helpers[j])
	}

	for count, helpers := range countMap {
		if count > 1 {
			return helpers
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
//   - error: An error of type *ErrNoElementsFound if no elements are found.
func (cs *ConflictSolver) FindConflicts(symbol string) ([]*Helper, error) {
	// 1. Find all helpers for the symbol.
	helpers, ok := cs.elems[symbol]
	if !ok {
		return nil, NewErrNoElementsFound(symbol)
	}

	// 2. Make a bucket for each position of the RHS.
	buckets := make(map[int][]*Helper)

	for _, h := range helpers {
		indices := h.IndicesOfRhs(symbol)

		for _, index := range indices {
			buckets[index] = append(buckets[index], h)
		}
	}

	if len(buckets) == 0 {
		return nil, nil
	}

	for limit, bucket := range buckets {
		conflicts := findConflict(limit, bucket)
		if conflicts != nil {
			return conflicts, nil
		}
	}

	return nil, nil
}

// GetElemsWithLhs is a method that returns all elements with a specific LHS.
//
// Parameters:
//   - rhs: The RHS to find elements for.
//
// Returns:
//   - []*Helper: The elements with the specified LHS.
func (cs *ConflictSolver) GetElemsWithLhs(rhs string) []*Helper {
	return slext.SliceFilter(cs.GetAllHelpers(), func(h *Helper) bool {
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

			err = cs.AppendHelper(newR)
			if err != nil {
				return false, NewErrHelper(newR, err)
			}
		}
	}

	return true, nil
}

// minimumUnique is a helper function that solves conflicts between a subgroup of helpers.
//
// Parameters:
//   - helpers: The helpers to solve conflicts for.
//   - limit: The limit of the subgroup.
func minimumUnique(helpers []*Helper, limit int) error {
	// Set all helpers to not done.
	elemsEval := make(map[*Helper]bool)

	for _, h := range helpers {
		elemsEval[h] = false
	}

	for i := limit; i >= 0; i-- {
		rhsPerLevel := make(map[string][]int)

		for j, h := range helpers {
			isDone, ok := elemsEval[h]
			if !ok {
				return NewErrHelper(h, errors.New("item not found in doneMap"))
			} else if isDone {
				continue
			}

			rhs, err := h.GetRhsAt(i)
			if err != nil {
				return NewErrHelper(h, err)
			}

			rhsPerLevel[rhs] = append(rhsPerLevel[rhs], j)
		}

		for rhs, indices := range rhsPerLevel {
			if len(indices) == 1 {
				// No conflict. Mark it as done.
				elemsEval[helpers[indices[0]]] = true
			}

			for _, index := range indices {
				currentH := helpers[index]

				// Add the RHS.
				err := currentH.Action.AppendRhs(rhs)
				if err != nil {
					return NewErrHelper(currentH, err)
				}
			}
		}
	}

	return nil
}

// solveSubgroup is a helper function that solves conflicts between a subgroup of helpers.
//
// Parameters:
//   - helpers: The helpers to solve conflicts for.
func solveSubgroup(helpers []*Helper) error {
	// 1. Bucket sort the items by their position.
	buckets := make(map[int][]*Helper)

	for _, h := range helpers {
		pos := h.GetPos()
		if pos == -1 {
			return NewErrHelper(h, errors.New("pos is -1"))
		}

		buckets[pos] = append(buckets[pos], h)
	}

	for limit, bucket := range buckets {
		err := minimumUnique(bucket, limit)
		if err != nil {
			return err
		}
	}

	// Now, find conflicts between buckets.
	// FIXME: Solve conflicts between buckets.

	return nil
}

// getLookAheadConflicts is a helper function that returns all look-ahead conflicts.
//
// Returns:
//   - map[string][]*Helper: A map of look-ahead conflicts.
func (cs *ConflictSolver) getLookAheadConflicts() map[string][]*Helper {
	laConflicts := make(map[string][]*Helper)

	todo := cs.GetAllHelpers()

	for _, h := range todo {
		if lookahead := h.GetLookahead(); lookahead != nil {
			laConflicts[*lookahead] = append(laConflicts[*lookahead], h)
		}
	}

	return laConflicts
}

// GetAllHelpers is a method that returns all helpers in the decision table.
//
// Returns:
//   - []*Helper: All helpers in the decision table.
func (cs *ConflictSolver) GetAllHelpers() []*Helper {
	result := make([]*Helper, 0)

	for _, helpers := range cs.elems {
		result = append(result, helpers...)
	}

	return result
}

// SolveConflictLookahead is a method that solves look-ahead conflicts in a decision table.
//
// To use in the Execute method.
//
// Returns:
//   - error: An error of type *ErrHelper if something goes wrong with a helper.
func SolveConflictLookahead(h *Helper) error {
	return h.SetLookahead()
}

// SolveConflicts is a method that solves conflicts in a decision table.
func (cs *ConflictSolver) SolveConflicts() error {
	err := cs.Execute(SolveConflictLookahead)
	if err != nil {
		return err
	}

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

	return nil
}

// ExecuteFunc is a function that executes a function on a helper.
//
// Parameters:
//   - h: The helper to execute the function on.
//
// Returns:
//   - error: An error if the operation failed.
type ExecuteFunc func(h *Helper) error

// Execute is a method that executes a function on all helpers in the decision table.
//
// Parameters:
//   - f: The function to execute.
//
// Returns:
//   - error: An error if the operation failed.
func (cs *ConflictSolver) Execute(f ExecuteFunc) error {
	todo := cs.GetAllHelpers()

	for _, h := range todo {
		err := f(h)
		if err != nil {
			return NewErrHelper(h, err)
		}
	}

	return nil
}
