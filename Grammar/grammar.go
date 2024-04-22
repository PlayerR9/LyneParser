package Grammar

import (
	"fmt"
	"strings"
)

const (
	// LeftToRight is the direction of a production from left to right.
	LeftToRight string = "->"

	// StartSymbolID is the identifier of the start symbol in the grammar.
	StartSymbolID string = "source"

	// EndSymbolID is the identifier of the end symbol in the grammar.
	EndSymbolID string = "EOF"
)

// Grammar represents a context-free grammar.
//
// A context-free grammar is a set of productions, each of which
// consists of a non-terminal symbol and a sequence of symbols.
//
// The non-terminal symbol is the left-hand side of the production,
// and the sequence of symbols is the right-hand side of the production.
//
// The grammar also contains a set of symbols, which are the
// non-terminal and terminal symbols in the grammar.
type Grammar struct {
	// Productions is a slice of Productions in the grammar.
	Productions []Productioner

	// LhsToSkip is a slice of productions to skip.
	LhsToSkip []string

	// Symbols is a slice of Symbols in the grammar.
	Symbols []string
}

// String is a method of Grammar that returns a string representation
// of a Grammar.
//
// It should only be used for debugging and logging purposes.
//
// Returns:
//
//   - string: A string representation of a Grammar.
func (g *Grammar) String() string {
	if g == nil {
		return "Grammar[nil]"
	}

	var builder strings.Builder

	builder.WriteString("Grammar[prouctions=[")

	if len(g.Productions) != 0 {
		fmt.Fprintf(&builder, "%v", g.Productions[0])

		for _, production := range g.Productions[1:] {
			fmt.Fprintf(&builder, ", %v", production)
		}
	}

	builder.WriteString("], symbols=[")

	if len(g.Symbols) != 0 {
		fmt.Fprintf(&builder, "%v", g.Symbols[0])

		for _, symbol := range g.Symbols[1:] {
			fmt.Fprintf(&builder, ", %v", symbol)
		}
	}

	builder.WriteString("], skipProductions=[")

	if len(g.LhsToSkip) != 0 {
		fmt.Fprintf(&builder, "%v", g.LhsToSkip[0])

		for _, production := range g.LhsToSkip[1:] {
			fmt.Fprintf(&builder, ", %v", production)
		}
	}

	builder.WriteString("]]")

	return builder.String()
}

// MatchedResult represents the result of a match operation.
type MatchedResult struct {
	// Matched is the matched token.
	Matched Tokener

	// RuleIndex is the index of the production that matched.
	RuleIndex int
}

// NewMatchResult is a constructor of MatchedResult.
//
// Parameters:
//
//   - matched: The matched token.
//   - ruleIndex: The index of the production that matched.
//
// Returns:
//
//   - MatchedResult: A new MatchedResult.
func NewMatchResult(matched Tokener, ruleIndex int) MatchedResult {
	return MatchedResult{Matched: matched, RuleIndex: ruleIndex}
}

// Match returns a slice of MatchedResult that match the input token.
//
// Parameters:
//
//   - at: The position in the input string.
//   - b: The input stream to match. Refers to Productioner.Match.
//
// Returns:
//
//   - []MatchedResult: A slice of MatchedResult that match the input token.
func (g *Grammar) Match(at int, b any) []MatchedResult {
	matches := make([]MatchedResult, 0)

	for i, p := range g.Productions {
		matched := p.Match(at, b)
		if matched != nil {
			matches = append(matches, NewMatchResult(matched, i))
		}
	}

	return matches
}
