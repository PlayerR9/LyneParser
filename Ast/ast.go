package Ast

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// SyntaxChecker is a helper struct for AST checks.
type SyntaxChecker[T uc.Enumer] struct {
	// table is the table of ASTs.
	table [][]T
}

// NewSyntaxChecker creates a new ASTer. The ASTer is used to check if a list of
// children matches an expected list of IDs.
//
// Parameters:
//   - rules: The rules to create the ASTer.
//
// Returns:
//   - ASTer: The new ASTer.
//
// Example:
//
// lhs -> rhs1 | rhs2
//
//	ast := NewSyntaxChecker([]T{
//		"rhs1",
//		"rhs2",
//	})
//
//	err := ast.Check([]gr.Token{
//		&gr.NonLeafToken{
//			ID:   "lhs",
//			Data: []gr.Token{
//				&gr.LeafToken{
//					ID:   "rhs1",
//					Data: "data1",
//				},
//			},
//		},
//	})
//
//	if err != nil {
//		// Handle error.
//	}
func NewSyntaxChecker[T uc.Enumer](rules [][]T) SyntaxChecker[T] {
	table := make([][]T, 0, len(rules))

	for _, rule := range rules {
		table = append(table, rule)
	}

	sc := SyntaxChecker[T]{
		table: table,
	}
	return sc
}

// isLastOfTable is a helper function to check if the last element
// of the table is reached.
//
// Parameters:
//   - top: The top of the table.
//   - i: The current index.
//
// Returns:
//   - bool: True if the last element is reached, false otherwise.
func (ast *SyntaxChecker[T]) isLastOfTable(top, i int) bool {
	return top == 0 && i == len(ast.table)-1
}

// filterMissingFields is a helper function to filter missing fields.
//
// Parameters:
//
//   - i: The current index.
//
// Returns:
//
//   - error: The error if a field is missing.
func (ast *SyntaxChecker[T]) filterMissingFields(i int) error {
	var top int

	allMissing := make([]T, 0, len(ast.table))

	for j, row := range ast.table {
		if i < len(row) {
			ast.table[top] = row
			top++
		} else {
			allMissing = append(allMissing, row[i])

			ok := ast.isLastOfTable(top, j)
			if ok {
				values := make([]string, 0, len(allMissing))
				for _, v := range allMissing {
					values = append(values, v.String())
				}

				return NewErrMissingFields(values...)
			}
		}
	}

	ast.table = ast.table[:top]

	return nil
}

// filterWrongFields is a helper function to filter wrong fields.
//
// Parameters:
//
//   - child: The child to check.
//   - i: The current index.
//
// Returns:
//
//   - error: The error if a field does not match the expected ID.
func (ast *SyntaxChecker[T]) filterWrongFields(child *gr.Token[T], i int) error {
	var top int

	allExpected := make([]T, 0, len(ast.table))

	for j, row := range ast.table {
		ok := IsToken(child, row[i])
		if ok {
			ast.table[top] = row
			top++
		} else {
			allExpected = append(allExpected, row[i])

			ok := ast.isLastOfTable(top, j)
			if ok {
				values := make([]string, 0, len(allExpected))
				for _, v := range allExpected {
					values = append(values, v.String())
				}

				return uc.NewErrUnexpected(child.GoString(), values...)
			}
		}
	}

	ast.table = ast.table[:top]

	return nil
}

// filterTooManyFields is a helper function to filter too many fields.
//
// Parameters:
//
//   - children: The list of children to check.
//   - i: The current index.
//
// Returns:
//
//   - error: The error if there are too many fields.
func (ast *SyntaxChecker[T]) filterTooManyFields(size int) error {
	top := 0

	for j := 0; j < len(ast.table); j++ {
		if size == len(ast.table[j]) {
			ast.table[top] = ast.table[j]
			top++
		} else {
			ok := ast.isLastOfTable(top, j)
			if ok {
				// DEBUG: Print the table
				// for _, row := range ast.table {
				// 	fmt.Println(row)
				// }
				// fmt.Println()

				return NewErrTooManyFields(size, len(ast.table[j]))
			}
		}
	}

	ast.table = ast.table[:top]

	return nil
}

// Check checks if a list of children matches an expected list of IDs.
//
// Errors:
//
//   - ErrMissingFields: If a field is missing.
//   - ers.ErrUnexpected: If a field does not match the expected ID.
//   - ErrTooManyFields: If there are too many fields.
//
// Parameters:
//
//   - children: The list of children to check.
//
// Returns:
//
//   - error: The error if the check fails.
func (ast *SyntaxChecker[T]) Check(children []*gr.Token[T]) error {
	// 1. Create a copy of the table.
	astCopy := ast.Copy().(*SyntaxChecker[T])

	// DEBUG: Print the table
	// for _, row := range astCopy.table {
	// 	fmt.Println(row)
	// }
	// fmt.Println()

	if len(children) == 0 {
		allExpected := make([]T, 0, len(astCopy.table))

		for _, row := range astCopy.table {
			allExpected = append(allExpected, row[0])
		}

		values := make([]string, 0, len(allExpected))
		for _, v := range allExpected {
			values = append(values, v.String())
		}

		return NewErrMissingFields(values...)
	}

	err := astCopy.filterTooManyFields(len(children))
	if err != nil {
		return err
	}

	for i := 0; i < len(children); i++ {
		err = astCopy.filterMissingFields(i)
		if err != nil {
			return err
		}

		err = astCopy.filterWrongFields(children[i], i)
		if err != nil {
			return err
		}
	}

	if len(astCopy.table) != 1 {
		return NewErrAmbiguousGrammar()
	}

	return nil
}

// Copy implements the common.Copier interface.
func (ast *SyntaxChecker[T]) Copy() uc.Copier {
	table := make([][]T, 0, len(ast.table))

	for _, row := range ast.table {
		rowCopy := make([]T, len(row))
		copy(rowCopy, row)

		table = append(table, rowCopy)
	}

	scCopy := &SyntaxChecker[T]{
		table: table,
	}

	return scCopy
}
