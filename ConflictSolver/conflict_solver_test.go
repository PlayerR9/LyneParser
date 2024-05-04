package ConflictSolver

import (
	"fmt"
	"testing"

	gr "github.com/PlayerR9/LyneParser/Grammar"
)

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

func TestConflictSolver(t *testing.T) {
	rules := ParserGrammar.GetProductions()

	cs, err := NewConflictSolver(ParserGrammar.Symbols, rules)
	if err != nil {
		t.Errorf("NewConflictSolver() returned an error: %s", err.Error())
	}

	lines := cs.FString(0)

	for _, line := range lines {
		fmt.Println(line)
	}

	err = cs.Solve()
	if err != nil {
		t.Errorf("ConflictSolver.Solve() returned an error: %s", err.Error())
	}

	lines = cs.FString(0)

	for _, line := range lines {
		fmt.Println(line)
	}

	t.Errorf("TestConflictSolver() is not implemented")
}
