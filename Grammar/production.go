package Grammar

import (
	"slices"
	"strings"

	ds "github.com/PlayerR9/MyGoLib/ListLike/DoubleLL"
	itf "github.com/PlayerR9/MyGoLib/Units/Iterators"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// Production represents a production in a grammar.
type Production struct {
	// Left-hand side of the production.
	lhs string

	// Right-hand side of the production.
	rhs []string
}

// String is a method of Production that returns a string representation
// of a Production.
//
// Returns:
//   - string: A string representation of a Production.
func (p *Production) String() string {
	var builder strings.Builder

	builder.WriteString(p.lhs)
	builder.WriteRune(' ')
	builder.WriteString(LeftToRight)
	builder.WriteRune(' ')

	if len(p.rhs) == 0 {
		builder.WriteString(EpsilonSymbolID)
	} else {
		str := strings.Join(p.rhs, " ")

		builder.WriteString(str)
	}

	return builder.String()
}

// Equals is a method of Production that returns whether the production
// is equal to another production. Two productions are equal if their
// left-hand sides are equal and their right-hand sides are equal.
//
// Parameters:
//   - other: The other production to compare to.
//
// Returns:
//   - bool: Whether the production is equal to the other production.
func (p *Production) Equals(other uc.Equaler) bool {
	if other == nil {
		return false
	}

	val, ok := other.(*Production)
	if !ok {
		return false
	}

	if val.lhs != p.lhs || len(val.rhs) != len(p.rhs) {
		return false
	}

	for i, symbol := range p.rhs {
		if symbol != val.rhs[i] {
			return false
		}
	}

	return true
}

// GetLhs is a method of Production that returns the left-hand side of
// the production.
//
// Returns:
//   - string: The left-hand side of the production.
func (p *Production) GetLhs() string {
	return p.lhs
}

// Iterator is a method of Production that returns an iterator for the
// production that iterates over the right-hand side of the production.
//
// Returns:
//   - itf.Iterater[string]: An iterator for the production.
func (p *Production) Iterator() itf.Iterater[string] {
	return itf.NewSimpleIterator(p.rhs)
}

// ReverseIterator is a method of Production that returns a reverse
// iterator for the production that iterates over the right-hand side of
// the production in reverse.
//
// Returns:
//   - itf.Iterater[string]: A reverse iterator for the production.
func (p *Production) ReverseIterator() itf.Iterater[string] {
	slice := make([]string, len(p.rhs))
	copy(slice, p.rhs)
	slices.Reverse(slice)

	return itf.NewSimpleIterator(slice)
}

// GetSymbols is a method of Production that returns a slice of symbols
// in the production. The slice contains the left-hand side of the
// production and the right-hand side of the production, with no
// duplicates.
//
// Returns:
//   - []string: A slice of symbols in the production.
func (p *Production) GetSymbols() []string {
	symbols := make([]string, len(p.rhs)+1)
	copy(symbols, p.rhs)

	symbols[len(symbols)-1] = p.lhs

	return us.Uniquefy(symbols, true)
}

// Match is a method of Production that returns a token that matches the
// production in the given stack. The token is a non-leaf token if the
// production is a non-terminal production, and a leaf token if the
// production is a terminal production.
//
// Parameters:
//   - at: The current index in the input stack.
//   - stack: The stack to match the production against.
//
// Returns:
//   - Tokener: A token that matches the production in the stack.
//
// Information:
//   - 'at' is the current index where the match is being attempted. It is
//     used by the lexer to specify the position of the token in the input
//     string. In parsers, however, it is not really used (at = 0). Despite
//     that, it can be used to provide additional information to the parser
//     for error reporting or debugging.
func (p *Production) Match(at int, stack *ds.DoubleStack[Tokener]) (*NonLeafToken, error) {
	solutions := make([]Tokener, 0)

	var reason error = nil

	for i := len(p.rhs) - 1; i >= 0; i-- {
		rhs := p.rhs[i]

		top, ok := stack.Pop()
		if !ok {
			reason = ue.NewErrUnexpected("", rhs)
			break
		}

		id := top.GetID()
		if id != rhs {
			reason = ue.NewErrUnexpected(top.GoString(), rhs)
			break
		}

		solutions = append(solutions, top)
	}

	stack.Refuse()

	if reason != nil {
		return nil, reason
	}

	slices.Reverse(solutions)

	tok := NewNonLeafToken(p.lhs, at, solutions...)

	return tok, nil
}

// Copy is a method of Production that returns a copy of the production.
//
// Returns:
//   - uc.Copier: A copy of the production.
func (p *Production) Copy() uc.Copier {
	pCopy := &Production{
		lhs: p.lhs,
		rhs: make([]string, len(p.rhs)),
	}
	copy(pCopy.rhs, p.rhs)

	return pCopy
}

// NewProduction is a function that returns a new Production with the
// given left-hand side and right-hand side.
//
// Parameters:
//   - lhs: The left-hand side of the production.
//   - rhs: The right-hand side of the production.
//
// Returns:
//   - *Production: A new Production with the given left-hand side and
//     right-hand side.
func NewProduction(lhs string, rhs string) *Production {
	return &Production{lhs: lhs, rhs: strings.Fields(rhs)}
}

// Size is a method of Production that returns the number of symbols in
// the right-hand side of the production.
//
// Returns:
//   - int: The number of symbols in the right-hand side of the
//     production.
func (p *Production) Size() int {
	return len(p.rhs)
}

// GetRhsAt is a method of Production that returns the symbol at the
// given index in the right-hand side of the production.
//
// Parameters:
//   - index: The index of the symbol to get.
//
// Returns:
//   - string: The symbol at the given index in the right-hand side of
//     the production.
//   - error: An error of type *ErrInvalidParameter if the index is
//     invalid.
func (p *Production) GetRhsAt(index int) (string, error) {
	if index < 0 || index >= len(p.rhs) {
		return "", ue.NewErrInvalidParameter(
			"index",
			ue.NewErrOutOfBounds(index, 0, len(p.rhs)),
		)
	}

	return p.rhs[index], nil
}

// IndicesOfRhs is a method of Production that returns the indices of the
// symbol in the right-hand side of the production.
//
// Parameters:
//   - rhs: The symbol to find the index of.
//
// Returns:
//   - []int: The indices of the symbol in the right-hand side of the
//     production.
func (p *Production) IndicesOfRhs(rhs string) []int {
	results := make([]int, 0)

	for i, symbol := range p.rhs {
		if symbol == rhs {
			results = append(results, i)
		}
	}

	return results
}

// ReplaceRhsAt is a method of Production that replaces the symbol at the
// given index in the right-hand side of the production with the right-hand
// side of another production.
//
// Parameters:
//   - index: The index of the symbol to replace.
//   - otherP: The other production to replace the symbol with.
//
// Returns:
//   - *Production: A new production with the symbol at the given index
//     replaced with the right-hand side of the other production.
//   - error: An error if the index is invalid or the other production is nil.
//
// Errors:
//   - *ue.ErrInvalidParameter: If the index is invalid or the other production is nil.
//   - *ErrLhsRhsMismatch: If the left-hand side of the other production does
//     not match the symbol at the given index in the right-hand side of the
//     production.
func (p *Production) ReplaceRhsAt(index int, rhs string) *Production {
	newP := p.Copy().(*Production)

	if index >= 0 && index < len(p.rhs) {
		newP.rhs[index] = rhs
	}

	return newP
}

// ReplaceRhsAt is a method of Production that replaces the symbol at the
// given index in the right-hand side of the production with the right-hand
// side of another production.
//
// Parameters:
//   - index: The index of the symbol to replace.
//   - otherP: The other production to replace the symbol with.
//
// Returns:
//   - *Production: A new production with the symbol at the given index
//     replaced with the right-hand side of the other production.
//   - error: An error if the index is invalid or the other production is nil.
//
// Errors:
//   - *ue.ErrInvalidParameter: If the index is invalid or the other production is nil.
//   - *ErrLhsRhsMismatch: If the left-hand side of the other production does
//     not match the symbol at the given index in the right-hand side of the
//     production.
func (p *Production) SubstituteRhsAt(index int, otherP *Production) *Production {
	newP := p.Copy().(*Production)

	if index < 0 || index >= len(p.rhs) || otherP == nil {
		return newP
	}

	if index == 0 {
		newP.rhs = append(otherP.rhs, newP.rhs[1:]...)
	} else if index == len(p.rhs)-1 {
		newP.rhs = append(newP.rhs[:index], otherP.rhs...)
	} else {
		newP.rhs = append(newP.rhs[:index], otherP.rhs...)
		newP.rhs = append(newP.rhs, newP.rhs[index+1:]...)
	}

	return newP
}

// HasRhs is a method of Production that returns whether the right-hand
// side of the production contains the given symbol.
//
// Parameters:
//   - rhs: The symbol to check for.
//
// Returns:
//   - bool: Whether the right-hand side of the production contains the
//     given symbol.
func (p *Production) HasRhs(rhs string) bool {
	for _, symbol := range p.rhs {
		if symbol == rhs {
			return true
		}
	}

	return false
}
