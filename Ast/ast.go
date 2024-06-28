package Ast

import (
	"strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// SyntaxChecker is a helper struct for AST checks.
type SyntaxChecker struct {
	// table is the table of ASTs.
	table [][]string
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
//	ast := NewSyntaxChecker([]string{
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
func NewSyntaxChecker(rules []string) SyntaxChecker {
	table := make([][]string, 0, len(rules))

	for _, rule := range rules {
		rows := strings.Fields(rule)
		table = append(table, rows)
	}

	sc := SyntaxChecker{
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
func (ast *SyntaxChecker) isLastOfTable(top, i int) bool {
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
func (ast *SyntaxChecker) filterMissingFields(i int) error {
	var top int

	allMissing := make([]string, 0, len(ast.table))

	for j, row := range ast.table {
		if i < len(row) {
			ast.table[top] = row
			top++
		} else {
			allMissing = append(allMissing, row[i])

			ok := ast.isLastOfTable(top, j)
			if ok {
				return NewErrMissingFields(allMissing...)
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
func (ast *SyntaxChecker) filterWrongFields(child gr.Token, i int) error {
	var top int

	allExpected := make([]string, 0, len(ast.table))

	for j, row := range ast.table {
		ok := IsToken(child, row[i])
		if ok {
			ast.table[top] = row
			top++
		} else {
			allExpected = append(allExpected, row[i])

			ok := ast.isLastOfTable(top, j)
			if ok {
				return uc.NewErrUnexpected(child.GoString(), allExpected...)
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
func (ast *SyntaxChecker) filterTooManyFields(size int) error {
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
func (ast *SyntaxChecker) Check(children []gr.Token) error {
	// 1. Create a copy of the table.
	astCopy := ast.Copy()

	// DEBUG: Print the table
	// for _, row := range astCopy.table {
	// 	fmt.Println(row)
	// }
	// fmt.Println()

	if len(children) == 0 {
		allExpected := make([]string, 0, len(astCopy.table))

		for _, row := range astCopy.table {
			allExpected = append(allExpected, row[0])
		}

		return NewErrMissingFields(allExpected...)
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

// Copy creates a copy of the SyntaxChecker.
//
// Returns:
//   - SyntaxChecker: A copy of the SyntaxChecker.
func (ast *SyntaxChecker) Copy() SyntaxChecker {
	table := make([][]string, 0, len(ast.table))

	for _, row := range ast.table {
		rowCopy := make([]string, len(row))
		copy(rowCopy, row)

		table = append(table, rowCopy)
	}

	scCopy := SyntaxChecker{
		table: table,
	}

	return scCopy
}
