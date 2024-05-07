package Lexer

import (
	com "github.com/PlayerR9/LyneParser/Common"
	gr "github.com/PlayerR9/LyneParser/Grammar"
	tr "github.com/PlayerR9/LyneParser/PlayerR9/Tree"

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
func convertBranchToTokenStream(branch []*helperToken) *com.TokenStream {
	ts := make([]*gr.LeafToken, 0, len(branch))

	for _, token := range branch {
		ts = append(ts, token.Tok)
	}

	ts = setEOFToken(ts)

	setLookahead(ts)

	return com.NewTokenStream(ts)
}

// addMatchLeaves adds the matches to a root tree as leaves.
//
// Parameters:
//   - root: The root of the tree to add the leaves to.
//   - matches: The matches to add to the lexer.
func addMatchLeaves(root *tr.Tree[*helperToken], matches []gr.MatchedResult[*gr.LeafToken]) {
	// Get the longest match.
	matches = getLongestMatches(matches)

	children := make([]*tr.Tree[*helperToken], 0, len(matches))

	for _, match := range matches {
		ht := newHelperToken(match.Matched)
		children = append(children, tr.NewTree(ht))
	}

	root.SetChildren(children)
}
