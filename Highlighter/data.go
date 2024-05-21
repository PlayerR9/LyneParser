package Highlighter

import (
	"fmt"

	cdd "github.com/PlayerR9/MyGoLib/ComplexData/Display/Table"
	"github.com/gdamore/tcell"
)

type Data struct {
	source []byte
	elems  []Texter

	// rules is a map of rules to apply.
	rules map[string]tcell.Style

	// defaultStyle is the default style to apply.
	defaultStyle tcell.Style

	// errorStyle is the style to apply to errors.
	errorStyle tcell.Style
}

func (d *Data) Draw(table cdd.DrawTable, x, y *int) error {
	for _, elem := range d.elems {
		switch elem := elem.(type) {
		case *ErrorText:
		case *ValidText:
		default:
			return fmt.Errorf("unknown Texter type: %T", elem)
		}
	}

	return nil
}

func NewData(source []byte) *Data {
	return &Data{
		elems:  make([]Texter, 0),
		source: source,
	}
}

func (d *Data) Add(elem Texter) {
	if elem == nil {
		return
	}

	d.elems = append(d.elems, elem)
}
