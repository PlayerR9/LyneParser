package Lexer

import (
	"slices"
	"strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// parseSingleRegexRule parses a single regex rule.
//
// Parameters:
//   - rule: The rule to parse.
//
// Returns:
//   - *RegProduction: The production.
//   - error: An error if there was a problem parsing the rule.
func parseSingleRegexRule(rule string) (*gr.RegProduction, error) {
	sides, err := gr.SplitByArrow(rule)
	if err != nil {
		return nil, err
	}

	regProd := gr.NewRegProduction(sides[0], sides[1])

	err = regProd.Compile()
	if err != nil {
		return nil, err
	}

	return regProd, nil
}

// parseRegexRules parses a string of regex rules.
//
// Parameters:
//   - rules: The rules to parse.
//
// Returns:
//   - []*RegProduction: A slice of productions.
//   - error: An error if there was a problem parsing the rules.
func parseRegexRules(rules string) ([]*gr.RegProduction, error) {
	lines := strings.Split(rules, "\n")
	lines = us.RemoveEmpty(lines)
	if len(lines) == 0 {
		return nil, nil
	}

	var productions []*gr.RegProduction

	for i, line := range lines {
		production, err := parseSingleRegexRule(line)
		if err != nil {
			return nil, uc.NewErrAt(i+1, "line", err)
		}

		if production != nil {
			productions = append(productions, production)
		}
	}

	return productions, nil
}

// Grammar represents a context-free grammar.
type Grammar struct {
	// productions is a slice of productions in the grammar.
	productions []*gr.RegProduction

	// lhsToSkip is a slice of productions to skip.
	lhsToSkip []string

	// symbols is a slice of symbols in the grammar.
	symbols []string
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
func NewGrammar(rules string, toSkip string) (*Grammar, error) {
	parsed, err := parseRegexRules(rules)
	if err != nil {
		return nil, uc.NewErrWhile("parsing regex rules", err)
	}

	if len(parsed) == 0 {
		return &Grammar{
			productions: make([]*gr.RegProduction, 0),
			lhsToSkip:   make([]string, 0),
			symbols:     make([]string, 0),
		}, nil
	}

	lhss := strings.Fields(toSkip)
	lhss = us.RemoveEmpty(lhss)
	lhss = us.Uniquefy(lhss, true)

	lhss = us.SliceFilter(
		lhss,
		func(lhs string) bool {
			filterProductionWithLHS := func(p *gr.RegProduction) bool {
				return p != nil && p.GetLhs() == lhs
			}

			return slices.ContainsFunc(parsed, filterProductionWithLHS)
		},
	)

	var symbols []string

	for _, p := range parsed {
		tmp := p.GetSymbols()

		for _, t := range tmp {
			pos, found := slices.BinarySearch(symbols, t)
			if !found {
				symbols = slices.Insert(symbols, pos, t)
			}
		}
	}

	return &Grammar{
		productions: parsed,
		lhsToSkip:   lhss,
		symbols:     symbols,
	}, nil
}

// GetSymbols returns a slice of symbols in the grammar.
//
// Returns:
//   - []string: A slice of symbols in the grammar.
func (g *Grammar) GetSymbols() []string {
	symbols := make([]string, len(g.symbols))
	copy(symbols, g.symbols)

	return symbols
}

// GetRegexProds returns a slice of RegProduction in the grammar.
//
// Returns:
//   - []*RegProduction: A slice of RegProduction in the grammar.
func (g *Grammar) GetRegexProds() []*gr.RegProduction {
	regProds := make([]*gr.RegProduction, len(g.productions))
	copy(regProds, g.productions)

	return regProds
}

// GetToSkip returns a slice of LHSs to skip.
//
// Returns:
//   - []string: A slice of LHSs to skip.
func (g *Grammar) GetToSkip() []string {
	toSkip := make([]string, len(g.lhsToSkip))
	copy(toSkip, g.lhsToSkip)

	return toSkip
}
