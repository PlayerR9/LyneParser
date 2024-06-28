package Lexer

import (
	"fmt"
	"slices"
	"strings"
	"unicode"

	stm "Ssalc/Parser/Stream"
)

type Lexer struct {
	is *stm.InputStream

	tokens []*stm.LeafToken
}

func (l *Lexer) peekNext() (*rune, error) {
	char, err := l.is.Peek()
	if err == nil {
		return &char, nil
	}

	ok := stm.IsExhausted(err)
	if ok {
		return nil, nil
	}

	return nil, fmt.Errorf("failed to peek next character: %w", err)
}

func (l *Lexer) lexNRepeat(leading rune, label1, label2 string, tails ...rune) (*stm.LeafToken, error) {
	next, err := l.peekNext()
	if err != nil {
		return nil, fmt.Errorf("after %c: %w", leading, err)
	}

	if next == nil {
		// leading -> label1
		pos := l.is.Pos()
		lt := stm.NewLeafToken(label1, string(leading), pos)

		return lt, nil
	}

	index := slices.Index(tails, *next)
	if index == -1 {
		// leading -> label1
		pos := l.is.Pos()
		lt := stm.NewLeafToken(label1, string(leading), pos)

		return lt, nil
	}

	tail := tails[index]

	// leading tail -> label2
	pos := l.is.Pos()
	lt := stm.NewLeafToken(label2, string(leading)+string(tail), pos)

	l.is.Accept()

	return lt, nil
}

func (l *Lexer) lex2Repeat(leading rune, label string, fn func(rune) error) (*stm.LeafToken, error) {
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
	lt := stm.NewLeafToken(label, string(leading)+string(*next), pos)

	l.is.Accept()

	return lt, nil
}

func (l *Lexer) lexRepeated(leading rune, label string, fn func(rune) (bool, error)) (*stm.LeafToken, error) {
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
	lt := stm.NewLeafToken(label, builder.String(), pos)

	return lt, nil
}

func (l *Lexer) lexOne() (*stm.LeafToken, error) {
	char, err := l.is.Next()
	if err != nil {
		ok := stm.IsExhausted(err)
		if ok {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to read next character: %w", err)
	}

	id, ok := singleTokens[char]
	if ok {
		pos := l.is.Pos()
		lt := stm.NewLeafToken(id, string(char), pos)

		return lt, nil
	}

	var lt *stm.LeafToken

	switch char {
	case '+':
		lt, err = l.lexNRepeat(char, "binary_operator", "unary_operator", '+')
	case '-':
		lt, err = l.lexNRepeat(char, "binary_operator", "unary_operator", '-')
	case '>':
		lt, err = l.lexNRepeat(char, "right_arrow", "unary_operator", '>')
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

		lt, err = l.lex2Repeat(char, "register", f)
	case '\n':
		f := func(r rune) (bool, error) {
			return r == '\n', nil
		}

		lt, err = l.lexRepeated(char, "newline", f)
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

		lt, err = l.lexRepeated(char, "immediate", f)
	}
	if err != nil {
		return nil, err
	}

	return lt, nil
}

func Lex(data []byte) ([]*stm.LeafToken, error) {
	is := stm.NewInputStream(data)

	l := &Lexer{
		is: is,
	}

	for {
		lt, err := l.lexOne()
		if err != nil {
			tokens := addEOF(l.tokens)

			return tokens, fmt.Errorf("failed to lex one: %w", err)
		}

		if lt == nil {
			break
		}

		if lt.ID != "ws" {
			l.tokens = append(l.tokens, lt)
		}
	}

	l.tokens = addEOF(l.tokens)

	return l.tokens, nil
}

func addEOF(tokens []*stm.LeafToken) []*stm.LeafToken {
	lt := stm.NewLeafToken("EOF", "", -1)

	tokens = append(tokens, lt)

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
