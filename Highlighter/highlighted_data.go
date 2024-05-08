package Highlighter

import (
	"strings"

	dtt "github.com/PlayerR9/MyGoLib/Safe/DtTable"
	"github.com/gdamore/tcell"
)

type HighlightedData struct {
	// data is the data of the highlighted data.
	data []*dtt.DtCell
}

func (hd *HighlightedData) String() string {
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

func NewHighlightedData() *HighlightedData {
	return &HighlightedData{
		data: make([]*dtt.DtCell, 0),
	}
}

func (hd *HighlightedData) FromBytes(data []byte, style tcell.Style) *HighlightedData {
	if len(data) == 0 {
		return NewHighlightedData()
	}

	hd = &HighlightedData{
		data: make([]*dtt.DtCell, 0, len(data)),
	}

	for _, b := range data {
		hd.data = append(hd.data, dtt.NewDtCell(rune(b), style))
	}

	return hd
}

func (hd *HighlightedData) AppendCell(cell *dtt.DtCell) {
	hd.data = append(hd.data, cell)
}

func (hd *HighlightedData) Size() int {
	return len(hd.data)
}

func (hd *HighlightedData) IsEmpty() bool {
	return len(hd.data) == 0
}

func (hd *HighlightedData) FieldsLine() ([]*HighlightedData, *HighlightedData) {
	result := make([]*HighlightedData, 0)
	line := NewHighlightedData()

	for _, cell := range hd.data {
		line.data = append(line.data, cell)

		if cell.First == '\n' {
			result = append(result, line)

			line = NewHighlightedData()
		}
	}

	return result, line
}

func (hd *HighlightedData) Merge(other *HighlightedData) {
	if other == nil {
		return
	}

	hd.data = append(hd.data, other.data...)
}
