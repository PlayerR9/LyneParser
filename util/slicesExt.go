package util

import (
	intf "github.com/PlayerR9/MyGoLibUnits/Interfaces"
)

// FindSubBytesFrom finds the first occurrence of a subslice in a byte
// slice starting from a given index.
//
// If the subslice is not found, the function returns -1.
//
// Parameters:
//
//   - s: The byte slice to search in.
//   - subS: The byte slice to search for.
//   - at: The index to start searching from.
//
// Returns:
//
//   - int: The index of the first occurrence of the subslice.
func FindSubsliceFrom[T intf.Comparable](s []T, subS []T, at int) int {
	// FIXME: Remove this once MyGoLib is updated.

	if len(subS) == 0 || len(s) == 0 || at+len(subS) > len(s) {
		return -1
	}

	if at < 0 {
		at = 0
	}

	possibleStarts := make([]int, 0)

	// Find all possible starting points.
	for i := at; i < len(s)-len(subS); i++ {
		if s[i] == subS[0] {
			possibleStarts = append(possibleStarts, i)
		}
	}

	// Check only the possible starting points that have enough space
	// to contain the subslice in full.
	top := 0

	for i := 0; i < len(possibleStarts)-1; i++ {
		if possibleStarts[i+1]-possibleStarts[i] >= len(subS) {
			possibleStarts[top] = possibleStarts[i]
			top++
		}
	}

	possibleStarts = possibleStarts[:top]

	// Check if the subslice is present at any of the possible starting points
	for _, start := range possibleStarts {
		found := true

		for j := 0; j < len(subS); j++ {
			if s[start+j] != subS[j] {
				found = false
				break
			}
		}

		if found {
			return start
		}
	}

	return -1
}
