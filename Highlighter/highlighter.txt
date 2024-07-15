package Highlighter

import (
	"unicode"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	lx "github.com/PlayerR9/LyneParser/Lexer"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	"github.com/gdamore/tcell"
)

// Highlighter is a highlighter that applies styles to tokens.
type Highlighter[T gr.TokenTyper] struct {
	// rules is a map of rules to apply.
	rules map[T]tcell.Style

	// default_style is the default style to apply.
	default_style tcell.Style

	// error_style is the style to apply to errors.
	error_style tcell.Style

	// data is the highlighted data.
	data *Data[T]

	// lexer is the lexer to use.
	lexer *lx.Lexer[T]

	// source is the source to use.
	source []byte
}

// NewHighlighter creates a new Highlighter.
//
// Returns:
//   - *Highlighter: The new Highlighter.
func NewHighlighter[T gr.TokenTyper](lexer *lx.Lexer[T], default_style tcell.Style) (*Highlighter[T], error) {
	if lexer == nil {
		return nil, uc.NewErrNilParameter("lexer")
	}

	h := &Highlighter[T]{
		rules:         make(map[T]tcell.Style),
		default_style: default_style,
		error_style:   default_style,
	}
	return h, nil
}

// SpecifyRule adds a rule to the highlighter.
//
// Parameters:
//   - style: The style to apply.
//   - ids: The IDs to apply the style to.
func (h *Highlighter[T]) SpecifyRule(style tcell.Style, ids ...T) {
	if h.rules == nil {
		h.rules = make(map[T]tcell.Style)
	}

	for _, id := range ids {
		h.rules[id] = style
	}
}

// ChangeErrorStyle sets the error style.
//
// Parameters:
//   - style: The style to apply to errors.
func (h *Highlighter[T]) ChangeErrorStyle(style tcell.Style) {
	h.error_style = style
}

func (h *Highlighter[T]) extractErrorSection(data []byte, firstInvalid int) int {
	// go until the first whitespace character
	for i := firstInvalid; i < len(data); i++ {
		ok := unicode.IsSpace(rune(data[i]))

		if ok {
			return i
		}
	}

	return -1
}

func (h *Highlighter[T]) makeData() *Data[T] {
	d := &Data[T]{
		elems:         make([]Texter, 0),
		source:        h.source,
		rules:         h.rules,
		default_style: h.default_style,
		error_style:   h.error_style,
	}

	return d
}

func (h *Highlighter[T]) Apply(data []byte) {
	h.data = h.makeData()

	v := lx.NewVerbose(true)
	defer v.Close()

	iter := h.lexer.Lex(data, v)

	for {
		branch, err := iter.Consume()
		var has_error bool

		if err != nil {
			ok := uc.Is[*uc.ErrExhaustedIter](err)
			if !ok {
				has_error = true
			}
		}

		// Find the most ideal token stream to use
		// As of now, we will use the first token stream
		token_items := branch.GetItems()

		txt, err := NewValidText(token_items)
		uc.AssertF(err == nil, "NewValidText failed: %s", err.Error())

		h.data.Add(txt)

		if !has_error {
			break
		}

		token_items = branch.GetItems()
		last_item := token_items[len(token_items)-1]

		first_invalid := last_item.At + len(last_item.Data.(string))

		// go until the first whitespace character
		index_of_ws := h.extractErrorSection(data, first_invalid)
		if index_of_ws == -1 {
			// Anything else is an error
			nt := NewNormalText(data[first_invalid:], h.error_style)
			h.data.Add(nt)

			return
		}

		// Extract the error section
		nt := NewNormalText(data[first_invalid:index_of_ws], h.error_style)
		h.data.Add(nt)

		// Create a new token stream for the rest of the data
		data = data[index_of_ws:]
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
func (h *Highlighter[T]) apply(stream *cds.Stream[*gr.Token[T]], source []byte) error {
	at_source := 0

	for at := 0; ; at++ {
		token, err := stream.GetOne(at)
		if err != nil {
			break
		}

		next_at_token := token.At

		if at_source < next_at_token {
			h.data.Add(NewNormalText(source[at_source:next_at_token], h.default_style))
			at_source = next_at_token
		}

		style, ok := h.rules[token.ID]
		if !ok {
			style = h.default_style
		}

		h.data.Add(NewNormalText([]byte(token.Data.(string)), style))
		at_source += len(token.Data.(string))
	}

	return nil
}

/*

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
func (h *Highlighter) SyntaxError(before, after *cds.Stream[*LeafToken], invalid []byte) {
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
