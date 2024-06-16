package Lexer

import (
	"sync"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

// Lex is the main function of the lexer. This can be parallelized.
//
// Parameters:
//   - lexer: The lexer to use.
//   - source: The source to lex.
//
// Returns:
//   - error: An error if lexing fails.
//
// Errors:
//   - *ErrNoTokensToLex: There are no tokens to lex.
//   - *ErrNoMatches: No matches are found in the source.
//   - *ErrAllMatchesFailed: All matches failed.
//   - *gr.ErrNoProductionRulesFound: No production rules are found in the grammar.
func Lex(lexer *Lexer, input []byte) ([]*cds.Stream[*gr.LeafToken], error) {
	if lexer == nil {
		return nil, ue.NewErrNilParameter("lexer")
	}

	lexer.mu.RLock()

	if len(lexer.productions) == 0 {
		lexer.mu.RUnlock()
		return nil, gr.NewErrNoProductionRulesFound()
	}

	prodCopy := make([]*gr.RegProduction, len(lexer.productions))
	copy(prodCopy, lexer.productions)
	toSkip := make([]string, len(lexer.toSkip))
	copy(toSkip, lexer.toSkip)

	lexer.mu.RUnlock()

	stream := cds.NewStream(input)
	tree, err := executeLexing(stream, prodCopy)
	if err != nil {
		tokenBranches, _ := getTokens(tree, toSkip)
		return tokenBranches, err
	}

	tokenBranches, err := getTokens(tree, toSkip)
	return tokenBranches, err
}

// FullLexer is a convenience function that creates a new lexer, lexes the content,
// and returns the token streams.
//
// Parameters:
//   - grammar: The grammar to use.
//   - input: The input to lex.
//
// Returns:
//   - []*cds.Stream[*LeafToken]: The tokens that have been lexed.
//   - error: An error if lexing fails.
func FullLexer(grammar *Grammar, input []byte) ([]*cds.Stream[*gr.LeafToken], error) {
	if grammar == nil {
		return nil, ue.NewErrNilParameter("grammar")
	}

	productions := grammar.GetRegexProds()
	toSkip := grammar.GetToSkip()

	if len(productions) == 0 {
		return nil, gr.NewErrNoProductionRulesFound()
	}

	stream := cds.NewStream(input)
	tree, err := executeLexing(stream, productions)
	if err != nil {
		tokenBranches, _ := getTokens(tree, toSkip)
		return tokenBranches, err
	}

	tokenBranches, err := getTokens(tree, toSkip)
	return tokenBranches, err
}

// Lexer is a lexer that uses a grammar to tokenize a string.
type Lexer struct {
	// grammar is the grammar used by the lexer.
	productions []*gr.RegProduction

	// toSkip is a list of LHSs to skip.
	toSkip []string

	// mu is a mutex to protect the lexer.
	mu sync.RWMutex
}

// NewLexer creates a new lexer.
//
// Parameters:
//   - grammar: The grammar to use.
//
// Returns:
//   - Lexer: The new lexer.
//
// Example:
//
//	lexer, err := NewLexer(grammar)
//	if err != nil {
//	    // Handle error.
//	}
//
//	branches, err := lexer.Lex(lexer, []byte("1 + 2"))
//	if err != nil {
//	    // Handle error.
//	}
//
// // Continue with parsing.
func NewLexer(grammar *Grammar) *Lexer {
	if grammar == nil {
		return &Lexer{
			productions: nil,
			toSkip:      nil,
		}
	}

	lex := &Lexer{
		productions: grammar.GetRegexProds(),
		toSkip:      grammar.GetToSkip(),
	}

	return lex
}
