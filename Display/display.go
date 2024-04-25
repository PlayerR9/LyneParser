package Display

import (
	"errors"
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

	// frameRate is the frame rate to display the table at.
	frameRate time.Duration
}

// NewDisplay creates a new display screen.
//
// Frame rate is the number of frames to display per second.
//
// Parameters:
//   - frameRate: The frame rate to display the table at.
//
// Returns:
//   - *Display: A pointer to the new display.
//   - error: An error if the display could not be created
//     or *ers.ErrInvalidParameter if the frame rate is less than or equal to 0.
//
// Example:
//
//	display, err := NewDisplay(60)
//	if err != nil {
//		panic(err)
//	}
//	defer display.Close()
//
//	errChan, err := display.Init(tcell.StyleDefault.Background(tcell.ColorBlack))
//	if err != nil {
//		panic(err)
//	}
//
//	table, err := display.GetTable()
//	if err != nil {
//		panic(err)
//	}
//
//	display.Start()
//
//	// Do something with the table.
func NewDisplay(frameRate float64) (*Display, error) {
	if frameRate <= 0 {
		return nil, ers.NewErrInvalidParameter(
			"frameRate",
			errors.New("value must be greater than 0"),
		)
	}

	d := &Display{
		frameRate: time.Duration(1/frameRate) * time.Second,
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
		go d.drawLoop()
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
	if d.shouldStop == nil {
		return
	}

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

// GetTable is a method of Display that returns the table.
//
// Returns:
//   - dtt.WriteOnlyDTer: The table in Write Only mode.
func (d *Display) GetTable() dtt.WriteOnlyDTer {
	return d.table
}
