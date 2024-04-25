package Grammar

// ErrNoMoreTokens is an error that indicates that
// there are no more tokens in the stream.
type ErrNoMoreTokens struct{}

// Error is a method of errors that returns the error message.
//
// Returns:
//  	- string: The error message.
func (e *ErrNoMoreTokens) Error() string {
	return "No more tokens in the stream."
}

// NewErrNoMoreTokens creates a new ErrNoMoreTokens error.
//
// Returns:
//  	- *ErrNoMoreTokens: A pointer to the new error.
func NewErrNoMoreTokens() *ErrNoMoreTokens {
	return &ErrNoMoreTokens{}
}
