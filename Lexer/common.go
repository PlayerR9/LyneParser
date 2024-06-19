package Lexer

import (
	"bytes"
	"errors"
	"strings"
	"sync"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

// Lex is the main function of the lexer. This can be parallelized.
//
// Parameters:
//   - lexer: The lexer to use.
//   - source: The source to lex.
//
// Returns:
//   - error: An error if lexing fails.
//
// Errors:
//   - *ErrNoTokensToLex: There are no tokens to lex.
//   - *ErrNoMatches: No matches are found in the source.
//   - *ErrAllMatchesFailed: All matches failed.
//   - *gr.ErrNoProductionRulesFound: No production rules are found in the grammar.
func Lex(lexer *Lexer, input []byte) ([]*cds.Stream[*gr.LeafToken], error) {
	if lexer == nil {
		return nil, ue.NewErrNilParameter("lexer")
	}

	if len(input) == 0 {
		return nil, errors.New("no tokens to lex")
	}

	lexer.mu.RLock()

	if len(lexer.productions) == 0 {
		lexer.mu.RUnlock()
		return nil, gr.NewErrNoProductionRulesFound()
	}

	prodCopy := make([]*gr.RegProduction, len(lexer.productions))
	copy(prodCopy, lexer.productions)
	toSkip := make([]string, len(lexer.toSkip))
	copy(toSkip, lexer.toSkip)

	lexer.mu.RUnlock()

	stream := cds.NewStream(input)
	tree, err := executeLexing(stream, prodCopy)
	if err != nil {
		tokenBranches, _ := getTokens(tree, toSkip)
		return tokenBranches, err
	}

	tokenBranches, err := getTokens(tree, toSkip)
	return tokenBranches, err
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
func FullLexer(grammar *Grammar, input []byte) ([]*cds.Stream[*gr.LeafToken], error) {
	if grammar == nil {
		return nil, ue.NewErrNilParameter("grammar")
	}

	productions := grammar.GetRegexProds()
	toSkip := grammar.GetToSkip()

	if len(productions) == 0 {
		return nil, gr.NewErrNoProductionRulesFound()
	}

	stream := cds.NewStream(input)
	tree, err := executeLexing(stream, productions)
	if err != nil {
		tokenBranches, _ := getTokens(tree, toSkip)
		return tokenBranches, err
	}

	tokenBranches, err := getTokens(tree, toSkip)
	return tokenBranches, err
}

// Lexer is a lexer that uses a grammar to tokenize a string.
type Lexer struct {
	// grammar is the grammar used by the lexer.
	productions []*gr.RegProduction

	// toSkip is a list of LHSs to skip.
	toSkip []string

	// mu is a mutex to protect the lexer.
	mu sync.RWMutex
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
//	branches, err := lexer.Lex(lexer, []byte("1 + 2"))
//	if err != nil {
//	    // Handle error.
//	}
//
// // Continue with parsing.
func NewLexer(grammar *Grammar) *Lexer {
	if grammar == nil {
		return &Lexer{
			productions: nil,
			toSkip:      nil,
		}
	}

	lex := &Lexer{
		productions: grammar.GetRegexProds(),
		toSkip:      grammar.GetToSkip(),
	}

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
