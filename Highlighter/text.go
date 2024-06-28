package Highlighter

import (
	"unicode/utf8"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	"github.com/gdamore/tcell"
)

// Texter is a text.
type Texter interface {
}

// NormalText is a highlighted text.
type NormalText struct {
	// data is the data of the highlighted data.
	data []rune

	// style is the style to apply.
	style tcell.Style
}

// NewNormalText creates a new highlighted data.
//
// Parameters:
//   - data: The data to create the highlighted data from.
//   - style: The style to apply.
//
// Returns:
//   - *NormalText: The new highlighted data.
func NewNormalText(data []byte, style tcell.Style) *NormalText {
	runes := make([]rune, 0, len(data))

	for len(data) > 0 {
		r, size := utf8.DecodeRune(data)
		data = data[size:]

		runes = append(runes, r)
	}

	nt := &NormalText{
		data:  runes,
		style: style,
	}

	return nt
}

// Runes returns the runes of the text.
//
// Returns:
//   - []rune: The runes of the text.
func (hd *NormalText) Runes() []rune {
	return hd.data
}

// ValidText is a highlighted text.
type ValidText struct {
	// data is the data of the highlighted data.
	data []gr.Token
}

// NewValidText creates a new highlighted data.
//
// Parameters:
//   - tokens: The tokens to create the highlighted data from.
//
// Returns:
//   - *ValidText: The new highlighted data.
//   - error: An error of type *uc.ErrInvalidParameter if the tokens are empty.
func NewValidText(tokens []gr.Token) (*ValidText, error) {
	if len(tokens) == 0 {
		return nil, uc.NewErrInvalidParameter(
			"tokens",
			uc.NewErrEmpty(tokens),
		)
	}

	vt := &ValidText{
		data: tokens,
	}

	return vt, nil
}

/*
// Draw is a method of cdd.TableDrawer that draws the unit to the table at the given x and y
// coordinates.
//
// Parameters:
//   - table: The table to draw the unit to.
//   - x: The x coordinate to draw the unit at.
//   - y: The y coordinate to draw the unit at.
//
// Returns:
//   - error: An error of type *uc.ErrInvalidParameter if the table is nil.
//
// Behaviors:
//   - Any value that would be drawn outside of the table is not drawn.
//   - Assumes that the table is not nil.
func (hd *ValidText) Draw(table *cdd.DrawTable, x, y *int) error {
	height := table.GetHeight()

	offsetX := *x

	for i, line := range hd.data {
		if *y >= height {
			break
		} else {
			offsetX = *x
		}

		sequence := make([]*cdd.ColoredUnit, 0)

		for j, cell := range line {
			sequence = append(sequence, cdd.NewColoredUnit(cell, hd.styles[i][j]))
		}

		table.WriteHorizontalSequence(&offsetX, y, sequence)
	}

	*x = offsetX

	return nil
}

// Append appends a rune and its style to the highlighted data.
//
// Parameters:
//   - str: The string to append.
//   - style: The style to append.
func (hd *ValidText) Append(str string, style tcell.Style) {
	if str == "" {
		return
	}

	for _, r := range str {
		if r == '\n' {
			hd.lastLine++

			hd.data = append(hd.data, make([]rune, 0))
			hd.styles = append(hd.styles, make([]tcell.Style, 0))
		} else {
			hd.data[hd.lastLine] = append(hd.data[hd.lastLine], r)
			hd.styles[hd.lastLine] = append(hd.styles[hd.lastLine], style)
		}
	}
}

/////////////////////////


func (hd *Text) String() string {
	if hd == nil || len(hd.data) == 0 {
		return ""
	}

	var builder strings.Builder

	for _, cell := range hd.data {
		if cell == nil {
			builder.WriteRune(' ')
		} else {
			builder.WriteRune(cell.First)
		}
	}

	return builder.String()
}


func (hd *Text) FromBytes(data []byte, style tcell.Style) *Text {
	if len(data) == 0 {
		return NewText()
	}

	hd = &Text{
		data: make([]*dtt.DtCell, 0, len(data)),
	}

	for _, b := range data {
		hd.data = append(hd.data, dtt.NewDtCell(rune(b), style))
	}

	return hd
}



func (hd *Text) Size() int {
	return len(hd.data)
}

func (hd *Text) IsEmpty() bool {
	return len(hd.data) == 0
}

func (hd *Text) FieldsLine() ([]*Text, *Text) {
	result := make([]*Text, 0)
	line := NewText()

	for _, cell := range hd.data {
		line.data = append(line.data, cell)

		if cell.First == '\n' {
			result = append(result, line)

			line = NewText()
		}
	}

	return result, line
}

func (hd *Text) Merge(other *Text) {
	if other == nil {
		return
	}

	hd.data = append(hd.data, other.data...)
}
*/
