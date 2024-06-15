package Grammar

import (
	"slices"
	"strings"

	ds "github.com/PlayerR9/MyGoLib/ListLike/DoubleLL"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// parseProductionRule is a helper function that parses a production rule.
//
// Parameters:
//   - str: The production rule to parse.
//
// Returns:
//   - []*Production: A slice of productions.
//   - error: An error if there was a problem parsing the production rule.
func parseProductionRule(str string) ([]*Production, error) {
	sides, err := splitByArrow(str)
	if err != nil {
		return nil, ue.NewErrWhile("parsing production rules", err)
	}

	lhs := sides[0]
	rhs := sides[1]

	rhss := strings.Split(rhs, "|")

	for i := 0; i < len(rhss); i++ {
		rhss[i] = strings.TrimSpace(rhss[i])
	}

	var productions []*Production

	for _, r := range rhss {
		prod := NewProduction(lhs, r)

		productions = append(productions, prod)
	}

	return productions, nil
}

// ParserGrammar represents a context-free grammar.
type ParserGrammar struct {
	// productions is a slice of productions in the grammar.
	productions []*Production

	// symbols is a slice of symbols in the grammar.
	symbols []string
}

// NewParserGrammar is a constructor of an empty ParserGrammar.
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
func NewParserGrammar(rules string) (*ParserGrammar, error) {
	parsed := strings.Split(rules, "\n")
	parsed = us.RemoveEmpty(parsed)
	if len(parsed) == 0 {
		return &ParserGrammar{
			productions: make([]*Production, 0),
			symbols:     make([]string, 0),
		}, nil
	}

	// Parse production rules
	var productions []*Production

	for _, rule := range parsed {
		tmp, err := parseProductionRule(rule)
		if err != nil {
			return nil, ue.NewErrWhile("parsing production rules", err)
		}

		productions = append(productions, tmp...)
	}

	productions = us.UniquefyEquals(productions, true)

	if productions == nil {
		return &ParserGrammar{
			productions: make([]*Production, 0),
			symbols:     make([]string, 0),
		}, nil
	}

	var symbols []string

	for _, p := range productions {
		tmp := p.GetSymbols()

		for _, s := range tmp {
			pos, ok := slices.BinarySearch(symbols, s)
			if !ok {
				symbols = slices.Insert(symbols, pos, s)
			}
		}
	}

	grammar := &ParserGrammar{
		productions: productions,
		symbols:     symbols,
	}

	return grammar, nil
}

// GetSymbols returns a slice of symbols in the grammar.
//
// Returns:
//   - []string: A slice of symbols in the grammar.
func (g *ParserGrammar) GetSymbols() []string {
	symbols := make([]string, len(g.symbols))
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
func (g *ParserGrammar) ProductionMatch(at int, stack *ds.DoubleStack[Tokener]) []*MatchedResult[*NonLeafToken] {
	matches := make([]*MatchedResult[*NonLeafToken], 0)

	for i, p := range g.productions {
		matched, err := p.Match(at, stack)
		if err != nil {
			matches = append(matches, NewMatchResult(matched, i))
		}
	}

	return matches
}

// GetProductions returns a slice of Production in the grammar.
//
// Returns:
//   - []*Production: A slice of Production in the grammar.
func (g *ParserGrammar) GetProductions() []*Production {
	prods := make([]*Production, len(g.productions))
	copy(prods, g.productions)

	return prods
}
