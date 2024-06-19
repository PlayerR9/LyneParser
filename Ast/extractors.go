package Ast

import (
	"fmt"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

// ExtractString extracts a string from a token.
//
// Parameters:
//   - root: The root to extract the string from.
//   - id: The ID of the token.
//
// Returns:
//   - string: The extracted string.
//   - error: The error if the extraction fails.
func ExtractString(root gr.Tokener, id string) (string, error) {
	ok, err := IsToken(root, id)
	if err != nil {
		return "", err
	} else if !ok {
		return "", ue.NewErrUnexpected(root.GoString(), id)
	}

	return root.(*gr.LeafToken).Data, nil
}

// CoreFunc is a function type for the core function of an extractor.
//
// Parameters:
//   - roots: The roots to extract the strings from.
//   - result: The result to append the extracted strings to.
//   - pos: The position of the extraction.
//
// Returns:
//   - O: The result with the extracted strings.
//   - error: The error if the extraction fails.
type CoreFunc[O any] func(roots []gr.Tokener, result O, pos int) (O, error)

// Extractor extracts strings from a rule.
type Extractor[O any] struct {
	// lhs is the left-hand side of the rule.
	lhs string

	// checker is the AST checker.
	checker SyntaxChecker

	// core is the core function.
	core CoreFunc[O]
}

// NewExtractor creates a new extractor.
//
// Parameters:
//   - lhs: The left-hand side of the rule.
//   - checker: The AST checker.
//   - core: The core function.
//
// Returns:
//   - Extractor: The new extractor.
//
// The extractor extracts strings from a rule with the following structure:
//
//	lhs -> rhs1
//	lhs -> rhs1 x y ... z lhs
//
// Where:
//   - lhs is the left-hand side of the rule.
//   - rhs1 is the first right-hand side of the rule. Either a terminal or a non-terminal
//     symbol.
//   - x, y, ..., z are the right-hand sides of the rule. (handled by the core function)
//   - lhs is the last right-hand side of the rule. It is the same as the left-hand side.
//
// Here are the assumptions:
//   - 1st RHS can be a terminal or a non-terminal symbol.
//   - last RHS is always a non-terminal symbol.
//   - the last RHS is the same as LHS.
//
// Example:
//
//	fieldCls1 -> ATTR
//	fieldCls1 -> ATTR SEP fieldCls1
//
// In core function, pos starts from 1 and the result is nil initially.
// Inside the core function, you have to extract the strings from the children
// provided in the roots parameter. The structure and syntax is handled by the
// checker parameter.
func NewExtractor[O any](lhs string, checker SyntaxChecker, core CoreFunc[O]) Extractor[O] {
	return Extractor[O]{
		lhs:     lhs,
		checker: checker,
		core:    core,
	}
}

// Apply applies the extractor to the root.
//
// Parameters:
//   - root: The root to apply the extractor to.
//
// Returns:
//   - O: The result with the extracted strings.
//   - error: The error if the extraction fails.
func (e *Extractor[O]) Apply(root gr.Tokener) (O, error) {
	var result O

	if e.core == nil {
		return result, ue.NewErrNilParameter("core")
	}

	for pos := 0; ; pos++ {
		ok, err := IsToken(root, e.lhs)
		if err != nil {
			return result, fmt.Errorf("could not extract %q: %w", e.lhs, err)
		}

		if !ok {
			rootID := root.GetID()

			return result, ue.NewErrUnexpected(rootID, e.lhs)
		}

		// ASSUMPTION: 1st RHS is a non-terminal symbol.
		rootToken, ok := root.(*gr.NonLeafToken)
		if !ok {
			return result, NewErrAssumptionViolated(
				fmt.Errorf("token must be a non-leaf token: %T", root),
			)
		}

		err = e.checker.Check(rootToken.Data)
		if err != nil {
			return result, fmt.Errorf("%q does not match the grammar: %w", e.lhs, err)
		}

		lastToken := rootToken.Data[len(rootToken.Data)-1]

		isNotBaseCase, err := IsToken(lastToken, e.lhs)
		if err != nil {
			return result, fmt.Errorf("could not check if %q is not the base case: %w", e.lhs, err)
		}

		var todo []gr.Tokener

		if isNotBaseCase {
			todo = rootToken.Data[:len(rootToken.Data)-1]
		} else {
			todo = rootToken.Data
		}

		result, err = e.core(todo, result, pos)
		if err != nil {
			return result, fmt.Errorf("could not apply the core function: %w", err)
		}

		if !isNotBaseCase {
			return result, nil
		}

		// ASSUMPTION: last RHS is always a non-terminal symbol.
		// ASSUMPTION: the last RHS is the same as LHS.
		root = lastToken
	}
}
