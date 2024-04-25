package Display

import (
	"errors"
	"sync"
	"time"

	"github.com/gdamore/tcell"

	rws "github.com/PlayerR9/MyGoLib/CustomData/Safe/RWSafe"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// Display represents a display.
type Display struct {
	// screen is the screen to display on.
	screen tcell.Screen

	// width and height are the width and height of the display,
	// respectively.
	width, height int

	// frameRate is the frame rate of the display.
	frameRate time.Duration

	// shouldStop is a flag that indicates if the display should stop.
	shouldStop rws.RWSafe[bool]

	// wg is a wait group that waits for the display to stop.
	wg sync.WaitGroup

	// eventChan is a channel that receives events.
	eventChan chan tcell.Event

	// errChan is a channel that receives errors.
	errChan chan error
}

// NewDisplay creates a new Display with the given frame rate.
//
// Parameters:
//   - frameRate: The frame rate of the display (number of draws per second).
//
// Returns:
//   - *Display: A pointer to the new Display.
//   - error: An error if the display could not be created or
//     typed *ers.ErrInvalidParameter if the frame rate is less than or equal to 0.
func NewDisplay(frameRate float64) (*Display, error) {
	if frameRate <= 0 {
		return nil, ers.NewErrInvalidParameter(
			"frameRate",
			errors.New("value must be greater than 0"),
		)
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	err = screen.Init()
	if err != nil {
		return nil, err
	}

	d := &Display{
		screen:    screen,
		frameRate: time.Duration(1/frameRate) * time.Second,
		eventChan: make(chan tcell.Event),
		errChan:   make(chan error),
	}

	d.width, d.height = d.screen.Size()

	return d, nil
}

// GetErrChan is a method of Display that returns the error channel.
//
// Returns:
//   - <-chan error: The error channel.
func (d *Display) GetErrChan() <-chan error {
	return d.errChan
}

// Start is a method of Display that starts the display.
//
// Parameters:
//   - table: The table to display.
func (d *Display) Start(table *DtTable) {
	d.shouldStop.Set(false)

	// Start the get event goroutine.
	go func() {
		for {
			ev := d.screen.PollEvent()
			if ev == nil {
				break
			}

			d.eventChan <- ev
		}
	}()

	d.wg.Add(2)

	go func() {
		defer d.wg.Done()

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
				}
			case *tcell.EventResize:
				d.width, d.height = ev.Size()
			}
		}
	}()

	go func() {
		defer d.wg.Done()

		for !d.shouldStop.Get() {
			d.screen.Clear()

			for y, row := range table.cells {
				if y >= d.height {
					break
				}

				for x, cell := range row {
					if x >= d.width {
						break
					}

					d.screen.SetContent(x, y, cell.Content, nil, cell.Style)
				}
			}

			d.screen.Show()

			time.Sleep(d.frameRate)
		}
	}()
}

// Stop is a method of Display that stops the display.
func (d *Display) Stop() {
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
