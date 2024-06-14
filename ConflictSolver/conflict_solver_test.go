package ConflictSolver

import (
	"strings"
	"testing"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
)

var (
	ParserGrammar *gr.ParserGrammar
)

func init() {
	var err error

	grammar, err := gr.NewParserGrammar(
		`source -> arrayObj EOF
		key -> WORD
		key -> key WORD
		arrayObj -> OP_SQUARE mapObj CL_SQUARE
		mapObj -> fieldCls OP_CURLY mapObj1 CL_CURLY
		mapObj1 -> fieldCls
		mapObj1 -> fieldCls mapObj1
		fieldCls -> key OP_PAREN fieldCls1 CL_PAREN
		fieldCls1 -> ATTR
		fieldCls1 -> ATTR SEP fieldCls1`,
	)
	if err != nil {
		panic(err)
	}

	ParserGrammar = grammar
}

func TestAmbiguousShifts(t *testing.T) {
	rules := ParserGrammar.GetProductions()

	cs := NewConflictSolver(ParserGrammar.GetSymbols(), rules)

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

	cs := NewConflictSolver(ParserGrammar.GetSymbols(), rules)

	err := cs.SolveAmbiguousShifts()
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

	t.Fatalf("TestConflictSolver() is not implemented")
}
