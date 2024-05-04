package Lexer

// ErrNoTokensToLex is an error that is returned when there are no tokens to lex.
type ErrNoTokensToLex struct{}

// Error returns the error message: "no tokens to lex".
//
// Returns:
//  	- string: The error message.
func (e *ErrNoTokensToLex) Error() string {
	return "no tokens to lex"
}

// NewErrNoTokensToLex creates a new error of type *ErrNoTokensToLex.
//
// Returns:
//  	- *ErrNoTokensToLex: The new error.
func NewErrNoTokensToLex() *ErrNoTokensToLex {
	return &ErrNoTokensToLex{}
}

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

// ErrInvalidToken is an error that is returned when an invalid token
// is found at a position.
type ErrInvalidToken struct{}

// Error returns the error message: "invalid token".
//
// Returns:
//  	- string: The error message.
func (e *ErrInvalidToken) Error() string {
	return "invalid token"
}

// NewErrInvalidTokencreates a new error of type *ErrInvalidToken.
//
// Returns:
//  	- *ErrInvalidToken: The new error.
func NewErrInvalidToken() *ErrInvalidToken {
	return &ErrInvalidToken{}
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
