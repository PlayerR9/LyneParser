package Ast

import (
	"fmt"

	gr "github.com/PlayerR9/LyneParser/Grammar"
)

// Aster is an interface for AST nodes.
type Aster interface {
	// AstOf converts the children to an AST node.
	//
	// Parameters:
	//   - children: The children to convert.
	//
	// Returns:
	//   - error: The error if the conversion fails.
	AstOf(children []gr.Token) error
}

// AstOf constructs the AST of a source.
//
// Parameters:
//   - tree: The token tree.
//   - source: The source to construct the AST of.
//
// Returns:
//   - T: The source with the AST.
//   - error: The error if the construction fails.
func AstOf[T Aster](tree *gr.TokenTree, source T) (T, error) {
	root := tree.GetRoot()

	ok, err := IsToken(root, "source")
	if err != nil {
		return *new(T), fmt.Errorf("failed to check if the root is a source token: %w", err)
	}

	if !ok {
		return *new(T), fmt.Errorf("the root is not a source token")
	}

	children := root.Data.([]gr.Token)

	err = source.AstOf(children)
	if err != nil {
		return *new(T), fmt.Errorf("failed to construct the AST: %w", err)
	}

	return source, nil
}
