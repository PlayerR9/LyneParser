package Ast

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// ExtractString extracts a string from a token.
//
// Parameters:
//
//   - root: The root to extract the string from.
//   - id: The ID of the token.
//
// Returns:
//
//   - string: The extracted string.
//   - error: The error if the extraction fails.
func ExtractString(root gr.Tokener, id string) (string, error) {
	if ok, err := IsToken(root, id); err != nil {
		return "", err
	} else if !ok {
		return "", ers.NewErrUnexpected(root, id)
	}

	return root.(*gr.LeafToken).Data, nil
}

// CoreFunc is a function type for the core function of an extractor.
//
// Parameters:
//
//   - roots: The roots to extract the strings from.
//   - result: The result to append the extracted strings to.
//   - pos: The position of the extraction.
//
// Returns:
//
//   - O: The result with the extracted strings.
//   - error: The error if the extraction fails.
type CoreFunc[O any] func(roots []gr.Tokener, result O, pos int) (O, error)

// Extractor extracts strings from a rule.
type Extractor[O any] struct {
	// lhs is the left-hand side of the rule.
	lhs string

	// checker is the AST checker.
	checker ASTer

	// core is the core function.
	core CoreFunc[O]
}

// NewExtractor creates a new extractor.
//
// Parameters:
//
//   - lhs: The left-hand side of the rule.
//   - checker: The AST checker.
//   - core: The core function.
//
// Returns:
//
//   - Extractor: The new extractor.
//
// The extractor extracts strings from a rule with the following structure:
//
//	lhs -> rhs1
//	lhs -> rhs1 x y ... z lhs
//
// Where:
//
//   - lhs is the left-hand side of the rule.
//   - rhs1 is the first right-hand side of the rule. (non-terminal)
//   - x, y, ..., z are the right-hand sides of the rule. (handled by the core function)
//   - lhs is the last right-hand side of the rule. It is the same as the left-hand side.
//
// Here are the assumptions:
//
//   - 1st RHS is a non-terminal symbol.
//   - last RHS is always a non-terminal symbol.
//   - the last RHS is the same as LHS.
//
// Example:
//
//	fieldCls1 -> ATTR
//	fieldCls1 -> ATTR SEP fieldCls1
func NewExtractor[O any](lhs string, checker ASTer, core CoreFunc[O]) Extractor[O] {
	return Extractor[O]{
		lhs:     lhs,
		checker: checker,
		core:    core,
	}
}

// Apply applies the extractor to the root.
//
// Parameters:
//
//   - root: The root to apply the extractor to.
//
// Returns:
//
//   - O: The result with the extracted strings.
//   - error: The error if the extraction fails.
func (e *Extractor[O]) Apply(root gr.Tokener) (O, error) {
	var result O

	var err error
	var ok bool

	for i := 1; root != nil; i++ {
		ok, err = IsToken(root, e.lhs)
		if err != nil {
			break
		} else if !ok {
			err = ers.NewErrUnexpected(root, e.lhs)
			break
		}

		// ASSUMPTION: 1st RHS is a non-terminal symbol.
		ruleRoot := root.(*gr.NonLeafToken)

		err = e.checker.Check(ruleRoot.Data)
		if err != nil {
			break
		}

		result, err = e.core(ruleRoot.Data, result, i)
		if err != nil {
			break
		}

		if len(ruleRoot.Data) == 1 {
			root = nil
		} else {
			// ASSUMPTION: last RHS is always a non-terminal symbol.
			// ASSUMPTION: the last RHS is the same as LHS.
			root = ruleRoot.Data[len(ruleRoot.Data)-1]
		}
	}

	return result, err
}
