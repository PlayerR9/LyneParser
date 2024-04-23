package Lexer

import (
	"errors"
	"fmt"
	"slices"

	gr "WppEditor/PlayerR9/LyneParser/Grammar"

	nd "github.com/PlayerR9/MyGoLib/CustomData/Node"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

// Lexer is a lexer that uses a grammar to tokenize a string
type Lexer struct {
	// grammar is the grammar used by the lexer
	grammar *gr.Grammar

	// toSkip is a list of LHSs to skip
	toSkip []string

	// root is the root node of the lexer
	root *nd.Node[helperToken]

	// leaves is a list of all the leaves in the lexer
	leaves []*nd.Node[helperToken]
}

// NewLexer creates a new lexer
//
// Parameters:
//
//   - grammar: The grammar to use
//
// Returns:
//
//   - *Lexer: The new lexer
//   - error: An error of type *ers.ErrInvalidParameter if the grammar is nil
func NewLexer(grammar *gr.Grammar) (*Lexer, error) {
	if grammar == nil {
		return nil, ers.NewErrNilParameter("grammar")
	}

	lex := &Lexer{
		grammar: grammar,
		toSkip:  grammar.LhsToSkip,
	}

	return lex, nil
}

// addFirstLeaves adds the first leaves to the lexer
//
// Parameters:
//
//   - b: The byte slice to lex
//
// Returns:
//
//   - error: An error if no matches are found at index 0
func (l *Lexer) addFirstLeaves(matches []gr.MatchedResult) {
	// Get the longest match
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

// processLeaf processes a leaf
//
// Parameters:
//
//   - leaf: The leaf to process
//   - b: The byte slice to lex
func (l *Lexer) processLeaf(leaf *nd.Node[helperToken], b []byte) {
	nextAt := leaf.Data.Tok.GetPos() + len(leaf.Data.Tok.Data)
	if nextAt >= len(b) {
		leaf.Data.SetStatus(TkComplete)
		return
	}
	subset := b[nextAt:]

	matches := l.grammar.Match(nextAt, subset)

	if len(matches) == 0 {
		// Branch is done but no match found
		leaf.Data.SetStatus(TkError)
		return
	}

	// Get the longest match
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

// Lex is the main function of the lexer
//
// Parameters:
//
//   - b: The byte slice to lex
//
// Returns:
//
//   - error: An error if lexing fails
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
		// Remove all the leaves that are completed
		todo := slext.SliceFilter(l.leaves, func(leaf *nd.Node[helperToken]) bool {
			return leaf.Data.Status != TkComplete
		})
		if len(todo) == 0 {
			// All leaves are complete
			break
		}

		// Remove all the leaves that are in error
		todo = slext.SliceFilter(todo, func(leaf *nd.Node[helperToken]) bool {
			return leaf.Data.Status != TkError
		})
		if len(todo) == 0 {
			// All leaves are in error
			break
		}

		// Remaining leaves are incomplete
		var newLeaves []*nd.Node[helperToken]

		for _, leaf := range todo {
			l.processLeaf(leaf, b)
			newLeaves = append(newLeaves, leaf.GetLeaves()...)
		}

		l.leaves = newLeaves
	}

	return nil
}

// removeSkippedTokens removes the skipped tokens from the token branches
//
// Parameters:
//
//   - tokenBranches: The token branches to filter
//
// Returns:
//
//   - [][]*LeafToken: The token branches with the skipped tokens removed
func (l *Lexer) removeSkippedTokens(tokenBranches [][]gr.LeafToken) [][]gr.LeafToken {
	// Remove the root token
	for i, branch := range tokenBranches {
		tokenBranches[i] = branch[1:]
	}

	for i, branch := range tokenBranches {
		tokenBranches[i] = slext.SliceFilter(branch, func(token gr.LeafToken) bool {
			return !slices.Contains(l.toSkip, token.ID)
		})
	}

	return tokenBranches
}

// GetTokens returns the tokens that have been lexed
//
// Returns:
//
//   - [][]*LeafToken: The tokens that have been lexed
//   - error: An error if the lexer has not been run yet
func (l *Lexer) GetTokens() ([][]gr.LeafToken, error) {
	if l.root == nil {
		return nil, errors.New("must call Lexer.Lex() first")
	}

	tokenBranches := l.root.SnakeTraversal()

	branches, invalidTokIndex := filterInvalidBranches(tokenBranches)

	// Convert the tokens to []*LeafToken
	result := make([][]gr.LeafToken, len(branches))

	for i, branch := range branches {
		result[i] = make([]gr.LeafToken, len(branch))

		for j, token := range branch {
			result[i][j] = *token.Tok
		}
	}

	result = l.removeSkippedTokens(result)
	result = slext.SliceFilter(result, emptyBranchFilter)

	if invalidTokIndex != -1 {
		return result, fmt.Errorf("invalid token at index %d", invalidTokIndex)
	}

	return result, nil
}

// FullLexer is a convenience function that creates a new lexer, lexes the content,
// and returns the tokens
//
// Parameters:
//
//   - grammar: The grammar to use
//   - content: The content to lex
//
// Returns:
//
//   - [][]*LeafToken: The tokens that have been lexed
//   - error: An error if lexing fails
func FullLexer(grammar *gr.Grammar, content string) ([][]gr.LeafToken, error) {
	lexer, err := NewLexer(grammar)
	if err != nil {
		return nil, err
	}

	err = lexer.Lex([]byte(content))
	tokens, _ := lexer.GetTokens()

	return tokens, err
}
