package Lexer

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"

	nd "github.com/PlayerR9/MyGoLib/CustomData/Node"
	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

// addFirstLeaves is a helper function that adds the first leaves to the lexer.
//
// Parameters:
//   - matches: The matches to add to the lexer.
func (l *Lexer) addFirstLeaves(matches []gr.MatchedResult) {
	// Get the longest match.
	matches = getLongestMatches(matches)
	for _, match := range matches {
		leafToken, ok := match.Matched.(*gr.LeafToken)
		if !ok {
			panic("this should not happen: match.Matched is not a *LeafToken")
		}

		l.root.AddChild(helperToken{
			Status: TkIncomplete,
			Tok:    leafToken,
		})
		l.leaves = l.root.GetLeaves()
	}
}

// processLeaf is a helper function that processes a leaf
// by adding children to it.
//
// Parameters:
//   - leaf: The leaf to process.
//   - b: The byte slice to lex.
func (l *Lexer) processLeaf(leaf *nd.Node[helperToken], b []byte) {
	nextAt := leaf.Data.Tok.GetPos() + len(leaf.Data.Tok.Data)
	if nextAt >= len(b) {
		leaf.Data.SetStatus(TkComplete)
		return
	}
	subset := b[nextAt:]

	matches := l.grammar.Match(nextAt, subset)

	if len(matches) == 0 {
		// Branch is done but no match found.
		leaf.Data.SetStatus(TkError)
		return
	}

	// Get the longest match.
	matches = getLongestMatches(matches)
	for _, match := range matches {
		leafToken, ok := match.Matched.(*gr.LeafToken)
		if !ok {
			leaf.Data.SetStatus(TkError)
			return
		}

		leaf.AddChild(helperToken{
			Status: TkIncomplete,
			Tok:    leafToken,
		})
	}

	leaf.Data.SetStatus(TkComplete)
}

// getLongestMatches returns the longest matches,
//
// Parameters:
//   - matches: A slice of matches to search through.
//
// Returns:
//
//   - []MatchedResult: A slice of the longest matches.
func getLongestMatches(matches []gr.MatchedResult) []gr.MatchedResult {
	weights := slext.ApplyWeightFunc(matches, func(match gr.MatchedResult) (float64, bool) {
		leaf, ok := match.Matched.(*gr.LeafToken)
		if !ok {
			return 0, false
		}

		return float64(len(leaf.Data)), true
	})
	if len(weights) == 0 {
		return []gr.MatchedResult{}
	}

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
func filterInvalidBranches(branches [][]helperToken) ([][]helperToken, int) {
	branches, ok := slext.SFSeparateEarly(branches, func(h []helperToken) bool {
		return len(h) != 0 && h[len(h)-1].Status == TkComplete
	})
	if ok {
		return branches, -1
	}

	// Return the longest branch.
	weights := slext.ApplyWeightFunc(branches, func(h []helperToken) (float64, bool) {
		return float64(len(h)), true
	})

	branches = slext.FilterByPositiveWeight(weights)

	return [][]helperToken{branches[0]}, len(branches[0])
}
