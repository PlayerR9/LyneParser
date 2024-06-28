package Lexer

import (
	"fmt"
	"testing"
)

func TestLexer(t *testing.T) {
	const (
		InputStr string = "$a $b |> $d"
	)

	tokens, err := Lex([]byte(InputStr))
	if err != nil {
		t.Fatalf("Expected no error, got %s", err.Error())
	}

	for _, token := range tokens {
		fmt.Printf("Token: %s\n", token.GoString())
	}

	t.Fatalf("Expected 3 tokens, got %d", len(tokens))
}
