package Grammar

import (
	"fmt"
	"strings"

	ds "github.com/PlayerR9/MyGoLib/ListLike/DoubleLL"
)

const (
	// LeftToRight is the direction of a production from left to right.
	LeftToRight string = "->"

	// StartSymbolID is the identifier of the start symbol in the grammar.
	StartSymbolID string = "source"

	// EndSymbolID is the identifier of the end symbol in the grammar.
	EndSymbolID string = "EOF"

	// EpsilonSymbolID is the identifier of the epsilon symbol in the grammar.
	EpsilonSymbolID string = "ε"
)

// Grammar represents a context-free grammar.
type Grammar struct {
	// Productions is a slice of Productions in the grammar.
	Productions []Productioner

	// LhsToSkip is a slice of productions to skip.
	LhsToSkip []string

	// Symbols is a slice of Symbols in the grammar.
	Symbols []string
}

// NewGrammar is a constructor of an empty Grammar.
//
// A context-free grammar is a set of productions, each of which
// consists of a non-terminal symbol and a sequence of symbols.
//
// The non-terminal symbol is the left-hand side of the production,
// and the sequence of symbols is the right-hand side of the production.
//
// The grammar also contains a set of symbols, which are the
// non-terminal and terminal symbols in the grammar.
//
// Returns:
//   - *Grammar: A new empty Grammar.
func NewGrammar() *Grammar {
	return &Grammar{
		Productions: make([]Productioner, 0),
		LhsToSkip:   make([]string, 0),
		Symbols:     make([]string, 0),
	}
}

// String is a method of Grammar that returns a string representation
// of a Grammar.
//
// It should only be used for debugging and logging purposes.
//
// Returns:
//   - string: A string representation of a Grammar.
func (g *Grammar) String() string {
	if g == nil {
		return "Grammar[nil]"
	}

	values := make([]string, 0, len(g.Productions))

	for _, production := range g.Productions {
		values = append(values, production.String())
	}

	return fmt.Sprintf(
		"Grammar[productions=[%s], symbols=[%s], skipProductions=[%s]]",
		strings.Join(values, ", "),
		strings.Join(g.Symbols, ", "),
		strings.Join(g.LhsToSkip, ", "),
	)
}

// RegexMatch returns a slice of MatchedResult that match the input token.
//
// Parameters:
//   - at: The position in the input string.
//   - b: The input stream to match. Refers to Productioner.Match.
//
// Returns:
//   - []MatchedResult: A slice of MatchedResult that match the input token.
func (g *Grammar) RegexMatch(at int, b []byte) []*MatchedResult[*LeafToken] {
	matches := make([]*MatchedResult[*LeafToken], 0)

	for i, p := range g.Productions {
		val, ok := p.(*RegProduction)
		if !ok {
			continue
		}

		matched := val.Match(at, b)
		if matched != nil {
			matches = append(matches, NewMatchResult(matched, i))
		}
	}

	return matches
}

// ProductionMatch returns a slice of MatchedResult that match the input token.
//
// Parameters:
//   - at: The position in the input string.
//   - stack: The input stream to match. Refers to Productioner.Match.
//
// Returns:
//   - []MatchedResult: A slice of MatchedResult that match the input token.
func (g *Grammar) ProductionMatch(at int, stack *ds.DoubleStack[Tokener]) []*MatchedResult[*NonLeafToken] {
	matches := make([]*MatchedResult[*NonLeafToken], 0)

	for i, p := range g.Productions {
		val, ok := p.(*Production)
		if !ok {
			continue
		}

		matched, err := val.Match(at, stack)
		if err != nil {
			matches = append(matches, NewMatchResult(matched, i))
		}
	}

	return matches
}

// GetRegProductions returns a slice of RegProduction in the grammar.
//
// Returns:
//   - []*RegProduction: A slice of RegProduction in the grammar.
func (g *Grammar) GetRegProductions() []*RegProduction {
	regProds := make([]*RegProduction, 0, len(g.Productions))

	for _, p := range g.Productions {
		if val, ok := p.(*RegProduction); ok {
			regProds = append(regProds, val)
		}
	}

	return regProds
}

// GetProductions returns a slice of Production in the grammar.
//
// Returns:
//   - []*Production: A slice of Production in the grammar.
func (g *Grammar) GetProductions() []*Production {
	prods := make([]*Production, 0, len(g.Productions))

	for _, p := range g.Productions {
		if val, ok := p.(*Production); ok {
			prods = append(prods, val)
		}
	}

	return prods
}

// Compile compiles the grammar.
//
// It should be called before using the grammar.
//
// Returns:
//   - error: An error if the grammar could not be compiled.
func (g *Grammar) compile() error {
	regProds := g.GetRegProductions()

	for _, p := range regProds {
		err := p.Compile()
		if err != nil {
			return err
		}
	}

	return nil
}
