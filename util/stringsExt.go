package util

import "strings"

// ByteSplitter splits a byte slice into multiple slices based on a separator byte.
// The separator byte is not included in the resulting slices.
//
// If the input slice is empty, the function returns nil.
//
// Parameters:
//
// 	- data: The byte slice to split.
// 	- sep: The separator byte.
//
// Returns:
//
// 	- [][]byte: A slice of byte slices.
func ByteSplitter(data []byte, sep byte) [][]byte {
	// FIXME: Remove this once MyGoLib is updated.
	if len(data) == 0 {
		return nil
	}

	slices := make([][]byte, 0)

	start := 0

	for i := 0; i < len(data); i++ {
		if data[i] == sep {
			slices = append(slices, data[start:i])
			start = i + 1
		}
	}

	slices = append(slices, data[start:])

	return slices
}

// JoinBytes joins multiple byte slices into a single string using a separator byte.
//
// If the input slice is empty, the function returns an empty string.
//
// Parameters:
//
// 	- slices: A slice of byte slices to join.
// 	- sep: The separator byte.
//
// Returns:
//
// 	- string: The joined string.
func JoinBytes(slices [][]byte, sep byte) string {
	// FIXME: Remove this once MyGoLib is updated.
	if len(slices) == 0 {
		return ""
	}

	var builder strings.Builder

	builder.Write(slices[0])

	for _, slice := range slices[1:] {
		builder.WriteByte(sep)
		builder.Write(slice)
	}

	return builder.String()
}
