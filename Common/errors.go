package Common

// ErrNoMatches is an error that is returned when there are no
// matches at a position.
type ErrNoMatches struct{}

// Error returns the error message: "no matches".
//
// Returns:
//  	- string: The error message.
func (e *ErrNoMatches) Error() string {
	return "no matches"
}

// NewErrNoMatches creates a new error of type *ErrNoMatches.
//
// Returns:
//  	- *ErrNoMatches: The new error.
func NewErrNoMatches() *ErrNoMatches {
	return &ErrNoMatches{}
}

// ErrCycleDetected is an error that is returned when a cycle is detected.
type ErrCycleDetected struct{}

// Error returns the error message: "cycle detected".
//
// Returns:
//  	- string: The error message.
func (e *ErrCycleDetected) Error() string {
	return "cycle detected"
}

// NewErrCycleDetected creates a new error of type *ErrCycleDetected.
//
// Returns:
//  	- *ErrCycleDetected: The new error.
func NewErrCycleDetected() *ErrCycleDetected {
	return &ErrCycleDetected{}
}
