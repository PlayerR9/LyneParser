package Display

import (
	"testing"

	"github.com/gdamore/tcell"
)

func TestStartAndClose(t *testing.T) {
	d, err := NewDisplay(1)
	if err != nil {
		t.Fatalf("error creating display: %v", err)
	}

	_, err = d.Init(tcell.StyleDefault)
	if err != nil {
		t.Fatalf("error initializing display: %v", err)
	}

	d.Start()

	d.Close()
}
