package ConflictSolver

import (
	"errors"

	gr "github.com/PlayerR9/LyneParser/Grammar"
)

// Pairing represents a pairing of a symbol and a production rule.
type pairing struct {
	// Symbol is the symbol of the pairing.
	Symbol string

	// Rule is the production rule of the pairing.
	Rule *gr.Production
}

// getIndices is a method that gets the indices of the right-hand side symbols of the pairing.
//
// Returns:
//   - []int: The indices of the right-hand side symbols.
func (p *pairing) getIndices() []int {
	return p.Rule.IndicesOfRhs(p.Symbol)
}

// groupProdsByRhss is a helper function that groups production rules by their right-hand side symbols.
//
// Parameters:
//   - symbols: The symbols to group by.
//   - rules: The production rules to group.
//
// Returns:
//   - pairs: The grouped production rules.
func groupProdsByRhss(symbols []string, rules []*gr.Production) (pairs []pairing) {
	for _, s := range symbols {
		for _, r := range rules {
			if r.HasRhs(s) {
				pairs = append(pairs, pairing{Symbol: s, Rule: r})
			}
		}
	}

	return
}

/////////////////////////////////////////////////////////////

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
func findConflictsPerSymbol(symbol string, helpers []*Helper) ([]*Helper, int) {
	buckets := make(map[int][]*Helper)

	for _, h := range helpers {
		indices := h.IndicesOfRhs(symbol)

		for _, index := range indices {
			buckets[index] = append(buckets[index], h)
		}
	}

	if len(buckets) == 0 {
		return nil, -1
	}

	for limit, bucket := range buckets {
		conflicts := findConflict(limit, bucket)
		if conflicts != nil {
			return conflicts, limit
		}
	}

	return nil, -1
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
