package Lexer

import (
	"fmt"
	"slices"
	"strings"
	"unicode"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	com "github.com/PlayerR9/LyneParser/SimpleLexer/Common"
	stm "github.com/PlayerR9/LyneParser/SimpleLexer/Stream"
)

var (
	singleTokens map[rune]com.TestTkType
)

func init() {
	singleTokens = map[rune]com.TestTkType{
		'|':  com.TkBinOp,
		'!':  com.TkBinOp,
		'&':  com.TkBinOp,
		'^':  com.TkBinOp,
		'(':  com.TkOpParen,
		')':  com.TkClParen,
		' ':  com.TkWs,
		'\t': com.TkWs,
		'0':  com.TkImmediate,
	}
}

type Lexer struct {
	is stm.Stream

	tokens []*gr.Token[com.TestTkType]
}

func (l *Lexer) peekNext() (*rune, error) {
	char, err := l.is.Peek()
	if err == nil {
		return &char, nil
	}

	ok := stm.IsStreamExhausted(err)
	if ok {
		return nil, nil
	}

	return nil, fmt.Errorf("failed to peek next character: %w", err)
}

func (l *Lexer) lexNRepeat(leading rune, label1, label2 com.TestTkType, tails ...rune) (*gr.Token[com.TestTkType], error) {
	next, err := l.peekNext()
	if err != nil {
		return nil, fmt.Errorf("after %c: %w", leading, err)
	}

	if next == nil {
		// leading -> label1
		pos := l.is.Pos()
		lt := gr.NewToken(label1, string(leading), pos, nil)

		return lt, nil
	}

	index := slices.Index(tails, *next)
	if index == -1 {
		// leading -> label1
		pos := l.is.Pos()
		lt := gr.NewToken(label1, string(leading), pos, nil)

		return lt, nil
	}

	tail := tails[index]

	// leading tail -> label2
	pos := l.is.Pos()
	lt := gr.NewToken(label2, string(leading)+string(tail), pos, nil)

	l.is.Accept()

	return lt, nil
}

func (l *Lexer) lex2Repeat(leading rune, label com.TestTkType, fn func(rune) error) (*gr.Token[com.TestTkType], error) {
	next, err := l.peekNext()
	if err != nil {
		return nil, fmt.Errorf("after %c: %w", leading, err)
	}

	if next == nil {
		return nil, fmt.Errorf("after %c: expected tail, got EOF", leading)
	}

	err = fn(*next)
	if err != nil {
		return nil, fmt.Errorf("after %c: %w", leading, err)
	}

	// leading tail -> label
	pos := l.is.Pos()
	lt := gr.NewToken(label, string(leading)+string(*next), pos, nil)

	l.is.Accept()

	return lt, nil
}

func (l *Lexer) lexRepeated(leading rune, label com.TestTkType, fn func(rune) (bool, error)) (*gr.Token[com.TestTkType], error) {
	var builder strings.Builder

	builder.WriteRune(leading)

	for {
		next, err := l.peekNext()
		if err != nil {
			return nil, fmt.Errorf("after %c: %w", leading, err)
		}

		if next == nil {
			break
		}

		ok, err := fn(*next)
		if err != nil {
			return nil, fmt.Errorf("after %c: %w", leading, err)
		}

		if !ok {
			break
		}

		builder.WriteRune(*next)
		l.is.Accept()
	}

	pos := l.is.Pos()
	lt := gr.NewToken(label, builder.String(), pos, nil)

	return lt, nil
}

func (l *Lexer) lexOne() (*gr.Token[com.TestTkType], error) {
	char, err := l.is.Next()
	if err != nil {
		return nil, err
	}

	id, ok := singleTokens[char]
	if ok {
		pos := l.is.Pos()
		lt := gr.NewToken(id, string(char), pos, nil)

		return lt, nil
	}

	var lt *gr.Token[com.TestTkType]

	switch char {
	case '+':
		lt, err = l.lexNRepeat(char, com.TkBinOp, com.TkUnaryOp, '+')
	case '-':
		lt, err = l.lexNRepeat(char, com.TkBinOp, com.TkUnaryOp, '-')
	case '>':
		lt, err = l.lexNRepeat(char, com.TkRightArrow, com.TkUnaryOp, '>')
	case '$':
		// register
		// "$" "a".."z"
		f := func(r rune) error {
			ok := unicode.IsLetter(r)
			if !ok {
				return fmt.Errorf("expected letter, got %c", r)
			}

			ok = unicode.IsLower(r)
			if !ok {
				return fmt.Errorf("expected lower case letter, got %c", r)
			}

			return nil
		}

		lt, err = l.lex2Repeat(char, com.TkRegister, f)
	case '\n':
		f := func(r rune) (bool, error) {
			return r == '\n', nil
		}

		lt, err = l.lexRepeated(char, com.TkNewline, f)
	default:
		ok := unicode.IsDigit(char)
		if !ok {
			pos := l.is.Pos()
			return nil, fmt.Errorf("invalid character %c at index %d", char, pos)
		}

		// digit | number
		// "1".."9"

		f := func(r rune) (bool, error) {
			ok := unicode.IsDigit(r)
			return ok, nil
		}

		lt, err = l.lexRepeated(char, com.TkImmediate, f)
	}
	if err != nil {
		return nil, err
	}

	return lt, nil
}

func Lex(data []byte) ([]*gr.Token[com.TestTkType], error) {
	is := stm.NewStream(data)

	l := &Lexer{
		is: is,
	}

	var lt *gr.Token[com.TestTkType]
	var err error

	for {
		lt, err = l.lexOne()
		if err != nil {
			break
		}

		if lt.ID != com.TkWs {
			l.tokens = append(l.tokens, lt)
		}
	}

	l.tokens = addEOF(l.tokens)

	ok := stm.IsStreamExhausted(err)
	if !ok {
		return l.tokens, fmt.Errorf("failed to lex one: %w", err)
	}

	return l.tokens, nil
}

func addEOF(tokens []*gr.Token[com.TestTkType]) []*gr.Token[com.TestTkType] {
	lt := gr.NewToken(com.TkEof, nil, 0, nil)

	tokens = append(tokens, lt)

	for i := 0; i < len(tokens)-1; i++ {
		token := tokens[i]
		token.SetLookahead(tokens[i+1])
	}

	return tokens
}

/*

```ebnf
Source = Statement { newline Statement } EOF .

Register = dollar letter .

Immediate
   = zero
   | digit { number }
   .

Statement
   = (
      UnaryInstruction
      | BinaryInstruction
      | LoadImmediate
   ) right_arrow Register
   .

Operand
   = Register
   | BinaryInstruction
   .

UnaryInstruction = Operand unary_operator .

BinaryInstruction = Operand Operand binary_operator .

LoadImmediate = op_paren Immediate cl_paren .
```
*/
