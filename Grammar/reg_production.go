package Grammar

import (
	"regexp"
	"strings"

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

// String is a method of fmt.Stringer that returns a string representation
// of a RegProduction.
//
// It should only be used for debugging and logging purposes.
//
// Returns:
//   - string: A string representation of a RegProduction.
func (r *RegProduction) GoString() string {
	var builder strings.Builder

	builder.WriteString("RegProduction{")
	builder.WriteString("lhs=")
	builder.WriteString(r.lhs)
	builder.WriteString(", rhs=")
	builder.WriteString(r.rhs)
	builder.WriteString(", rxp=")

	if r.rxp == nil {
		builder.WriteString("N/A")
	} else {
		str := r.rxp.String()

		builder.WriteString(str)
	}

	builder.WriteRune('}')

	return builder.String()
}

// Equals is a method of RegProduction that returns whether the production
// is equal to another production. Two productions are equal if their
// left-hand sides are equal and their right-hand sides are equal.
//
// Parameters:
//   - other: The other production to compare to.
//
// Returns:
//   - bool: Whether the production is equal to the other production.
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
//   - Tokener: A token that matches the production in the stack. nil if
//     there is no match.
func (p *RegProduction) Match(at int, b []byte) *LeafToken {
	data := p.rxp.Find(b)
	if data == nil {
		return nil
	}

	return NewLeafToken(p.lhs, string(data), at)
}

// Copy is a method of RegProduction that returns a copy of the production.
//
// Returns:
//   - uc.Copier: A copy of the production.
func (p *RegProduction) Copy() uc.Copier {
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
	return &RegProduction{
		lhs: lhs,
		rhs: "^" + regex,
	}
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
