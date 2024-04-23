package Grammar

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	itf "github.com/PlayerR9/MyGoLib/CustomData/Iterators"
	Stacker "github.com/PlayerR9/MyGoLib/ListLike/Common"
	impl "github.com/PlayerR9/MyGoLib/ListLike/Stack"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
	itff "github.com/PlayerR9/MyGoLib/Units/Interfaces"
)

// Productioner is an interface that defines methods for a production in a grammar.
type Productioner interface {
	// Equals returns whether the production is equal to another production.
	// Two productions are equal if their left-hand sides are equal and their
	// right-hand sides are equal.
	//
	// Parameters:
	//
	//   - other: The other production to compare to.
	//
	// Returns:
	//
	//   - bool: Whether the production is equal to the other production.
	Equals(other Productioner) bool

	// GetLhs returns the left-hand side of the production.
	//
	// Returns:
	//
	//   - string: The left-hand side of the production.
	GetLhs() string

	// GetSymbols returns a slice of symbols in the production. The slice
	// contains the left-hand side of the production and the right-hand side
	// of the production, with no duplicates.
	//
	// Returns:
	//
	//   - []string: A slice of symbols in the production.
	GetSymbols() []string

	// Match returns a token that matches the production at the given index
	// in the given stack. The token is a non-leaf token if the production is
	// a non-terminal production, and a leaf token if the production is a
	// terminal production.
	//

	// Match returns a token that matches the production in the given stack.
	// The token is a non-leaf token if the production is a non-terminal
	// production, and a leaf token if the production is a terminal production.
	//
	// Parameters:
	//
	//   - at: The current index in the input stack.
	//   - b: The input stream or stack to match the production against.
	//
	// Returns:
	//
	//   - Tokener: A token that matches the production in the input stream or stack.
	// 	nil if there is no match.
	//
	// Information:
	//
	//  - 'at' is the current index where the match is being attempted.
	//   It is used by the lexer to specify the position of the token in the
	//   input string. In parsers, however, it is not really used (at = 0).
	//   Despite that, it can be used to provide additional information to
	//   the parser for error reporting or debugging.
	Match(at int, b any) Tokener

	fmt.Stringer
	itff.Copier
}

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
//
//   - string: A string representation of a Production.
func (p *Production) String() string {
	if p == nil {
		return "Production[nil]"
	}

	if len(p.rhs) == 0 {
		return fmt.Sprintf("Production[%s %s]", p.lhs, LeftToRight)
	} else {
		return fmt.Sprintf("Production[%s %s %s]", p.lhs, LeftToRight, strings.Join(p.rhs, " "))
	}
}

