package Ast

import (
	"fmt"
	"strings"

	fs "github.com/PlayerR9/MyGoLib/Formatting/Strings"

	gr "github.com/PlayerR9/LyneParser/Grammar"
)

// ErrAssumptionViolated is an error for when an assumption is violated.
type ErrAssumptionViolated struct {
	// Reason is the violated assumption.
	Reason error
}

// Error implements the error interface.
//
// Message: "assumption violation: (assumption)".
// If the assumption is nil, the message is "an assumption was violated".
func (e *ErrAssumptionViolated) Error() string {
	if e.Reason == nil {
		return "an assumption was violated"
	}

	var builder strings.Builder

	builder.WriteString("assumption violation: ")
	builder.WriteString(e.Reason.Error())

	return builder.String()
}

// NewErrAssumptionViolated creates a new ErrAssumptionViolated.
//
// Parameters:
//   - reason: The violated assumption.
//
// Returns:
//   - *ErrAssumptionViolated: A pointer to the new error.
func NewErrAssumptionViolated(reason error) *ErrAssumptionViolated {
	return &ErrAssumptionViolated{Reason: reason}
}

// ErrExpectedNonNil is an error for when a non-nil value is expected.
type ErrExpectedNonNil struct {
	// Expected is the expected type.
	Expected string
}

// Error is a method of error interface.
//
// Returns:
//
//   - string: The error message.
func (e *ErrExpectedNonNil) Error() string {
	return fmt.Sprintf("expected non-nil %s", e.Expected)
}

// NewErrExpectedNonNil creates a new ErrExpectedNonNil.
//
// Parameters:
//
//   - expected: The expected type.
//
// Returns:
//
//   - *ErrExpectedNonNil: A pointer to the new error.
func NewErrExpectedNonNil(expected string) *ErrExpectedNonNil {
	return &ErrExpectedNonNil{Expected: expected}
}

// ErrInvalidParsing is an error for when a parsing is invalid.
type ErrInvalidParsing struct {
	// Root is the root of the parsing where the error occurred.
	Root gr.Tokener
}

// Error is a method of error interface.
//
// Message: "invalid parsing: (root)".
func (e *ErrInvalidParsing) Error() string {
	return fmt.Sprintf("invalid parsing: %+v", e.Root)
}

// NewErrInvalidParsing creates a new ErrInvalidParsing.
//
// Parameters:
//   - root: The root of the parsing where the error occurred.
//
// Returns:
//   - *ErrInvalidParsing: A pointer to the new error.
func NewErrInvalidParsing(root gr.Tokener) *ErrInvalidParsing {
	return &ErrInvalidParsing{Root: root}
}

// ErrMissingFields is an error for when fields are missing.
type ErrMissingFields struct {
	// Missings is the list of missing fields.
	Missings []string
}

// Error is a method of error interface.
//
// Returns:
//
//   - string: The error message.
func (e *ErrMissingFields) Error() string {
	switch len(e.Missings) {
	case 0:
		return "no field is missing"
	case 1:
		return fmt.Sprintf("missing %q field", e.Missings[0])
	default:
		return fmt.Sprintf("missing %q fields", fs.OrString(e.Missings...))
	}
}

// NewErrMissingFields creates a new ErrMissingFields.
//
// Parameters:
//
//   - missings: The list of missing fields.
//
// Returns:
//
//   - *ErrMissingFields: A pointer to the new error.
func NewErrMissingFields(missings ...string) *ErrMissingFields {
	return &ErrMissingFields{Missings: missings}
}

// ErrTooManyFields is an error for when too many fields are present.
type ErrTooManyFields struct {
	// Got is the number of fields present.
	Got int

	// Wanted is the number of fields expected.
	Wanted int
}

// Error is a method of error interface.
//
// Returns:
//
//   - string: The error message.
func (e *ErrTooManyFields) Error() string {
	return fmt.Sprintf("expected %d fields, got %d instead", e.Wanted, e.Got)
}

// NewErrTooManyFields creates a new ErrTooManyFields.
//
// Parameters:
//
//   - got: The number of fields present.
//   - wanted: The number of fields expected.
//
// Returns:
//
//   - *ErrTooManyFields: A pointer to the new error.
func NewErrTooManyFields(got, wanted int) *ErrTooManyFields {
	return &ErrTooManyFields{Got: got, Wanted: wanted}
}

// ErrAmbiguousGrammar is an error for when a grammar is ambiguous.
type ErrAmbiguousGrammar struct{}

// Error is a method of error interface.
//
// Returns:
//
//   - string: The error message.
func (e *ErrAmbiguousGrammar) Error() string {
	return "ambiguous grammar"
}

// NewErrAmbiguousGrammar creates a new ErrAmbiguousGrammar.
//
// Returns:
//
//   - *ErrAmbiguousGrammar: A pointer to the new error.
func NewErrAmbiguousGrammar() *ErrAmbiguousGrammar {
	return &ErrAmbiguousGrammar{}
}
