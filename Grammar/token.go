package Grammar

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	// EOFTokenID is the identifier of the end-of-file token.
	EOFTokenID string = "EOF"
)

// NewEOFToken creates a new end-of-file token.
//
// Returns:
//
//   - *LeafToken: A pointer to the new end-of-file token.
func NewEOFToken() *LeafToken {
	tok := &LeafToken{
		ID:   EOFTokenID,
		Data: "",
		At:   -1,
	}

	return tok
}

func NewRootToken() *LeafToken {
	tok := &LeafToken{
		ID:   "ROOT",
		Data: "",
		At:   -1,
	}

	return tok
}

// IsTerminal checks if the given identifier is a terminal. Terminals are identifiers
// that start with an uppercase letter.
//
// Parameters:
//
//   - identifier: The identifier to check.
//
// Returns:
//
//   - bool: True if the identifier is a terminal, false otherwise.
func IsTerminal(identifier string) bool {
	firstLetter := []rune(identifier)[0]

	return unicode.IsLetter(firstLetter) && unicode.IsUpper(firstLetter)
}

// Tokener is an interface that defines the methods that a token must implement.
type Tokener interface {
	// GetID returns the identifier of the token.
	//
	// Returns:
	//
	//   - string: The identifier of the token.
	GetID() string

	// GetData returns the data of the token.
	//
	// Returns:
	//
	//   - any: The data of the token.
	GetData() any

	// GetPos returns the position of the token in the input string.
	//
	// Returns:
	//
	//   - int: The position of the token in the input string.
	GetPos() int

	// GetLookahead returns the next token in the input string.
	//
	// Returns:
	//
	//   - LeafToken: The next token in the input string.
	GetLookahead() *LeafToken

	// SetLookahead sets the next token in the input string.
	//
	// Parameters:
	//
	//   - lookahead: The next token in the input string.
	SetLookahead(lookahead *LeafToken)

	fmt.Stringer
}

// LeafToken represents a token that contains a single piece of data.
type LeafToken struct {
	// ID is the identifier of the token.
	ID string

	// Data is the data of the token.
	Data string

	// At is the position of the token in the input string.
	At int

	// Lookahead is the next token in the input string.
	Lookahead *LeafToken
}

// String is a method of fmt.Stringer interface.
//
// It should only be used for debugging and logging purposes.
//
// Returns:
//
//   - string: A string representation of the leaf token.
func (t *LeafToken) String() string {
	if t == nil {
		return "LeafToken[nil]"
	}

	return fmt.Sprintf("LeafToken[id=%s, data='%s', at=%d]",
		t.ID,
		t.Data,
		t.At,
	)
}

// GetID returns the identifier of the token.
//
// Returns:
//
//   - string: The identifier of the token.
func (t *LeafToken) GetID() string {
	return t.ID
}

// GetData returns the data of the token.
//
// Returns:
//
//   - any: The data of the token.
func (t *LeafToken) GetData() any {
	return t.Data
}

// GetPos returns the position of the token in the input string.
//
// Returns:
//
//   - int: The position of the token in the input string.
func (t *LeafToken) GetPos() int {
	return t.At
}

// GetLookahead returns the next token in the input string.
//
// Returns:
//
//   - LeafToken: The next token in the input string.
func (t *LeafToken) GetLookahead() *LeafToken {
	return t.Lookahead
}

// SetLookahead sets the next token in the input string.
//
// Parameters:
//
//   - lookahead: The next token in the input string.
func (t *LeafToken) SetLookahead(lookahead *LeafToken) {
	t.Lookahead = lookahead
}

// NewLeafToken creates a new leaf token with the given identifier, data, and position.
//
// Parameters:
//
//   - id: The identifier of the token.
//   - data: The data of the token.
//   - at: The position of the token in the input string.
//
// Returns:
//
//   - *LeafToken: A pointer to the new leaf token.
func NewLeafToken(id string, data string, at int) *LeafToken {
	return &LeafToken{
		id,
		data,
		at,
		nil,
	}
}

// NonLeafToken represents a token that contains multiple pieces of data.
type NonLeafToken struct {
	// ID is the identifier of the token.
	ID string

	// Data is the data of the token.
	Data []Tokener

	// At is the position of the token in the input string.
	At int

	// Lookahead is the next token in the input string.
	Lookahead *LeafToken
}

// String is a method of fmt.Stringer interface.
//
// It should only be used for debugging and logging purposes.
//
// Returns:
//
//   - string: A string representation of the non-leaf token.
func (t *NonLeafToken) String() string {
	if t == nil {
		return "NonLeafToken[nil]"
	}

	if len(t.Data) == 0 {
		return fmt.Sprintf("NonLeafToken[id=%s, data=[], at=%d]",
			t.ID,
			t.At,
		)
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "NonLeafToken[id=%s, data=[%v", t.ID, t.Data[0])

	for _, token := range t.Data[1:] {
		fmt.Fprintf(&builder, ", %v", token)
	}

	fmt.Fprintf(&builder, "], at=%d]", t.At)

	return builder.String()
}

// GetID returns the identifier of the token.
//
// Returns:
//
//   - string: The identifier of the token.
func (t *NonLeafToken) GetID() string {
	return t.ID
}

// GetData returns the data of the token.
//
// Returns:
//
//   - any: The data of the token.
func (t *NonLeafToken) GetData() any {
	return t.Data
}

// GetPos returns the position of the token in the input string.
//
// Returns:
//
//   - int: The position of the token in the input string.
func (t *NonLeafToken) GetPos() int {
	return t.At
}

// GetLookahead returns the next token in the input string.
//
// Returns:
//
//   - LeafToken: The next token in the input string.
func (t *NonLeafToken) GetLookahead() *LeafToken {
	return t.Lookahead
}

// SetLookahead sets the next token in the input string.
//
// Parameters:
//
//   - lookahead: The next token in the input string.
func (t *NonLeafToken) SetLookahead(lookahead *LeafToken) {
	t.Lookahead = lookahead
}

// NewNonLeafToken creates a new non-leaf token with the given identifier, data, and position.
//
// Parameters:
//
//   - id: The identifier of the token.
//   - data: The data of the token.
//   - at: The position of the token in the input string.
//
// Returns:
//
//   - *NonLeafToken: A pointer to the new non-leaf token.
func NewNonLeafToken(id string, at int, data ...Tokener) *NonLeafToken {
	return &NonLeafToken{
		ID:        id,
		At:        at,
		Data:      data,
		Lookahead: nil,
	}
}
