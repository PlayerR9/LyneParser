package FrontierEvaluation

// ErrNoAcceptance is an error type for when the result is not accepted.
type ErrNoAcceptance struct{}

// Error returns the error message: "result is not accepted".
//
// Returns:
// 	- string: the error message
func (e *ErrNoAcceptance) Error() string {
	return "result is not accepted"
}

// NewErrNoAcceptance creates a new ErrNoAcceptance error.
//
// Returns:
// 	- *ErrNoAcceptance: the new error
func NewErrNoAcceptance() *ErrNoAcceptance {
	return &ErrNoAcceptance{}
}
