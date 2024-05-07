package Interfaces

import (
	"fmt"
	"strings"
)

// StringOf converts any type to a string.
//
// The function converts the input element to a string using the following rules:
//   - If the element is a string, it is returned as is.
//   - If the element implements the fmt.Stringer interface, the String method is called.
//   - If the element is an error, the Error method is called.
//   - If the element is a byte slice, it is converted to a string.
//   - If the element is a rune slice, it is converted to a string.
//   - If the element implements the ffs.FStringer interface, the FString function is called.
//   - Otherwise, the element is converted to a string using the fmt.Sprintf function.
//
// Parameters:
//   - elem: The element to convert to a string.
//
// Returns:
//   - string: The string representation of the element.
func StringOf(elem any) string {
	if elem == nil {
		return ""
	}

	switch elem := elem.(type) {
	case string:
		return elem
	case fmt.Stringer:
		return elem.String()
	case error:
		return elem.Error()
	case []byte:
		return string(elem)
	case []rune:
		return string(elem)
	case FStringer:
		return strings.Join(elem.FString(0), "\n")
	default:
		return fmt.Sprintf("%v", elem)
	}
}

// FStringer is an interface that defines the behavior of a type that can be
// converted to a string representation.
type FStringer interface {
	// FString returns a string representation of the object.
	//
	// Parameters:
	//   - int: The current indentation level.
	//
	// Returns:
	//   - []string: A slice of strings that represent the object.
	FString(int) []string
}

// FString is a function that returns a string representation of an object that
// implements the FStringer interface.
//
// It joins the strings returned by the FString method of the object using a newline
// character with no indentation at the beginning.
//
// Parameters:
//   - elem: The object that implements the FStringer interface.
//
// Returns:
//   - string: A string representation of the object.
func FStringOf(elem any) string {
	if elem == nil {
		return ""
	}

	switch elem := elem.(type) {
	case FStringer:
		return strings.Join(elem.FString(0), "\n")
	case string:
		return elem
	case fmt.Stringer:
		return elem.String()
	case error:
		return elem.Error()
	case []byte:
		return string(elem)
	case []rune:
		return string(elem)
	default:
		return fmt.Sprintf("%v", elem)
	}
}
