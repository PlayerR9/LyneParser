package Lexer0

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
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
	FilterEmptyBranch us.PredicateFilter[[]*CurrentEval]

	// RemoveToSkipTokens removes tokens that are marked as to skip in the grammar.
	//
	// Parameters:
	//   - branches: The branches to remove tokens from.
	//
	// Returns:
	//   - []gr.TokenStream: The branches with the tokens removed.
	RemoveToSkipTokens FilterBranchesFunc
)

func init() {
	MatchWeightFunc = func(elem *gr.MatchedResult[*gr.LeafToken]) (float64, bool) {
		return float64(len(elem.Matched.Data)), true
	}

	FilterEmptyBranch = func(branch []*CurrentEval) bool {
		return len(branch) != 0
	}

	RemoveToSkipTokens = func(branches [][]*CurrentEval) ([][]*CurrentEval, error) {
		var newBranches [][]*CurrentEval

		for _, branch := range branches {
			if len(branch) != 0 {
				newBranches = append(newBranches, branch[1:])
			}
		}

		return newBranches, nil
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
	s := cds.NewStream(source)

	lexer := &TreeEvaluator{
		source: s,
	}

	lexer.root = tr.NewTree(NewCurrentEval(gr.NewRootToken()))

	matches, err := MatchFrom(lexer.source, 0, Productions)
	if err != nil {
		branches, _ := lexer.GetBranches()

		return branches, ue.NewErrAt(0, "position", err)
	}

	lexer.addMatchLeaves(lexer.root, matches)

	lexer.root.Root().Data.SetStatus(EvalComplete)

	for {
		err := lexer.root.ProcessLeaves(lexer.processLeaves())
		if err != nil {
			branches, _ := lexer.GetBranches()

			return branches, err
		}

		for {
			target := lexer.root.SearchNodes(FilterErrorLeaves)
			if target == nil {
				break
			}

			err = lexer.root.DeleteBranchContaining(target)
			if err != nil {
				branches, _ := lexer.GetBranches()
				return branches, err
			}
		}

		if lexer.root.Size() == 0 {
			branches, _ := lexer.GetBranches()

			return branches, NewErrAllMatchesFailed()
		}

		if !lexer.canContinue() {
			break
		}
	}

	var reason error

	for {
		target := lexer.root.SearchNodes(FilterIncompleteLeaves)
		if target == nil {
			break
		}

		err = lexer.root.DeleteBranchContaining(target)
		if err != nil {
			reason = err
			break
		}
	}

	branches, err := lexer.GetBranches()
	if reason != nil {
		return branches, reason
	}

	return branches, err
}

// FilterBranchesFunc is a function that filters branches.
//
// Parameters:
//   - branches: The branches to filter.
//
// Returns:
//   - [][]*CurrentEval: The filtered branches.
//   - error: An error if the branches are invalid.
type FilterBranchesFunc func(branches [][]*CurrentEval) ([][]*CurrentEval, error)

// FilterErrorLeaves is a filter that filters out leaves that are in error.
//
// Parameters:
//   - leaf: The leaf to filter.
//
// Returns:
//   - bool: True if the leaf is in error, false otherwise.
func FilterErrorLeaves(h *CurrentEval) bool {
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
func filterInvalidBranches(branches [][]*CurrentEval) ([][]*CurrentEval, int) {
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

	return [][]*CurrentEval{elems}, len(elems)
}
