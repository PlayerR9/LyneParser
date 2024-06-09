package ConflictSolver

import (
	uts "github.com/PlayerR9/MyGoLib/Utility/Sorting"
)

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
func getRhssAt(bucket *uts.Bucket[*Helper], index int) (map[string][]*Helper, error) {
	groups := make(map[string][]*Helper)

	iter := bucket.Iterator()

	for {
		h, err := iter.Consume()
		if err != nil {
			break
		}

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
func minimumUnique(bucket *uts.Bucket[*Helper], limit int) error {
	todo := make(map[*Helper]bool)

	iter := bucket.Iterator()

	for {
		h, err := iter.Consume()
		if err != nil {
			break
		}

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
func solveSubgroup(bucket *uts.Bucket[*Helper]) error {
	// 1. Bucket sort the items by their position.

	// TODO: Modify this once MyGoLib is updated.
	buckets := make(map[int]*uts.Bucket[*Helper])

	iter := bucket.Iterator()

	for {
		h, err := iter.Consume()
		if err != nil {
			break
		}

		pos := h.GetPos()

		_, ok := buckets[pos]
		if !ok {
			buckets[pos] = uts.NewBucket([]*Helper{h})
		} else {
			buckets[pos].Add(h)
		}
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
func findConflict(limit int, bucket *uts.Bucket[*Helper]) *uts.Bucket[*Helper] {
	if bucket.GetSize() < 2 {
		return nil
	}

	// 1. Fill the matrix of symbols.
	matrixOfSymbols := make([][]*string, limit+1)
	for i := range matrixOfSymbols {
		matrixOfSymbols[i] = make([]*string, bucket.GetSize())
	}

	iter := bucket.Iterator()

	for i := 0; ; i++ {
		h, err := iter.Consume()
		if err != nil {
			break
		}

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
	conflictMatrix := make([][]int, 0, bucket.GetSize())

	iter = bucket.Iterator()

	for {
		_, err := iter.Consume()
		if err != nil {
			break
		}

		conflictMatrix = append(conflictMatrix, make([]int, bucket.GetSize()))
	}

	// 3. Evaluate the matrix of conflicts.
	for _, row := range matrixOfSymbols {
		for i := 0; i < bucket.GetSize(); i++ {
			for j := 0; j < bucket.GetSize(); j++ {
				if *row[i] == *row[j] {
					conflictMatrix[i][j]++
				}
			}
		}
	}

	// 4. Find conflicts.
	countMap := make(map[int]*uts.Bucket[*Helper])

	columnIter := bucket.Iterator()

	for j := 0; ; j++ {
		h1, err := columnIter.Consume()
		if err != nil {
			break
		}

		count := 0

		for i := 0; i < bucket.GetSize(); i++ {
			count += conflictMatrix[i][j]
		}

		prev, ok := countMap[count]
		if !ok {
			countMap[count] = uts.NewBucket([]*Helper{h1})
		} else {
			prev.Add(h1)
		}
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
func findConflictsPerSymbol(symbol string, bucket *uts.Bucket[*Helper]) (*uts.Bucket[*Helper], int) {
	if bucket.GetSize() < 2 {
		return nil, -1
	}

	buckets := make(map[int]*uts.Bucket[*Helper])

	iter := bucket.Iterator()

	for {
		h, err := iter.Consume()
		if err != nil {
			break
		}

		indices := h.IndicesOfRhs(symbol)

		for _, index := range indices {
			prev, ok := buckets[index]
			if !ok {
				buckets[index] = uts.NewBucket([]*Helper{h})
			} else {
				prev.Add(h)
			}
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
