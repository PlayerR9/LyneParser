package Lexer

import (
	"bytes"
	"strings"
	"sync"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
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

// findInvalidTokenIndex finds the index of the first invalid token in the data.
// The function returns -1 if no invalid token is found.
//
// Parameters:
//   - branch: The branch of tokens to search for.
//   - data: The data to search in.
//
// Returns:
//   - int: The index of the first invalid token.
func findInvalidTokenIndex(branch []*gr.LeafToken, data []byte) int {
	pos := 0

	for _, token := range branch {
		b := []byte(token.Data)

		startIndex := us.FindSubsliceFrom(data, b, pos)
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

// writeArrow writes an arrow pointing to the position in the data.
//
// Parameters:
//   - pos: The position to write the arrow to.
//
// Returns:
//   - string: The arrow.
func writeArrow(pos int) string {
	var builder strings.Builder

	leftStr := strings.Repeat(" ", pos)
	builder.WriteString(leftStr)
	builder.WriteRune('^')

	return builder.String()
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
func FormatSyntaxError(branch *cds.Stream[*gr.LeafToken], data []byte) []string {
	if branch == nil {
		lines := bytes.Split(data, []byte{'\n'})

		var formatted []string
		for _, line := range lines {
			formatted = append(formatted, string(line))
		}

		return formatted
	}

	items := branch.GetItems()

	firstInvalid := findInvalidTokenIndex(items, data)
	if firstInvalid == -1 {
		lines := bytes.Split(data, []byte{'\n'})

		var formatted []string
		for _, line := range lines {
			formatted = append(formatted, string(line))
		}

		return formatted
	}

	var lines []string

	before := data[:firstInvalid]

	after := data[firstInvalid:]

	// Write all lines before the one containing the invalid token.
	beforeLines := bytes.Split(before, []byte{'\n'})
	if len(beforeLines) > 1 {
		joined := bytes.Join(beforeLines[:len(beforeLines)-1], []byte{'\n'})

		lines = append(lines, string(joined))
	}

	// Write the faulty line.
	faultyLine := beforeLines[len(beforeLines)-1]

	afterLines := bytes.Split(after, []byte{'\n'})

	var builder strings.Builder

	builder.Write(faultyLine)

	if len(afterLines) > 0 {
		builder.Write(afterLines[0])
	}

	lines = append(lines, builder.String())

	// Write the caret.
	arrow := writeArrow(len(faultyLine))
	lines = append(lines, arrow)

	if len(afterLines) == 0 {
		return lines
	}

	for _, line := range afterLines[1:] {
		lines = append(lines, string(line))
	}

	return lines
}
