package TreeExplorer

import (
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// FilterBranchesFunc is a function that filters branches.
//
// Parameters:
//   - branches: The branches to filter.
//
// Returns:
//   - [][]*CurrentEval: The filtered branches.
//   - error: An error if the branches are invalid.
type FilterBranchesFunc[O any] func(branches [][]*uc.Pair[EvalStatus, O]) ([][]*uc.Pair[EvalStatus, O], error)

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

// filterInvalidBranches filters out invalid branches.
//
// Parameters:
//   - branches: The branches to filter.
//
// Returns:
//   - [][]helperToken: The filtered branches.
//   - int: The index of the last invalid token. -1 if no invalid token is found.
func filterInvalidBranches[O any](branches [][]*uc.Pair[EvalStatus, O]) ([][]*uc.Pair[EvalStatus, O], int) {
	branches, ok := us.SFSeparateEarly(branches, FilterCompleteTokens)
	if ok {
		return branches, -1
	} else if len(branches) == 0 {
		return nil, -1
	}

	// Return the longest branch.
	weights := us.ApplyWeightFunc(branches, HelperWeightFunc)
	weights = us.FilterByPositiveWeight(weights)

	elems := weights[0].GetData().First

	return [][]*uc.Pair[EvalStatus, O]{elems}, len(elems)
}
