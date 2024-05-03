package Parser

import (
	"fmt"
	"testing"

	gr "github.com/PlayerR9/LyneParser/Grammar"
)

var ParserGrammar *gr.Grammar = func() *gr.Grammar {
	var builder gr.GrammarBuilder

	// EOF arrayObj -> source
	builder.AddProduction(gr.NewProduction("source", "arrayObj EOF"))

	// WORD -> key
	// WORD key -> key
	builder.AddProduction(gr.NewProduction("key", "WORD"))
	builder.AddProduction(gr.NewProduction("key", "key WORD"))

	// CL_SQUARE mapObj OP_SQUARE -> arrayObj
	builder.AddProduction(gr.NewProduction("arrayObj", "OP_SQUARE mapObj CL_SQUARE"))

	// CL_CURLY mapObj1 OP_CURLY fieldCls -> mapObj
	builder.AddProduction(gr.NewProduction("mapObj", "fieldCls OP_CURLY mapObj1 CL_CURLY"))

	// fieldCls -> mapObj1
	// mapObj1 fieldCls -> mapObj1
	builder.AddProduction(gr.NewProduction("mapObj1", "fieldCls"))
	builder.AddProduction(gr.NewProduction("mapObj1", "fieldCls mapObj1"))

	// CL_PAREN fieldCls1 OP_PAREN key -> fieldCls
	builder.AddProduction(gr.NewProduction("fieldCls", "key OP_PAREN fieldCls1 CL_PAREN"))

	// ATTR -> fieldCls1
	// fieldCls1 SEP ATTR -> fieldCls1
	builder.AddProduction(gr.NewProduction("fieldCls1", "ATTR"))
	builder.AddProduction(gr.NewProduction("fieldCls1", "ATTR SEP fieldCls1"))

	grammar, err := builder.Build()
	if err != nil {
		panic(err)
	}

	return grammar
}()

func TestDecisionTable(t *testing.T) {
	rules := make([]*gr.Production, 0)

	for _, p := range ParserGrammar.Productions {
		rule, ok := p.(*gr.Production)
		if !ok {
			t.Errorf("Production is not a *gr.Production")
		}

		rules = append(rules, rule)
	}

	dt := NewDecisionTable()

	err := dt.GenerateItems(rules)
	if err != nil {
		t.Errorf("GenerateItems() returned an error: %s", err.Error())
	}

	lines := dt.FString(0)

	for _, line := range lines {
		fmt.Println(line)
	}

	err = dt.FixConflicts()
	if err != nil {
		t.Errorf("Conflict: %s", err.Error())
	}

	lines = dt.FString(0)

	for _, line := range lines {
		fmt.Println(line)
	}

	t.Errorf("TestDecisionTable() is not implemented")
}
