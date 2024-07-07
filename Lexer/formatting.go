package Lexer

import (
	"strconv"
	"strings"
)

// SyntaxErrorer is an interface that represents a syntax error.
type SyntaxErrorer interface {
	// Display displays the syntax error.
	//
	// Returns:
	//   - []string: The lines to display.
	Display() []string

	// GetPosition gets the position of the syntax error.
	//
	// Returns:
	//   - int: The position.
	GetPosition() int
}

// ClosestSyntaxError is a syntax error that represents a syntax error
// where the closest keyword is suggested.
type ClosestSyntaxError struct {
	// unexpected_char is the unexpected character.
	reason *ErrUnexpectedChar

	// end_idx is the end index.
	end_idx int

	// closest is the closest keyword.
	closest string

	// actual is the actual keyword.
	actual string
}

// GetPosition implements the SyntaxErrorer interface.
func (cse *ClosestSyntaxError) GetPosition() int {
	return cse.end_idx
}

// Display implements the SyntaxErrorer interface.
func (cse *ClosestSyntaxError) Display() []string {
	var reason_msg string

	if cse.reason == nil {
		reason_msg = "found an unexpected character"
	} else {
		reason_msg = cse.reason.Error()
	}

	var error_detail strings.Builder

	error_detail.WriteRune('\t')
	error_detail.WriteString(reason_msg)

	var suggestion strings.Builder

	suggestion.WriteString("\tDid you mean ")
	suggestion.WriteString(strconv.Quote(cse.closest))
	suggestion.WriteString(" instead of ")
	suggestion.WriteString(strconv.Quote(cse.actual))
	suggestion.WriteString("?")

	lines := []string{
		"Syntax Error:",
		error_detail.String(),
		"",
		"",
		"Suggestion:",
		suggestion.String(),
	}

	return lines
}

// NewClosestSyntaxError creates a new ClosestSyntaxError.
//
// Parameters:
//   - unexpected_char: The unexpected character.
//   - previous_char: The previous character.
//   - got_char: The got character.
//   - end_idx: The end index.
//
// Returns:
//   - *ClosestSyntaxError: The new ClosestSyntaxError.
func NewClosestSyntaxError(reason *ErrUnexpectedChar, closest, actual string, end_idx int) *ClosestSyntaxError {
	cse := &ClosestSyntaxError{
		reason:  reason,
		end_idx: end_idx,
		closest: closest,
		actual:  actual,
	}

	// closest, err := laven_table.GetClosest(f.data[f.idx:cse.end_idx])
	// uc.Assert(err == nil, "In GetClosest, no error should occur")

	// actual := f.data[f.idx:cse.end_idx]

	return cse
}

// UnrecognizedSyntaxError is a syntax error that represents an unrecognized syntax error.
type UnrecognizedSyntaxError struct {
	// char is the character.
	char rune

	// at is the position.
	at int
}

// GetPosition implements the SyntaxErrorer interface.
func (use *UnrecognizedSyntaxError) GetPosition() int {
	return use.at
}

// Display implements the SyntaxErrorer interface.
func (use *UnrecognizedSyntaxError) Display() []string {

	var errorDetail strings.Builder

	errorDetail.WriteRune('\t')
	errorDetail.WriteString(strconv.QuoteRune(use.char))
	errorDetail.WriteString(" is not a recognized character")

	lines := []string{
		"Syntax Error:",
		errorDetail.String(),
		"",
		"",
		"Suggestion:",
		"\tYou may want to check your code for any typos",
	}

	return lines
}

// NewUnrecognizedSyntaxError creates a new UnrecognizedSyntaxError.
//
// Parameters:
//   - char: The character.
//   - at: The position.
//
// Returns:
//   - *UnrecognizedSyntaxError: The new UnrecognizedSyntaxError.
func NewUnrecognizedSyntaxError(char rune, at int) *UnrecognizedSyntaxError {
	use := &UnrecognizedSyntaxError{
		char: char,
		at:   at,
	}
	return use
}

// GenericSyntaxError is a syntax error that represents a generic syntax error.
type GenericSyntaxError struct {
	// at is the position.
	at int

	// Message is the error message.
	Message string

	// Suggestion is the suggestion.
	Suggestion string
}

// GetPosition implements the SyntaxErrorer interface.
func (gse *GenericSyntaxError) GetPosition() int {
	return gse.at
}

// Display implements the SyntaxErrorer interface.
func (gse *GenericSyntaxError) Display() []string {
	var errorDetail strings.Builder

	errorDetail.WriteRune('\t')
	errorDetail.WriteString(gse.Message)

	lines := []string{
		"Syntax Error:",
		errorDetail.String(),
	}

	if gse.Suggestion != "" {
		var suggestion strings.Builder

		suggestion.WriteString("\t")
		suggestion.WriteString(gse.Suggestion)

		lines = append(lines, "", "", "Suggestion:", suggestion.String())
	}

	return lines
}

// NewGenericSyntaxError creates a new GenericSyntaxError.
//
// Parameters:
//   - at: The position.
//   - message: The error message.
//   - suggestion: The suggestion.
//
// Returns:
//   - *GenericSyntaxError: The new GenericSyntaxError.
func NewGenericSyntaxError(at int, message, suggestion string) *GenericSyntaxError {
	gse := &GenericSyntaxError{
		at:         at,
		Message:    message,
		Suggestion: suggestion,
	}
	return gse
}
