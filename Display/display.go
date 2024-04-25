package Display

import (
	"fmt"
	"sync"
	"time"

	dtt "github.com/PlayerR9/LyneParser/DtTable"
	rws "github.com/PlayerR9/MyGoLib/CustomData/Safe/RWSafe"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
	"github.com/gdamore/tcell"
)

// Display represents a display on a screen.
type Display struct {
	// screen is the screen to display on.
	screen tcell.Screen

	// table is the table to display.
	table *dtt.DtTable

	// eventChan is the channel to send events to.
	eventChan chan tcell.Event

	// errChan is the channel to send errors to.
	errChan chan error

	// shouldStop is a flag to stop the display.
	shouldStop *rws.RWSafe[bool]

	// wg is a wait group to wait for the display to stop.
	wg sync.WaitGroup

	// once is a sync.Once to ensure that the display is only started once.
	once sync.Once

	// runFunc is the function that runs the display.
	rf RunFunc
}

// NewDisplayScreen creates a new display screen.
//
// Parameters:
//   - rf: The function that runs the display.
//
// Returns:
//   - *Display: A pointer to the new display.
//   - error: An error if the display could not be created
//     or *ers.ErrInvalidParameter if rf is nil.
func NewDisplayScreen(rf RunFunc) (*Display, error) {
	if rf == nil {
		return nil, ers.NewErrNilParameter("rf")
	}

	d := &Display{
		rf:        rf,
		eventChan: make(chan tcell.Event),
		errChan:   make(chan error),
	}

	var err error

	d.screen, err = tcell.NewScreen()

	return d, err
}

// Init initializes the display.
//
// Parameters:
//   - bgStyle: The background style of the display.
//
// Returns:
//   - <-chan error: The error channel.
//   - error: An error if the display could not be initialized.
func (d *Display) Init(bgStyle tcell.Style) (<-chan error, error) {
	err := d.screen.Init()
	if err != nil {
		return d.errChan, err
	}

	width, height := d.screen.Size()

	d.table, err = dtt.NewDtTable(width, height)
	if err != nil {
		panic(fmt.Errorf("error creating table: %w", err))
	}

	d.screen.EnableMouse()
	d.screen.SetStyle(bgStyle)

	d.shouldStop = rws.NewRWSafe(false)

	return d.errChan, nil
}

// GetErrChan returns the error channel.
//
// Returns:
//   - <-chan error: The error channel.
func (d *Display) GetErrorChan() <-chan error {
	return d.errChan
}

// Start starts the display.
func (d *Display) Start() {
	d.once.Do(func() {
		d.shouldStop.Set(false)

		go d.eventListener()

		d.wg.Add(2)

		go d.handleEvents()
		go d.rf(d)
	})
}

// Wait waits for the display to stop.
// WARNING: This method will block until the display stops;
// may cause deadlock.
func (d *Display) Wait() {
	d.wg.Wait()
}

// Close closes the display.
// WARNING: This method will block until the display stops;
// may cause deadlock.
// It cleans up the resources used by the display.
func (d *Display) Close() {
	d.shouldStop.Set(true)

	d.wg.Wait()

	if d.eventChan != nil {
		close(d.eventChan)
		d.eventChan = nil
	}

	if d.errChan != nil {
		close(d.errChan)
		d.errChan = nil
	}

	if d.screen != nil {
		d.screen.Fini()
		d.screen = nil
	}
}

// eventListener is a helper method of Display that listens for events.
func (d *Display) eventListener() {
	for {
		ev := d.screen.PollEvent()
		if ev == nil {
			break
		}

		d.eventChan <- ev
	}
}

// displayLoop is a helper method of Display that displays the table.
func (d *Display) displayOnce() {
	d.screen.Clear()

	height := d.table.GetHeight()
	width := d.table.GetWidth()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if cell := d.table.GetCellAt(x, y); cell != nil {
				d.screen.SetContent(x, y, cell.Content, nil, cell.Style)
			}
		}
	}

	d.screen.Show()
}

// setTable is a helper method of Display that sets the table.
//
// Parameters:
//   - table: The table to set.
//
// Returns:
//   - error: An error if the table could not be set.
func (d *Display) setTable(table *dtt.DtTable) error {
	err := table.ResizeWidth(d.table.GetWidth())
	if err != nil {
		return err
	}

	err = table.ResizeHeight(d.table.GetHeight())
	if err != nil {
		return err
	}

	d.table = table

	return nil
}

// GetTable is a method of Display that returns the table.
//
// Returns:
//   - dtt.WriteOnlyDTer: The table in Write Only mode.
func (d *Display) GetTable() dtt.WriteOnlyDTer {
	return d.table
}

// handleEvents is an helper method of Display that handles events.
func (d *Display) handleEvents() {
	defer d.wg.Done()

	for {
		select {
		case <-time.After(100 * time.Millisecond):
			if d.shouldStop.Get() {
				return
			}
		case ev := <-d.eventChan:
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyCtrlC {
					d.errChan <- NewErrESCPressed()
				} else {
					d.errChan <- fmt.Errorf("unknown key: %v", ev.Key())
				}
			case *tcell.EventResize:
				width, height := ev.Size()

				err := d.table.ResizeHeight(height)
				if err != nil {
					d.errChan <- err
				}

				err = d.table.ResizeWidth(width)
				if err != nil {
					d.errChan <- err
				}
			default:
				d.errChan <- fmt.Errorf("unknown event: %v", ev)
			}
		}
	}
}
