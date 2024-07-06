package Ast

import (
	"fmt"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// IsToken checks if a root is a token with a specific ID.
//
// This works regardless of whether the root is a leaf or non-leaf token.
//
// Parameters:
//   - root: The root to check.
//   - id: The ID to check.
//
// Returns:
//   - bool: True if the root is a token with the ID, false otherwise.
func IsToken[T uc.Enumer](root *gr.Token[T], id T) bool {
	rootID := root.GetID()
	if rootID != id {
		return false
	}

	isTerminal := gr.IsTerminal(id.String())
	return isTerminal
}

// Aster is an interface for AST nodes.
type Aster[T uc.Enumer] interface {
	// AstOf converts the children to an AST node.
	//
	// Parameters:
	//   - children: The children to convert.
	//
	// Returns:
	//   - error: The error if the conversion fails.
	AstOf(children []*gr.Token[T]) error
}

// AstOf constructs the AST of a source.
//
// Parameters:
//   - tree: The token tree.
//   - source: The source to construct the AST of.
//
// Returns:
//   - E: The source with the AST.
//   - error: The error if the construction fails.
func AstOf[E Aster[T], T uc.Enumer](tree *gr.TokenTree[T], source E) (E, error) {
	root := tree.GetRoot()

	// ok := IsToken(root, "source")
	// if !ok {
	// 	return *new(E), fmt.Errorf("the root is not a source token")
	// }

	children := root.Data.([]*gr.Token[T])

	err := source.AstOf(children)
	if err != nil {
		return *new(E), fmt.Errorf("failed to construct the AST: %w", err)
	}

	return source, nil
}
