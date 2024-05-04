package Lexer

import (
	"fmt"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// SourceStream is a stream of bytes that can be matched against production rules.
type SourceStream struct {
	// bytes is the byte stream.
	bytes []byte
}

// NewSourceStream creates a new source stream from a given source.
//
// Parameters:
//   - source: The source to create the stream from.
//
// Returns:
//   - *SourceStream: The new source stream.
//
// Behaviors:
//   - If the source is nil, the source stream will be created from a empty byte slice.
//   - If the source is a string, the source stream will be created from the string.
//   - If the source is a []byte, the source stream will be created from the byte slice.
//   - If the source is a fmt.Stringer, the source stream will be created from the stringer.
//   - If the source is a *SourceStream, the source stream will return the source stream as is.
//   - Otherwise, the source stream will be created from the string representation of the source.
func NewSourceStream(source any) *SourceStream {
	if source == nil {
		return &SourceStream{
			bytes: []byte{},
		}
	}

	var b []byte

	switch source := source.(type) {
	case *SourceStream:
		return source
	case []byte:
		b = source
	case fmt.Stringer:
		b = []byte(source.String())
	case string:
		b = []byte(source)
	default:
		b = []byte(fmt.Sprintf("%v", source))
	}

	return &SourceStream{
		bytes: b,
	}
}

// FromString sets the source stream to a string.
//
// Parameters:
//   - str: The string to set the source stream to.
//
// Returns:
//   - *SourceStream: The source stream.
func (s *SourceStream) FromString(str string) *SourceStream {
	s.bytes = []byte(str)

	return s
}

// FromBytes sets the source stream to a byte slice.
//
// Parameters:
//   - b: The byte slice to set the source stream to.
//
// Returns:
//   - *SourceStream: The source stream.
func (s *SourceStream) FromBytes(b []byte) *SourceStream {
	if b == nil {
		b = []byte{}
	}

	s.bytes = b

	return s
}

// IsEmpty checks if the source stream is empty.
//
// Returns:
//   - bool: True if the source stream is empty, false otherwise.
func (s *SourceStream) IsEmpty() bool {
	return len(s.bytes) == 0
}

// MatchFrom matches the source stream from a given index with a list of production rules.
//
// Parameters:
//   - from: The index to start matching from.
//   - ps: The production rules to match.
//
// Returns:
//   - matches: A slice of MatchedResult that match the input token.
//   - reason: An error if no matches are found.
//
// Errors:
//   - *ers.ErrInvalidParameter: The from index is out of bounds.
//   - *ErrNoMatches: No matches are found.
func (s *SourceStream) MatchFrom(from int, ps []*gr.RegProduction) (matches []gr.MatchedResult[*gr.LeafToken], reason error) {
	if from < 0 || from >= len(s.bytes) {
		reason = ers.NewErrInvalidParameter(
			"from",
			ers.NewErrOutOfBounds(from, 0, len(s.bytes)),
		)

		return
	}

	subSet := s.bytes[from:]

	for i, p := range ps {
		matched := p.Match(from, subSet)
		if matched != nil {
			matches = append(matches, gr.NewMatchResult(matched, i))
		}
	}

	if len(matches) == 0 {
		reason = NewErrNoMatches()
	}

	return
}

// IsDone checks if the source stream is done.
//
// Parameters:
//   - from: The index to check if the source stream is done.
//
// Returns:
//   - bool: True if the source stream is done, false otherwise.
func (s *SourceStream) IsDone(from int) bool {
	return from >= len(s.bytes)
}
