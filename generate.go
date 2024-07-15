package LyneParser

import _ "github.com/PlayerR9/treenode"

// go:generate go run github.com/PlayerR9/treenode/cmd/treenode -name=HelperNode -fields=h/*Helper[T] -g=T/gr.TokenTyper -o=ConflictSolver/helper_treenode.go
//go:generate go run github.com/PlayerR9/treenode/cmd/treenode -name=TokenNode -fields=Token/*gr.Token[T],Status/EvalStatus -g=T/gr.TokenTyper -o=Lexer/token_treenode.go
