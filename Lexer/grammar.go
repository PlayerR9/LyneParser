package Lexer

import (
	"slices"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// Grammar represents a context-free grammar.
type Grammar[T gr.TokenTyper] struct {
	// productions is a slice of productions in the grammar.
	productions []*gr.RegProduction[T]

	// lhs_to_skip is a slice of productions to skip.
	lhs_to_skip []T

	// symbols is a slice of symbols in the grammar.
	symbols []T
}

// Fix implements the object.Fixer interface.
//
// Never errors.
func (g *Grammar[T]) Fix() error {
	g.lhs_to_skip = us.SliceFilter(
		g.lhs_to_skip,
		func(lhs T) bool {
			filter_production_with_lhs := func(p *gr.RegProduction[T]) bool {
				return p != nil && p.GetLhs() == lhs
			}

			return slices.ContainsFunc(g.productions, filter_production_with_lhs)
		},
	)

	return nil
}

// NewGrammar is a constructor of an empty LexerGrammar.
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
//   - *LexerGrammar: A new empty LexerGrammar.
func NewGrammar[T gr.TokenTyper](to_skip []T) *Grammar[T] {
	to_skip = us.Uniquefy(to_skip, true)

	g := &Grammar[T]{
		lhs_to_skip: to_skip,
	}

	return g
}

// AddRule adds a new rule to the grammar.
//
// Parameters:
//   - lhs: The left-hand side of the production.
//   - regex: The regular expression of the production.
//
// Returns:
//   - error: An error if there was a problem adding the rule.
func (g *Grammar[T]) AddRule(lhs T, regex string) error {
	production := gr.NewRegProduction(lhs, regex)

	err := production.Compile()
	if err != nil {
		return err
	}

	g.productions = append(g.productions, production)

	tmp := production.GetSymbols()

	for _, t := range tmp {
		pos, found := slices.BinarySearch(g.symbols, t)
		if !found {
			g.symbols = slices.Insert(g.symbols, pos, t)
		}
	}

	return nil
}

// GetSymbols returns a slice of symbols in the grammar.
//
// Returns:
//   - []T: A slice of symbols in the grammar.
func (g *Grammar[T]) GetSymbols() []T {
	symbols := make([]T, len(g.symbols))
	copy(symbols, g.symbols)

	return symbols
}

// GetRegexProds returns a slice of RegProduction in the grammar.
//
// Returns:
//   - []*RegProduction: A slice of RegProduction in the grammar.
func (g *Grammar[T]) GetRegexProds() []*gr.RegProduction[T] {
	reg_prods := make([]*gr.RegProduction[T], len(g.productions))
	copy(reg_prods, g.productions)

	return reg_prods
}

// GetToSkip returns a slice of LHSs to skip.
//
// Returns:
//   - []T: A slice of LHSs to skip.
func (g *Grammar[T]) GetToSkip() []T {
	to_skip := make([]T, len(g.lhs_to_skip))
	copy(to_skip, g.lhs_to_skip)

	return to_skip
}
