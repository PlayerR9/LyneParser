package Lexer

import (
	"errors"
	"fmt"

	gr "github.com/PlayerR9/LyneParser/Grammar"

	nd "github.com/PlayerR9/MyGoLib/CustomData/Node"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

// Lexer is a lexer that uses a grammar to tokenize a string.
type Lexer struct {
	// grammar is the grammar used by the lexer.
	grammar *gr.Grammar

	// toSkip is a list of LHSs to skip.
	toSkip []string

	// root is the root node of the lexer.
	root *nd.Node[helperToken]

	// leaves is a list of all the leaves in the lexer.
	leaves []*nd.Node[helperToken]
}

// NewLexer creates a new lexer.
//
// Parameters:
//   - grammar: The grammar to use.
//
// Returns:
//
//   - Lexer: The new lexer.
//   - error: An error of type *ers.ErrInvalidParameter if the grammar is nil.
func NewLexer(grammar *gr.Grammar) (Lexer, error) {
	var lex Lexer

	if grammar == nil {
		return lex, ers.NewErrNilParameter("grammar")
	}

	lex.grammar = grammar
	lex.toSkip = grammar.LhsToSkip

	return lex, nil
}

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

// Lex is the main function of the lexer.
//
// Parameters:
//   - b: The byte slice to lex.
//
// Returns:
//   - error: An error if lexing fails.
func (l *Lexer) Lex(b []byte) error {
	if len(b) == 0 {
		return errors.New("no tokens to parse")
	}

	tok := gr.NewLeafToken("root", "", -1)
	root := nd.NewNode(newHelperToken(&tok))

	l.root = &root

	matches := l.grammar.Match(0, b)
	if len(matches) == 0 {
		return errors.New("no matches found at index 0")
	}

	l.addFirstLeaves(matches)

	l.root.Data.SetStatus(TkComplete)

	for {
		// Remove all the leaves that are completed.
		todo := slext.SliceFilter(l.leaves, func(leaf *nd.Node[helperToken]) bool {
			return leaf.Data.Status != TkComplete
		})
		if len(todo) == 0 {
			// All leaves are complete.
			break
		}

		// Remove all the leaves that are in error.
		todo = slext.SliceFilter(todo, func(leaf *nd.Node[helperToken]) bool {
			return leaf.Data.Status != TkError
		})
		if len(todo) == 0 {
			// All leaves are in error.
			break
		}

		// Remaining leaves are incomplete.
		var newLeaves []*nd.Node[helperToken]

		for _, leaf := range todo {
			l.processLeaf(leaf, b)
			newLeaves = append(newLeaves, leaf.GetLeaves()...)
		}

		l.leaves = newLeaves
	}

	return nil
}

// GetTokens returns the tokens that have been lexed.
//
// Remember to use Lexer.RemoveToSkipTokens() to remove tokens that
// are not needed for the parser (i.e., marked as to skip in the grammar).
//
// Returns:
//   - []gr.TokenStream: The tokens that have been lexed.
//   - error: An error if the lexer has not been run yet.
func (l *Lexer) GetTokens() ([]gr.TokenStream, error) {
	if l.root == nil {
		return nil, errors.New("must call Lexer.Lex() first")
	}

	tokenBranches := l.root.SnakeTraversal()

	branches, invalidTokIndex := filterInvalidBranches(tokenBranches)

	// Convert the tokens to gr.TokenStream.
	result := make([]gr.TokenStream, len(branches))

	for i, branch := range branches {
		if len(branch) == 0 {
			// Skip empty branches.
			continue
		}

		ts := make([]gr.LeafToken, len(branch)-1)

		// branch[1:] to skip the root token.
		for j, token := range branch[1:] {
			ts[j] = *token.Tok
		}

		result[i] = gr.NewTokenStream(ts)
	}

	if invalidTokIndex != -1 {
		return result, fmt.Errorf("invalid token at index %d", invalidTokIndex)
	}

	return result, nil
}

// RemoveToSkipTokens removes tokens that are marked as to skip in the grammar.
//
// Parameters:
//   - branches: The branches to remove tokens from.
//
// Returns:
//   - []gr.TokenStream: The branches with the tokens removed.
func (l *Lexer) RemoveToSkipTokens(branches []gr.TokenStream) []gr.TokenStream {
	for _, toSkip := range l.toSkip {
		if len(branches) == 0 {
			break
		}

		top := 0

		for i := 0; i < len(branches); i++ {
			branches[i].RemoveByTokenID(toSkip)

			if !branches[i].IsEmpty() {
				branches[top] = branches[i]
				top++
			}
		}

		branches = branches[:top]
	}

	return branches
}

// FullLexer is a convenience function that creates a new lexer, lexes the content,
// and returns the token streams.
//
// Parameters:
//   - grammar: The grammar to use.
//   - content: The content to lex.
//
// Returns:
//   - []gr.TokenStream: The tokens that have been lexed.
//   - error: An error if lexing fails.
func FullLexer(grammar *gr.Grammar, content string) ([]gr.TokenStream, error) {
	lexer, err := NewLexer(grammar)
	if err != nil {
		return nil, err
	}

	err = lexer.Lex([]byte(content))
	tokens, _ := lexer.GetTokens()

	return tokens, err
}
