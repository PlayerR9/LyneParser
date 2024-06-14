package LexerN

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	teval "github.com/PlayerR9/MyGoLib/TreeLike/Explorer"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

var (
	// MatchWeightFunc is a weight function that returns the length of the match.
	//
	// Parameters:
	//   - match: The match to weigh.
	//
	// Returns:
	//   - float64: The weight of the match.
	//   - bool: True if the weight is valid, false otherwise.
	MatchWeightFunc us.WeightFunc[*gr.MatchedResult[*gr.LeafToken]]

	// FilterEmptyBranch is a filter that filters out empty branches.
	//
	// Parameters:
	//   - branch: The branch to filter.
	//
	// Returns:
	//   - bool: True if the branch is not empty, false otherwise.
	FilterEmptyBranch us.PredicateFilter[[]*teval.CurrentEval[*gr.LeafToken]]

	// RemoveToSkipTokens removes tokens that are marked as to skip in the grammar.
	//
	// Parameters:
	//   - branches: The branches to remove tokens from.
	//
	// Returns:
	//   - []gr.TokenStream: The branches with the tokens removed.
	RemoveToSkipTokens teval.FilterBranchesFunc[*gr.LeafToken]
)

func init() {
	MatchWeightFunc = func(elem *gr.MatchedResult[*gr.LeafToken]) (float64, bool) {
		return float64(len(elem.Matched.Data)), true
	}

	FilterEmptyBranch = func(branch []*teval.CurrentEval[*gr.LeafToken]) bool {
		return len(branch) != 0
	}

	RemoveToSkipTokens = func(branches [][]*teval.CurrentEval[*gr.LeafToken]) ([][]*teval.CurrentEval[*gr.LeafToken], error) {
		var newBranches [][]*teval.CurrentEval[*gr.LeafToken]
		var reason error

		for _, branch := range branches {
			if len(branch) != 0 {
				newBranches = append(newBranches, branch[1:])
			}
		}

		for _, toSkip := range ToSkip {
			newBranches = us.SliceFilter(newBranches, FilterEmptyBranch)
			if len(newBranches) == 0 {
				reason = teval.NewErrAllMatchesFailed()

				return newBranches, reason
			}

			filterTokenDifferentID := func(h *teval.CurrentEval[*gr.LeafToken]) bool {
				return h.GetElem().ID != toSkip
			}

			for i := 0; i < len(newBranches); i++ {
				newBranches[i] = us.SliceFilter(newBranches[i], filterTokenDifferentID)
			}
		}

		return newBranches, reason
	}
}

// MatchFrom matches the source stream from a given index with a list of production rules.
//
// Parameters:
//   - s: The source stream to match.
//   - from: The index to start matching from.
//   - ps: The production rules to match.
//
// Returns:
//   - matches: A slice of MatchedResult that match the input token.
//   - reason: An error if no matches are found.
//
// Errors:
//   - *ue.ErrInvalidParameter: The from index is out of bounds.
//   - *ErrNoMatches: No matches are found.
func MatchFrom(s *cds.Stream[byte], from int, ps []*gr.RegProduction) (matches []*gr.MatchedResult[*gr.LeafToken], reason error) {
	size := s.Size()

	if from < 0 || from >= size {
		reason = ue.NewErrInvalidParameter(
			"from",
			ue.NewErrOutOfBounds(from, 0, size),
		)

		return
	}

	subSet, err := s.Get(from, size)
	if err != nil {
		panic(err)
	}

	for i, p := range ps {
		matched := p.Match(from, subSet)
		if matched != nil {
			matches = append(matches, gr.NewMatchResult(matched, i))
		}
	}

	if len(matches) == 0 {
		reason = NewErrNoMatches()
	}

	return
}

// Lex is the main function of the lexer.
//
// Parameters:
//   - source: The source to lex.
//
// Returns:
//   - []*cds.Stream[*LeafToken]: The tokens that have been lexed.
//   - error: An error if lexing fails.
//
// Errors:
//   - *ErrNoTokensToLex: There are no tokens to lex.
//   - *ErrNoMatches: No matches are found in the source.
//   - *ErrAllMatchesFailed: All matches failed.
//
// Example:
//
//	err := Lex([]byte("1 + 2")))
//	if err != nil {
//	    // Handle error.
//	}
//
//	tokenBranches, err := lexer.GetTokens()
//	if err != nil {
//	    // Handle error.
//	} else if len(tokenBranches) == 0 {
//	    // No tokens found.
//	}
//
//	tokenBranches = lexer.RemoveToSkipTokens(tokenBranches) // prepare for parsing
//
// // DEBUG: Print tokens.
//
//	for _, branch := range tokenBranches {
//	    for _, token := range branch {
//	        fmt.Println(token)
//	    }
//	}
//
// // Continue with parsing.
func Lex(source []byte) ([]*cds.Stream[*gr.LeafToken], error) {
	lexer := teval.NewTreeEvaluator[*gr.MatchedResult[*gr.LeafToken], *LexerMatcher](
		RemoveToSkipTokens,
	)

	matcher := &LexerMatcher{
		source: cds.NewStream(source),
	}

	err := lexer.Evaluate(matcher, gr.NewRootToken())

	branches, err2 := lexer.GetBranches()

	var result []*cds.Stream[*gr.LeafToken]
	for _, branch := range branches {
		conv := convertBranchToTokenStream(branch)

		result = append(result, conv)
	}

	if err != nil {
		return result, err
	}

	return result, err2
}

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
	branches, ok := us.SFSeparateEarly(branches, FilterIncompleteTokens)
	if ok {
		return branches, -1
	} else if len(branches) == 0 {
		return nil, -1
	}

	// Return the longest branch.
	weights := us.ApplyWeightFunc(branches, HelperWeightFunc)
	weights = us.FilterByPositiveWeight(weights)

	elems := weights[0].GetData().First

	return [][]*CurrentEval[O]{elems}, len(elems)
}
