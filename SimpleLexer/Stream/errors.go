package Stream

import "errors"

// ErrStreamExhausted is the error returned when the stream is exhausted.
type ErrStreamExhausted struct{}

// Error implements the error interface.
//
// Message: "stream is exhausted"
func (e *ErrStreamExhausted) Error() string {
	return "stream is exhausted"
}

// NewStreamExhausted creates a new stream exhausted error.
//
// Returns:
//   - *ErrStreamExhausted: The stream exhausted error.
func NewStreamExhausted() *ErrStreamExhausted {
	e := &ErrStreamExhausted{}
	return e
}

// IsStreamExhausted checks if the error is a stream exhausted error.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - bool: True if the error is a stream exhausted error.
//     False otherwise.
func IsStreamExhausted(err error) bool {
	if err == nil {
		return false
	}

	var streamExhaustedErr *ErrStreamExhausted

	ok := errors.As(err, &streamExhaustedErr)
	return ok
}
