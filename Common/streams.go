package Common

import (
	"fmt"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"

	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
)

// ByteStream is a stream of bytes that can be matched against production rules.
type ByteStream struct {
	*cds.Stream[byte]
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
func NewSourceStream(source any) *ByteStream {
	if source == nil {
		return &ByteStream{cds.NewStream([]byte{})}
	}

	var b []byte

	switch source := source.(type) {
	case *ByteStream:
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

	return &ByteStream{cds.NewStream(b)}
}

// FromString sets the source stream to a string.
//
// Parameters:
//   - str: The string to set the source stream to.
//
// Returns:
//   - *SourceStream: The source stream.
func (s *ByteStream) FromString(str string) *ByteStream {
	return &ByteStream{cds.NewStream([]byte(str))}
}

// FromBytes sets the source stream to a byte slice.
//
// Parameters:
//   - b: The byte slice to set the source stream to.
//
// Returns:
//   - *SourceStream: The source stream.
func (s *ByteStream) FromBytes(b []byte) *ByteStream {
	if b == nil {
		b = []byte{}
	}

	return &ByteStream{cds.NewStream(b)}
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
func (s *ByteStream) MatchFrom(from int, ps []*gr.RegProduction) (matches []gr.MatchedResult[*gr.LeafToken], reason error) {
	size := s.Size()

	if from < 0 || from >= size {
		reason = ers.NewErrInvalidParameter(
			"from",
			ers.NewErrOutOfBounds(from, 0, size),
		)

		return
	}

	subSet, err := s.Get(from, size)
	if err != nil {
		panic(err)
	}

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

// TokenStream is a stream of tokens that have been lexed.
type TokenStream struct {
	*cds.Stream[*gr.LeafToken]
}

// NewTokenStream creates a new token stream from a given slice of tokens.
//
// Parameters:
//   - tokens: The tokens to create the stream from.
//
// Returns:
//   - *TokenStream: The new token stream.
//
// Behaviors:
//   - If the tokens are nil, the token stream will be created from a empty token slice.
//   - If the tokens are a *TokenStream, the token stream will return the token stream as is.
//   - Otherwise, the token stream will be created from the token slice.
func NewTokenStream(tokens []*gr.LeafToken) *TokenStream {
	if tokens == nil {
		tokens = []*gr.LeafToken{}
	}

	return &TokenStream{cds.NewStream(tokens)}
}
