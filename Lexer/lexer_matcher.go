package Lexer

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	ers "github.com/PlayerR9/MyGoLib/Units/errors"

	com "github.com/PlayerR9/LyneParser/Common"
	hlp "github.com/PlayerR9/MyGoLib/Utility/Helpers"
)

// LexerMatcher is a struct that represents a lexer matcher.
type LexerMatcher struct {
	// source is the source to match.
	source *com.ByteStream

	// productions is the list of productions to match.
	productions []*gr.RegProduction
}

// IsDone is a function that checks if the matcher is done.
//
// Parameters:
//   - from: The starting position of the match.
//
// Returns:
//   - bool: True if the matcher is done, false otherwise.
func (lm *LexerMatcher) IsDone(from int) bool {
	return lm.source.IsDone(from, 1)
}

// Match is a function that matches the element.
//
// Parameters:
//   - from: The starting position of the match.
//
// Returns:
//   - []Matcher: The list of matchers.
//   - error: An error if the matchers cannot be created.
func (lm *LexerMatcher) Match(from int) ([]*gr.MatchedResult[*gr.LeafToken], error) {
	return lm.source.MatchFrom(from, lm.productions)
}

// SelectBestMatches selects the best matches from the list of matches.
// Usually, the best matches' euristic is the longest match.
//
// Parameters:
//   - matches: The list of matches.
//
// Returns:
//   - []Matcher: The best matches.
func (lm *LexerMatcher) SelectBestMatches(matches []*gr.MatchedResult[*gr.LeafToken]) []*gr.MatchedResult[*gr.LeafToken] {
	weights := hlp.ApplyWeightFunc(matches, MatchWeightFunc)
	pairs := hlp.FilterByPositiveWeight(weights)

	return hlp.ExtractResults(pairs)
}

// GetNext is a function that returns the next position of an element.
//
// Parameters:
//   - elem: The element to get the next position of.
//
// Returns:
//   - int: The next position of the element.
func (lm *LexerMatcher) GetNext(elem *gr.LeafToken) int {
	return elem.GetPos() + len(elem.Data)
}

// NewLexerMatcher creates a new lexer matcher.
//
// Parameters:
//   - source: The source to match.
//
// Returns:
//   - *LexerMatcher: The new lexer matcher.
//   - error: An error if the source is nil.
func NewLexerMatcher(source *com.ByteStream) (*LexerMatcher, error) {
	if source == nil {
		return nil, ers.NewErrNilParameter("source")
	}

	return &LexerMatcher{
		source: source,
	}, nil
}
