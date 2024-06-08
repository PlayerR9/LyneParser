package Parser

import (
	"testing"

	gr "github.com/PlayerR9/LyneParser/Grammar"

	com "github.com/PlayerR9/LyneParser/Common"
)

var ParserGrammar *gr.Grammar

func init() {
	var builder gr.GrammarBuilder

	// EOF arrayObj -> source
	builder.AddProductions(gr.NewProduction("source", "arrayObj EOF"))

	// WORD -> key
	// WORD key -> key
	builder.AddProductions(gr.NewProduction("key", "WORD"))
	builder.AddProductions(gr.NewProduction("key", "key WORD"))

	// CL_SQUARE mapObj OP_SQUARE -> arrayObj
	builder.AddProductions(gr.NewProduction("arrayObj", "OP_SQUARE mapObj CL_SQUARE"))

	// CL_CURLY mapObj1 OP_CURLY fieldCls -> mapObj
	builder.AddProductions(gr.NewProduction("mapObj", "fieldCls OP_CURLY mapObj1 CL_CURLY"))

	// fieldCls -> mapObj1
	// mapObj1 fieldCls -> mapObj1
	builder.AddProductions(gr.NewProduction("mapObj1", "fieldCls"))
	builder.AddProductions(gr.NewProduction("mapObj1", "fieldCls mapObj1"))

	// CL_PAREN fieldCls1 OP_PAREN key -> fieldCls
	builder.AddProductions(gr.NewProduction("fieldCls", "key OP_PAREN fieldCls1 CL_PAREN"))

	// ATTR -> fieldCls1
	// fieldCls1 SEP ATTR -> fieldCls1
	builder.AddProductions(gr.NewProduction("fieldCls1", "ATTR"))
	builder.AddProductions(gr.NewProduction("fieldCls1", "ATTR SEP fieldCls1"))

	grammar, err := builder.Build()
	if err != nil {
		panic(err)
	}

	ParserGrammar = grammar
}

var LexedContent *com.TokenStream

func init() {
	LexedContent = com.NewTokenStream(
		[]*gr.LeafToken{
			{
				ID:   "OP_SQUARE",
				Data: "[",
				At:   0,
			},
			{
				ID:   "WORD",
				Data: "char",
				At:   1,
			},
			{
				ID:   "OP_PAREN",
				Data: "(",
				At:   5,
			},
			{
				ID:   "ATTR",
				Data: "\"Mark\"",
				At:   6,
			},
			{
				ID:   "CL_PAREN",
				Data: ")",
				At:   12,
			},
			{
				ID:   "OP_CURLY",
				Data: "{",
				At:   13,
			},
			{
				ID:   "WORD",
				Data: "Species",
				At:   16,
			},
			{
				ID:   "OP_PAREN",
				Data: "(",
				At:   23,
			},
			{
				ID:   "ATTR",
				Data: "\"Human\"",
				At:   24,
			},
			{
				ID:   "CL_PAREN",
				Data: ")",
				At:   31,
			},
			{
				ID:   "WORD",
				Data: "Personality",
				At:   34,
			},
			{
				ID:   "OP_PAREN",
				Data: "(",
				At:   45,
			},
			{
				ID:   "ATTR",
				Data: "\"Kind\"",
				At:   46,
			},
			{
				ID:   "SEP",
				Data: "+",
				At:   52,
			},
			{
				ID:   "ATTR",
				Data: "\"Caring\"",
				At:   53,
			},
			{
				ID:   "CL_PAREN",
				Data: ")",
				At:   61,
			},
			{
				ID:   "CL_CURLY",
				Data: "}",
				At:   63,
			},
			{
				ID:   "EOF",
				Data: "",
				At:   -1,
			},
		},
	)
}

func TestParsing(t *testing.T) {
	p, err := NewParser(ParserGrammar)
	if err != nil {
		t.Fatalf("NewParser() returned an error: %s", err.Error())
	}

	err = p.Parse(LexedContent)
	if err != nil {
		t.Fatalf("Parser.Parse() returned an error: %s", err.Error())
	}

	forest, err := p.GetParseTree()
	if err != nil {
		t.Fatalf("Parser.GetParseTree() returned an error: %s", err.Error())
	}

	if len(forest) == 0 {
		t.Fatalf("no parse trees were found")
	}

	t.Log(forest[0].DebugString())

	t.Error("TestParsing() is not implemented")
}
