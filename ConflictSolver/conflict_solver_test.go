package ConflictSolver

import (
	"strings"
	"testing"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
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

	// builder.AddProductions(gr.NewProduction("test", "e1 e2 X"))
	// builder.AddProductions(gr.NewProduction("test", "e2 e2 X"))

	grammar, err := builder.Build()
	if err != nil {
		panic(err)
	}

	ParserGrammar = grammar
}

func TestAmbiguousShifts(t *testing.T) {
	rules := ParserGrammar.GetProductions()

	cs, err := NewConflictSolver(ParserGrammar.Symbols, rules)
	if err != nil {
		t.Fatalf("NewConflictSolver() returned an error: %s", err.Error())
	}

	// DEBUG: Display the decision table before solving ambiguous shifts.
	doc, err := ffs.SprintFString(ffs.NewFormatter(ffs.NewIndentConfig("   ", 0)), cs)
	if err != nil {
		t.Fatalf("ffs.SprintFString() returned an error: %s", err.Error())
	}

	pages := strings.Join(ffs.Stringfy(doc), "\f")

	t.Log(pages)

	err = cs.SolveAmbiguousShifts()
	if err != nil {
		t.Fatalf("ConflictSolver.SolveAmbiguousShifts() returned an error: %s", err.Error())
	}

	// DEBUG: Display the decision table after solving ambiguous shifts.
	doc, err = ffs.SprintFString(ffs.NewFormatter(ffs.NewIndentConfig("   ", 0)), cs)
	if err != nil {
		t.Fatalf("ffs.SprintFString() returned an error: %s", err.Error())
	}

	pages = strings.Join(ffs.Stringfy(doc), "\f")

	t.Log(pages)
}

func TestConflictSolver(t *testing.T) {
	rules := ParserGrammar.GetProductions()

	cs, err := NewConflictSolver(ParserGrammar.Symbols, rules)
	if err != nil {
		t.Fatalf("NewConflictSolver() returned an error: %s", err.Error())
	}

	err = cs.SolveAmbiguousShifts()
	if err != nil {
		t.Fatalf("ConflictSolver.SolveAmbiguousShifts() returned an error: %s", err.Error())
	}

	// DEBUG: Display the decision table before solving conflicts.
	doc, err := ffs.SprintFString(ffs.NewFormatter(ffs.NewIndentConfig("   ", 0)), cs)
	if err != nil {
		t.Fatalf("ffs.SprintFString() returned an error: %s", err.Error())
	}

	pages := strings.Join(ffs.Stringfy(doc), "\f")

	t.Log(pages)

	err = cs.Solve()
	if err != nil {
		t.Fatalf("ConflictSolver.Solve() returned an error: %s", err.Error())
	}

	// DEBUG: Display the decision table after solving conflicts.
	doc, err = ffs.SprintFString(ffs.NewFormatter(ffs.NewIndentConfig("   ", 0)), cs)
	if err != nil {
		t.Fatalf("ffs.SprintFString() returned an error: %s", err.Error())
	}

	pages = strings.Join(ffs.Stringfy(doc), "\n")

	t.Log(pages)
}
