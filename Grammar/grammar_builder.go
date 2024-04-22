package Grammar

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

// GrammarBuilder represents a builder for a grammar.
//
// The default direction of the productions is LeftToRight.
type GrammarBuilder struct {
	// Slice of productions to add to the grammar.
	productions []Productioner

	// Slice of productions to skip.
	skipProductions []string
}

// String is a method of GrammarBuilder that returns a string
// representation of a GrammarBuilder.
//
// It should only be used for debugging and logging purposes.
//
// Returns:
//
//   - string: A string representation of a GrammarBuilder.
func (b *GrammarBuilder) String() string {
	if b.productions == nil {
		return "GrammarBuilder[nil]"
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "GrammarBuilder[total=%d, productions=[", len(b.productions))

	if len(b.productions) != 0 {
		fmt.Fprintf(&builder, "%v", b.productions[0])

		for _, production := range b.productions[1:] {
			fmt.Fprintf(&builder, ", %v", production)
		}
	}

	builder.WriteString("], skipProductions=[")

	if len(b.skipProductions) != 0 {
		fmt.Fprintf(&builder, "%v", b.skipProductions[0])

		for _, production := range b.skipProductions[1:] {
			fmt.Fprintf(&builder, ", %v", production)
		}
	}

	builder.WriteString("]]")

	return builder.String()
}

// AddProduction is a method of GrammarBuilder that adds a production to
// the GrammarBuilder.
//
// Parameters:
//
//   - p: The production to add to the GrammarBuilder.
func (b *GrammarBuilder) AddProduction(p ...Productioner) {
	for _, production := range p {
		if production == nil {
			continue
		}

		if b.productions == nil {
			b.productions = []Productioner{production}
		} else {
			b.productions = append(b.productions, production)
		}
	}
}

// SetToSkip is a method of GrammarBuilder that sets the productions to skip
// in the GrammarBuilder.
//
// Parameters:
//
//   - lhss: The left-hand sides of the productions to skip.
func (b *GrammarBuilder) SetToSkip(lhss ...string) {
	if b.skipProductions == nil {
		b.skipProductions = lhss
	} else {
		b.skipProductions = append(b.skipProductions, lhss...)
	}
}

// Build is a method of GrammarBuilder that builds a Grammar from the
// GrammarBuilder.
//
// Returns:
//
//   - *Grammar: A Grammar built from the GrammarBuilder.
func (b *GrammarBuilder) Build() (*Grammar, error) {
	if b.productions == nil {
		return &Grammar{
			Productions: make([]Productioner, 0),
			Symbols:     make([]string, 0),
			LhsToSkip:   make([]string, 0),
		}, nil
	}

	grammar := &Grammar{
		Symbols: make([]string, 0),
	}

	// 1. Remove duplicates
	for i := 0; i < len(b.productions); {
		index := slices.IndexFunc(b.productions[i+1:], func(p Productioner) bool {
			return p.Equals(b.productions[i])
		})

		if index != -1 {
			b.productions = append(b.productions[:index], b.productions[index+1:]...)
		} else {
			i++
		}
	}

	for i := 0; i < len(b.skipProductions); {
		index := slices.Index(b.skipProductions[i+1:], b.skipProductions[i])

		if index != -1 {
			b.skipProductions = append(b.skipProductions[:index], b.skipProductions[index+1:]...)
		} else {
			i++
		}
	}

	// 2. Remove LHS to skip if they don't exist in productions
	b.skipProductions = slext.SliceFilter(b.skipProductions, func(lhs string) bool {
		return slices.ContainsFunc(b.productions, func(p Productioner) bool {
			return p.GetLhs() == lhs
		})
	})

	// 3. Add productions to grammar
	grammar.Productions = make([]Productioner, len(b.productions))
	copy(grammar.Productions, b.productions)

	grammar.LhsToSkip = make([]string, len(b.skipProductions))
	copy(grammar.LhsToSkip, b.skipProductions)

	// 4. Add symbols to grammar
	for _, p := range b.productions {
		for _, symbol := range p.GetSymbols() {
			if !slices.Contains(grammar.Symbols, symbol) {
				grammar.Symbols = append(grammar.Symbols, symbol)
			}
		}
	}

	// 4. Compile all regular expressions
	var err error

	for i, p := range grammar.Productions {
		val, ok := p.(*RegProduction)
		if !ok {
			continue
		}

		val.rxp, err = regexp.Compile(val.rhs)
		if err != nil {
			return nil, err
		}

		grammar.Productions[i] = val
	}

	// 4. Clear builder
	for i := range b.productions {
		b.productions[i] = nil
	}

	b.productions = nil
	b.skipProductions = nil

	return grammar, nil
}

// Reset is a method of GrammarBuilder that resets a GrammarBuilder.
func (b *GrammarBuilder) Reset() {
	for i := range b.productions {
		b.productions[i] = nil
	}

	b.productions = nil
	b.skipProductions = nil
}
