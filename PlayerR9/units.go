package PlayerR9

import (
	"strings"
)

// Runer is an interface that provides a method to get the runes of a string.
type Runer interface {
	// Runes returns the runes of the string.
	//
	// Returns:
	//   - []rune: The runes of the string.
	Runes() []rune
}

// StringToLines splits a string into lines.
//
// Parameters:
//   - str: The string to split into lines.
//
// Returns:
//   - []string: The lines of the string.
func StringToLines(str string) []string {
	var lines []string
	var builder strings.Builder

	for _, c := range str {
		if c == '\n' {
			lines = append(lines, builder.String())
			builder.Reset()
		} else {
			builder.WriteRune(c)
		}
	}

	if builder.Len() > 0 {
		lines = append(lines, builder.String())
	}

	return lines
}

// StringToLines splits a string into lines.
//
// Parameters:
//   - elem: The element to split into lines.
//   - f: The function to execute on each rune to convert it to the desired type.
//
// Returns:
//   - [][]O: The lines of the elements.
//   - error: An error if it occurs during the conversion.
func AnyToLines[I Runer, O any](elem I, f func(rune) (O, error)) ([][]O, error) {
	lines := [][]O{make([]O, 0)}
	lastLine := 0

	chars := elem.Runes()

	for _, c := range chars {
		if c == '\n' {
			lines = append(lines, make([]O, 0))
			lastLine++
		} else {
			newO, err := f(c)
			if err != nil {
				return lines, err
			}

			lines[lastLine] = append(lines[lastLine], newO)
		}
	}

	return lines, nil
}
