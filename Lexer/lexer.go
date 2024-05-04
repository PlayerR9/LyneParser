package Lexer

import (
	"errors"

	gr "github.com/PlayerR9/LyneParser/Grammar"

	nd "github.com/PlayerR9/MyGoLib/CustomData/Node"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

// Lexer is a lexer that uses a grammar to tokenize a string.
type Lexer struct {
	// grammar is the grammar used by the lexer.
	productions []*gr.RegProduction

	// toSkip is a list of LHSs to skip.
	toSkip []string

	// root is the root node of the lexer.
	root *nd.Node[*helperToken]

	// leaves is a list of all the leaves in the lexer.
	leaves []*nd.Node[*helperToken]
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
func NewLexer(grammar *gr.Grammar) (*Lexer, error) {
	if grammar == nil {
		return nil, ers.NewErrNilParameter("grammar")
	}

	lex := &Lexer{
		productions: grammar.GetRegProductions(),
		toSkip:      grammar.LhsToSkip,
		root:        nil,
		leaves:      nil,
	}

	if len(lex.productions) == 0 {
		return lex, gr.NewErrNoProductionRulesFound()
	}

	return lex, nil
}

// addFirstLeaves is a helper function that adds the first leaves to the lexer.
//
// Parameters:
//   - matches: The matches to add to the lexer.
func (l *Lexer) addFirstLeaves(matches []gr.MatchedResult[*gr.LeafToken]) {
	// Get the longest match.
	matches = getLongestMatches(matches)
	for _, match := range matches {
		l.root.AddChild(newHelperToken(match.Matched))
		l.leaves = l.root.GetLeaves()
	}
}

// processLeaf is a helper function that processes a leaf
// by adding children to it.
//
// Parameters:
//   - leaf: The leaf to process.
//   - b: The byte slice to lex.
func (l *Lexer) processLeaf(leaf *nd.Node[*helperToken], source *SourceStream) {
	nextAt := leaf.Data.GetPos() + len(leaf.Data.GetData())
	if source.IsDone(nextAt) {
		leaf.Data.SetStatus(TkComplete)
		return
	}

	matches, err := source.MatchFrom(nextAt, l.productions)
	if err != nil {
		leaf.Data.SetStatus(TkError)
		return
	}

	// Get the longest match.
	matches = getLongestMatches(matches)
	for _, match := range matches {
		leaf.AddChild(newHelperToken(match.Matched))
	}

	leaf.Data.SetStatus(TkComplete)
}

// filterLeaves filters out leaves that are incomplete or in error.
//
// Returns:
//   - []*nd.Node[*helperToken]: The filtered leaves.
//   - error: An error of type *ErrAllMatchesFailed if all matches failed.
func (l *Lexer) filterLeaves() ([]*nd.Node[*helperToken], error) {
	todo := slext.SliceFilter(l.leaves, FilterIncompleteLeaves)
	if len(todo) == 0 {
		return nil, nil
	}

	todo = slext.SliceFilter(todo, FilterErrorLeaves)
	if len(todo) == 0 {
		return nil, NewErrAllMatchesFailed()
	}

	return todo, nil
}

// generateNewLeaves generates new leaves for the lexer.
//
// Parameters:
//   - todo: The leaves to generate new leaves from.
//
// Returns:
//   - newLeaves: The new leaves generated.
func (l *Lexer) generateNewLeaves(todo []*nd.Node[*helperToken], source *SourceStream) (newLeaves []*nd.Node[*helperToken]) {
	for _, leaf := range todo {
		l.processLeaf(leaf, source)
		newLeaves = append(newLeaves, leaf.GetLeaves()...)
	}

	return
}

// processLeaves processes the leaves in the lexer.
//
// Returns:
//   - bool: True if all leaves are complete, false otherwise.
//   - error: An error of type *ErrAllMatchesFailed if all matches failed.
func (l *Lexer) processLeaves(source *SourceStream) (bool, error) {
	todo, err := l.filterLeaves()
	if err != nil {
		return false, err
	} else if len(todo) == 0 {
		return true, nil
	}

	l.leaves = l.generateNewLeaves(todo, source)

	return false, nil
}

// Lex is the main function of the lexer.
//
// Returns:
//   - error: An error if lexing fails.
//
// Errors:
//   - *ErrNoTokensToLex: There are no tokens to lex.
//   - *ErrNoMatches: No matches are found in the source.
//   - *ErrAllMatchesFailed: All matches failed.
func (l *Lexer) Lex(source *SourceStream) error {
	if source == nil || source.IsEmpty() {
		return NewErrNoTokensToLex()
	}

	root := nd.NewNode(newHelperToken(gr.NewRootToken()))
	l.root = &root

	matches, err := source.MatchFrom(0, l.productions)
	if err != nil {
		return ers.NewErrAt(0, err)
	}

	l.addFirstLeaves(matches)

	l.root.Data.SetStatus(TkComplete)

	for {
		isDone, err := l.processLeaves(source)
		if err != nil {
			return err
		} else if isDone {
			return nil
		}
	}
}

// GetTokens returns the tokens that have been lexed.
//
// Remember to use Lexer.RemoveToSkipTokens() to remove tokens that
// are not needed for the parser (i.e., marked as to skip in the grammar).
//
// Returns:
//   - result: The tokens that have been lexed.
//   - reason: An error if the lexer has not been run yet.
func (l *Lexer) GetTokens() (result []*gr.TokenStream, reason error) {
	if l.root == nil {
		reason = errors.New("must call Lexer.Lex() first")
		return
	}

	tokenBranches := l.root.SnakeTraversal()

	branches, invalidTokIndex := filterInvalidBranches(tokenBranches)
	if invalidTokIndex != -1 {
		reason = ers.NewErrAt(invalidTokIndex, NewErrInvalidToken())
	}

	branches, err := l.removeToSkipTokens(branches)
	if err != nil {
		reason = err
		return
	}

	for _, branch := range branches {
		result = append(result, convertBranchToTokenStream(branch))
	}

	return
}

// removeToSkipTokens removes tokens that are marked as to skip in the grammar.
//
// Parameters:
//   - branches: The branches to remove tokens from.
//
// Returns:
//   - []gr.TokenStream: The branches with the tokens removed.
func (l *Lexer) removeToSkipTokens(branches [][]*helperToken) (newBranches [][]*helperToken, reason error) {
	for _, branch := range branches {
		if len(branch) != 0 {
			newBranches = append(newBranches, branch[1:])
		}
	}

	for _, toSkip := range l.toSkip {
		newBranches = slext.SliceFilter(newBranches, FilterEmptyBranch)
		if len(newBranches) == 0 {
			reason = NewErrAllMatchesFailed()

			return
		}

		filterTokenDifferentID := func(h *helperToken) bool {
			return h.GetID() != toSkip
		}

		for i := 0; i < len(newBranches); i++ {
			newBranches[i] = slext.SliceFilter(newBranches[i], filterTokenDifferentID)
		}
	}

	return
}
