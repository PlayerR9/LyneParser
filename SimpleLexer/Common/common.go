package Common

type TestTkType int

const (
	TkEof TestTkType = iota
	TkBinOp
	TkUnaryOp
	TkOpParen
	TkClParen
	TkWs
	TkNewline
	TkImmediate
	TkRightArrow
	TkRegister
	TkUnaryInstruction
	TkLoadImmediate
	TkOperand
	TkSource
	TkSource1
	TkStatement
	TkBinInstruction
)

func (t TestTkType) String() string {
	return [...]string{
		"End of File",
		"binary operator",
		"unary operator",
		"open parenthesis",
		"close parenthesis",
		"whitespace",
		"newline",
		"immediate",
		"right arrow",
		"register",
		"unary instruction",
		"load immediate",
		"operand",
		"source",
		"source (I)",
		"statement",
		"binary instruction",
	}[t]
}
