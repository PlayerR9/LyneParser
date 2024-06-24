package Lexer

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

var (
	TestGrammar *Grammar
)

func init() {
	grammar, err := NewGrammar(
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

	TestGrammar = grammar
}

func TestLex(t *testing.T) {
	const (
		Source string = "[char(\"Mark\"){\n\tSpecies(\"Human\")\n\tPersonality(\"Kind\"+\"Caring\")\n}]"
	)

	lexer := NewLexer(TestGrammar)

	for i := 0; ; i++ {
		alx := Lex(lexer, []byte(Source))

		branch, err := alx.GetBranch()
		if err != nil {
			t.Fatalf("Lex() returned an error: %s", err.Error())
		}

		if branch == nil {
			break
		}

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

func TestSyntaxError(t *testing.T) {
	const (
		Source string = "[char(\"Mark\"){\n\tSpecies(\"Human\")\n\tPersonality(\"Kind\"+\")\n}]"
	)

	lexer := NewLexer(TestGrammar)

	alx := Lex(lexer, []byte(Source))

	branch, err := alx.GetBranch()
	if err != nil {
		t.Fatalf("Lex() returned an error: %s", err.Error())
	}

	if branch == nil {
		t.Fatalf("Lex() returned no branches")
	}

	// DEBUG: Print syntax error
	line := FormatSyntaxError(branch, []byte(Source))
	fmt.Println(line)

	t.Fatalf("Test failed")
}
