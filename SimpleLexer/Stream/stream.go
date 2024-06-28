package Stream

import (
	"unicode/utf8"
)

// lookahead is a lookahead rune.
type lookahead struct {
	// char is the rune.
	char rune

	// size is the size of the rune.
	size int
}

// Stream is a stream of runes.
type Stream struct {
	// data is the data of the stream.
	data []byte

	// pos is the current position of the stream.
	pos int

	// lastPos is the last position of the stream.
	lastPos int

	// lookahead is the next rune.
	lookahead *lookahead
}

// NewStream creates a new stream.
//
// Parameters:
//   - data: The data of the stream.
//
// Returns:
//   - Stream: The stream.
func NewStream(data []byte) Stream {
	s := Stream{
		data:      data,
		pos:       0,
		lastPos:   0,
		lookahead: nil,
	}
	return s
}

// Peek peeks the next rune without consuming it.
//
// Returns:
//   - rune: The next rune.
//   - error: The error if any.
func (s *Stream) Peek() (rune, error) {
	if s.lookahead != nil {
		return s.lookahead.char, nil
	}

	if len(s.data) == 0 {
		return 0, NewStreamExhausted()
	}

	r, size := utf8.DecodeRune(s.data)
	s.data = s.data[size:]

	la := &lookahead{
		char: r,
		size: size,
	}

	s.lookahead = la

	return r, nil
}

// Next returns the next rune and consumes it.
//
// Returns:
//   - rune: The next rune.
//   - error: The error if any.
func (s *Stream) Next() (rune, error) {
	if s.lookahead != nil {
		r := s.lookahead.char
		s.lastPos = s.pos

		s.lookahead = nil

		return r, nil
	}

	if len(s.data) == 0 {
		return 0, NewStreamExhausted()
	}

	r, size := utf8.DecodeRune(s.data)
	s.data = s.data[size:]
	s.pos += size

	return r, nil
}

// Pos returns the current position of the stream.
//
// Returns:
//   - int: The current position.
func (s *Stream) Pos() int {
	return s.lastPos
}

// Accept accepts the lookahead rune.
//
// Does nothing if there is no lookahead rune.
func (s *Stream) Accept() {
	if s.lookahead == nil {
		return
	}

	s.lastPos += s.lookahead.size

	s.lookahead = nil
}
