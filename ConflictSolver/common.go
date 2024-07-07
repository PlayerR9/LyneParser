package ConflictSolver

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
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
func SolveConflicts[T gr.TokenTyper](symbols []T, rules []*gr.Production[T]) (*ConflictSolver[T], error) {
	if len(rules) == 0 {
		return nil, uc.NewErrInvalidParameter("rules", uc.NewErrEmpty(rules))
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
