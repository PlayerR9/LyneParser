package Lexer

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// Lexer is a lexer that uses a grammar to tokenize a string.
type Lexer[T gr.TokenTyper] struct {
	// productions are the production rules to use.
	productions []*gr.RegProduction[T]

	// to_skip are the tokens to skip.
	to_skip []T
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
func NewLexer[T gr.TokenTyper](grammar *Grammar[T]) *Lexer[T] {
	lex := new(Lexer[T])

	if grammar == nil {
		return lex
	}

	lex.productions = grammar.GetRegexProds()
	lex.to_skip = grammar.GetToSkip()

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
	prod_copy := make([]*gr.RegProduction[T], len(l.productions))
	copy(prod_copy, l.productions)
	to_skip := make([]T, len(l.to_skip))
	copy(to_skip, l.to_skip)

	stream := cds.NewStream(input)

	si := newSourceIterator(stream, prod_copy, logger)

	lr := &leaves_result[T]{
		leaves: nil,
	}

	li := &LexerIterator[T]{
		to_skip:          to_skip,
		source_iter:      si,
		completed_leaves: lr,
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
func FullLexer[T gr.TokenTyper](grammar *Grammar[T], input []byte, logger *Verbose) *LexerIterator[T] {
	lexer := NewLexer(grammar)
	if lexer == nil {
		return nil
	}

	iter := lexer.Lex(input, logger)

	return iter
}

// index_of_newline_char is a helper function to find the index of a newline character.
//
// Parameters:
//   - data: The data.
//   - from: The starting index.
//
// Returns:
//   - int: The index of the newline character. -1 if not found or data is empty.
//
// Asserts:
//   - from >= 0 only if data is not empty.
func index_of_newline_char(data []rune, from int) int {
	if len(data) == -1 {
		return -1
	}

	uc.AssertParam("from", from >= 0, uc.NewErrGTE(0))

	for i := from; i < len(data); i++ {
		if data[i] == '\n' {
			return i
		}
	}

	return -1
}

// split_into_lines is a helper function to split data into lines.
//
// Parameters:
//   - data: The data.
//
// Returns:
//   - [][]rune: The lines. Nil if data is empty.
func split_into_lines(data []rune) [][]rune {
	if len(data) == 0 {
		return nil
	}

	var table [][]rune
	var from int

	for from < len(data) {
		index := index_of_newline_char(data, from)

		if index == -1 {
			table = append(table, data[from:])
			break
		}

		extract := data[from:index]

		table = append(table, extract)

		from = index
	}

	return table
}

// write_lines writes lines to a string builder.
//
// Parameters:
//   - b: The string builder.
//   - lines: The lines.
//
// Asserts:
//   - b != nil
func write_lines(b *strings.Builder, lines [][]rune) {
	uc.AssertParam("b", b != nil, errors.New("b should not be nil"))

	if len(lines) == 0 {
		return
	}

	var values []string

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		values = append(values, string(line))
	}

	joined_str := strings.Join(values, "\n")

	b.WriteString(joined_str)
}

// determine_faulty_line determines the faulty line.
//
// Parameters:
//   - data: The data.
//   - end_idx: The end index.
//   - before_line: The line before the end index.
//
// Returns:
//   - []rune: The faulty line.
//   - [][]rune: The lines after the faulty line. Nil if there are no after lines.
//
// Asserts:
//   - len(data) > 0
//   - end_idx >= 0 && end_idx <= len(data)
func determine_faulty_line(data []rune, end_idx int, before_line []rune) ([]rune, [][]rune) {
	uc.AssertParam("data", len(data) > 0, errors.New("data should not be empty"))
	uc.AssertParam("end_idx", end_idx >= 0 && end_idx <= len(data),
		uc.NewErrOutOfBounds(end_idx, 0, len(data)).WithUpperBound(true),
	)

	if end_idx >= len(data) {
		return before_line, nil
	}

	after_lines := split_into_lines(data[end_idx:])
	uc.Assert(len(after_lines) > 0, "after_lines should not be empty")

	var faulty_line []rune

	if len(after_lines[0]) == 0 {
		faulty_line = before_line
	} else {
		faulty_line = append(before_line, after_lines[0]...)
	}

	return faulty_line, after_lines[1:]
}

// LastInstanceOfWS finds the last instance of whitespace in the characters.
//
// Parameters:
//   - chars: The characters.
//   - from_idx: The starting index. (inclusive)
//   - to_idx: The ending index. (exclusive)
//
// Returns:
//   - int: The index of the last whitespace character. -1 if not found.
//
// Behaviors:
//   - If from_idx < 0, from_idx is set to 0.
//   - If to_idx >= len(chars), to_idx is set to len(chars) - 1.
//   - If from_idx > to_idx, from_idx and to_idx are swapped.
//
// FIXME: Remove this function once MyGoLib is updated.
func LastInstanceOfWS(chars []rune, from_idx, to_idx int) int {
	if len(chars) == 0 {
		return -1
	}

	if from_idx < 0 {
		from_idx = 0
	}

	if to_idx >= len(chars) {
		to_idx = len(chars)
	}

	if from_idx > to_idx {
		from_idx, to_idx = to_idx, from_idx
	}

	for i := to_idx - 1; i >= from_idx; i-- {
		ok := unicode.IsSpace(chars[i])
		if ok {
			return i
		}
	}

	return -1
}

// FirstInstanceOfWS finds the first instance of whitespace in the characters.
//
// Parameters:
//   - chars: The characters.
//   - from_idx: The starting index. (inclusive)
//   - to_idx: The ending index. (exclusive)
//
// Returns:
//   - int: The index of the first whitespace character. -1 if not found.
//
// Behaviors:
//   - If from_idx < 0, from_idx is set to 0.
//   - If to_idx >= len(chars), to_idx is set to len(chars) - 1.
//   - If from_idx > to_idx, from_idx and to_idx are swapped.
//
// FIXME: Remove this function once MyGoLib is updated.
func FirstInstanceOfWS(chars []rune, from_idx, to_idx int) int {
	if len(chars) == 0 {
		return -1
	}

	if from_idx < 0 {
		from_idx = 0
	}

	if to_idx >= len(chars) {
		to_idx = len(chars)
	}

	if from_idx > to_idx {
		from_idx, to_idx = to_idx, from_idx
	}

	for i := from_idx; i < to_idx; i++ {
		ok := unicode.IsSpace(chars[i])
		if ok {
			return i
		}
	}

	return -1
}

// write_arrow writes an arrow from the from index to the to index.
//
// Parameters:
//   - b: The string builder.
//   - from_idx: The from index.
//   - to_idx: The to index.
//
// Asserts:
//   - b != nil
//   - to_idx >= from_idx
func write_arrow(b *strings.Builder, from_idx, to_idx int) {
	uc.AssertParam("b", b != nil, errors.New("b should not be nil"))
	uc.AssertParam("to_idx", to_idx >= from_idx, errors.New("to_idx should be greater than from_idx"))

	for i := 0; i < from_idx; i++ {
		b.WriteRune(' ')
	}

	for i := from_idx; i < to_idx-1; i++ {
		b.WriteRune('^')
	}
}

// highlight_line highlights a line.
//
// Parameters:
//   - b: The string builder.
//   - faulty_line: The faulty line.
//   - idx: The index to highlight.
//
// Asserts:
//   - b != nil
//   - len(faulty_line) > 0
//   - idx >= 0 && idx <= len(faulty_line)
func highlight_line(b *strings.Builder, faulty_line []rune, idx int) {
	uc.AssertParam("b", b != nil, errors.New("b should not be nil"))
	uc.AssertParam("faulty_line", len(faulty_line) > 0, uc.NewErrEmpty(faulty_line))
	uc.AssertParam("idx", idx >= 0 && idx <= len(faulty_line), uc.NewErrOutOfBounds(idx, 0, len(faulty_line)).WithUpperBound(true))

	left_idx := LastInstanceOfWS(faulty_line, 0, idx)

	right_idx := FirstInstanceOfWS(faulty_line, idx, len(faulty_line))
	if right_idx == -1 {
		right_idx = len(faulty_line)
	}

	b.WriteString(string(faulty_line))
	b.WriteString("\n")

	write_arrow(b, left_idx+1, right_idx) // +1 to exclude the ws character.

	b.WriteString("\n")
}

// PrintCode prints the code but highlights the faulty line.
//
// A faulty line is defined as the line containing the last token as, supposedly, the last token
// is the one that caused the error.
//
// Parameters:
//   - data: The original data read.
//   - tokens: The tokens lexed.
//
// Returns:
//   - string: The formatted code.
//
// Example:
//
//	data := []rune("Hello, word!")
//	tokens := []*gr.Token{
//	  {ID: TkWord, Data: "Hello", At: 0},
//	  {ID: TkComma, Data: ",", At: 6},
//	  {ID: TkWord, Data: "word", At: 8}, // Invalid token (expected "world")
//	  {ID: gr.TkEof, Data: nil, At: -1},
//	}
//
//	str := PrintCode(data, tokens)
//	fmt.Println(str)
//
// Output:
//
//	Hello, word!
//	       ^^^^
func PrintCode[T gr.TokenTyper](data []rune, tokens []*gr.Token[T]) string {
	if len(data) == 0 {
		return ""
	}

	var builder strings.Builder

	if len(tokens) < 2 {
		after_lines := split_into_lines(data)
		// uc.Assert(len(after_lines) > 0, "after_lines should not be empty")

		if len(after_lines[0]) > 0 {
			highlight_line(&builder, after_lines[0], 0)

			after_lines = after_lines[1:]
		}

		write_lines(&builder, after_lines)
	} else {
		tokens = tokens[:len(tokens)-1] // Ignore the 'Enf of File' token.

		last_token := tokens[len(tokens)-1]
		tok_size := utf8.RuneCountInString(last_token.Data.(string))

		end_idx := last_token.At + tok_size

		before_lines := split_into_lines(data[:end_idx])

		var last_line []rune

		if len(before_lines) > 0 {
			valid_lines := before_lines[:len(before_lines)-1]

			write_lines(&builder, valid_lines)

			last_line = before_lines[len(before_lines)-1]
		} else {
			last_line = before_lines[0]
		}

		faulty_line, after_lines := determine_faulty_line(data, end_idx, last_line)

		highlight_line(&builder, faulty_line, end_idx)

		write_lines(&builder, after_lines)
	}

	str := builder.String()

	return str
}
