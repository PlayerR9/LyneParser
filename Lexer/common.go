package Lexer

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
)

// LexString is a function that, given an input string, returns a slice of tokens.
//
// Parameters:
//   - input: The input string.
//
// Returns:
//   - []gr.TokenStream: A slice of slices of tokens.
//   - error: An error if the input string cannot be lexed.
func LexString(lexer *Lexer, input string) ([]*gr.TokenStream, error) {
	lexer.SetSource([]byte(input))

	err := lexer.Lex()
	if err != nil {
		return nil, err
	}

	return lexer.GetTokens()
}

// LexBytes is a function that, given an input byte slice, returns a slice of tokens.
//
// Parameters:
//   - input: The input byte slice.
//
// Returns:
//
//   - []gr.TokenStream: A slice of slices of tokens.
//   - error: An error if the input byte slice cannot be lexed.
func LexBytes(lexer *Lexer, input []byte) ([]*gr.TokenStream, error) {
	lexer.SetSource(input)

	err := lexer.Lex()
	if err != nil {
		return nil, err
	}

	return lexer.GetTokens()
}
