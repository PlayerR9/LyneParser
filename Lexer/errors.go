package Lexer

import (
	"strconv"
	"strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// ErrUnexpectedChar is an error that is returned when an unexpected character is found.
type ErrUnexpectedChar struct {
	// Expected is the expected rune.
	Expected rune

	// Previous is the previous rune.
	Previous rune

	// Got is the got rune.
	Got *rune
}

// Error implements the error interface.
//
// Message: "expected <expected> after <previous>, found <got> instead".
func (e *ErrUnexpectedChar) Error() string {
	var got_str string

	if e.Got != nil {
		got_str = strconv.QuoteRune(*e.Got)
	} else {
		got_str = "nothing"
	}

	var builder strings.Builder

	builder.WriteString("expected ")
	builder.WriteString(strconv.QuoteRune(e.Expected))
	builder.WriteString(" after ")
	builder.WriteString(strconv.QuoteRune(e.Previous))
	builder.WriteString(", found ")
	builder.WriteString(got_str)
	builder.WriteString(" instead")

	msg := builder.String()
	return msg
}

// NewErrUnexpectedChar creates a new ErrUnexpectedChar.
//
// Parameters:
//   - expected: The expected rune.
//   - previous: The previous rune.
//   - got: The got rune.
//
// Returns:
//   - *ErrUnexpectedChar: The new ErrUnexpectedChar.
func NewErrUnexpectedChar(expected, previous rune, got *rune) *ErrUnexpectedChar {
	e := &ErrUnexpectedChar{
		Expected: expected,
		Previous: previous,
		Got:      got,
	}
	return e
}

type ErrLexerError[T gr.TokenTyper] struct {
	At   int
	Prev []*gr.Token[T]
}

func (e *ErrLexerError[T]) Error() string {
	var builder strings.Builder

	builder.WriteString("no matches found at ")
	builder.WriteString(strconv.Itoa(e.At))

	return builder.String()
}

func NewErrLexerError[T gr.TokenTyper](at int, prev []*gr.Token[T]) *ErrLexerError[T] {
	e := &ErrLexerError[T]{
		At:   at,
		Prev: prev,
	}
	return e
}

// ErrNoMatches is an error that is returned when there are no
// matches at a position.
type ErrNoMatches struct{}

// Error returns the error message: "no matches".
//
// Returns:
//   - string: The error message.
func (e *ErrNoMatches) Error() string {
	return "no matches"
}

// NewErrNoMatches creates a new error of type *ErrNoMatches.
//
// Returns:
//   - *ErrNoMatches: The new error.
func NewErrNoMatches() *ErrNoMatches {
	return &ErrNoMatches{}
}

// ErrAllMatchesFailed is an error that is returned when all matches
// fail.
type ErrAllMatchesFailed struct{}

// Error returns the error message: "all matches failed".
//
// Returns:
//   - string: The error message.
func (e *ErrAllMatchesFailed) Error() string {
	return "all matches failed"
}

// NewErrAllMatchesFailed creates a new error of type *ErrAllMatchesFailed.
//
// Returns:
//   - *ErrAllMatchesFailed: The new error.
func NewErrAllMatchesFailed() *ErrAllMatchesFailed {
	return &ErrAllMatchesFailed{}
}

// ErrInvalidElement is an error that is returned when an invalid element
// is found.
type ErrInvalidElement struct{}

// Error returns the error message: "invalid element".
//
// Returns:
//   - string: The error message.
func (e *ErrInvalidElement) Error() string {
	return "invalid element"
}

// NewErrInvalidElement creates a new error of type *ErrInvalidElement.
//
// Returns:
//   - *ErrInvalidElement: The new error.
func NewErrInvalidElement() *ErrInvalidElement {
	return &ErrInvalidElement{}
}

// IsDone checks if an error is a completion error or nil.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - bool: True if the error is a completion error or nil.
//     False otherwise.
func IsDone(err error) bool {
	if err == nil {
		return true
	}

	ok := uc.Is[*uc.ErrExhaustedIter](err)
	return ok
}
