package Grammar

import (
	"fmt"
	"regexp"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

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

// GoString implements the fmt.GoStringer interface.
func (r *RegProduction) GoString() string {
	str := fmt.Sprintf("%+v", *r)

	return str
}

// Equals implements the common.Equaler interface.
//
// Two productions are equal if their left-hand sides are equal and their
// right-hand sides are equal.
func (p *RegProduction) Equals(other uc.Equaler) bool {
	if other == nil {
		return false
	}

	val, ok := other.(*RegProduction)
	if !ok {
		return false
	}

	return val.lhs == p.lhs && val.rhs == p.rhs
}

// Copy implements the common.Copier interface.
func (p *RegProduction) Copy() uc.Copier {
	pCopy := &RegProduction{
		lhs: p.lhs,
		rhs: p.rhs,
		rxp: p.rxp,
	}
	return pCopy
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
func NewRegProduction(lhs string, regex string) *RegProduction {
	p := &RegProduction{
		lhs: lhs,
		rhs: "^" + regex,
	}
	return p
}

// GetLhs is a method of RegProduction that returns the left-hand side of
// the production.
//
// Returns:
//   - string: The left-hand side of the production.
func (p *RegProduction) GetLhs() string {
	return p.lhs
}

// GetSymbols is a method of RegProduction that returns a slice of symbols
// in the production. The slice contains the left-hand side of the
// production.
//
// Returns:
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
//   - at: The current index in the input stack.
//   - b: The slice of bytes to match the production against.
//
// Returns:
//   - Token: A token that matches the production in the stack.
//   - bool: True if the production matches the input stack, false
//     otherwise.
func (p *RegProduction) Match(at int, b []byte) (Token, bool) {
	data := p.rxp.Find(b)
	if data == nil {
		return Token{}, false
	}

	// Must be an exact match.
	if len(data) != len(b) {
		return Token{}, false
	}

	lt := NewToken(p.lhs, string(data), at, nil)

	return lt, true
}

// Compile is a method of RegProduction that compiles the regular
// expression of the production.
//
// Returns:
//   - error: An error if the regular expression cannot be compiled.
func (r *RegProduction) Compile() error {
	rxp, err := regexp.Compile(r.rhs)
	if err != nil {
		return err
	}

	r.rxp = rxp
	return nil
}
