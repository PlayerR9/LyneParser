package Grammar

// TokenStream is a stream of tokens.
type TokenStream struct {
	// tokens is the slice of tokens in the stream.
	tokens []*LeafToken

	// currentIndex is the current index of the token stream.
	// It indicates the first non-consumed token.
	currentIndex int
}

// NewTokenStream creates a new token stream with the given tokens.
//
// Parameters:
//   - tokens: The tokens to add to the stream.
//
// Returns:
//   - TokenStream: The new token stream.
func NewTokenStream(tokens []*LeafToken) TokenStream {
	return TokenStream{
		tokens:       tokens,
		currentIndex: 0,
	}
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

// RemoveByTokenID removes tokens by their token ID.
//
// Parameters:
//   - id: The token ID to remove.
func (ts *TokenStream) RemoveByTokenID(id string) {
	if len(ts.tokens) == 0 {
		return
	}

	top := 0

	for _, token := range ts.tokens {
		if token.ID != id {
			ts.tokens[top] = token
			top++
		}
	}

	ts.tokens = ts.tokens[:top]
}

// SetEOFToken sets the end-of-file token in the token stream.
//
// If the end-of-file token is already present, it will not be added again.
func (ts *TokenStream) SetEOFToken() {
	if len(ts.tokens) != 0 && ts.tokens[len(ts.tokens)-1].ID == EOFTokenID {
		// EOF token is already present
		return
	}

	tok := NewEOFToken()
	ts.tokens = append(ts.tokens, tok)
}

// SetLookahead sets the lookahead token for all the tokens in the stream.
func (ts *TokenStream) SetLookahead() {
	for i, token := range ts.tokens[:len(ts.tokens)-1] {
		token.SetLookahead(ts.tokens[i+1])
	}
}

// Peek returns the next token in the stream without consuming it.
// It panics if there are no more tokens in the stream.
//
// Returns:
//   - *LeafToken: A pointer to the next token in the stream.
func (ts *TokenStream) Peek() *LeafToken {
	if ts.currentIndex >= len(ts.tokens) {
		panic(NewErrNoMoreTokens())
	}

	return ts.tokens[ts.currentIndex]
}

// Consume consumes the next token in the stream.
// It panics if there are no more tokens in the stream.
//
// Returns:
//   - *LeafToken: A pointer to the consumed token.
func (ts *TokenStream) Consume() *LeafToken {
	if ts.currentIndex >= len(ts.tokens) {
		panic(NewErrNoMoreTokens())
	}

	token := ts.tokens[ts.currentIndex]
	ts.currentIndex++

	return token
}

// Reset resets the token stream to the beginning.
func (ts *TokenStream) Reset() {
	ts.currentIndex = 0
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
