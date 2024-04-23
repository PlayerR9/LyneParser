package Parser

// ErrNoAccept is an error that is returned when the parser reaches the end of the
// input stream without accepting the input stream.
type ErrNoAccept struct{}

// Error is a method of the error interface.
//
// Returns:
//
//   - string: The error message.
func (e *ErrNoAccept) Error() string {
	return "reached end of input stream without accepting"
}

// NewErrNoAccept creates a new ErrNoAccept error.
//
// Returns:
//
//   - *ErrNoAccept: A pointer to the new ErrNoAccept error.
func NewErrNoAccept() *ErrNoAccept {
	return &ErrNoAccept{}
}
