package Common

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"

	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
)

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
