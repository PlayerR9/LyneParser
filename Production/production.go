package Production

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	lxr "github.com/PlayerR9/LyneParser/Lexer"
	prs "github.com/PlayerR9/LyneParser/Parser"
)

var (
	LexerGrammar  *gr.LexerGrammar
	ParserGrammar *gr.ParserGrammar
)

func init() {
	grammar1, err := gr.NewLexerGrammar(
		`LHS -> [a-z][a-zA-Z0-9_]*
		WORD -> [a-zA-Z0-9_]+
		ARROW -> ->
		SYMBOL -> [\*\+\?]
		PIPE -> \|
		OP_PAREN -> \(
		CL_PAREN -> \)
		WS -> [ \t\n\r]+`,
		`WS`,
	)
	if err != nil {
		panic(err)
	}

	LexerGrammar = grammar1

	grammar2, err := gr.NewParserGrammar(
		`source -> LHS ARROW rhsCls EOF

		rhsCls -> rhs
		rhsCls -> rhs PIPE rhs
		rhsCls -> rhs PIPE rhs rhsCls

		rhs -> WORD
		rhs -> WORD SYMBOL
		rhs -> OP_PAREN rhsCls CL_PAREN
		rhs -> OP_PAREN rhsCls1 CL_PAREN SYMBOL
		
		rhs -> WORD rhs
		rhs -> WORD SYMBOL rhs
		rhs -> OP_PAREN rhsCls CL_PAREN rhs
		rhs -> OP_PAREN rhsCls1 CL_PAREN SYMBOL rhs`,
	)
	if err != nil {
		panic(err)
	}

	ParserGrammar = grammar2
}

var (
	Lexer  *lxr.Lexer
	Parser *prs.Parser
)

func init() {
	lexer, err := lxr.NewLexer(LexerGrammar)
	if err != nil {
		panic(err)
	}

	Lexer = lexer

	parser, err := prs.NewParser(ParserGrammar)
	if err != nil {
		panic(err)
	}

	Parser = parser
}

func Lex(source string) error {
	return nil
}
