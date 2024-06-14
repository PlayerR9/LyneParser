package Lexer1

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// LexerMatcher is a struct that represents a lexer matcher.
type LexerMatcher struct {
	// source is the source to match.
	source *cds.Stream[byte]
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
	matched, err := MatchFrom(lm.source, from, Productions)
	if err != nil {
		return nil, err
	}

	return matched, nil
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
	weights := us.ApplyWeightFunc(matches, MatchWeightFunc)
	pairs := us.FilterByPositiveWeight(weights)

	return us.ExtractResults(pairs)
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
