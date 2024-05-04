package Grammar

import (
	"github.com/PlayerR9/MyGoLib/CustomData/Stream"
)

// TokenStream is a stream of tokens.
type TokenStream struct {
	// tokens is the slice of tokens in the stream.
	tokens []*LeafToken

	// currentIndex is the current index of the token stream.
	// It indicates the first non-consumed token.
	currentIndex int
}

// Size returns the number of tokens in the token stream.
//
// Returns:
//   - int: The number of tokens in the token stream.
func (ts *TokenStream) Size() int {
	return len(ts.tokens)
}

// IsEmpty returns true if the token stream is empty.
//
// Returns:
//   - bool: True if the token stream is empty.
func (ts *TokenStream) IsEmpty() bool {
	return len(ts.tokens) == 0
}

// Reset resets the token stream to the beginning.
func (ts *TokenStream) Reset() {
	ts.currentIndex = 0
}

// Peek returns the next token in the stream without consuming it.
//
// Returns:
//   - *LeafToken: A pointer to the next token in the stream.
//   - error: An error if there are no more tokens in the stream.
func (ts *TokenStream) Peek() (*LeafToken, error) {
	if ts.currentIndex >= len(ts.tokens) {
		return nil, Stream.NewErrNoMoreItems()
	}

	return ts.tokens[ts.currentIndex], nil
}

// Consume consumes the next token in the stream.
//
// Returns:
//   - *LeafToken: A pointer to the consumed token.
//   - error: An error if there are no more tokens in the stream.
func (ts *TokenStream) Consume() (*LeafToken, error) {
	if ts.currentIndex >= len(ts.tokens) {
		return nil, Stream.NewErrNoMoreItems()
	}

	token := ts.tokens[ts.currentIndex]
	ts.currentIndex++

	return token, nil
}

// IsDone returns true if the token stream has been fully consumed.
//
// Returns:
//   - bool: True if the token stream has been fully consumed.
func (ts *TokenStream) IsDone() bool {
	return ts.currentIndex >= len(ts.tokens)
}

// GetTokens returns the tokens in the stream.
// It ignores the end-of-file token if present.
//
// Returns:
//   - []*LeafToken: The tokens in the stream.
func (ts *TokenStream) GetTokens() []*LeafToken {
	result := make([]*LeafToken, len(ts.tokens))
	copy(result, ts.tokens)

	if len(result) > 0 && result[len(result)-1].ID == EOFTokenID {
		result = result[:len(result)-1]
	}

	return result
}

// GetLeftoverTokens returns the tokens that have not been consumed.
//
// Returns:
//   - []*LeafToken: The leftover tokens.
func (ts *TokenStream) GetLeftoverTokens() []*LeafToken {
	return ts.tokens[ts.currentIndex:]
}

// NewTokenStream creates a new token stream with the given tokens.
//
// Parameters:
//   - tokens: The tokens to add to the stream.
//
// Returns:
//   - TokenStream: The new token stream.
func NewTokenStream(tokens []*LeafToken) *TokenStream {
	return &TokenStream{
		tokens:       tokens,
		currentIndex: 0,
	}
}
