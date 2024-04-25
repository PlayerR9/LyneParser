package Lexer

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"

	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

// getLongestMatches returns the longest matches,
//
// Parameters:
//   - matches: A slice of matches to search through.
//
// Returns:
//
//   - []MatchedResult: A slice of the longest matches.
func getLongestMatches(matches []gr.MatchedResult) []gr.MatchedResult {
	return slext.FilterByPositiveWeight(matches, func(match gr.MatchedResult) (int, bool) {
		leaf, ok := match.Matched.(*gr.LeafToken)
		if !ok {
			return 0, false
		}

		return len(leaf.Data), true
	})
}

// filterInvalidBranches filters out invalid branches.
//
// Parameters:
//   - branches: The branches to filter.
//
// Returns:
//   - [][]helperToken: The filtered branches.
//   - int: The length of the longest branch.
func filterInvalidBranches(branches [][]helperToken) ([][]helperToken, int) {
	branches, ok := slext.SFSeparateEarly(branches, func(h []helperToken) bool {
		return len(h) != 0 && h[len(h)-1].Status == TkComplete
	})
	if ok {
		return branches, -1
	}

	// Return the longest branch.
	branches = slext.FilterByPositiveWeight(branches, func(h []helperToken) (int, bool) {
		return len(h), true
	})

	return [][]helperToken{branches[0]}, len(branches[0])
}

// findInvalidTokenIndex finds the index of the first invalid token in the data.
// The function returns -1 if no invalid token is found.
//
// Parameters:
//   - branch: The branch of tokens to search for.
//   - data: The data to search in.
//
// Returns:
//   - int: The index of the first invalid token.
func findInvalidTokenIndex(branch []*gr.LeafToken, data []byte) int {
	pos := 0

	for _, token := range branch {
		b := []byte(token.Data)

		startIndex := slext.FindSubsliceFrom(data, b, pos)
		if startIndex == -1 {
			return -1
		}

		pos += startIndex + len(token.Data)
	}

	if pos >= len(data) {
		return -1
	}

	return pos
}
