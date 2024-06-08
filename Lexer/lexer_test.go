package Lexer

import (
	"strconv"
	"strings"
	"testing"

	gr "github.com/PlayerR9/LyneParser/Grammar"

	com "github.com/PlayerR9/LyneParser/Common"
)

var LexerGrammar *gr.Grammar = func() *gr.Grammar {
	var builder gr.GrammarBuilder

	// Fragments
	builder.AddRegProductions(gr.NewRegProduction("WORD", `[a-zA-Z]+`))

	// Literals
	builder.AddRegProductions(gr.NewRegProduction("ATTR", `".*?"`))

	// Brackets
	builder.AddRegProductions(gr.NewRegProduction("OP_PAREN", `\(`))
	builder.AddRegProductions(gr.NewRegProduction("CL_PAREN", `\)`))
	builder.AddRegProductions(gr.NewRegProduction("OP_SQUARE", `\[`))
	builder.AddRegProductions(gr.NewRegProduction("CL_SQUARE", `\]`))
	builder.AddRegProductions(gr.NewRegProduction("OP_CURLY", `\{`))
	builder.AddRegProductions(gr.NewRegProduction("CL_CURLY", `\}`))

	// Operators
	builder.AddRegProductions(gr.NewRegProduction("SEP", `[+]`))

	// Whitespace
	builder.AddRegProductions(gr.NewRegProduction("WS", `[ \t\r\n]+`))

	builder.SetToSkip("WS")

	grammar, err := builder.Build()
	if err != nil {
		panic(err)
	}

	return grammar
}()

func TestLex(t *testing.T) {
	const (
		Source string = "[char(\"Mark\"){\n\tSpecies(\"Human\")\n\tPersonality(\"Kind\"+\"Caring\")\n}]"
	)

	lexer, err := NewLexer(LexerGrammar)
	if err != nil {
		t.Fatalf("NewLexer() returned an error: %s", err.Error())
	}

	err = lexer.Lex(new(com.ByteStream).FromString(Source))
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
			builder.WriteString("&gr.Token{\n")
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

		t.Logf("[]*gr.Token{\n\t%s\n}", strings.Join(values, ",\n\t"))
	}

	t.Fatalf("Test failed")
}
