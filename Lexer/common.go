package Lexer

import (
	"strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	sext "github.com/PlayerR9/MyGoLib/Utility/StringExt"
)

// LexString is a function that, given an input string, returns a slice of tokens.
//
// Parameters:
//
//   - input: The input string.
//
// Returns:
//
//   - [][]gr.LeafToken: A slice of slices of tokens.
//   - error: An error if the input string cannot be lexed.
func LexString(lexer *Lexer, input string) ([][]gr.LeafToken, error) {
	err := lexer.Lex([]byte(input))
	if err != nil {
		return nil, err
	}

	return lexer.GetTokens()
}

// LexBytes is a function that, given an input byte slice, returns a slice of tokens.
//
// Parameters:
//
//   - input: The input byte slice.
//
// Returns:
//
//   - [][]gr.LeafToken: A slice of slices of tokens.
//   - error: An error if the input byte slice cannot be lexed.
func LexBytes(lexer *Lexer, input []byte) ([][]gr.LeafToken, error) {
	err := lexer.Lex(input)
	if err != nil {
		return nil, err
	}

	return lexer.GetTokens()
}

// FormatSyntaxError formats a syntax error in the data.
// The function returns a string with the faulty line and a caret pointing to the invalid token.
//
// Parameters:
//
//   - branch: The branch of tokens to search for.
//   - data: The data to search in.
//
// Returns:
//
//   - string: The formatted syntax error.
//
// Example:
//
//	data := []byte("Hello, word!")
//	branch := []gr.LeafToken{
//	  {Data: "Hello", At: 0},
//	  {Data: ",", At: 6},
//	  {Data: "word", At: 8}, // Invalid token (expected "world")
//	  {Data: "!", At: 12},
//	}
//
//	fmt.Println(FormatSyntaxError(branch, data))
//
// Output:
//
//	Hello, word!
//	       ^
func FormatSyntaxError(branch []gr.LeafToken, data []byte) string {
	firstInvalid := findInvalidTokenIndex(branch, data)
	if firstInvalid == -1 {
		return string(data)
	}

	var builder strings.Builder

	before := data[:firstInvalid]
	after := data[firstInvalid:]

	// Write all lines before the one containing the invalid token

	beforeLines := sext.ByteSplitter(before, '\n')

	if len(beforeLines) > 1 {
		builder.WriteString(sext.JoinBytes(beforeLines[:len(beforeLines)-1], '\n'))
		builder.WriteRune('\n')
	}

	// Write the faulty line
	faultyLine := beforeLines[len(beforeLines)-1]
	afterLines := sext.ByteSplitter(after, '\n')

	builder.WriteString(string(faultyLine))

	if len(afterLines) > 0 {
		builder.WriteString(string(afterLines[0]))
	}

	builder.WriteRune('\n')

	// Write the caret
	builder.WriteString(strings.Repeat(" ", len(faultyLine)))
	builder.WriteRune('^')
	builder.WriteRune('\n')

	if len(afterLines) > 1 {
		builder.WriteString(sext.JoinBytes(afterLines[1:], '\n'))
	}

	return builder.String()
}
