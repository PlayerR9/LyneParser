package Common

// ErrCycleDetected is an error that is returned when a cycle is detected.
type ErrCycleDetected struct{}

// Error returns the error message: "cycle detected".
//
// Returns:
//   - string: The error message.
func (e *ErrCycleDetected) Error() string {
	return "cycle detected"
}

// NewErrCycleDetected creates a new error of type *ErrCycleDetected.
//
// Returns:
//   - *ErrCycleDetected: The new error.
func NewErrCycleDetected() *ErrCycleDetected {
	return &ErrCycleDetected{}
}
