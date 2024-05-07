package Pair

import (
	"fmt"

	intf "github.com/PlayerR9/LyneParser/PlayerR9/Common"
)

// Pair is a pair of values.
type Pair[A any, B any] struct {
	// The first value.
	First A

	// The second value.
	Second B
}

func (p *Pair[A, B]) String() string {
	return fmt.Sprintf("(%v, %v)", p.First, p.Second)
}

func (p *Pair[A, B]) Copy() intf.Copier {
	return &Pair[A, B]{
		First:  intf.CopyOf(p.First).(A),
		Second: intf.CopyOf(p.Second).(B),
	}
}

// NewPair creates a new pair.
//
// Parameters:
//   - first: The first value.
//   - second: The second value.
//
// Returns:
//   - Pair[A, B]: The new pair.
func NewPair[A any, B any](first A, second B) *Pair[A, B] {
	return &Pair[A, B]{
		First:  first,
		Second: second,
	}
}
