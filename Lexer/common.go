package Lexer

import (
	"bytes"
	"strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/StatusTree"
)

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
func FullLexer(grammar *Grammar, input []byte) *Lexer {
	lex := new(Lexer)

	if grammar == nil {
		return lex
	}

	lex.productions = grammar.GetRegexProds()
	lex.toSkip = grammar.GetToSkip()

	if len(input) == 0 || len(lex.productions) == 0 {
		return nil
	}

	lex.source = cds.NewStream(input)

	rootNode := gr.NewRootToken()

	lex.tree = tr.NewTreeWithHistory(EvalIncomplete, rootNode)
	lex.isFirst = true
	lex.canContinue = true

	return lex
}

// writeArrow writes an arrow pointing to the position in the data.
//
// Parameters:
//   - pos: The position to write the arrow to.
//
// Returns:
//   - string: The arrow.
func writeArrow(faultyLine []byte) string {
	var builder strings.Builder

	for _, char := range faultyLine {
		if char == '\t' {
			builder.WriteRune('\t')
		} else {
			builder.WriteRune(' ')
		}
	}

	builder.WriteRune('^')

	return builder.String()
}

// splitBytesFromEnd splits the data into two parts from the end.
//
// Parameters:
//   - data: The data to split.
//
// Returns:
//   - [][]byte: The split data.
func splitBytesFromEnd(data []byte) [][]byte {
	for i := len(data) - 1; i >= 0; i-- {
		if data[i] == '\n' {
			return [][]byte{data[:i], data[i:]}
		}
	}

	return [][]byte{data}
}

// FormatSyntaxError formats a syntax error in the data.
// The function returns a string with the faulty line and a caret pointing to the invalid token.
//
// Parameters:
//   - branch: The branch of tokens to search for.
//   - data: The data to search in.
//
// Returns:
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
func FormatSyntaxError(branch *cds.Stream[*gr.LeafToken], data []byte) string {
	if branch == nil {
		return string(data)
	}

	items := branch.GetItems()
	lastToken := items[len(items)-2]

	firstInvalid := lastToken.At + len(lastToken.Data)
	if firstInvalid == len(data) {
		return string(data)
	}

	var builder strings.Builder

	before := data[:firstInvalid]
	after := data[firstInvalid:]

	// Write all lines before the one containing the invalid token.
	beforeLines := splitBytesFromEnd(before)

	var faultyLine []byte

	if len(beforeLines) > 1 {
		builder.Write(beforeLines[0])

		faultyLine = beforeLines[1]
	} else {
		faultyLine = beforeLines[0]
	}

	builder.Write(faultyLine)

	afterLines := bytes.SplitN(after, []byte{'\n'}, 2)
	if len(afterLines) > 0 {
		builder.Write(afterLines[0])
	}

	// Write the caret.
	arrow := writeArrow(faultyLine)
	builder.WriteRune('\n')
	builder.WriteString(arrow)

	if len(afterLines) == 0 {
		return builder.String()
	}

	builder.WriteRune('\n')

	builder.Write(afterLines[1])

	return builder.String()
}
