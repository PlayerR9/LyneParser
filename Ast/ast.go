package Ast

import (
	"fmt"
	"strings"

	ers "github.com/PlayerR9/MyGoLibUnits/Errors"

	gr "github.com/PlayerR9/LyneParser/Grammar"
)

// ASTNoder is an interface for AST nodes.
type ASTNoder interface {
	// ToAST converts the parser node to an AST node.
	//
	// Parameters:
	//
	//   - root: The root to convert.
	//
	// Returns:
	//
	//   - error: The error if the conversion fails.
	ToAST(root gr.Tokener) error

	fmt.Stringer
}

// IsToken checks if a root is a token with a specific ID.
//
// Parameters:
//
//   - root: The root to check.
//   - id: The ID to check.
//
// Returns:
//
//   - bool: True if the root is a token with the ID, false otherwise.
//   - error: The error if the root is not a token or the ID does not match.
func IsToken(root gr.Tokener, id string) (bool, error) {
	if root == nil {
		return false, NewErrExpectedNonNil("root")
	} else if root.GetID() != id {
		return false, nil
	}

	var ok bool

	if gr.IsTerminal(id) {
		_, ok = root.(*gr.LeafToken)
	} else {
		_, ok = root.(*gr.NonLeafToken)
	}

	if !ok {
		return false, NewErrInvalidParsing(root)
	} else {
		return true, nil
	}
}

// ASTer is a helper struct for AST checks.
type ASTer struct {
	// table is the table of ASTs.
	table [][]string
}

// NewASTer creates a new ASTer.
//
// Parameters:
//
//   - rules: The rules to create the ASTer.
//
// Returns:
//
//   - ASTer: The new ASTer.
func NewASTer(rules []string) ASTer {
	ast := ASTer{
		table: make([][]string, 0, len(rules)),
	}

	for _, rule := range rules {
		ast.table = append(ast.table, strings.Fields(rule))
	}

	return ast
}

// isLastOfTable is a helper function to check if the last element
// of the table is reached.
//
// Parameters:
//
//   - top: The top of the table.
//   - i: The current index.
//
// Returns:
//
//   - bool: True if the last element is reached, false otherwise.
func (ast *ASTer) isLastOfTable(top, i int) bool {
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
func (ast *ASTer) filterMissingFields(i int) error {
	top := 0

	allMissing := make([]string, 0, len(ast.table))

	for j, row := range ast.table {
		if i < len(row) {
			ast.table[top] = row
			top++
		} else if ast.isLastOfTable(top, j) {
			allMissing = append(allMissing, row[i])

			return NewErrMissingFields(allMissing...)
		} else {
			allMissing = append(allMissing, row[i])
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
func (ast *ASTer) filterWrongFields(child gr.Tokener, i int) error {
	top := 0

	allExpected := make([]string, 0, len(ast.table))

	for j, row := range ast.table {
		ok, err := IsToken(child, row[i])
		if err == nil && ok {
			ast.table[top] = row
			top++
		} else if ast.isLastOfTable(top, j) {
			allExpected = append(allExpected, row[i])

			return ers.NewErrUnexpected(child, allExpected...)
		} else {
			allExpected = append(allExpected, row[i])
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
func (ast *ASTer) filterTooManyFields(children []gr.Tokener) error {
	top := 0

	for j := 0; j < len(ast.table); j++ {
		if len(children) <= len(ast.table[j]) {
			ast.table[top] = ast.table[j]
			top++
		} else if ast.isLastOfTable(top, j) {
			return NewErrTooManyFields(len(children), len(ast.table[j]))
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
func (ast *ASTer) Check(children []gr.Tokener) error {
	if len(children) == 0 {
		allExpected := make([]string, 0, len(ast.table))

		for _, row := range ast.table {
			allExpected = append(allExpected, row[0])
		}

		return NewErrMissingFields(allExpected...)
	}

	var err error

	for i := 0; i < len(children); i++ {
		err = ast.filterMissingFields(i)
		if err != nil {
			return err
		}

		err = ast.filterWrongFields(children[i], i)
		if err != nil {
			return err
		}

		err = ast.filterTooManyFields(children)
		if err != nil {
			return err
		}
	}

	if len(ast.table) != 1 {
		return NewErrAmbiguousGrammar()
	}

	return nil
}
