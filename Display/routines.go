package Display

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell"
)

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

// drawLoop is an helper method of Display that draws the table.
func (d *Display) drawLoop() {
	defer d.wg.Done()

	for !d.shouldStop.Get() {
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

		time.Sleep(d.frameRate)
	}
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
