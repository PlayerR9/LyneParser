package Lexer

import (
	"testing"

	gr "github.com/PlayerR9/LyneParser/Grammar"
)

type TestTokenType int

const (
	TkEof TestTokenType = iota
	TkWord
	TkComma
	TkExclamation
)

func (t TestTokenType) String() string {
	return [...]string{
		gr.EOFTokenID,
		"word",
		"comma",
		"exclamation",
	}[t]
}

func (t TestTokenType) IsTerminal() bool {
	return true
}

func TestPrintCode(t *testing.T) {
	var (
		DataToTest   []rune                     = []rune("Hello, word!")
		TokensToTest []*gr.Token[TestTokenType] = []*gr.Token[TestTokenType]{
			{
				ID:   TkWord,
				Data: "Hello",
				At:   0,
			},
			{
				ID:   TkComma,
				Data: ",",
				At:   6,
			},
			{
				ID:   TkWord,
				Data: "word",
				At:   8,
			}, // Invalid token (expected "world")
			{
				ID:   TkEof,
				Data: nil,
				At:   -1,
			},
		}

		ExpectedOutput string = "Hello, word!\n       ^^^^\n"
	)

	output := PrintCode(DataToTest, TokensToTest)
	if output != ExpectedOutput {
		t.Errorf("PrintCode() =\n%s, want\n%s", output, ExpectedOutput)
	}
}
