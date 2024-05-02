package Grammar

// ErrNoProductionRulesFound is an error that is returned when no production rules
// are found in a grammar.
type ErrNoProductionRulesFound struct{}

// Error returns the error message: "no production rules found".
//
// Returns:
//  	- string: The error message.
func (e *ErrNoProductionRulesFound) Error() string {
	return "no production rules found"
}

// NewErrNoProductionRulesFound creates a new error of type *ErrNoProductionRulesFound.
//
// Returns:
//  	- *ErrNoProductionRulesFound: The new error.
func NewErrNoProductionRulesFound() *ErrNoProductionRulesFound {
	return &ErrNoProductionRulesFound{}
}
