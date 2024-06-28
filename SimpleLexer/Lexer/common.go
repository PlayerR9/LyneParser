package Lexer

var (
	singleTokens map[rune]string
)

func init() {
	singleTokens = map[rune]string{
		'|':  "binary_operator",
		'!':  "binary_operator",
		'&':  "binary_operator",
		'^':  "binary_operator",
		'(':  "op_paren",
		')':  "cl_paren",
		' ':  "ws",
		'\t': "ws",
		'0':  "immediate",
	}
}
