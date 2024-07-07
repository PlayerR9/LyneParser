package Highlighter

import (
	"fmt"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	p9 "github.com/PlayerR9/LyneParser/PlayerR9"
	cdd "github.com/PlayerR9/MyGoLib/Display/drawtable"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	"github.com/gdamore/tcell"
)

// Data is a highlighted data.
type Data[T gr.TokenTyper] struct {
	// source is the source of the data.
	source []byte

	// elems is the elements of the data.
	elems []Texter

	// rules is a map of rules to apply.
	rules map[T]tcell.Style

	// default_style is the default style to apply.
	default_style tcell.Style

	// error_style is the style to apply to errors.
	error_style tcell.Style
}

// Draw draws the data.
//
// Parameters:
//   - table: The table to draw to.
//   - x: The x position to draw to.
//   - y: The y position to draw to.
//
// Returns:
//   - error: An error if there was a problem drawing the data.
func (d *Data[T]) Draw(table cdd.DrawTable, x, y *int) error {
	for _, elem := range d.elems {
		switch elem := elem.(type) {
		case *NormalText:
			sequences, err := p9.AnyToLines(elem, func(r rune) (*cdd.ColoredUnit, error) {
				return cdd.NewColoredUnit(r, elem.style), nil
			})
			uc.AssertF(err == nil, "AnyToLines failed: %s", err.Error())

			// FINISH THIS
			for _, sequence := range sequences {
				table.WriteHorizontalSequence(x, y, sequence)
			}
		case *ValidText[T]:
		default:
			return fmt.Errorf("unknown Texter type: %T", elem)
		}
	}

	return nil
}

// Add adds an element to the data.
//
// Parameters:
//   - elem: The element to add.
func (d *Data[T]) Add(elem Texter) {
	if elem == nil {
		return
	}

	d.elems = append(d.elems, elem)
}
