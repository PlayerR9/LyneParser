package Common

type TestTkType int

const (
	// Terminal tokens

	// TkEof is the end of file token.
	TkEof TestTkType = iota

	// TkBinOp is a binary operator token.
	TkBinOp

	// TkUnaryOp is a unary operator token.
	TkUnaryOp

	// TkOpParen is an open parenthesis token.
	TkOpParen

	// TkClParen is a close parenthesis token.
	TkClParen

	// TkWs is a whitespace token.
	TkWs

	// TkNewline is a newline token.
	TkNewline

	// TkImmediate is an immediate token.
	TkImmediate

	// TkRightArrow is a right arrow token.
	TkRightArrow

	// TkRegister is a register token.
	TkRegister

	// TkUnaryInstruction is a unary instruction token.
	TkUnaryInstruction

	// TkLoadImmediate is a load immediate token.
	TkLoadImmediate

	// TkOperand is an operand token.
	TkOperand

	// TkSource is a source token.
	TkSource

	// TkSource1 is a source (I) token.
	TkSource1

	// TkStatement is a statement token.
	TkStatement

	// TkBinInstruction is a binary instruction token.
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

func (t TestTkType) IsTerminal() bool {
	switch t {
	case TkEof, TkBinOp, TkUnaryOp, TkOpParen, TkClParen, TkWs, TkNewline,
		TkImmediate, TkRightArrow, TkRegister:
		return true
	case TkUnaryInstruction,
		TkLoadImmediate, TkOperand, TkSource, TkSource1, TkStatement,
		TkBinInstruction:
		return false
	}

	panic("unreachable")
}
