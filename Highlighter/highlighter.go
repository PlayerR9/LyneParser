package Highlighter

import (
	"fmt"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	"github.com/gdamore/tcell"

	ds "github.com/PlayerR9/LyneParser/Display"
)

// Highlighter is a highlighter that applies styles to tokens.
type Highlighter struct {
	// rules is a map of rules to apply.
	rules map[string]tcell.Style
}

// Apply is a method of Highlighter that applies the rules to the input stream.
//
// Parameters:
//   - inputStream: The input stream to apply the rules to.
//
// Returns:
//   - []ds.DtCell: The cells with the applied styles.
//   - error: An error if the rules could not be applied.
func (h *Highlighter) Apply(inputStream gr.TokenStream) ([]ds.DtCell, error) {
	result := make([]ds.DtCell, 0)

	tokens := inputStream.GetTokens()

	for _, token := range tokens {
		style, ok := h.rules[token.ID]
		if !ok {
			return result, fmt.Errorf("no style found for token ID %s", token.ID)
		}

		for _, c := range token.Data {
			result = append(result, ds.NewDtCell(c, style))
		}
	}

	return result, nil
}
