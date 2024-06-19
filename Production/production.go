package Production

import (
	lxr "github.com/PlayerR9/LyneParser/Lexer"
	prs "github.com/PlayerR9/LyneParser/Parser"
)

var (
	LexerGrammar  *lxr.Grammar
	ParserGrammar *prs.Grammar
)

func init() {
	grammar1, err := lxr.NewGrammar(
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

	grammar2, err := prs.NewGrammar(
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
	Lexer = lxr.NewLexer(LexerGrammar)

	parser, err := prs.NewParser(ParserGrammar)
	if err != nil {
		panic(err)
	}

	Parser = parser
}

func Lex(source string) error {
	return nil
}
