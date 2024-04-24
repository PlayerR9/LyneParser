package Grammar

// TokenStream is a stream of tokens.
type TokenStream struct {
	// tokens is the slice of tokens in the stream.
	tokens []LeafToken
}

func NewTokenStream(tokens []LeafToken) TokenStream {
	return TokenStream{tokens: tokens}
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
//		- id: The token ID to remove.
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
