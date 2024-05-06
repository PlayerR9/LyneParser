package Lexer

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"

	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"

	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
)

// getLongestMatches returns the longest matches,
//
// Parameters:
//   - matches: A slice of matches to search through.
//
// Returns:
//
//   - []MatchedResult: A slice of the longest matches.
func getLongestMatches(matches []gr.MatchedResult[*gr.LeafToken]) []gr.MatchedResult[*gr.LeafToken] {
	weights := slext.ApplyWeightFunc(matches, MatchWeightFunc)
	return slext.FilterByPositiveWeight(weights)
}

// filterInvalidBranches filters out invalid branches.
//
// Parameters:
//   - branches: The branches to filter.
//
// Returns:
//   - [][]helperToken: The filtered branches.
//   - int: The index of the last invalid token. -1 if no invalid token is found.
func filterInvalidBranches(branches [][]*helperToken) ([][]*helperToken, int) {
	branches, ok := slext.SFSeparateEarly(branches, FilterIncompleteTokens)
	if ok {
		return branches, -1
	}

	// Return the longest branch.
	weights := slext.ApplyWeightFunc(branches, HelperWeightFunc)
	branches = slext.FilterByPositiveWeight(weights)

	return [][]*helperToken{branches[0]}, len(branches[0])
}

// SetEOFToken sets the end-of-file token in the token stream.
//
// If the end-of-file token is already present, it will not be added again.
func setEOFToken(tokens []*gr.LeafToken) []*gr.LeafToken {
	if len(tokens) != 0 && tokens[len(tokens)-1].ID == gr.EOFTokenID {
		// EOF token is already present
		return tokens
	}

	return append(tokens, gr.NewEOFToken())
}

// SetLookahead sets the lookahead token for all the tokens in the stream.
func setLookahead(tokens []*gr.LeafToken) {
	for i, token := range tokens[:len(tokens)-1] {
		token.SetLookahead(tokens[i+1])
	}
}

// convertBranchToTokenStream converts a branch to a token stream.
//
// Parameters:
//   - branch: The branch to convert.
//
// Returns:
//   - *gr.TokenStream: The token stream.
func convertBranchToTokenStream(branch []*helperToken) *cds.Stream[*gr.LeafToken] {
	ts := make([]*gr.LeafToken, 0, len(branch))

	for _, token := range branch {
		ts = append(ts, token.Tok)
	}

	ts = setEOFToken(ts)

	setLookahead(ts)

	return cds.NewStream(ts)
}
