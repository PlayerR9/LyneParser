package LyneLexer

import "fmt"

type ErrInvalidTokenName struct {
	name string
}

func (e *ErrInvalidTokenName) Error() string {
	return fmt.Sprintf("expected a non-empty string starting with a lowercase letter, got %s instead", e.name)
}

func NewErrInvalidTokenName(name string) *ErrInvalidTokenName {
	return &ErrInvalidTokenName{name: name}
}
