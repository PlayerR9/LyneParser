package DtTable

import (
	"errors"
	"sync"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// DtRow represents a row in a data table.
// It is safe for concurrent use.
type DtRow struct {
	// cells is a slice of cells in the row.
	cells []*DtCell

	// width represents the width of the row.
	width int

	// mu is a mutex that protects the row.
	mu sync.RWMutex
}

// NewDtRow creates a new DtRow with the given width all initialized to nil.
//
// Parameters:
//   - width: The width of the row.
//
// Returns:
//   - *DtRow: A pointer to the new DtRow.
//   - error: An error of type *ers.ErrInvalidParameter if the width
//     is less than 0.
func NewDtRow(width int) (*DtRow, error) {
	if width < 0 {
		return nil, ers.NewErrInvalidParameter(
			"width",
			errors.New("value must be non-negative"),
		)
	}

	return &DtRow{
		cells: make([]*DtCell, width),
		width: width,
	}, nil
}

// GetWidth returns the width of the row.
//
// Returns:
//   - int: The width of the row.
func (r *DtRow) GetWidth() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.width
}

// SetCell sets the cell at the given index.
//
// Parameters:
//   - cell: The cell to set.
//   - index: The index of the cell.
//
// Returns:
//   - error: An error of type *ers.ErrInvalidParameter if the index
//     is out of bounds.
func (r *DtRow) SetCell(cell *DtCell, index int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if index < 0 || index >= r.width {
		return ers.NewErrInvalidParameter(
			"index",
			errors.New("value must be between 0 and width"),
		)
	}

	r.cells[index] = cell

	return nil
}

// SetCells sets the cells at the given index.
//
// Parameters:
//   - cells: The cells to set.
//   - from: The index to start setting the cells.
//
// Returns:
//   - error: An error of type *ers.ErrInvalidParameter if cells
//     cannot be set at the given index.
func (r *DtRow) SetCells(cells []*DtCell, from int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if from < 0 || from+len(cells) > r.width {
		return ers.NewErrInvalidParameter(
			"from",
			ers.NewErrOutOfBounds(from, 0, r.width-len(cells)),
		)
	}

	for i, cell := range cells {
		r.cells[from+i] = cell
	}

	return nil
}

// Append appends the given cells to the row.
//
// Parameters:
//   - cells: The cells to append.
func (r *DtRow) Append(cells ...*DtCell) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cells = append(r.cells, cells...)
	r.width += len(cells)
}

// GetCellAt returns the cell at the given index.
// If the index is out of bounds, it returns nil.
//
// Parameters:
//   - index: The index of the cell.
func (r *DtRow) GetCellAt(index int) *DtCell {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if index < 0 || index >= r.width {
		return nil
	}

	return r.cells[index]
}

// Resize resizes the row to the given width.
//
// Parameters:
//   - newWidth: The new width of the row.
//
// Returns:
//   - error: An error of type *ers.ErrInvalidParameter if the new width
//     is less than 0.
func (r *DtRow) Resize(newWidth int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if newWidth < 0 {
		return ers.NewErrInvalidParameter(
			"newWidth",
			errors.New("value must be greater than 0"),
		)
	}

	if newWidth == r.width {
		return nil
	}

	if newWidth < r.width {
		r.cells = r.cells[:newWidth]
	} else {
		r.cells = append(r.cells, make([]*DtCell, newWidth-r.width)...)
	}

	r.width = newWidth

	return nil
}
