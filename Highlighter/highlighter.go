package Highlighter

import (
	"unicode"

	"github.com/gdamore/tcell"

	com "github.com/PlayerR9/LyneParser/Common"
	lx "github.com/PlayerR9/LyneParser/Lexer"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// Highlighter is a highlighter that applies styles to tokens.
type Highlighter struct {
	// rules is a map of rules to apply.
	rules map[string]tcell.Style

	// defaultStyle is the default style to apply.
	defaultStyle tcell.Style

	// errorStyle is the style to apply to errors.
	errorStyle tcell.Style

	// data is the highlighted data.
	data *Data

	// lexer is the lexer to use.
	lexer *lx.Lexer

	// source is the source to use.
	source []byte
}

// NewHighlighter creates a new Highlighter.
//
// Returns:
//   - *Highlighter: The new Highlighter.
func NewHighlighter(lexer *lx.Lexer, defaultStyle tcell.Style) (*Highlighter, error) {
	if lexer == nil {
		return nil, ers.NewErrNilParameter("lexer")
	}

	return &Highlighter{
		rules:        make(map[string]tcell.Style),
		defaultStyle: defaultStyle,
		errorStyle:   defaultStyle,
	}, nil
}

// SpecifyRule adds a rule to the highlighter.
//
// Parameters:
//   - style: The style to apply.
//   - ids: The IDs to apply the style to.
func (h *Highlighter) SpecifyRule(style tcell.Style, ids ...string) {
	if h.rules == nil {
		h.rules = make(map[string]tcell.Style)
	}

	for _, id := range ids {
		h.rules[id] = style
	}
}

// ChangeErrorStyle sets the error style.
//
// Parameters:
//   - style: The style to apply to errors.
func (h *Highlighter) ChangeErrorStyle(style tcell.Style) {
	h.errorStyle = style
}

func (h *Highlighter) extractErrorSection(source *com.ByteStream, firstInvalid int) int {
	// go until the first whitespace character
	bytes := source.GetItems()

	for i := firstInvalid; i < len(bytes); i++ {
		if unicode.IsSpace(rune(bytes[i])) {
			return i
		}
	}

	return -1
}

func (h *Highlighter) makeData() *Data {
	return &Data{
		elems:        make([]Texter, 0),
		source:       h.source,
		rules:        h.rules,
		defaultStyle: h.defaultStyle,
		errorStyle:   h.errorStyle,
	}
}

func (h *Highlighter) Apply(source *com.ByteStream) {
	h.source = source.GetItems()
	h.data = h.makeData()

	for {
		hasError := h.lexer.Lex(source) != nil

		tokens, err := h.lexer.GetTokens()
		if err != nil {
			panic(err)
		} else if len(tokens) == 0 {
			break
		}

		// Find the most ideal token stream to use
		// As of now, we will use the first token stream
		h.data.Add(NewValidText(tokens[0].GetItems()))

		if !hasError {
			break
		}

		items := tokens[0].GetItems()
		lastItem := items[len(items)-1]

		firstInvalid := lastItem.At + len(lastItem.Data)

		// go until the first whitespace character
		bytes := source.GetItems()

		indexOfWS := h.extractErrorSection(source, firstInvalid)
		if indexOfWS == -1 {
			// Anything else is an error
			h.data.Add(NewErrorText(bytes[firstInvalid:]))

			return
		}

		// Extract the error section
		h.data.Add(NewErrorText(bytes[firstInvalid:indexOfWS]))

		// Create a new token stream for the rest of the data
		source = com.NewSourceStream(bytes[indexOfWS:])
	}
}

// apply applies the rules to the input stream using the source for context.
//
// Parameters:
//   - stream: The stream to apply the rules to.
//   - source: The source to apply the rules to.
//
// Returns:
//   - error: An error if the rules could not be applied.
func (h *Highlighter) apply(stream *com.TokenStream, source []byte) error {
	atSource := 0

	for at := 0; ; at++ {
		token, err := stream.GetOne(at)
		if err != nil {
			break
		}

		nextAtToken := token.At

		if atSource < nextAtToken {
			h.data.Append(string(source[atSource:nextAtToken]), h.defaultStyle)
			atSource = nextAtToken
		}

		style, ok := h.rules[token.ID]
		if !ok {
			style = h.defaultStyle
		}

		h.data.Append(token.Data, style)
		atSource += len(token.Data)
	}

	return nil
}

// GetHighlight returns the highlighted data.
//
// Returns:
//   - *HighlightedData: The highlighted data.
func (h *Highlighter) GetHighlight() *ValidText {
	if h.data == nil {
		h.data = NewText()
	}

	return h.data
}

/*

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
func (h *Highlighter) SyntaxError(before, after *com.TokenStream, invalid []byte) {
	// 1. Highlight the tokens and split them into lines.
	h.Apply(before)
	beforeHighlight := h.data

	bhs, leftBhs := beforeHighlight.FieldsLine()

	errHighlight := new(Text).FromBytes(invalid, h.errorStyle)
	ehs, leftEhs := errHighlight.FieldsLine()

	h.Apply(after)
	afterHighlight := h.data

	ahs, leftAhs := afterHighlight.FieldsLine()

	// 2. Initialize the result.
	result := NewText()

	var faultyLine *Text
	faultIndex := 0

	// 2. Write all lines before the one containing the invalid token.
	if !leftBhs.IsEmpty() {
		for _, line := range bhs {
			result.Merge(line)
		}

		faultyLine = leftBhs
		faultIndex = faultyLine.Size()

		if len(ehs) != 0 {
			result.Merge(ehs[0])
			ehs = ehs[1:]
		}
	} else {
		for _, line := range bhs[:len(bhs)-1] {
			result.Merge(line)
		}

		faultyLine = bhs[len(bhs)-1]
		faultIndex = faultyLine.Size()
	}

	// 3. Write the faulty line.
	result.Merge(faultyLine)

	// 4. Write the caret.
	for i := 0; i < faultIndex; i++ {
		result.AppendCell(nil)
	}

	result.AppendCell(dtt.NewDtCell('^', h.errorStyle))
	result.AppendCell(dtt.NewDtCell('\n', h.errorStyle))

	// 5. Write all invalid tokens.
	for _, line := range ehs[1:] {
		result.Merge(line)
	}

	if !leftEhs.IsEmpty() {
		result.Merge(leftEhs)
		result.Merge(ahs[0])
		ahs = ahs[1:]
	}

	// 6. Write all lines after the one containing the invalid token.
	for _, line := range ahs {
		result.Merge(line)
	}

	result.Merge(leftAhs)

	h.data = result
}
*/
