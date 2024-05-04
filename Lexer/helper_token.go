package Lexer

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
)

// TokenStatus represents the status of a token.
type TokenStatus int

const (
	// TkComplete represents a token that has been fully lexed.
	TkComplete TokenStatus = iota

	// TkIncomplete represents a token that has not been fully lexed.
	TkIncomplete

	// TkError represents a token that has an error.
	TkError
)

// String is a method of fmt.Stringer that returns the string
// representation of the TokenStatus.
//
// Returns:
//   - string: The string representation of the TokenStatus
func (s TokenStatus) String() string {
	return [...]string{
		"complete",
		"incomplete",
		"error",
	}[s]
}

// helperToken is a wrapper around a *gr.LeafToken that adds a status field.
type helperToken struct {
	// Status is the status of the token.
	Status TokenStatus

	// Tok is the *gr.LeafToken to wrap.
	Tok *gr.LeafToken
}

// newHelperToken creates a new helperToken, marking it as incomplete.
//
// Parameters:
//   - tok: The *gr.LeafToken to wrap.
//
// Returns:
//   - helperToken: The new helperToken.
func newHelperToken(tok *gr.LeafToken) *helperToken {
	return &helperToken{
		Status: TkIncomplete,
		Tok:    tok,
	}
}

// SetStatus sets the status of the token.
//
// Parameters:
//   - status: The status to set.
func (ht *helperToken) SetStatus(status TokenStatus) {
	ht.Status = status
}

// GetPos returns the position of the token in the input string.
//
// Returns:
//   - int: The position of the token in the input string. If the token is nil, -1 is returned.
func (ht *helperToken) GetPos() int {
	if ht.Tok == nil {
		return -1
	}

	return ht.Tok.GetPos()
}

// GetData returns the data of the token.
//
// Returns:
//   - string: The data of the token. If the token is nil, an empty string is returned.
func (ht *helperToken) GetData() string {
	if ht.Tok == nil {
		return ""
	}

	return ht.Tok.Data
}

func (ht *helperToken) GetID() string {
	if ht.Tok == nil {
		return ""
	}

	return ht.Tok.ID
}
