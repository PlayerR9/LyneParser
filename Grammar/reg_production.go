package Grammar

import (
	"fmt"
	"regexp"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// RegProduction represents a production in a grammar that matches a
// regular expression.
type RegProduction[T TokenTyper] struct {
	// Left-hand side of the production.
	lhs T

	// Right-hand side of the production.
	rhs string

	// Regular expression to match the right-hand side of the production.
	rxp *regexp.Regexp
}

// GoString implements the fmt.GoStringer interface.
func (r *RegProduction[T]) GoString() string {
	str := fmt.Sprintf("%+v", *r)
	return str
}

// Equals implements the common.Equaler interface.
//
// Two productions are equal if their left-hand sides are equal and their
// right-hand sides are equal.
func (p *RegProduction[T]) Equals(other uc.Equaler) bool {
	if other == nil {
		return false
	}

	val, ok := other.(*RegProduction[T])
	if !ok {
		return false
	}

	return val.lhs == p.lhs && val.rhs == p.rhs
}

// Copy implements the common.Copier interface.
func (p *RegProduction[T]) Copy() uc.Copier {
	p_copy := &RegProduction[T]{
		lhs: p.lhs,
		rhs: p.rhs,
		rxp: p.rxp,
	}
	return p_copy
}

// NewRegProduction is a function that returns a new RegProduction with the
// given left-hand side and regular expression.
//
// It adds the '^' character to the beginning of the regular expression to
// match the beginning of the input string.
//
// Parameters:
//   - lhs: The left-hand side of the production.
//   - regex: The regular expression to match the right-hand side of the
//     production.
//
// Returns:
//   - *RegProduction: A new RegProduction with the given left-hand side
//     and regular expression.
//
// Information:
//   - Must call Compile() on the returned RegProduction to compile the
//     regular expression.
func NewRegProduction[T TokenTyper](lhs T, regex string) *RegProduction[T] {
	p := &RegProduction[T]{
		lhs: lhs,
		rhs: "^" + regex,
	}
	return p
}

// GetLhs is a method of RegProduction that returns the left-hand side of
// the production.
//
// Returns:
//   - T: The left-hand side of the production.
func (p *RegProduction[T]) GetLhs() T {
	return p.lhs
}

// GetSymbols is a method of RegProduction that returns a slice of symbols
// in the production. The slice contains the left-hand side of the
// production.
//
// Returns:
//   - []T: A slice of symbols in the production.
func (p *RegProduction[T]) GetSymbols() []T {
	return []T{p.lhs}
}

// Match is a method of RegProduction that returns a token that matches the
// production in the given stack. The token is a non-leaf token if the
// production is a non-terminal production, and a leaf token if the
// production is a terminal production.
//
// Parameters:
//   - at: The current index in the input stack.
//   - b: The slice of bytes to match the production against.
//
// Returns:
//   - Token: A token that matches the production in the stack.
//   - bool: True if the production matches the input stack, false
//     otherwise.
func (p *RegProduction[T]) MatchRegProd(at int, b []byte) (*Token[T], bool) {
	data := p.rxp.Find(b)
	if data == nil {
		return nil, false
	}

	// Must be an exact match.
	if len(data) != len(b) {
		return nil, false
	}

	lt := NewToken(p.lhs, string(data), at, nil)

	return lt, true
}

// Compile is a method of RegProduction that compiles the regular
// expression of the production.
//
// Returns:
//   - error: An error if the regular expression cannot be compiled.
func (r *RegProduction[T]) Compile() error {
	rxp, err := regexp.Compile(r.rhs)
	if err != nil {
		return err
	}

	r.rxp = rxp
	return nil
}
