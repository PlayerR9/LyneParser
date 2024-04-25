package Display

import "github.com/gdamore/tcell"

// DtTable represents a table of cells.
type DtTable struct {
	// height and width represent the height and width of the table, respectively.
	height, width int

	// cells is a 2D slice of cells.
	cells [][]DtCell
}

// TransformIntoTable transforms a slice of cells into a table.
//
// Parameters:
//   - highlights: The slice of cells to transform.
//   - bgStyle: The background style to use for empty cells.
//
// Returns:
//   - DtTable: The new table.
func TransformIntoTable(highlights []DtCell, bgStyle tcell.Style) *DtTable {
	if len(highlights) == 0 {
		return &DtTable{
			height: 0,
			width:  0,
			cells:  make([][]DtCell, 0),
		}
	}

	table := &DtTable{
		cells: make([][]DtCell, 0),
	}

	row := make([]DtCell, 0)

	for _, hl := range highlights {
		if hl.Content == '\n' {
			table.cells = append(table.cells, row)
			row = make([]DtCell, 0)
		} else {
			row = append(row, hl)
		}
	}

	if len(row) > 0 {
		table.cells = append(table.cells, row)
	}

	table.height = len(table.cells)

	// Fix the sizes of the table.
	table.width = 0

	for _, row := range table.cells {
		if len(row) > table.width {
			table.width = len(row)
		}
	}

	for i, row := range table.cells {
		if len(row) == table.width {
			continue
		}

		newRow := make([]DtCell, len(row))
		copy(newRow, row)

		for j := len(row); j < table.width; j++ {
			newRow = append(newRow, NewDtCell(' ', bgStyle))
		}

		table.cells[i] = newRow
	}

	return table
}
