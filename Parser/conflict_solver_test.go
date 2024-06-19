package Parser

import (
	"strings"
	"testing"

	cs "github.com/PlayerR9/LyneParser/ConflictSolver"
	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
)

var (
	TestGrammar *Grammar
)

func init() {
	var err error

	grammar, err := NewGrammar(
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

	TestGrammar = grammar
}

func TestAmbiguousShifts(t *testing.T) {
	rules := TestGrammar.GetProductions()

	tmp := cs.NewConflictSolver(TestGrammar.GetSymbols(), rules)

	// DEBUG: Display the decision table before solving ambiguous shifts.
	doc, err := ffs.SprintFString(ffs.NewFormatter(ffs.NewIndentConfig("   ", 0)), tmp)
	if err != nil {
		t.Fatalf("ffs.SprintFString() returned an error: %s", err.Error())
	}

	pages := strings.Join(ffs.Stringfy(doc), "\f")

	t.Log(pages)

	err = tmp.SolveAmbiguousShifts()
	if err != nil {
		t.Fatalf("ConflictSolver.SolveAmbiguousShifts() returned an error: %s", err.Error())
	}

	// DEBUG: Display the decision table after solving ambiguous shifts.
	doc, err = ffs.SprintFString(ffs.NewFormatter(ffs.NewIndentConfig("   ", 0)), tmp)
	if err != nil {
		t.Fatalf("ffs.SprintFString() returned an error: %s", err.Error())
	}

	pages = strings.Join(ffs.Stringfy(doc), "\f")

	t.Log(pages)
}

func TestConflictSolver(t *testing.T) {
	rules := TestGrammar.GetProductions()

	tmp := cs.NewConflictSolver(TestGrammar.GetSymbols(), rules)

	err := tmp.SolveAmbiguousShifts()
	if err != nil {
		t.Fatalf("ConflictSolver.SolveAmbiguousShifts() returned an error: %s", err.Error())
	}

	// DEBUG: Display the decision table before solving conflicts.
	doc, err := ffs.SprintFString(ffs.NewFormatter(ffs.NewIndentConfig("   ", 0)), tmp)
	if err != nil {
		t.Fatalf("ffs.SprintFString() returned an error: %s", err.Error())
	}

	pages := strings.Join(ffs.Stringfy(doc), "\f")

	t.Log(pages)

	err = tmp.Solve()
	if err != nil {
		t.Fatalf("ConflictSolver.Solve() returned an error: %s", err.Error())
	}

	// DEBUG: Display the decision table after solving conflicts.
	doc, err = ffs.SprintFString(ffs.NewFormatter(ffs.NewIndentConfig("   ", 0)), tmp)
	if err != nil {
		t.Fatalf("ffs.SprintFString() returned an error: %s", err.Error())
	}

	pages = strings.Join(ffs.Stringfy(doc), "\n")

	t.Log(pages)

	t.Fatalf("TestConflictSolver() is not implemented")
}

var LexedContent *cds.Stream[*gr.LeafToken]

func init() {
	tokens := []*gr.LeafToken{
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
			ID:   "CL_SQUARE",
			Data: "]",
			At:   64,
		},
		{
			ID:   "EOF",
			Data: "",
			At:   -1,
		},
	}

	for i := 0; i < len(tokens)-1; i++ {
		tokens[i].SetLookahead(tokens[i+1])
	}

	LexedContent = cds.NewStream(tokens)
}

func TestParsing(t *testing.T) {
	p, err := NewParser(TestGrammar)
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
}
