package Lexer

import (
	"bytes"
	"strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// Lexer is a lexer that uses a grammar to tokenize a string.
type Lexer[T uc.Enumer] struct {
	// productions are the production rules to use.
	productions []*gr.RegProduction[T]

	// toSkip are the tokens to skip.
	toSkip []T
}

// NewLexer creates a new lexer.
//
// Parameters:
//   - grammar: The grammar to use.
//
// Returns:
//   - Lexer: The new lexer.
//
// Example:
//
//	lexer, err := NewLexer(grammar)
//	if err != nil {
//	    // Handle error.
//	}
//
//	iter := lexer.Lex([]byte("1 + 2"))
//
//	var branch *cds.Stream[gr.Token]
//	var err error
//
//	for {
//	    branch, err = iter.Consume()
//	    if err != nil {
//	        break
//	    }
//
//	    // Parse the branch.
//	}
//
//	ok := IsDone(err)
//	if !ok {
//	    // Handle error.
//	}
//
//	// Finished successfully.
func NewLexer[T uc.Enumer](grammar *Grammar[T]) *Lexer[T] {
	lex := new(Lexer[T])

	if grammar == nil {
		return lex
	}

	lex.productions = grammar.GetRegexProds()
	lex.toSkip = grammar.GetToSkip()

	return lex
}

// Lex is the main function of the lexer. This can be parallelized.
//
// Parameters:
//   - source: The source to lex.
//   - logger: A verbose logger.
//
// Returns:
//   - Lexer: The active lexer. Nil if there are no tokens to lex or
//     grammar is invalid.
//
// Errors:
//   - *ErrNoTokensToLex: There are no tokens to lex.
//   - *ErrNoMatches: No matches are found in the source.
//   - *ErrAllMatchesFailed: All matches failed.
//   - *gr.ErrNoProductionRulesFound: No production rules are found in the grammar.
func (l *Lexer[T]) Lex(input []byte, logger *Verbose) *LexerIterator[T] {
	prodCopy := make([]*gr.RegProduction[T], len(l.productions))
	copy(prodCopy, l.productions)
	toSkip := make([]T, len(l.toSkip))
	copy(toSkip, l.toSkip)

	stream := cds.NewStream(input)

	si := newSourceIterator(stream, prodCopy, logger)

	li := &LexerIterator[T]{
		toSkip:     toSkip,
		sourceIter: si,
		completedLeaves: &leavesResult[T]{
			leaves: nil,
		},
	}

	return li
}

// FullLexer is a convenience function that creates a new lexer, lexes the content,
// and returns the token streams.
//
// Parameters:
//   - grammar: The grammar to use.
//   - input: The input to lex.
//   - logger: A verbose logger.
//
// Returns:
//   - *LexerIterator: The lexer iterator.
func FullLexer[T uc.Enumer](grammar *Grammar[T], input []byte, logger *Verbose) *LexerIterator[T] {
	lexer := NewLexer(grammar)
	if lexer == nil {
		return nil
	}

	iter := lexer.Lex(input, logger)

	return iter
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
func FormatSyntaxError[T uc.Enumer](branch *cds.Stream[*gr.Token[T]], data []byte) string {
	if branch == nil {
		return string(data)
	}

	items := branch.GetItems()
	lastToken := items[len(items)-2]

	firstInvalid := lastToken.At + len(lastToken.Data.(string))
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
