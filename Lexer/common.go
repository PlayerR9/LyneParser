package Lexer

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// Lex is a shorthand function that creates a new lexer, sets the source, lexes the content,
// and returns the token streams.
//
// Parameters:
//   - lexer: The lexer to use.
//   - input: The input to lex.
//
// Returns:
//   - []*gr.TokenStream: The tokens that have been lexed.
//   - error: An error if lexing fails.
//
// Errors:
//   - *ers.ErrInvalidParameter: The lexer or input is nil.
//   - *ErrNoTokensToLex: There are no tokens to lex.
//   - *ErrNoMatches: No matches are found in the source.
//   - *ErrAllMatchesFailed: All matches failed.
func Lex(lexer *Lexer, input any) ([]*gr.TokenStream, error) {
	if lexer == nil {
		return nil, ers.NewErrNilParameter("lexer")
	}

	err := lexer.Lex(NewSourceStream(input))
	if err != nil {
		return nil, err
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
//   - []*gr.TokenStream: The tokens that have been lexed.
//   - error: An error if lexing fails.
func FullLexer(grammar *gr.Grammar, input any) ([]*gr.TokenStream, error) {
	lexer, err := NewLexer(grammar)
	if err != nil {
		return nil, err
	}

	err = lexer.Lex(NewSourceStream(input))
	tokens, _ := lexer.GetTokens()

	return tokens, err
}
