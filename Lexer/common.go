package Lexer

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	teval "github.com/PlayerR9/MyGoLib/TreeLike/Explorer"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

var (
	// MatchWeightFunc is a weight function that returns the length of the match.
	//
	// Parameters:
	//   - match: The match to weigh.
	//
	// Returns:
	//   - float64: The weight of the match.
	//   - bool: True if the weight is valid, false otherwise.
	MatchWeightFunc us.WeightFunc[*gr.MatchedResult[*gr.LeafToken]]

	// FilterEmptyBranch is a filter that filters out empty branches.
	//
	// Parameters:
	//   - branch: The branch to filter.
	//
	// Returns:
	//   - bool: True if the branch is not empty, false otherwise.
	FilterEmptyBranch us.PredicateFilter[[]*teval.CurrentEval[*gr.LeafToken]]
)

func init() {
	MatchWeightFunc = func(elem *gr.MatchedResult[*gr.LeafToken]) (float64, bool) {
		return float64(len(elem.Matched.Data)), true
	}

	FilterEmptyBranch = func(branch []*teval.CurrentEval[*gr.LeafToken]) bool {
		return len(branch) != 0
	}
}

// Lex is a shorthand function that creates a new lexer, sets the source, lexes the content,
// and returns the token streams.
//
// Parameters:
//   - lexer: The lexer to use.
//   - input: The input to lex.
//
// Returns:
//   - []*cds.Stream[*LeafToken]: The tokens that have been lexed.
//   - error: An error if lexing fails.
//
// Errors:
//   - *ue.ErrInvalidParameter: The lexer or input is nil.
//   - *ErrNoTokensToLex: There are no tokens to lex.
//   - *ErrNoMatches: No matches are found in the source.
//   - *ErrAllMatchesFailed: All matches failed.
func Lex(lexer *Lexer, input []byte) ([]*cds.Stream[*gr.LeafToken], error) {
	if lexer == nil {
		return nil, ue.NewErrNilParameter("lexer")
	}

	err := lexer.Lex(input)
	if err != nil {
		tokens, _ := lexer.GetTokens()

		return tokens, err
	}

	return lexer.GetTokens()
}

// FullLexer is a convenience function that creates a new lexer, lexes the content,
// and returns the token streams.
//
// Parameters:
//   - grammar: The grammar to use.
//   - input: The input to lex.
//
// Returns:
//   - []*cds.Stream[*LeafToken]: The tokens that have been lexed.
//   - error: An error if lexing fails.
func FullLexer(grammar *gr.LexerGrammar, input []byte) ([]*cds.Stream[*gr.LeafToken], error) {
	lexer, err := NewLexer(grammar)
	if err != nil {
		return nil, err
	}

	err = lexer.Lex(input)
	if err != nil {
		tokens, _ := lexer.GetTokens()

		return tokens, err
	}

	return lexer.GetTokens()
}

// MatchFrom matches the source stream from a given index with a list of production rules.
//
// Parameters:
//   - s: The source stream to match.
//   - from: The index to start matching from.
//   - ps: The production rules to match.
//
// Returns:
//   - matches: A slice of MatchedResult that match the input token.
//   - reason: An error if no matches are found.
//
// Errors:
//   - *ue.ErrInvalidParameter: The from index is out of bounds.
//   - *ErrNoMatches: No matches are found.
func MatchFrom(s *cds.Stream[byte], from int, ps []*gr.RegProduction) (matches []*gr.MatchedResult[*gr.LeafToken], reason error) {
	size := s.Size()

	if from < 0 || from >= size {
		reason = ue.NewErrInvalidParameter(
			"from",
			ue.NewErrOutOfBounds(from, 0, size),
		)

		return
	}

	subSet, err := s.Get(from, size)
	if err != nil {
		panic(err)
	}

	for i, p := range ps {
		matched := p.Match(from, subSet)
		if matched != nil {
			matches = append(matches, gr.NewMatchResult(matched, i))
		}
	}

	if len(matches) == 0 {
		reason = NewErrNoMatches()
	}

	return
}