// Equals is a method of Production that returns whether the production
// is equal to another production. Two productions are equal if their
// left-hand sides are equal and their right-hand sides are equal.
//
// Parameters:
//
//   - other: The other production to compare to.
//
// Returns:
//
//   - bool: Whether the production is equal to the other production.
func (p *Production) Equals(other Productioner) bool {
	if p == nil || other == nil || other.GetLhs() != p.lhs {
		return false
	}

	val, ok := other.(*Production)
	if !ok || len(val.rhs) != len(p.rhs) {
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
//
//   - string: The left-hand side of the production.
func (p *Production) GetLhs() string {
	return p.lhs
}

// Iterator is a method of Production that returns an iterator for the
// production that iterates over the right-hand side of the production.
//
// Returns:
//
//   - itf.Iterater[string]: An iterator for the production.
func (p *Production) Iterator() itf.Iterater[string] {
	return itf.IteratorFromSlice(p.rhs)
}

// ReverseIterator is a method of Production that returns a reverse
// iterator for the production that iterates over the right-hand side of
// the production in reverse.
//
// Returns:
//
//   - itf.Iterater[string]: A reverse iterator for the production.
func (p *Production) ReverseIterator() itf.Iterater[string] {
	slice := make([]string, len(p.rhs))
	copy(slice, p.rhs)
	slices.Reverse(slice)

	return itf.IteratorFromSlice(slice)
}

// GetSymbols is a method of Production that returns a slice of symbols
// in the production. The slice contains the left-hand side of the
// production and the right-hand side of the production, with no
// duplicates.
//
// Returns:
//
//   - []string: A slice of symbols in the production.
func (p *Production) GetSymbols() []string {
	symbols := make([]string, len(p.rhs)+1)
	copy(symbols, p.rhs)

	symbols[len(symbols)-1] = p.lhs

	// Remove duplicates
	for i := 0; i < len(symbols); {
		index := slices.Index(symbols[i+1:], symbols[i])

		if index != -1 {
			symbols = append(symbols[:index], symbols[index+1:]...)
		} else {
			i++
		}
	}

	return symbols
}

// Match is a method of Production that returns a token that matches the
// production in the given stack. The token is a non-leaf token if the
// production is a non-terminal production, and a leaf token if the
// production is a terminal production.
//
// Parameters:
//
//   - at: The current index in the input stack.
//   - b: The stack to match the production against.
//
// Returns:
//
//   - Tokener: A token that matches the production in the stack.
//
// Information:
//
//   - 'at' is the current index where the match is being attempted. It is
//     used by the lexer to specify the position of the token in the input
//     string. In parsers, however, it is not really used (at = 0). Despite
//     that, it can be used to provide additional information to the parser
//     for error reporting or debugging.
//   - as of now, only Stack.Stacker[Tokener] is supported as the type of
//     'b'.
func (p *Production) Match(at int, b any) Tokener {
	switch val := b.(type) {
	case Stacker.Stacker[Tokener]:
		popped := impl.NewArrayStack[Tokener]()

		for i := len(p.rhs) - 1; i >= 0; i-- {
			if val.IsEmpty() {
				// Push back the popped symbols.
				for !popped.IsEmpty() {
					val.Push(popped.Pop())
				}

				return nil
			}

			top := val.Peek()

			// FIXME: Implement a way to check if the top of the stack is a
			// non-leaf token.

			popped.Push(top)
			val.Pop()
		}

		slice := popped.Slice()
		slices.Reverse(slice)

		// Push back the popped symbols.
		for !popped.IsEmpty() {
			val.Push(popped.Pop())
		}

		tok := NewNonLeafToken(p.lhs, at, slice...)

		return &tok
	default:
		return nil
	}
}

// Copy is a method of Production that returns a copy of the production.
//
// Returns:
//
//   - itff.Copier: A copy of the production.
func (p *Production) Copy() itff.Copier {
	rhsCopy := make([]string, len(p.rhs))
	copy(rhsCopy, p.rhs)

	return &Production{
		lhs: p.lhs,
		rhs: rhsCopy,
	}
}

// NewProduction is a function that returns a new Production with the
// given left-hand side and right-hand side.
//
// Parameters:
//
//   - lhs: The left-hand side of the production.
//   - rhs: The right-hand side of the production.
//
// Returns:
//
//   - *Production: A new Production with the given left-hand side and
//     right-hand side.
func NewProduction(lhs string, rhs string) *Production {
	return &Production{lhs: lhs, rhs: strings.Fields(rhs)}
}

// Size is a method of Production that returns the number of symbols in
// the right-hand side of the production.
//
// Returns:
//
//   - int: The number of symbols in the right-hand side of the
//     production.
func (p *Production) Size() int {
	return len(p.rhs)
}

// GetRhsAt is a method of Production that returns the symbol at the
// given index in the right-hand side of the production.
//
// Parameters:
//
//   - index: The index of the symbol to get.
//
// Returns:
//
//   - string: The symbol at the given index in the right-hand side of
//     the production.
//   - error: An error of type *ErrInvalidParameter if the index is
//     invalid.
func (p *Production) GetRhsAt(index int) (string, error) {
	if index < 0 || index >= len(p.rhs) {
		return "", ers.NewErrInvalidParameter(
			"index",
			ers.NewErrOutOfBounds(index, 0, len(p.rhs)),
		)
	}

	return p.rhs[index], nil
}

// IndexOfRhs is a method of Production that returns the index of the
// given symbol in the right-hand side of the production.
//
// Parameters:
//
//   - rhs: The symbol to find the index of.
//
// Returns:
//
//   - int: The index of the symbol in the right-hand side of the
//     production. Returns -1 if the symbol is not found.
func (p *Production) IndexOfRhs(rhs string) int {
	return slices.Index(p.rhs, rhs)
}

// RegProduction represents a production in a grammar that matches a
// regular expression.
type RegProduction struct {
	// Left-hand side of the production.
	lhs string

	// Right-hand side of the production.
	rhs string

	// Regular expression to match the right-hand side of the production.
	rxp *regexp.Regexp
}

// String is a method of fmt.Stringer that returns a string representation
// of a RegProduction.
//
// It should only be used for debugging and logging purposes.
//
// Returns:
//
//   - string: A string representation of a RegProduction.
func (r *RegProduction) String() string {
	if r == nil {
		return "RegProduction[nil]"
	}

	if r.rxp == nil {
		return fmt.Sprintf("RegProduction[lhs=%s, rhs=%s, rxp=N/A]", r.lhs, r.rhs)
	}

	return fmt.Sprintf("RegProduction[lhs=%s, rhs=%s, rxp=%v]", r.lhs, r.rhs, r.rxp)
}

// Equals is a method of RegProduction that returns whether the production
// is equal to another production. Two productions are equal if their
// left-hand sides are equal and their right-hand sides are equal.
//
// Parameters:
//
//   - other: The other production to compare to.
//
// Returns:
//
//   - bool: Whether the production is equal to the other production.
func (p *RegProduction) Equals(other Productioner) bool {
	if p == nil || other == nil || other.GetLhs() != p.lhs {
		return false
	}

	val, ok := other.(*RegProduction)
	return ok && val.rhs == p.rhs
}

// GetLhs is a method of RegProduction that returns the left-hand side of
// the production.
//
// Returns:
//
//   - string: The left-hand side of the production.
func (p *RegProduction) GetLhs() string {
	return p.lhs
}

// GetSymbols is a method of RegProduction that returns a slice of symbols
// in the production. The slice contains the left-hand side of the
// production.
//
// Returns:
//
//   - []string: A slice of symbols in the production.
func (p *RegProduction) GetSymbols() []string {
	return []string{p.lhs}
}

// Match is a method of RegProduction that returns a token that matches the
// production in the given stack. The token is a non-leaf token if the
// production is a non-terminal production, and a leaf token if the
// production is a terminal production.
//
// Parameters:
//
//   - at: The current index in the input stack.
//   - b: The slice of bytes to match the production against.
//
// Returns:
//
//   - Tokener: A token that matches the production in the stack. nil if
//     there is no match.
func (p *RegProduction) Match(at int, b any) Tokener {
	val, ok := b.([]byte)
	if !ok {
		return nil
	}

	data := p.rxp.Find(val)
	if data == nil {
		return nil
	}

	tok := NewLeafToken(p.lhs, string(data), at)

	return &tok
}

// Copy is a method of RegProduction that returns a copy of the production.
//
// Returns:
//
//   - itff.Copier: A copy of the production.
func (p *RegProduction) Copy() itff.Copier {
	return &RegProduction{
		lhs: p.lhs,
		rhs: p.rhs,
		rxp: p.rxp,
	}
}

// NewRegProduction is a function that returns a new RegProduction with the
// given left-hand side and regular expression.
//
// It adds the '^' character to the beginning of the regular expression to
// match the beginning of the input string.
//
// Parameters:
//
//   - lhs: The left-hand side of the production.
//   - regex: The regular expression to match the right-hand side of the
//     production.
//
// Returns:
//
//   - *RegProduction: A new RegProduction with the given left-hand side
//     and regular expression.
//
// Information:
//
//   - Must call Compile() on the returned RegProduction to compile the
//     regular expression.
func NewRegProduction(lhs string, regex string) *RegProduction {
	return &RegProduction{
		lhs: lhs,
		rhs: "^" + regex,
	}
}

// Compile is a method of RegProduction that compiles the regular
// expression of the production.
//
// Returns:
//
//   - error: An error if the regular expression cannot be compiled.
func (r *RegProduction) Compile() error {
	rxp, err := regexp.Compile(r.rhs)
	if err != nil {
		return err
	}

	r.rxp = rxp
	return nil
}
