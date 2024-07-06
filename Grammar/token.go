package Grammar

import (
	"errors"
	"fmt"
	"unicode"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

const (
	// EOFTokenID is the identifier of the end-of-file token.
	EOFTokenID string = "EOF"

	// RootTokenID is the identifier of the root token.
	RootTokenID string = "ROOT"
)

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
	uc.Assert(identifier != "", "In IsTerminal: identifier is empty")

	firstLetter := []rune(identifier)[0]

	ok := unicode.IsUpper(firstLetter)
	return ok
}

// Token is the information about a token.
type Token[T uc.Enumer] struct {
	// ID is the identifier of the token.
	ID T

	// At is the position of the token in the input string.
	At int

	// Lookahead is the next token in the input string.
	Lookahead *Token[T]

	// Data is the data of the token.
	// If data is a string, it is the data of a leaf token.
	// If data is a slice of Token, it is the data of a non-leaf token.
	// any other type of data is not supported.
	//
	// Only EofToken and RootToken have nil data.
	Data any
}

// GoString is a method of fmt.GoStringer interface.
func (tok *Token[T]) GoString() string {
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
func NewToken[T uc.Enumer](id T, data any, at int, lookahead *Token[T]) *Token[T] {
	uc.AssertParam("data", data != nil, errors.New("in NewToken: data is nil"))

	switch data.(type) {
	case string, []*Token[T]:
		tok := &Token[T]{
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
//   - T: The identifier of the token.
func (tok *Token[T]) GetID() T {
	return tok.ID
}

// GetPos returns the position of the token in the input string.
//
// Returns:
//   - int: The position of the token in the input string.
func (tok *Token[T]) GetPos() int {
	return tok.At
}

// GetLookahead returns the next token in the input string.
//
// Returns:
//   - *Token: The next token in the input string.
func (tok *Token[T]) GetLookahead() *Token[T] {
	return tok.Lookahead
}

// SetLookahead sets the next token in the input string.
//
// Parameters:
//   - lookahead: The next token in the input string.
func (tok *Token[T]) SetLookahead(lookahead *Token[T]) {
	tok.Lookahead = lookahead
}

// IsLeaf checks if the token is a leaf token.
//
// Returns:
//   - bool: True if the token is a leaf token, false otherwise.
func (tok *Token[T]) IsLeaf() bool {
	if tok.Data == nil {
		return true
	}

	_, ok := tok.Data.(string)
	return ok
}

// IsNonLeaf checks if the token is a non-leaf token.
//
// Returns:
//   - bool: True if the token is a non-leaf token, false otherwise.
func (tok *Token[T]) IsNonLeaf() bool {
	if tok.Data == nil {
		return false
	}

	_, ok := tok.Data.([]*Token[T])
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
func (tok *Token[T]) GetData() any {
	return tok.Data
}

// Copy implements common.Copier interface.
func (tok *Token[T]) Copy() uc.Copier {
	lt := &Token[T]{
		ID: tok.ID,
		At: tok.At,
	}

	switch data := tok.Data.(type) {
	case string:
		lt.Data = data
	case []*Token[T]:
		sliceCopy := make([]*Token[T], 0, len(data))
		for _, elem := range data {
			elemCopy := elem.Copy().(*Token[T])
			sliceCopy = append(sliceCopy, elemCopy)
		}

		lt.Data = sliceCopy
	default:
		panic("In Token.Copy: unsupported data type")
	}

	return lt
}
