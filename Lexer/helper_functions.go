package Lexer

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/StatusTree"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

var (
	// matchWeightFunc is a weight function that returns the length of the match.
	//
	// Parameters:
	//   - match: The match to weigh.
	//
	// Returns:
	//   - float64: The weight of the match.
	//   - bool: True if the weight is valid, false otherwise.
	matchWeightFunc us.WeightFunc[*gr.MatchedResult[*gr.LeafToken]]

	// filterEmptyBranch is a filter that filters out empty branches.
	//
	// Parameters:
	//   - branch: The branch to filter.
	//
	// Returns:
	//   - bool: True if the branch is not empty, false otherwise.
	filterEmptyBranch us.PredicateFilter[[]uc.Pair[EvalStatus, *gr.LeafToken]]

	// filterIncompleteLeaves is a filter that filters out incomplete leaves.
	//
	// Parameters:
	//   - leaf: The leaf to filter.
	//
	// Returns:
	//   - bool: True if the leaf is incomplete, false otherwise.
	filterIncompleteLeaves us.PredicateFilter[uc.Pair[EvalStatus, *gr.LeafToken]]

	// filterCompleteTokens is a filter that filters complete helper tokens.
	//
	// Parameters:
	//   - h: The helper tokens to filter.
	//
	// Returns:
	//   - bool: True if the helper tokens are incomplete, false otherwise.
	filterCompleteTokens us.PredicateFilter[[]uc.Pair[EvalStatus, *gr.LeafToken]]

	// filterErrorLeaves is a filter that filters out leaves that are in error.
	//
	// Parameters:
	//   - leaf: The leaf to filter.
	//
	// Returns:
	//   - bool: True if the leaf is in error, false otherwise.
	filterErrorLeaves us.PredicateFilter[uc.Pair[EvalStatus, *gr.LeafToken]]

	// selectBestMatches selects the best matches from the list of matches.
	// Usually, the best matches' euristic is the longest match.
	//
	// Parameters:
	//   - matches: The list of matches.
	//
	// Returns:
	//   - []*gr.MatchedResult[*gr.LeafToken]: The best matches.
	selectBestMatches func(matches []*gr.MatchedResult[*gr.LeafToken]) []*gr.MatchedResult[*gr.LeafToken]
)

func init() {
	matchWeightFunc = func(elem *gr.MatchedResult[*gr.LeafToken]) (float64, bool) {
		return float64(len(elem.Matched.Data)), true
	}

	filterEmptyBranch = func(branch []uc.Pair[EvalStatus, *gr.LeafToken]) bool {
		return len(branch) != 0
	}

	filterCompleteTokens = func(h []uc.Pair[EvalStatus, *gr.LeafToken]) bool {
		status := h[len(h)-1].First
		return status == EvalComplete
	}

	filterIncompleteLeaves = func(h uc.Pair[EvalStatus, *gr.LeafToken]) bool {
		return h.First == EvalIncomplete
	}

	filterErrorLeaves = func(h uc.Pair[EvalStatus, *gr.LeafToken]) bool {
		return h.First == EvalError
	}

	selectBestMatches = func(matches []*gr.MatchedResult[*gr.LeafToken]) []*gr.MatchedResult[*gr.LeafToken] {
		weights := us.ApplyWeightFunc(matches, matchWeightFunc)
		pairs := us.FilterByPositiveWeight(weights)

		return us.ExtractResults(pairs)
	}

}

var (
	// SortFunc is a function that sorts the token stream.
	sortFunc func(a, b []uc.Pair[EvalStatus, *gr.LeafToken]) int
)

func init() {
	sortFunc = func(a, b []uc.Pair[EvalStatus, *gr.LeafToken]) int {
		return len(b) - len(a)
	}
}

// SetEOFToken sets the end-of-file token in the token stream.
//
// If the end-of-file token is already present, it will not be added again.
func setEOFToken(tokens []*gr.LeafToken) []*gr.LeafToken {
	if len(tokens) != 0 && tokens[len(tokens)-1].ID == gr.EOFTokenID {
		// EOF token is already present
		return tokens
	}

	tok := gr.NewEOFToken()

	return append(tokens, tok)
}

// SetLookahead sets the lookahead token for all the tokens in the stream.
func setLookahead(tokens []*gr.LeafToken) {
	for i := 0; i < len(tokens)-1; i++ {
		tokens[i].SetLookahead(tokens[i+1])
	}
}

// convertBranchToTokenStream converts a branch to a token stream.
//
// Parameters:
//   - branch: The branch to convert.
//   - toSkip: The tokens to skip.
//
// Returns:
//   - *cds.Stream[*LeafToken]: The token stream.
func convertBranchToTokenStream(branch []*tr.TreeNode[EvalStatus, *gr.LeafToken], toSkip []string) *cds.Stream[*gr.LeafToken] {
	branch = branch[1:]

	for _, elem := range toSkip {
		filterTokenDifferentID := func(h *tr.TreeNode[EvalStatus, *gr.LeafToken]) bool {
			id := h.Data.ID

			return id != elem
		}

		branch = us.SliceFilter(branch, filterTokenDifferentID)
	}

	var ts []*gr.LeafToken

	for _, elem := range branch {
		ts = append(ts, elem.Data)
	}

	ts = setEOFToken(ts)

	setLookahead(ts)

	stream := cds.NewStream(ts)

	return stream
}

// matchFrom matches the source stream from a given index with a list of production rules.
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
func matchFrom(s *cds.Stream[byte], from int, ps []*gr.RegProduction) (matches []*gr.MatchedResult[*gr.LeafToken], reason error) {
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
			res := gr.NewMatchResult(matched, i)
			matches = append(matches, res)
		}
	}

	if len(matches) == 0 {
		reason = NewErrNoMatches()
	}

	return
}

/*
// getTokens returns the tokens that have been lexed.
//
// Returns:
//   - result: The tokens that have been lexed.
//   - reason: An error if the lexer has not been run yet.
func getTokens(tree *tr.Tree[EvalStatus, *gr.LeafToken], toSkip []string) ([]*cds.Stream[*gr.LeafToken], error) {
	branches := tree.SnakeTraversal()

	branches = removeToSkipTokens(toSkip, branches)
	if len(branches) == 0 {
		return nil, errors.New("all tokens were skipped")
	}

	branches, ok := us.SFSeparateEarly(branches, filterCompleteTokens)

	// Sort the branches by length (descending order)
	slices.SortStableFunc(branches, sortFunc)

	if ok {
		result := make([]*cds.Stream[*gr.LeafToken], 0, len(branches))

		for _, branch := range branches {
			conv := convertBranchToTokenStream(branch)
			result = append(result, conv)
		}

		return result, nil
	}

	// Assume that the longest branch is the one with the
	// most likely error
	firstBranch := branches[0]
	size := len(firstBranch)

	lastToken := firstBranch[size-1].Second

	result := make([]*cds.Stream[*gr.LeafToken], 0, len(branches))

	for _, branch := range branches {
		conv := convertBranchToTokenStream(branch)
		result = append(result, conv)
	}

	return result, ue.NewErrPossibleError(
		NewErrAllMatchesFailed(),
		fmt.Errorf("after token %q, at index %d, there is no valid continuation",
			lastToken.Data,
			lastToken.At,
		),
	)
}
*/
