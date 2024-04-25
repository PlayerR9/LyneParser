package Display

import (
	"errors"
	"time"

	dtt "github.com/PlayerR9/LyneParser/DtTable"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// RunFunc is a function that runs a display.
//
// Parameters:
//
//   - d: The display to run.
type RunFunc func(d *Display)

// ContinuousDisplay creates a display that continuously displays
// the table at the given frame rate.
//
// Frame rate is the number of frames to display per second.
//
// Parameters:
//
//   - frameRate: The frame rate to display the table at.
//
// Returns:
//
//   - RunFunc: The function that runs the display.
//   - error: An error if the frame rate is invalid.
func ContinuousDisplay(frameRate float64) (RunFunc, error) {
	if frameRate <= 0 {
		return nil, ers.NewErrInvalidParameter(
			"frameRate",
			errors.New("value must be greater than 0"),
		)
	}

	sleepTime := time.Duration(1/frameRate) * time.Second

	return func(d *Display) {
		defer d.wg.Done()

		for !d.shouldStop.Get() {
			d.displayOnce()
			time.Sleep(sleepTime)
		}
	}, nil
}

// StaticDisplay creates a display that displays the table once.
//
// Parameters:
//
//   - table: The table to display.
//
// Returns:
//
//   - RunFunc: The function that runs the display.
//   - error: An error if the table is nil.
func StaticDisplay(table *dtt.DtTable) (RunFunc, error) {
	if table == nil {
		return nil, ers.NewErrNilParameter("table")
	}

	return func(d *Display) {
		defer d.wg.Done()

		err := d.setTable(table)
		if err != nil {
			d.errChan <- err
			return
		}

		d.displayOnce()
	}, nil
}
