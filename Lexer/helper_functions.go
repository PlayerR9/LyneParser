package Lexer

import (
	com "github.com/PlayerR9/LyneParser/Common"
	gr "github.com/PlayerR9/LyneParser/Grammar"
	teval "github.com/PlayerR9/MyGoLib/CustomData/TreeExplorer"
)

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
func convertBranchToTokenStream(branch []*teval.CurrentEval[*gr.LeafToken]) *com.TokenStream {
	ts := make([]*gr.LeafToken, 0, len(branch))

	for _, leaf := range branch {
		ts = append(ts, leaf.GetElem())
	}

	ts = setEOFToken(ts)

	setLookahead(ts)

	return com.NewTokenStream(ts)
}
