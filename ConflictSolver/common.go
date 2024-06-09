package ConflictSolver

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"

	ers "github.com/PlayerR9/MyGoLib/Units/errors"
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
func SolveConflicts(symbols []string, rules []*gr.Production) (*ConflictSolver, error) {
	if len(rules) == 0 {
		return nil, ers.NewErrInvalidParameter("rules", ers.NewErrEmpty(rules))
	}

	cs := NewConflictSolver(symbols, rules)

	err := cs.SolveAmbiguousShifts()
	if err != nil {
		return cs, err
	}

	err = cs.Solve()
	if err != nil {
		return cs, err
	}

	return cs, nil
}
