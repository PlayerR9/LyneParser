package ConflictSolver

import "fmt"

// ErrNoElementsFound is an error that is returned when no
// elements are found for a symbol.
type ErrNoElementsFound struct {
	// Symbol is the symbol for which no elements were found.
	Symbol string
}

// Error returns the error message: "no elements found for symbol (symbol)".
//
// Returns:
//   - string: The error message.
func (e *ErrNoElementsFound) Error() string {
	return fmt.Sprintf("no elements found for symbol (%s)", e.Symbol)
}

// NewErrNoElementsFound creates a new error of type *ErrNoElementsFound.
//
// Parameters:
//   - symbol: The symbol for which no elements were found.
//
// Returns:
//   - *ErrNoElementsFound: A pointer to the new error.
func NewErrNoElementsFound(symbol string) *ErrNoElementsFound {
	return &ErrNoElementsFound{
		Symbol: symbol,
	}
}

// ErrItemIsNil is an error that is returned when an item is nil.
type ErrItemIsNil struct{}

// Error returns the error message: "item is nil".
//
// Returns:
//   - string: The error message.
func (e *ErrItemIsNil) Error() string {
	return "item is nil"
}

// NewErrItemIsNil creates a new error of type *ErrItemIsNil.
//
// Returns:
//   - *ErrItemIsNil: A pointer to the new error.
func NewErrItemIsNil() *ErrItemIsNil {
	return &ErrItemIsNil{}
}

// ErrHelper is an error that is returned when something goes wrong
// with a helper.
type ErrHelper struct {
	// Elem is the helper that caused the error.
	Elem *Helper

	// Reason is the reason for the error.
	Reason error
}

// Error returns the error message: "helper (elem) error: (reason)".
//
// Returns:
//   - string: The error message.
func (e *ErrHelper) Error() string {
	var elem string

	if e.Elem == nil {
		elem = "no helper"
	} else {
		elem = e.Elem.String()
	}

	if e.Reason == nil {
		return fmt.Sprintf("something went wrong with helper (%s)", elem)
	} else {
		return fmt.Sprintf("helper (%s) error: %s", elem, e.Reason.Error())
	}
}

// NewErrHelper creates a new error of type *ErrHelper.
//
// Parameters:
//   - elem: The helper that caused the error.
//   - reason: The reason for the error.
//
// Returns:
//   - *ErrHelper: A pointer to the new error.
func NewErrHelper(elem *Helper, reason error) *ErrHelper {
	return &ErrHelper{
		Elem:   elem,
		Reason: reason,
	}
}
