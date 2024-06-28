package Grammar

import (
	"fmt"
	"unicode"

	util "github.com/PlayerR9/LyneParser/Util"
)

const (
	// EOFTokenID is the identifier of the end-of-file token.
	EOFTokenID string = "EOF"

	// RootTokenID is the identifier of the root token.
	RootTokenID string = "ROOT"
)

// EofToken creates the end-of-file token.
//
// Returns:
//   - Token: The end-of-file token.
func EofToken() Token {
	tok := Token{
		ID: EOFTokenID,
		At: -1,
	}
	return tok
}

// RootToken creates the root token.
//
// Returns:
//   - Token: The root token.
func RootToken() Token {
	tok := Token{
		ID: RootTokenID,
		At: -1,
	}
	return tok
}

// IsTerminal checks if the given identifier is a terminal. Terminals are identifiers
// that start with an uppercase letter.
//
// Parameters:
//   - identifier: The identifier to check.
//
// Returns:
//   - bool: True if the identifier is a terminal, false otherwise.
//
// Asserts:
//   - The identifier is not empty.
func IsTerminal(identifier string) bool {
	util.Assert(identifier != "", "In IsTerminal: identifier is empty")

	firstLetter := []rune(identifier)[0]

	ok := unicode.IsUpper(firstLetter)
	return ok
}

// Token is the information about a token.
type Token struct {
	// ID is the identifier of the token.
	ID string

	// At is the position of the token in the input string.
	At int

	// Lookahead is the next token in the input string.
	Lookahead *Token

	// Data is the data of the token.
	// If data is a string, it is the data of a leaf token.
	// If data is a slice of Token, it is the data of a non-leaf token.
	// any other type of data is not supported.
	//
	// Only EofToken and RootToken have nil data.
	Data any
}

// GoString is a method of fmt.GoStringer interface.
func (tok *Token) GoString() string {
	str := fmt.Sprintf("%+v", *tok)
	return str
}

// NewToken creates a new token info with the given identifier, data, position,
// and lookahead token.
//
// Parameters:
//   - id: The identifier of the token.
//   - data: The data of the token.
//   - at: The position of the token in the input string.
//   - lookahead: The next token in the input string.
//
// Returns:
//   - *Token: A pointer to the new token info. Nil if the data is nil
//     or not a string or a slice of Token.
func NewToken(id string, data any, at int, lookahead *Token) Token {
	if data == nil {
		panic("In NewToken: data is nil")
	}

	switch data.(type) {
	case string, []Token:
		tok := Token{
			ID:        id,
			Data:      data,
			At:        at,
			Lookahead: lookahead,
		}
		return tok
	default:
		panic("In NewToken: data is not a string or a slice of Token")
	}
}

// GetID returns the identifier of the token.
//
// Returns:
//   - string: The identifier of the token.
func (tok *Token) GetID() string {
	return tok.ID
}

// GetPos returns the position of the token in the input string.
//
// Returns:
//   - int: The position of the token in the input string.
func (tok *Token) GetPos() int {
	return tok.At
}

// GetLookahead returns the next token in the input string.
//
// Returns:
//   - *Token: The next token in the input string.
func (tok *Token) GetLookahead() *Token {
	return tok.Lookahead
}

// SetLookahead sets the next token in the input string.
//
// Parameters:
//   - lookahead: The next token in the input string.
func (tok *Token) SetLookahead(lookahead *Token) {
	tok.Lookahead = lookahead
}

// IsLeaf checks if the token is a leaf token.
//
// Returns:
//   - bool: True if the token is a leaf token, false otherwise.
func (tok *Token) IsLeaf() bool {
	_, ok := tok.Data.(string)
	return ok
}

// IsNonLeaf checks if the token is a non-leaf token.
//
// Returns:
//   - bool: True if the token is a non-leaf token, false otherwise.
func (tok *Token) IsNonLeaf() bool {
	_, ok := tok.Data.([]Token)
	return ok
}

// GetData returns the data of the token.
//
// Data can only be a string or a slice of Token. Unless
// the token is the EofToken or the RootToken, the data
// should not be nil.
//
// Returns:
//   - any: The data of the token.
func (tok *Token) GetData() any {
	return tok.Data
}

// Copy creates a copy of the token.
//
// It does not copy the lookahead token.
//
// Returns:
//   - Token: A copy of the token.
func (tok *Token) Copy() Token {
	lt := Token{
		ID: tok.ID,
		At: tok.At,
	}

	if tok.Data == nil {
		return lt
	}

	switch data := tok.Data.(type) {
	case string:
		lt.Data = data
	case []Token:
		sliceCopy := make([]Token, 0, len(data))
		for _, elem := range data {
			elemCopy := elem.Copy()
			sliceCopy = append(sliceCopy, elemCopy)
		}

		lt.Data = sliceCopy
	default:
		panic("In Token.Copy: unsupported data type")
	}

	return lt
}
