package Parser

import (
	"slices"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	ud "github.com/PlayerR9/MyGoLib/Units/Debugging"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// Grammar represents a context-free grammar.
type Grammar[T uc.Enumer] struct {
	// productions is a slice of productions in the grammar.
	productions []*gr.Production[T]

	// symbols is a slice of symbols in the grammar.
	symbols []T
}

// NewGrammar is a constructor of an empty ParserGrammar.
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
//   - *ParserGrammar: A new empty ParserGrammar.
func NewGrammar[T uc.Enumer]() (*Grammar[T], error) {
	g := &Grammar[T]{
		productions: make([]*gr.Production[T], 0),
		symbols:     make([]T, 0),
	}
	return g, nil

}

// AddRule adds a new rule to the grammar.
//
// Parameters:
//   - lhs: The left-hand side of the production.
//   - rhss: The right-hand side of the production.
//
// Returns:
//   - error: An error if there was a problem adding the rule.
func (g *Grammar[T]) AddRule(lhs T, rhss []T) error {
	production := gr.NewProduction(lhs, rhss)

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

// ProductionMatch returns a slice of MatchedResult that match the input token.
//
// Parameters:
//   - at: The position in the input string.
//   - stack: The input stream to match. Refers to Productioner.Match.
//
// Returns:
//   - []MatchedResult: A slice of MatchedResult that match the input token.
func (g *Grammar[T]) ProductionMatch(at int, stack *ud.History[lls.Stacker[*gr.Token[T]]]) []*gr.MatchedResult[T] {
	var matches []*gr.MatchedResult[T]

	for i, p := range g.productions {
		matched, err := p.Match(at, stack)
		if err != nil {
			mr := gr.NewMatchResult(matched, i)
			matches = append(matches, mr)
		}
	}

	return matches
}

// GetProductions returns a slice of Production in the grammar.
//
// Returns:
//   - []*Production: A slice of Production in the grammar.
func (g *Grammar[T]) GetProductions() []*gr.Production[T] {
	prods := make([]*gr.Production[T], len(g.productions))
	copy(prods, g.productions)

	return prods
}
