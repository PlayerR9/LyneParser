package LyneParser

import "fmt"

type ErrUnexpectedToken struct {
	expected string
	before   string
	got      *Node
}

func (e *ErrUnexpectedToken) Error() string {
	if e.got == nil {
		return fmt.Sprintf("expected %s before %s. Got Nothing instead", e.expected, e.before)
	} else {
		return fmt.Sprintf("expected %s before %s. Got %v instead", e.expected, e.before, e.got)
	}
}

func NewErrUnexpectedToken(expected, before string, got *Node) *ErrUnexpectedToken {
	return &ErrUnexpectedToken{expected, before, got}
}

type ErrParsing struct {
	line, column int
	reason       error
}

func NewErrParsing(line, column int, reason error) *ErrParsing {
	return &ErrParsing{
		line:   line,
		column: column,
		reason: reason,
	}
}

func (e *ErrParsing) Error() string {
	return fmt.Sprintf("parse error at line %d, column %d: %s", e.line, e.column, e.reason.Error())
}

func (e *ErrParsing) Unwrap() error {
	return e.reason
}

func (e *ErrParsing) Wrap(reason error) *ErrParsing {
	e.reason = reason
	return e
}
