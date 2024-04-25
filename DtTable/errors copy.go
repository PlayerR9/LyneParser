package DtTable

import (
	"fmt"
)

// ErrInvalidCharacter is an error that indicates that an invalid character
// was found in the input.
type ErrInvalidCharacter struct {
	Character rune
}

// Error is a method of errors that returns the error message.
//
// Returns:
//   - string: The error message.
func (e *ErrInvalidCharacter) Error() string {
	return fmt.Sprintf("invalid character: %c", e.Character)
}

// NewErrInvalidCharacter creates a new ErrInvalidCharacter error.
//
// Parameters:
//   - character: The invalid character.
//
// Returns:
//   - *ErrInvalidCharacter: A pointer to the new error.
func NewErrInvalidCharacter(character rune) *ErrInvalidCharacter {
	return &ErrInvalidCharacter{Character: character}
}
