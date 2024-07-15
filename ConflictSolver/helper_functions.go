package ConflictSolver

import gr "github.com/PlayerR9/LyneParser/Grammar"

// get_rhss_at is a helper function that groups helpers by their right-hand side
// symbols at a specific index.
//
// Parameters:
//   - helpers: The helpers to group.
//   - index: The index to group by.
//
// Returns:
//   - map[T][]*Helper: The grouped helpers.
//   - error: An error of type *ErrHelpersConflictingSize if the helpers have conflicting sizes.
func get_rhss_at[T gr.TokenTyper](bucket []*HelperNode[T], index int) (map[T][]*HelperNode[T], error) {
	groups := make(map[T][]*HelperNode[T])

	for _, h := range bucket {
		rhs, err := h.GetRhsAt(index)
		if err != nil {
			return nil, NewErrHelpersConflictingSize()
		}

		groups[rhs] = append(groups[rhs], h)
	}

	return groups, nil
}

// minimum_unique is a helper function that, given a set of helpers and the limit
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
func minimum_unique[T gr.TokenTyper](bucket []*HelperNode[T], limit int) error {
	todo := make(map[*HelperNode[T]]bool)

	for _, h := range bucket {
		todo[h] = true
	}

	for i := limit; i >= 0; i-- {
		rhs_per_level, err := get_rhss_at(bucket, i)
		if err != nil {
			return NewErrHelpersConflictingSize()
		}

		for rhs, helpers := range rhs_per_level {
			// Add the rhs to the helpers.
			for _, h := range helpers {
				h.AppendRhs(rhs)
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

// solve_subgroup is a helper function that solves conflicts between a subgroup of helpers.
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
func solve_subgroup[T gr.TokenTyper](bucket []*HelperNode[T]) error {
	// 1. Bucket sort the items by their position.

	buckets := make(map[int][]*HelperNode[T])

	for _, h := range bucket {
		pos := h.GetPos()

		prev, ok := buckets[pos]
		if !ok {
			prev = []*HelperNode[T]{h}
		} else {
			prev = append(prev, h)
		}

		buckets[pos] = prev
	}

	// 2. Solve conflicts for each bucket.
	for limit, bucket := range buckets {
		err := minimum_unique(bucket, limit)
		if err != nil {
			return err
		}
	}

	return nil
}

// find_conflict is a helper function that finds conflicts between a subgroup of helpers.
//
// Parameters:
//   - limit: The index of the shared symbol for all helpers.
//   - helpers: The helpers to find conflicts for.
//
// Returns:
//   - []*Helper: The conflicting helpers.
func find_conflict[T gr.TokenTyper](limit int, bucket []*HelperNode[T]) []*HelperNode[T] {
	if len(bucket) < 2 {
		return nil
	}

	// 1. Fill the matrix of symbols.
	matrix_of_symbols := make([][]*T, limit+1)
	for i := range matrix_of_symbols {
		matrix_of_symbols[i] = make([]*T, len(bucket))
	}

	for i, h := range bucket {
		for l := 0; l <= limit; l++ {
			rhs, err := h.GetRhsAt(l)
			if err != nil {
				matrix_of_symbols[l][i] = nil
			} else {
				matrix_of_symbols[l][i] = &rhs
			}
		}
	}

	// 2. Create a conflict matrix.
	conflict_matrix := make([][]int, 0, len(bucket))

	for range bucket {
		row := make([]int, len(bucket))
		conflict_matrix = append(conflict_matrix, row)
	}

	// 3. Evaluate the matrix of conflicts.
	for _, row := range matrix_of_symbols {
		for i := 0; i < len(bucket); i++ {
			for j := 0; j < len(bucket); j++ {
				if *row[i] == *row[j] {
					conflict_matrix[i][j]++
				}
			}
		}
	}

	// 4. Find conflicts.
	count_map := make(map[int][]*HelperNode[T])

	for j, h1 := range bucket {
		count := 0

		for i := 0; i < len(bucket); i++ {
			count += conflict_matrix[i][j]
		}

		prev, ok := count_map[count]
		if !ok {
			prev = []*HelperNode[T]{h1}
		} else {
			prev = append(prev, h1)
		}

		count_map[count] = prev
	}

	for count, helpers := range count_map {
		if count > 1 {
			return helpers
		}
	}

	return nil
}

// find_conflicts_per_symbol is a helper function that finds conflicts for a specific symbol.
//
// Parameters:
//   - symbol: The symbol to find conflicts for.
//   - helpers: The helpers to find conflicts for.
//
// Returns:
//   - []*Helper: The conflicting helpers.
//   - int: The position of the conflict.
func find_conflicts_per_symbol[T gr.TokenTyper](symbol T, bucket []*HelperNode[T]) ([]*HelperNode[T], int) {
	if len(bucket) < 2 {
		return nil, -1
	}

	buckets := make(map[int][]*HelperNode[T])

	for _, h := range bucket {
		indices := h.IndicesOfRhs(symbol)

		for _, index := range indices {
			prev, ok := buckets[index]
			if !ok {
				prev = []*HelperNode[T]{h}
			} else {
				prev = append(prev, h)
			}

			buckets[index] = prev
		}
	}

	for limit, bucket := range buckets {
		conflicts := find_conflict(limit, bucket)
		if conflicts != nil {
			return conflicts, limit
		}
	}

	return nil, -1
}
