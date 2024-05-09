package TreeExplorer

// ErrEmptyInput is an error that is returned when the input is empty.
type ErrEmptyInput struct{}

// Error returns the error message: "empty input".
//
// Returns:
//  	- string: The error message.
func (e *ErrEmptyInput) Error() string {
	return "empty input"
}

// NewErrEmptyInput creates a new error of type *ErrEmptyInput.
//
// Returns:
//  	- *ErrEmptyInput: The new error.
func NewErrEmptyInput() *ErrEmptyInput {
	return &ErrEmptyInput{}
}

// ErrAllMatchesFailed is an error that is returned when all matches
// fail.
type ErrAllMatchesFailed struct{}

// Error returns the error message: "all matches failed".
//
// Returns:
//  	- string: The error message.
func (e *ErrAllMatchesFailed) Error() string {
	return "all matches failed"
}

// NewErrAllMatchesFailed creates a new error of type *ErrAllMatchesFailed.
//
// Returns:
//  	- *ErrAllMatchesFailed: The new error.
func NewErrAllMatchesFailed() *ErrAllMatchesFailed {
	return &ErrAllMatchesFailed{}
}

// ErrNilRoot is an error that is returned when the root is nil.
type ErrNilRoot struct{}

// Error returns the error message: "root is nil".
//
// Returns:
//  	- string: The error message.
func (e *ErrNilRoot) Error() string {
	return "root is nil"
}

// NewErrNilRoot creates a new error of type *ErrNilRoot.
//
// Returns:
//  	- *ErrNilRoot: The new error.
func NewErrNilRoot() *ErrNilRoot {
	return &ErrNilRoot{}
}

// ErrInvalidElement is an error that is returned when an invalid element
// is found.
type ErrInvalidElement struct{}

// Error returns the error message: "invalid element".
//
// Returns:
//  	- string: The error message.
func (e *ErrInvalidElement) Error() string {
	return "invalid element"
}

// NewErrInvalidElement creates a new error of type *ErrInvalidElement.
//
// Returns:
//  	- *ErrInvalidElement: The new error.
func NewErrInvalidElement() *ErrInvalidElement {
	return &ErrInvalidElement{}
}
