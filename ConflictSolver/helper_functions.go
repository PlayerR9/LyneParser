package ConflictSolver

// getRhssAt is a helper function that groups helpers by their right-hand side
// symbols at a specific index.
//
// Parameters:
//   - helpers: The helpers to group.
//   - index: The index to group by.
//
// Returns:
//   - map[string][]*Helper: The grouped helpers.
//   - error: An error of type *ErrHelpersConflictingSize if the helpers have conflicting sizes.
func getRhssAt(bucket []*Helper, index int) (map[string][]*Helper, error) {
	groups := make(map[string][]*Helper)

	for _, h := range bucket {
		rhs, err := h.GetRhsAt(index)
		if err != nil {
			return nil, NewErrHelpersConflictingSize()
		}

		groups[rhs] = append(groups[rhs], h)
	}

	return groups, nil
}

// minimumUnique is a helper function that, given a set of helpers and the limit
// (i.e., the position of the shared symbol for all helpers), solves conflicts
// by finding the least number of rhs symbols under the limit that are unique
// to each helper.
//
// For example:
//
//	Helper 1: A -> B C [D]
//	Helper 2: A -> A C [D]
//	Helper 3: A -> E F [D]
//
//	Limit: 2 (i.e., the position of the shared symbol [D])
//
//	Then, to distinguish Helper 1 and Helper 2, we need up to (B) for Helper 1,
//	and up to (A) for Helper 2, but (F) for Helper 3.
//
//	Helper 1: A -> B C [D]
//	Helper 2: A -> A C [D]
//	Helper 3: A -> F [D]
//
// In this way, we optimize the numbers of checks needed to make an informed decision.
//
// Parameters:
//   - helpers: The helpers to solve conflicts for.
//   - limit: The index of the shared symbol for all helpers.
//
// Returns:
//   - error: An error if the operation failed.
//
// Errors:
//   - *ErrHelpersConflictingSize: If the helpers have conflicting sizes.
//   - *ErrHelper: If there is an error appending the right-hand side to the helper.
func minimumUnique(bucket []*Helper, limit int) error {
	todo := make(map[*Helper]bool)

	for _, h := range bucket {
		todo[h] = true
	}

	for i := limit; i >= 0; i-- {
		rhsPerLevel, err := getRhssAt(bucket, i)
		if err != nil {
			return NewErrHelpersConflictingSize()
		}

		for rhs, helpers := range rhsPerLevel {
			// Add the rhs to the helpers.
			for _, h := range helpers {
				err := h.AppendRhs(rhs)
				if err != nil {
					return NewErrHelper(h, err)
				}
			}

			// However, if there is only one helper, then there is no conflict.
			// Therefore, remove it from the todo list.
			if len(helpers) == 1 {
				delete(todo, helpers[0])
			}
		}
	}

	return nil
}

// solveSubgroup is a helper function that solves conflicts between a subgroup of helpers.
//
// Parameters:
//   - helpers: The helpers to solve conflicts for.
//
// Returns:
//   - error: An error if the operation failed.
//
// Errors:
//   - *ErrHelpersConflictingSize: If the helpers have conflicting sizes.
//   - *ErrHelper: If there is an error appending the right-hand side to the helper.
func solveSubgroup(bucket []*Helper) error {
	// 1. Bucket sort the items by their position.

	buckets := make(map[int][]*Helper)

	for _, h := range bucket {
		pos := h.GetPos()

		prev, ok := buckets[pos]
		if !ok {
			prev = []*Helper{h}
		} else {
			prev = append(prev, h)
		}

		buckets[pos] = prev
	}

	// 2. Solve conflicts for each bucket.
	for limit, bucket := range buckets {
		err := minimumUnique(bucket, limit)
		if err != nil {
			return err
		}
	}

	return nil
}

// findConflict is a helper function that finds conflicts between a subgroup of helpers.
//
// Parameters:
//   - limit: The index of the shared symbol for all helpers.
//   - helpers: The helpers to find conflicts for.
//
// Returns:
//   - []*Helper: The conflicting helpers.
func findConflict(limit int, bucket []*Helper) []*Helper {
	if len(bucket) < 2 {
		return nil
	}

	// 1. Fill the matrix of symbols.
	matrixOfSymbols := make([][]*string, limit+1)
	for i := range matrixOfSymbols {
		matrixOfSymbols[i] = make([]*string, len(bucket))
	}

	for i, h := range bucket {
		for l := 0; l <= limit; l++ {
			rhs, err := h.GetRhsAt(l)
			if err != nil {
				matrixOfSymbols[l][i] = nil
			} else {
				matrixOfSymbols[l][i] = &rhs
			}
		}
	}

	// 2. Create a conflict matrix.
	conflictMatrix := make([][]int, 0, len(bucket))

	for range bucket {
		row := make([]int, len(bucket))
		conflictMatrix = append(conflictMatrix, row)
	}

	// 3. Evaluate the matrix of conflicts.
	for _, row := range matrixOfSymbols {
		for i := 0; i < len(bucket); i++ {
			for j := 0; j < len(bucket); j++ {
				if *row[i] == *row[j] {
					conflictMatrix[i][j]++
				}
			}
		}
	}

	// 4. Find conflicts.
	countMap := make(map[int][]*Helper)

	for j, h1 := range bucket {
		count := 0

		for i := 0; i < len(bucket); i++ {
			count += conflictMatrix[i][j]
		}

		prev, ok := countMap[count]
		if !ok {
			prev = []*Helper{h1}
		} else {
			prev = append(prev, h1)
		}

		countMap[count] = prev
	}

	for count, helpers := range countMap {
		if count > 1 {
			return helpers
		}
	}

	return nil
}

// findConflictsPerSymbol is a helper function that finds conflicts for a specific symbol.
//
// Parameters:
//   - symbol: The symbol to find conflicts for.
//   - helpers: The helpers to find conflicts for.
//
// Returns:
//   - []*Helper: The conflicting helpers.
//   - int: The position of the conflict.
func findConflictsPerSymbol(symbol string, bucket []*Helper) ([]*Helper, int) {
	if len(bucket) < 2 {
		return nil, -1
	}

	buckets := make(map[int][]*Helper)

	for _, h := range bucket {
		indices := h.IndicesOfRhs(symbol)

		for _, index := range indices {
			prev, ok := buckets[index]
			if !ok {
				prev = []*Helper{h}
			} else {
				prev = append(prev, h)
			}

			buckets[index] = prev
		}
	}

	for limit, bucket := range buckets {
		conflicts := findConflict(limit, bucket)
		if conflicts != nil {
			return conflicts, limit
		}
	}

	return nil, -1
}
