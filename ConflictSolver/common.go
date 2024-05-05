package ConflictSolver

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// SolveConflicts solves conflicts in a decision table.
//
// Parameters:
//   - symbols: The symbols in the decision table.
//   - rules: The rules in the decision table.
//
// Returns:
//   - map[string][]*Helper: The elements in the decision table with conflicts solved.
//   - error: An error if the operation failed.
func SolveConflicts(symbols []string, rules []*gr.Production) (map[string][]*Helper, error) {
	if len(rules) == 0 {
		return nil, ers.NewErrInvalidParameter("rules", ers.NewErrEmptySlice())
	}

	cs, err := NewConflictSolver(symbols, rules)
	if err != nil {
		return nil, err
	}

	err = cs.SolveAmbiguousShifts()
	if err != nil {
		return cs.table, err
	}

	err = cs.Solve()
	if err != nil {
		return cs.table, err
	}

	return cs.table, nil
}
