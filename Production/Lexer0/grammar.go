package Lexer0

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
)

var (
	Productions []*gr.RegProduction
)

func init() {
	grammar, err := gr.NewLexerGrammar(
		`LHS -> [a-z][a-zA-Z0-9_]*
		WORD -> [a-zA-Z0-9_]+
		ARROW -> ->
		SYMBOL -> [\*\+\?]
		PIPE -> \|
		OP_PAREN -> \(
		CL_PAREN -> \)
		WS -> [ \t\n\r]+`,
		"",
	)
	if err != nil {
		panic(err)
	}

	Productions = grammar.GetRegProductions()
	if len(Productions) == 0 {
		panic("No productions found")
	}
}
