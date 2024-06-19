package Lexer

import (
	"slices"

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
	filterEmptyBranch us.PredicateFilter[[]*uc.Pair[EvalStatus, *gr.LeafToken]]

	// filterInvalidBranches filters out invalid branches.
	//
	// Parameters:
	//   - branches: The branches to filter.
	//
	// Returns:
	//   - [][]helperToken: The filtered branches.
	//   - int: The index of the last invalid token. -1 if no invalid token is found.
	filterInvalidBranches func(branches [][]*uc.Pair[EvalStatus, *gr.LeafToken]) ([][]*uc.Pair[EvalStatus, *gr.LeafToken], int)

	// filterIncompleteLeaves is a filter that filters out incomplete leaves.
	//
	// Parameters:
	//   - leaf: The leaf to filter.
	//
	// Returns:
	//   - bool: True if the leaf is incomplete, false otherwise.
	filterIncompleteLeaves us.PredicateFilter[*uc.Pair[EvalStatus, *gr.LeafToken]]

	// filterCompleteTokens is a filter that filters complete helper tokens.
	//
	// Parameters:
	//   - h: The helper tokens to filter.
	//
	// Returns:
	//   - bool: True if the helper tokens are incomplete, false otherwise.
	filterCompleteTokens us.PredicateFilter[[]*uc.Pair[EvalStatus, *gr.LeafToken]]

	// helperWeightFunc is a weight function that returns the length of the helper tokens.
	//
	// Parameters:
	//   - h: The helper tokens to weigh.
	//
	// Returns:
	//   - float64: The weight of the helper tokens.
	//   - bool: True if the weight is valid, false otherwise.
	helperWeightFunc us.WeightFunc[[]*uc.Pair[EvalStatus, *gr.LeafToken]]

	// filterErrorLeaves is a filter that filters out leaves that are in error.
	//
	// Parameters:
	//   - leaf: The leaf to filter.
	//
	// Returns:
	//   - bool: True if the leaf is in error, false otherwise.
	filterErrorLeaves us.PredicateFilter[*uc.Pair[EvalStatus, *gr.LeafToken]]

	// selectBestMatches selects the best matches from the list of matches.
	// Usually, the best matches' euristic is the longest match.
	//
	// Parameters:
	//   - matches: The list of matches.
	//
	// Returns:
	//   - []*gr.MatchedResult[*gr.LeafToken]: The best matches.
	selectBestMatches func(matches []*gr.MatchedResult[*gr.LeafToken]) []*gr.MatchedResult[*gr.LeafToken]

	// removeToSkipTokens removes tokens that are marked as to skip in the grammar.
	//
	// Parameters:
	//   - toSkip: The tokens to skip.
	//   - branches: The branches to remove tokens from.
	//
	// Returns:
	//   - []gr.TokenStream: The branches with the tokens removed.
	removeToSkipTokens func(toSkip []string, branches [][]*uc.Pair[EvalStatus, *gr.LeafToken]) ([][]*uc.Pair[EvalStatus, *gr.LeafToken], error)
)

