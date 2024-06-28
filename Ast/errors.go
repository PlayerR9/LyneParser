package Ast

import (
	"strconv"
	"strings"

	fs "github.com/PlayerR9/MyGoLib/Formatting/Strings"
)

// ErrMissingFields is an error for when fields are missing.
type ErrMissingFields struct {
	// Missings is the list of missing fields.
	Missings []string
}

// Error implements the error interface.
//
// Messages:
//   - "missing <field> fields".
//   - "no field is missing" if there are no missing fields.
func (e *ErrMissingFields) Error() string {
	if len(e.Missings) == 0 {
		return "no field is missing"
	}

	var builder strings.Builder

	builder.WriteString("missing ")

	if len(e.Missings) == 1 {
		builder.WriteString(strconv.Quote(e.Missings[0]))
		builder.WriteString(" field")
	} else {
		vals := fs.OrString(e.Missings...)
		builder.WriteString(strconv.Quote(vals))
		builder.WriteString(" fields")
	}

	return builder.String()
}

// NewErrMissingFields creates a new ErrMissingFields.
//
// Parameters:
//   - missings: The list of missing fields.
//
// Returns:
//   - *ErrMissingFields: A pointer to the new error.
func NewErrMissingFields(missings ...string) *ErrMissingFields {
	e := &ErrMissingFields{
		Missings: missings,
	}
	return e
}

// ErrTooManyFields is an error for when too many fields are present.
type ErrTooManyFields struct {
	// Got is the number of fields present.
	Got int

	// Wanted is the number of fields expected.
	Wanted int
}

// Error implements the error interface.
//
// Message: "expected <wanted> fields, got <got> instead".
func (e *ErrTooManyFields) Error() string {
	var builder strings.Builder

	builder.WriteString("expected ")
	builder.WriteString(strconv.Itoa(e.Wanted))
	builder.WriteString(" fields, got ")
	builder.WriteString(strconv.Itoa(e.Got))
	builder.WriteString(" instead")

	return builder.String()
}

// NewErrTooManyFields creates a new ErrTooManyFields.
//
// Parameters:
//   - got: The number of fields present.
//   - wanted: The number of fields expected.
//
// Returns:
//   - *ErrTooManyFields: A pointer to the new error.
func NewErrTooManyFields(got, wanted int) *ErrTooManyFields {
	e := &ErrTooManyFields{
		Got:    got,
		Wanted: wanted,
	}
	return e
}

// ErrAmbiguousGrammar is an error for when a grammar is ambiguous.
type ErrAmbiguousGrammar struct{}

// Error implements the error interface.
//
// Message: "ambiguous grammar".
func (e *ErrAmbiguousGrammar) Error() string {
	return "ambiguous grammar"
}

// NewErrAmbiguousGrammar creates a new ErrAmbiguousGrammar.
//
// Returns:
//   - *ErrAmbiguousGrammar: A pointer to the new error.
func NewErrAmbiguousGrammar() *ErrAmbiguousGrammar {
	e := &ErrAmbiguousGrammar{}
	return e
}
