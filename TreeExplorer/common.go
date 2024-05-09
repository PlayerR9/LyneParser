package TreeExplorer

import slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"

// FilterBranchesFunc is a function that filters branches.
//
// Parameters:
//   - branches: The branches to filter.
//
// Returns:
//   - [][]*CurrentEval: The filtered branches.
//   - error: An error if the branches are invalid.
type FilterBranchesFunc[O any] func(branches [][]*CurrentEval[O]) ([][]*CurrentEval[O], error)

// MatchResult is an interface that represents a match result.
type MatchResulter[O any] interface {
	// GetMatch returns the match.
	//
	// Returns:
	//   - O: The match.
	GetMatch() O
}

// Matcher is an interface that represents a matcher.
type Matcher[R MatchResulter[O], O any] interface {
	// IsDone is a function that checks if the matcher is done.
	//
	// Parameters:
	//   - from: The starting position of the match.
	//
	// Returns:
	//   - bool: True if the matcher is done, false otherwise.
	IsDone(from int) bool

	// Match is a function that matches the element.
	//
	// Parameters:
	//   - from: The starting position of the match.
	//
	// Returns:
	//   - []R: The list of matched results.
	//   - error: An error if the matchers cannot be created.
	Match(from int) ([]R, error)

	// SelectBestMatches selects the best matches from the list of matches.
	// Usually, the best matches' euristic is the longest match.
	//
	// Parameters:
	//   - matches: The list of matches.
	//
	// Returns:
	//   - []T: The best matches.
	SelectBestMatches(matches []R) []R

	// GetNext is a function that returns the next position of an element.
	//
	// Parameters:
	//   - elem: The element to get the next position of.
	//
	// Returns:
	//   - int: The next position of the element.
	GetNext(elem O) int
}

// FilterErrorLeaves is a filter that filters out leaves that are in error.
//
// Parameters:
//   - leaf: The leaf to filter.
//
// Returns:
//   - bool: True if the leaf is in error, false otherwise.
func FilterErrorLeaves[O any](h *CurrentEval[O]) bool {
	return h == nil || h.Status == EvalError
}

// filterInvalidBranches filters out invalid branches.
//
// Parameters:
//   - branches: The branches to filter.
//
// Returns:
//   - [][]helperToken: The filtered branches.
//   - int: The index of the last invalid token. -1 if no invalid token is found.
func filterInvalidBranches[O any](branches [][]*CurrentEval[O]) ([][]*CurrentEval[O], int) {
	branches, ok := slext.SFSeparateEarly(branches, FilterIncompleteTokens)
	if ok {
		return branches, -1
	}

	// Return the longest branch.
	weights := slext.ApplyWeightFunc(branches, HelperWeightFunc)
	branches = slext.FilterByPositiveWeight(weights)

	return [][]*CurrentEval[O]{branches[0]}, len(branches[0])
}
