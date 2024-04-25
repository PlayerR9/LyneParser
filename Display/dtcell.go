package Display

import (
	"github.com/gdamore/tcell"
)

// DtCell represents a cell in a data table.
type DtCell struct {
	// Content is the content of the cell.
	Content rune

	// Style is the style of the cell.
	Style tcell.Style
}

// NewDtCell creates a new DtCell with the given content and style.
//
// Parameters:
//   - content: The content of the cell.
//   - style: The style of the cell.
//
// Returns:
//   - DtCell: The new DtCell.
func NewDtCell(content rune, style tcell.Style) DtCell {
	return DtCell{
		Content: content,
		Style:   style,
	}
}
