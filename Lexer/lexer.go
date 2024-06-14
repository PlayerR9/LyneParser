package Lexer

import (
	com "github.com/PlayerR9/LyneParser/Common"
	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	teval "github.com/PlayerR9/MyGoLib/TreeLike/Explorer"
	ers "github.com/PlayerR9/MyGoLib/Units/errors"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// Lexer is a lexer that uses a grammar to tokenize a string.
type Lexer struct {
	// te is the tree evaluator used by the lexer.
	te *teval.TreeEvaluator[*gr.MatchedResult[*gr.LeafToken], *LexerMatcher, *gr.LeafToken]

	// grammar is the grammar used by the lexer.
	productions []*gr.RegProduction

	// toSkip is a list of LHSs to skip.
	toSkip []string
}

// NewLexer creates a new lexer.
//
// Parameters:
//   - grammar: The grammar to use.
//
// Returns:
//   - Lexer: The new lexer.
//   - error: An error if the lexer cannot be created.
//
// Errors:
//   - *ers.ErrInvalidParameter: The grammar is nil.
//   - *gr.ErrNoProductionRulesFound: No production rules are found in the grammar.
//
// Example:
//
//	lexer, err := NewLexer(grammar)
//	if err != nil {
//	    // Handle error.
//	}
//
//	lexer.SetSource([]byte("1 + 2"))
//
//	err = lexer.Lex()
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
func NewLexer(grammar *gr.LexerGrammar) (*Lexer, error) {
	if grammar == nil {
		return nil, ers.NewErrNilParameter("grammar")
	}

	lex := &Lexer{
		productions: grammar.GetRegProductions(),
		toSkip:      grammar.GetToSkip(),
	}

	if len(lex.productions) == 0 {
		return lex, gr.NewErrNoProductionRulesFound()
	}

	lex.te = teval.NewTreeEvaluator[*gr.MatchedResult[*gr.LeafToken], *LexerMatcher](
		lex.removeToSkipTokens(),
	)

	return lex, nil
}

// Lex is the main function of the lexer.
//
// Parameters:
//   - source: The source to lex.
//
// Returns:
//   - error: An error if lexing fails.
//
// Errors:
//   - *ErrNoTokensToLex: There are no tokens to lex.
//   - *ErrNoMatches: No matches are found in the source.
//   - *ErrAllMatchesFailed: All matches failed.
func (l *Lexer) Lex(source []byte) error {
	matcher := &LexerMatcher{
		source: cds.NewStream(source),
	}

	matcher.productions = l.productions

	err := l.te.Evaluate(matcher, gr.NewRootToken())
	if err != nil {
		return err
	}

	return nil
}

// GetTokens returns the tokens that have been lexed.
//
// Remember to use Lexer.RemoveToSkipTokens() to remove tokens that
// are not needed for the parser (i.e., marked as to skip in the grammar).
//
// Returns:
//   - result: The tokens that have been lexed.
//   - reason: An error if the lexer has not been run yet.
func (l *Lexer) GetTokens() ([]*com.TokenStream, error) {
	branches, err := l.te.GetBranches()
	if err != nil {
		return nil, err
	}

	var result []*com.TokenStream

	for _, branch := range branches {
		result = append(result, convertBranchToTokenStream(branch))
	}

	return result, nil
}

// removeToSkipTokens removes tokens that are marked as to skip in the grammar.
//
// Parameters:
//   - branches: The branches to remove tokens from.
//
// Returns:
//   - []gr.TokenStream: The branches with the tokens removed.
func (l *Lexer) removeToSkipTokens() teval.FilterBranchesFunc[*gr.LeafToken] {
	return func(branches [][]*teval.CurrentEval[*gr.LeafToken]) ([][]*teval.CurrentEval[*gr.LeafToken], error) {
		var newBranches [][]*teval.CurrentEval[*gr.LeafToken]
		var reason error

		for _, branch := range branches {
			if len(branch) != 0 {
				newBranches = append(newBranches, branch[1:])
			}
		}

		for _, toSkip := range l.toSkip {
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
