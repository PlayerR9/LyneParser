package Ast

import (
	"fmt"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
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
func ExtractString[T uc.Enumer](root *gr.Token[T], id T) (string, error) {
	ok := IsToken(root, id)
	if !ok {
		return "", uc.NewErrUnexpected(root.GoString(), id.String())
	}

	return root.Data.(string), nil
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
type CoreFunc[O any, T uc.Enumer] func(roots []*gr.Token[T], result O, pos int) (O, error)

// Extractor extracts strings from a rule.
type Extractor[O any, T uc.Enumer] struct {
	// lhs is the left-hand side of the rule.
	lhs T

	// checker is the AST checker.
	checker *SyntaxChecker[T]

	// core is the core function.
	core CoreFunc[O, T]
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
func NewExtractor[O any, T uc.Enumer](lhs T, checker *SyntaxChecker[T], core CoreFunc[O, T]) Extractor[O, T] {
	extr := Extractor[O, T]{
		lhs:     lhs,
		checker: checker,
		core:    core,
	}
	return extr
}

// Apply applies the extractor to the root.
//
// Parameters:
//   - root: The root to apply the extractor to.
//
// Returns:
//   - O: The result with the extracted strings.
//   - error: The error if the extraction fails.
func (e *Extractor[O, T]) Apply(root *gr.Token[T]) (O, error) {
	var result O

	if e.core == nil {
		return result, uc.NewErrNilParameter("core")
	} else if e.checker == nil {
		return result, uc.NewErrNilParameter("checker")
	}

	for pos := 0; ; pos++ {
		ok := IsToken(root, e.lhs)
		if !ok {
			rootID := root.GetID()

			return result, uc.NewErrUnexpected(rootID.String(), e.lhs.String())
		}

		// ASSUMPTION: 1st RHS is a non-terminal symbol.
		uc.Assert(root.IsNonLeaf(), "token must be a non-leaf token")

		data := root.Data.([]*gr.Token[T])

		err := e.checker.Check(data)
		if err != nil {
			return result, fmt.Errorf("%q does not match the grammar: %w", e.lhs, err)
		}

		lastToken := data[len(data)-1]

		isNotBaseCase := IsToken(lastToken, e.lhs)

		var todo []*gr.Token[T]

		if isNotBaseCase {
			todo = data[:len(data)-1]
		} else {
			todo = data
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
