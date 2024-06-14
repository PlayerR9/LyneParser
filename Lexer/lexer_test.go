package Lexer

import (
	"strconv"
	"strings"
	"testing"

	gr "github.com/PlayerR9/LyneParser/Grammar"
)

var (
	LexerGrammar *gr.LexerGrammar
)

func init() {
	grammar, err := gr.NewLexerGrammar(
		`WORD -> [a-zA-Z]+
		ATTR -> ".*?"
		OP_PAREN -> \(
		CL_PAREN -> \)
		OP_SQUARE -> \[
		CL_SQUARE -> \]
		OP_CURLY -> \{
		CL_CURLY -> \}
		SEP -> [\+]
		WS -> [ \t\r\n]+`,
		"WS",
	)
	if err != nil {
		panic(err)
	}

	LexerGrammar = grammar
}

func TestLex(t *testing.T) {
	const (
		Source string = "[char(\"Mark\"){\n\tSpecies(\"Human\")\n\tPersonality(\"Kind\"+\"Caring\")\n}]"
	)

	lexer, err := NewLexer(LexerGrammar)
	if err != nil {
		t.Fatalf("NewLexer() returned an error: %s", err.Error())
	}

	err = lexer.Lex([]byte(Source))
	if err != nil {
		t.Fatalf("Lexer.Lex() returned an error: %s", err.Error())
	}

	tokenBranches, err := lexer.GetTokens()
	if err != nil {
		t.Fatalf("Lexer.GetTokens() returned an error: %s", err.Error())
	}

	// DEBUG: Print token branches
	for i, branch := range tokenBranches {
		t.Logf("Branch %d", i)

		var values []string
		var builder strings.Builder

		for _, token := range branch.GetItems() {
			builder.WriteString("&Token{\n")
			builder.WriteString("\t\tID: ")
			builder.WriteString(strconv.Quote(token.ID))
			builder.WriteString(",\n")
			builder.WriteString("\t\tData: ")
			builder.WriteString(strconv.Quote(token.Data))
			builder.WriteString(",\n")
			builder.WriteString("\t\tAt: ")
			builder.WriteString(strconv.Itoa(token.At))
			builder.WriteString(",\n}")

			values = append(values, builder.String())
			builder.Reset()
		}

		t.Logf("[]*Token{\n\t%s\n}", strings.Join(values, ",\n\t"))
	}

	t.Fatalf("Test failed")
}
