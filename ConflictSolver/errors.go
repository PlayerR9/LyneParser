package ConflictSolver

import "fmt"

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

// ErrInvalidPosition is an error that is returned when a position is invalid (i.e., less than 0).
type ErrInvalidPosition struct{}

// Error returns the error message: "invalid position".
//
// Returns:
//   - string: The error message.
func (e *ErrInvalidPosition) Error() string {
	return "invalid position"
}

// NewErrInvalidPosition creates a new error of type *ErrInvalidPosition.
//
// Returns:
//   - *ErrInvalidPosition: A pointer to the new error.
func NewErrInvalidPosition() *ErrInvalidPosition {
	return &ErrInvalidPosition{}
}

// ErrCannotCreateItem is an error that is returned when an item cannot be created.
type ErrCannotCreateItem struct{}

// Error returns the error message: "cannot create item".
//
// Returns:
//   - string: The error message.
func (e *ErrCannotCreateItem) Error() string {
	return "cannot create item"
}

// NewErrCannotCreateItem creates a new error of type *ErrCannotCreateItem.
//
// Returns:
//   - *ErrCannotCreateItem: A pointer to the new error.
func NewErrCannotCreateItem() *ErrCannotCreateItem {
	return &ErrCannotCreateItem{}
}

// ErrHelpersConflictingSize is an error that is returned when helpers have conflicting sizes.
type ErrHelpersConflictingSize struct{}

// Error returns the error message: "helpers have conflicting sizes".
//
// Returns:
//   - string: The error message.
func (e *ErrHelpersConflictingSize) Error() string {
	return "helpers have conflicting sizes"
}

// NewErrHelpersConflictingSize creates a new error of type *ErrHelpersConflictingSize.
//
// Returns:
//   - *ErrHelpersConflictingSize: A pointer to the new error.
func NewErrHelpersConflictingSize() *ErrHelpersConflictingSize {
	return &ErrHelpersConflictingSize{}
}

// ErrNoActionProvided is an error that is returned when no action is provided.
type ErrNoActionProvided struct{}

// Error returns the error message: "no action provided".
//
// Returns:
//   - string: The error message.
func (e *ErrNoActionProvided) Error() string {
	return "no action provided"
}

// NewErrNoActionProvided creates a new error of type *ErrNoActionProvided.
//
// Returns:
//   - *ErrNoActionProvided: A pointer to the new error.
func NewErrNoActionProvided() *ErrNoActionProvided {
	return &ErrNoActionProvided{}
}

// Err0thRhsNotSet is an error that is returned when the 0th right-hand side is not set.
type Err0thRhsNotSet struct{}

// Error returns the error message: "0th RHS not set".
//
// Returns:
//   - string: The error message.
func (e *Err0thRhsNotSet) Error() string {
	return "0th RHS not set"
}

// NewErr0thRhsNotSet creates a new error of type *Err0thRhsNotSet.
//
// Returns:
//   - *Err0thRhsNotSet: A pointer to the new error.
func NewErr0thRhsNotSet() *Err0thRhsNotSet {
	return &Err0thRhsNotSet{}
}

/////////////////////////////////////////////////////////////

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
