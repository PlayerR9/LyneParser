package Highlighter

import (
	"fmt"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	"github.com/gdamore/tcell"

	dtt "github.com/PlayerR9/MyGoLib/Safe/DtTable"

	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
)

// Highlighter is a highlighter that applies styles to tokens.
type Highlighter struct {
	// rules is a map of rules to apply.
	rules map[string]tcell.Style

	// errorStyle is the style to apply to errors.
	errorStyle tcell.Style
}

// Apply is a method of Highlighter that applies the rules to the input stream.
//
// Parameters:
//   - inputStream: The input stream to apply the rules to.
//
// Returns:
//   - []ds.DtCell: The cells with the applied styles.
//   - error: An error if the rules could not be applied.
func (h *Highlighter) Apply(inputStream *cds.Stream[*gr.LeafToken]) (*HighlightedData, error) {
	result := NewHighlightedData()

	tokens := inputStream.GetItems()

	for _, token := range tokens {
		style, ok := h.rules[token.ID]
		if !ok {
			return result, fmt.Errorf("no style found for token ID %s", token.ID)
		}

		for _, c := range token.Data {
			result.AppendCell(dtt.NewDtCell(c, style))
		}
	}

	return result, nil
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
func (h *Highlighter) SyntaxError(before, after *cds.Stream[*gr.LeafToken], invalid []byte) (*HighlightedData, error) {
	// 1. Highlight the tokens and split them into lines.
	beforeHighlight, err := h.Apply(before)
	if err != nil {
		return nil, err
	}
	bhs, leftBhs := beforeHighlight.FieldsLine()

	errHighlight := new(HighlightedData).FromBytes(invalid, h.errorStyle)
	ehs, leftEhs := errHighlight.FieldsLine()

	afterHighlight, err := h.Apply(after)
	if err != nil {
		return nil, err
	}
	ahs, leftAhs := afterHighlight.FieldsLine()

	// 2. Initialize the result.
	result := &HighlightedData{
		data: make([]*dtt.DtCell, 0),
	}

	var faultyLine *HighlightedData
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

	return result, nil
}
