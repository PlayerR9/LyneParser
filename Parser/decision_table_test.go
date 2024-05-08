package Parser

import (
	"fmt"
	"testing"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	lx "github.com/PlayerR9/LyneParser/Lexer"
	hlp "github.com/PlayerR9/MyGoLib/CustomData/Helpers"

	com "github.com/PlayerR9/LyneParser/Common"
)

var LexerGrammar *gr.Grammar = func() *gr.Grammar {
	var builder gr.GrammarBuilder

	// Fragments
	builder.AddProductions(gr.NewRegProduction("WORD", `[a-zA-Z]+`))

	// Literals
	builder.AddProductions(gr.NewRegProduction("ATTR", `".*?"`))

	// Brackets
	builder.AddProductions(gr.NewRegProduction("OP_PAREN", `\(`))
	builder.AddProductions(gr.NewRegProduction("CL_PAREN", `\)`))
	builder.AddProductions(gr.NewRegProduction("OP_SQUARE", `\[`))
	builder.AddProductions(gr.NewRegProduction("CL_SQUARE", `\]`))
	builder.AddProductions(gr.NewRegProduction("OP_CURLY", `\{`))
	builder.AddProductions(gr.NewRegProduction("CL_CURLY", `\}`))

	// Operators
	builder.AddProductions(gr.NewRegProduction("SEP", `[+]`))

	// Whitespace
	builder.AddProductions(gr.NewRegProduction("WS", `[ \t\r\n]+`))

	builder.SetToSkip("WS")

	grammar, err := builder.Build()
	if err != nil {
		panic(err)
	}

	return grammar
}()

var ParserGrammar *gr.Grammar = func() *gr.Grammar {
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

	return grammar
}()

var LexedContents []*com.TokenStream = func() []*com.TokenStream {
	l, err := lx.NewLexer(LexerGrammar)
	if err != nil {
		panic(err)
	}

	const (
		Source string = "[char(\"Mark\"){\n\tSpecies(\"Human\")\n\tPersonality(\"Kind\"+\"Caring\")\n}]"
	)

	err = l.Lex(new(com.ByteStream).FromString(Source))
	if err != nil {
		panic(err)
	}

	tokenBranches, err := l.GetTokens()
	if err != nil {
		panic(err)
	} else if len(tokenBranches) == 0 {
		panic("no token branches found")
	}

	return tokenBranches
}()

var TestParser *Parser = func() *Parser {
	p, err := NewParser(ParserGrammar)
	if err != nil {
		panic(err)
	}

	return p
}()

func TestParsing(t *testing.T) {
	results := make([]hlp.HResult[*com.TokenStream], 0)

	forest := make([]*com.TokenTree, 0)

	for _, branch := range LexedContents {
		err := TestParser.Parse(branch)
		if err != nil {
			results = append(results, hlp.HResult[*com.TokenStream]{First: branch, Second: err})

			continue
		}

		tmp, err := TestParser.GetParseTree()
		if err != nil {
			results = append(results, hlp.HResult[*com.TokenStream]{First: branch, Second: err})

			continue
		}

		forest = append(forest, tmp...)
	}

	if len(forest) == 0 {
		for _, result := range results {
			if result.First != nil {
				fmt.Printf("Failed to parse: %s\n", result.Second.Error())
			}
		}

		t.Errorf("no parse trees were found")
	} else {
		fmt.Println(forest[0].DebugString())
	}

	t.Errorf("TestParsing() is not implemented")
}