func init() {
	matchWeightFunc = func(elem *gr.MatchedResult[*gr.LeafToken]) (float64, bool) {
		return float64(len(elem.Matched.Data)), true
	}

	filterEmptyBranch = func(branch []*uc.Pair[EvalStatus, *gr.LeafToken]) bool {
		return len(branch) != 0
	}

	filterInvalidBranches = func(branches [][]*uc.Pair[EvalStatus, *gr.LeafToken]) ([][]*uc.Pair[EvalStatus, *gr.LeafToken], int) {
		branches, ok := us.SFSeparateEarly(branches, filterCompleteTokens)
		if ok {
			return branches, -1
		} else if len(branches) == 0 {
			return nil, -1
		}

		// Return the longest branch.
		weights := us.ApplyWeightFunc(branches, helperWeightFunc)
		weights = us.FilterByPositiveWeight(weights)

		elems := weights[0].GetData().First

		return [][]*uc.Pair[EvalStatus, *gr.LeafToken]{elems}, len(elems)
	}

	filterIncompleteLeaves = func(h *uc.Pair[EvalStatus, *gr.LeafToken]) bool {
		if h == nil {
			return true
		}

		return h.First == EvalIncomplete
	}

	filterCompleteTokens = func(h []*uc.Pair[EvalStatus, *gr.LeafToken]) bool {
		if len(h) == 0 {
			return false
		}

		status := h[len(h)-1].First
		return status == EvalComplete
	}

	helperWeightFunc = func(h []*uc.Pair[EvalStatus, *gr.LeafToken]) (float64, bool) {
		return float64(len(h)), true
	}

	filterErrorLeaves = func(h *uc.Pair[EvalStatus, *gr.LeafToken]) bool {
		if h == nil {
			return true
		}

		return h.First == EvalError
	}

	selectBestMatches = func(matches []*gr.MatchedResult[*gr.LeafToken]) []*gr.MatchedResult[*gr.LeafToken] {
		weights := us.ApplyWeightFunc(matches, matchWeightFunc)
		pairs := us.FilterByPositiveWeight(weights)

		return us.ExtractResults(pairs)
	}

	removeToSkipTokens = func(toSkip []string, branches [][]*uc.Pair[EvalStatus, *gr.LeafToken]) ([][]*uc.Pair[EvalStatus, *gr.LeafToken], error) {
		var newBranches [][]*uc.Pair[EvalStatus, *gr.LeafToken]
		var reason error

		for _, branch := range branches {
			if len(branch) != 0 {
				newBranches = append(newBranches, branch[1:])
			}
		}

		for _, elem := range toSkip {
			newBranches = us.SliceFilter(newBranches, filterEmptyBranch)
			if len(newBranches) == 0 {
				reason = NewErrAllMatchesFailed()

				return newBranches, reason
			}

			filterTokenDifferentID := func(h *uc.Pair[EvalStatus, *gr.LeafToken]) bool {
				id := h.Second.ID

				return id != elem
			}

			for i := 0; i < len(newBranches); i++ {
				newBranches[i] = us.SliceFilter(newBranches[i], filterTokenDifferentID)
			}
		}

		return newBranches, reason
	}
}

var (
	// SortFunc is a function that sorts the token stream.
	sortFunc func(a, b *cds.Stream[*gr.LeafToken]) int
)

func init() {
	sortFunc = func(a, b *cds.Stream[*gr.LeafToken]) int {
		return b.Size() - a.Size()
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
//
// Returns:
//   - *cds.Stream[*LeafToken]: The token stream.
func convertBranchToTokenStream(branch []*uc.Pair[EvalStatus, *gr.LeafToken]) *cds.Stream[*gr.LeafToken] {
	var ts []*gr.LeafToken

	// +1 for removing the ROOT token
	for i := 1; i < len(branch); i++ {
		ts = append(ts, branch[i].Second)
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

// getTokens returns the tokens that have been lexed.
//
// Remember to use Lexer.RemoveToSkipTokens() to remove tokens that
// are not needed for the parser (i.e., marked as to skip in the grammar).
//
// Returns:
//   - result: The tokens that have been lexed.
//   - reason: An error if the lexer has not been run yet.
func getTokens(tree *tr.Tree[EvalStatus, *gr.LeafToken], toSkip []string) ([]*cds.Stream[*gr.LeafToken], error) {
	tokenBranches := tree.SnakeTraversal()

	branches, invalidTokIndex := filterInvalidBranches(tokenBranches)
	if invalidTokIndex != -1 {
		var result []*cds.Stream[*gr.LeafToken]

		for _, branch := range branches {
			conv := convertBranchToTokenStream(branch)
			result = append(result, conv)
		}

		// Sort the result by size. (descending order)
		slices.SortStableFunc(result, sortFunc)

		return result, ue.NewErrAt(invalidTokIndex, "token", NewErrInvalidElement())
	}

	branches, err := removeToSkipTokens(toSkip, branches)

	var result []*cds.Stream[*gr.LeafToken]

	for _, branch := range branches {
		conv := convertBranchToTokenStream(branch)
		result = append(result, conv)
	}

	slices.SortStableFunc(result, sortFunc)

	return result, err
}
