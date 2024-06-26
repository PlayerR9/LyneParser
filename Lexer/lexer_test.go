package Lexer

import (
	"fmt"
	"strings"
	"testing"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
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

	v := NewVerbose(true)
	defer v.Close()

	iter := lexer.Lex([]byte(Source), v)

	var branch *cds.Stream[*gr.LeafToken]
	var err error

	for i := 0; ; i++ {
		branch, err = iter.Consume()
		if err != nil {
			break
		}

		t.Logf("Branch %d", i)

		items := branch.GetItems()
		values := make([]string, 0, len(items))

		for _, token := range items {
			str := fmt.Sprintf("%+v", token)
			values = append(values, str)
		}

		joinedStr := strings.Join(values, "\n")

		t.Logf("Branch %d:\n%s", i, joinedStr)
	}

	ok := IsDone(err)
	if !ok {
		t.Errorf("Lex() returned an error: %s", err.Error())
	}

	t.Fatalf("Done")
}

func TestSyntaxError(t *testing.T) {
	const (
		Source string = "[char!(\"Mark\"){\n\tSpecies(\"Human\")\n\tPersonality(\"Kind\"+\")\n}]"
	)

	lexer := NewLexer(TestGrammar)

	v := NewVerbose(true)
	defer v.Close()

	iter := lexer.Lex([]byte(Source), v)

	branch, err := iter.Consume()
	if err != nil {
		ok := IsDone(err)
		if !ok {
			t.Fatalf("Lex() returned an error: %s", err.Error())
		}
	}

	// DEBUG: Print syntax error
	line := FormatSyntaxError(branch, []byte(Source))
	fmt.Println(line)

	t.Fatalf("Syntax error:")
}
