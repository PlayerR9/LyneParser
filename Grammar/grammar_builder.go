package Grammar

import (
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

// AddProductions is a method of GrammarBuilder that adds a production to
// the GrammarBuilder.
//
// Parameters:
//   - ps: The productions to add to the GrammarBuilder.
func (b *GrammarBuilder) AddProductions(ps ...Productioner) {
	b.productions = append(b.productions, slext.SliceFilter(ps, FilterNilProduction)...)
}

// SetToSkip is a method of GrammarBuilder that sets the productions to skip
// in the GrammarBuilder.
//
// Parameters:
//   - lhss: The left-hand sides of the productions to skip.
func (b *GrammarBuilder) SetToSkip(lhss ...string) {
	b.skipProductions = append(b.skipProductions, lhss...)
}

// Build is a method of GrammarBuilder that builds a Grammar from the
// GrammarBuilder.
//
// Returns:
//   - *Grammar: A Grammar built from the GrammarBuilder.
//   - error: An error if the GrammarBuilder could not build a Grammar.
func (b *GrammarBuilder) Build() (*Grammar, error) {
	if b.productions == nil {
		return NewGrammar(), nil
	}

	b.productions = slext.RemoveDuplicatesFunc(b.productions)
	b.skipProductions = slext.RemoveDuplicates(b.skipProductions)
	b.skipProductions = slext.SliceFilter(b.skipProductions, b.FilterProductionsWithoutLHS)

	grammar := &Grammar{
		Symbols:     make([]string, 0),
		Productions: make([]Productioner, len(b.productions)),
		LhsToSkip:   make([]string, len(b.skipProductions)),
	}
	copy(grammar.Productions, b.productions)
	copy(grammar.LhsToSkip, b.skipProductions)

	for _, p := range b.productions {
		grammar.Symbols = append(grammar.Symbols, p.GetSymbols()...)
	}

	grammar.Symbols = slext.RemoveDuplicates(grammar.Symbols)

	err := grammar.compile()
	if err != nil {
		return nil, err
	}

	b.Reset()

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
