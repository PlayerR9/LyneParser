package Lexer

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// matchWeightFunc is a weight function that returns the length of the match.
//
// Parameters:
//   - match: The match to weigh.
//
// Returns:
//   - float64: The weight of the match.
//   - bool: True if the weight is valid, false otherwise.
func matchWeightFunc[T uc.Enumer](elem *gr.MatchedResult[T]) (float64, bool) {
	return float64(len(elem.Matched.Data.(string))), true
}

// filterErrorLeaves is a filter that filters out leaves that are in error.
//
// Parameters:
//   - leaf: The leaf to filter.
//
// Returns:
//   - bool: True if the leaf is in error, false otherwise.
func filterErrorLeaves[T uc.Enumer](h *tr.StatusInfo[EvalStatus, *gr.Token[T]]) bool {
	status := h.GetStatus()
	return status == EvalError
}

// selectBestMatches selects the best matches from the list of matches.
// Usually, the best matches' euristic is the longest match.
//
// Parameters:
//   - matches: The list of matches.
//   - logger: A verbose logger.
//
// Returns:
//   - []*gr.MatchedResult: The best matches.
func selectBestMatches[T uc.Enumer](matches []*gr.MatchedResult[T], logger *Verbose) []*gr.MatchedResult[T] {
	logger.DoIf(func(p *Printer) {
		p.Print("Selecting best matches...")
	})

	weights := us.ApplyWeightFunc(matches, matchWeightFunc)
	pairs := us.FilterByPositiveWeight(weights)

	results := us.ExtractResults(pairs)

	logger.DoIf(func(p *Printer) {
		p.Print("The best matches are:")

		for _, elem := range results {
			p.Printf("\t%+v", elem.Matched.Data)
		}
	})

	return results
}

// SetEOFToken sets the end-of-file token in the token stream.
//
// If the end-of-file token is already present, it will not be added again.
func setEOFToken[T uc.Enumer](tokens []*gr.Token[T]) []*gr.Token[T] {
	if len(tokens) != 0 && tokens[len(tokens)-1].ID.String() == gr.EOFTokenID {
		// EOF token is already present
		return tokens
	}

	// tok := gr.EofToken()
	//
	// return append(tokens, tok)
	panic("FIXME: EOF token is not implemented")
}

// SetLookahead sets the lookahead token for all the tokens in the stream.
func setLookahead[T uc.Enumer](tokens []*gr.Token[T]) {
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
func convertBranchToTokenStream[T uc.Enumer](branch []*tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.Token[T]]], toSkip []T) *cds.Stream[*gr.Token[T]] {
	branch = branch[1:]

	for _, elem := range toSkip {
		filterTokenDifferentID := func(h *tr.TreeNode[*tr.StatusInfo[EvalStatus, *gr.Token[T]]]) bool {
			data := h.Data.GetData()
			return data.ID != elem
		}

		branch = us.SliceFilter(branch, filterTokenDifferentID)
	}

	var ts []*gr.Token[T]

	for _, elem := range branch {
		data := elem.Data.GetData()
		ts = append(ts, data)
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
//   - *uc.ErrInvalidParameter: The from index is out of bounds.
//   - *ErrNoMatches: No matches are found.
func matchFrom[T uc.Enumer](s *cds.Stream[byte], from int, ps []*gr.RegProduction[T]) ([]*gr.MatchedResult[T], error) {
	size := s.Size()

	if from < 0 || from >= size {
		return nil, uc.NewErrInvalidParameter(
			"from",
			uc.NewErrOutOfBounds(from, 0, size),
		)
	}

	type Result struct {
		subset  []byte
		matches []*gr.MatchedResult[T]
	}

	var prevResult *Result

	for i := 1; i <= size; i++ {
		subset, _ := s.Get(from, i)

		var matches []*gr.MatchedResult[T]

		for i, p := range ps {
			matched, ok := p.MatchRegProd(from, subset)
			if ok {
				res := gr.NewMatchResult(matched, i)
				matches = append(matches, res)
			}
		}

		if len(matches) != 0 {
			result := &Result{
				subset:  subset,
				matches: matches,
			}

			prevResult = result
		}
	}

	if prevResult == nil {
		return nil, NewErrNoMatches()
	}

	return prevResult.matches, nil
}

// filterLeaves processes the leaves in the tree evaluator.
//
// Parameters:
//   - source: The source stream to match.
//   - productions: The production rules to match.
//   - logger: A verbose logger.
//
// Returns:
//   - bool: True if all leaves are complete, false otherwise.
//   - error: An error of type *ErrAllMatchesFailed if all matches failed.
func filterLeaves[T uc.Enumer](source *cds.Stream[byte], productions []*gr.RegProduction[T], logger *Verbose) uc.EvalManyFunc[*tr.StatusInfo[EvalStatus, *gr.Token[T]], *tr.StatusInfo[EvalStatus, *gr.Token[T]]] {
	filterFunc := func(ld *tr.StatusInfo[EvalStatus, *gr.Token[T]]) ([]*tr.StatusInfo[EvalStatus, *gr.Token[T]], error) {
		data := ld.GetData()

		var nextAt int

		if data.ID.String() == gr.RootTokenID {
			nextAt = 0
		} else {
			nextAt = data.GetPos() + len(data.Data.(string))
		}

		if nextAt >= source.Size() {
			ld.ChangeStatus(EvalComplete)
			return nil, nil
		}

		matches, err := matchFrom(source, nextAt, productions)
		if err != nil {
			ld.ChangeStatus(EvalError)
			return nil, nil
		}

		children := generateEvalTrees(matches, logger)
		ld.ChangeStatus(EvalComplete)

		return children, nil
	}

	return filterFunc
}

/*
// getTokens returns the tokens that have been lexed.
//
// Returns:
//   - result: The tokens that have been lexed.
//   - reason: An error if the lexer has not been run yet.
func getTokens(tree *tr.Tree[*tr.StatusInfo[EvalStatus, gr.Token]], toSkip []string) ([]*cds.Stream[gr.Token], error) {
	branches := tree.SnakeTraversal()

	branches = removeToSkipTokens(toSkip, branches)
	if len(branches) == 0 {
		return nil, errors.New("all tokens were skipped")
	}

	branches, ok := us.SFSeparateEarly(branches, filterCompleteTokens)

	// Sort the branches by length (descending order)
	slices.SortStableFunc(branches, sortFunc)

	if ok {
		result := make([]*cds.Stream[gr.Token], 0, len(branches))

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

	result := make([]*cds.Stream[gr.Token], 0, len(branches))

	for _, branch := range branches {
		conv := convertBranchToTokenStream(branch)
		result = append(result, conv)
	}

	return result, uc.NewErrPossibleError(
		NewErrAllMatchesFailed(),
		fmt.Errorf("after token %q, at index %d, there is no valid continuation",
			lastToken.Data,
			lastToken.At,
		),
	)
}
*/
