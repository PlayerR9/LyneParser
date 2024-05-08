package Parser

import (
	"fmt"

	com "github.com/PlayerR9/LyneParser/Common"
	cs "github.com/PlayerR9/LyneParser/ConflictSolver"
	gr "github.com/PlayerR9/LyneParser/Grammar"
)

// FullParse parses the input stream using the given grammar and decision
// function. It is a convenience function intended for simple parsing tasks.
//
// Parameters:
//
//   - grammar: The grammar that the parser will use.
//   - inputStream: The input stream that the parser will parse.
//   - decisionFunc: The decision function that the parser will use.
//
// Returns:
//
//   - []gr.NonLeafToken: The parse tree.
//   - error: An error if the input stream could not be parsed.
func FullParse(grammar *gr.Grammar, source *com.TokenStream, dt *cs.ConflictSolver) ([]*com.TokenTree, error) {
	parser, err := NewParser(grammar)
	if err != nil {
		return nil, fmt.Errorf("could not create parser: %s", err.Error())
	}

	err = parser.Parse(source)
	if err != nil {
		return nil, fmt.Errorf("parse error: %s", err.Error())
	}

	roots, err := parser.GetParseTree()
	if err != nil {
		return nil, fmt.Errorf("could not get parse tree: %s", err.Error())
	}

	return roots, nil
}

/////////////////////////////////////////////////////////////

const (
	Indentation string = "|  "
)

/*
// findInvalidTokenIndex finds the index of the first invalid token in the data.
// The function returns -1 if no invalid token is found.
//
// Parameters:
//
//   - branch: The branch of tokens to search for.
//   - data: The data to search in.
//
// Returns:
//
//   - int: The index of the first invalid token.
func findInvalidTokenIndex(branch []gr.LeafToken, data []byte) int {
	pos := 0

	for _, token := range branch {
		b := []byte(token.Data)

		startIndex := slext.FindSubsliceFrom(data, b, pos)
		if startIndex == -1 {
			return -1
		}

		pos += startIndex + len(token.Data)
	}

	if pos >= len(data) {
		return -1
	}

	return pos
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
func FormatSyntaxError(root gr.Tokener, data []byte) string {
	return TokenerString(root)

	/*
		root.Data

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
*/
