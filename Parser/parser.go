package Parser

import (
	"errors"

	cs "github.com/PlayerR9/LyneParser/ConflictSolver"
	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// Parser is a parser that uses a stack to parse a stream of tokens.
type Parser[T uc.Enumer] struct {
	// evals is a list of evaluations that the parser will use.
	evals []*CurrentEval[T]

	// decisionFunc represents the function that the parser will use to determine
	// the next action to take.
	dt *cs.ConflictSolver[T]
}

/////////////////////////////////////////////////////////////

// NewParser creates a new parser with the given grammar.
//
// Parameters:
//   - grammar: The grammar that the parser will use.
//
// Returns:
//   - *Parser: A pointer to the new parser.
//   - error: An error if the parser could not be created.
//
// Errors:
//   - *uc.ErrInvalidParameter: The grammar is nil.
//   - *gr.ErrNoProductionRulesFound: No production rules are found in the grammar.
func NewParser[T uc.Enumer](grammar *Grammar[T]) (*Parser[T], error) {
	if grammar == nil {
		return nil, uc.NewErrNilParameter("grammar")
	}

	productions := grammar.GetProductions()
	if len(productions) == 0 {
		return nil, gr.NewErrNoProductionRulesFound()
	}

	table, err := cs.SolveConflicts(grammar.GetSymbols(), productions)
	if err != nil {
		return nil, err
	}

	p := &Parser[T]{
		dt: table,
	}

	return p, nil
}

// Parse parses the input stream using the parser's decision function.
//
// Parameters:
//   - p: The parser to use.
//   - source: The input stream to parse.
//
// Returns:
//   - error: An error if the input stream could not be parsed.
func Parse[T uc.Enumer](p *Parser[T], source *cds.Stream[*gr.Token[T]]) error {
	if p == nil {
		return uc.NewErrNilParameter("parser")
	}

	if p.dt == nil {
		return errors.New("no grammar was set")
	}

	if source == nil || source.IsEmpty() {
		return errors.New("source is empty")
	}

	ceRoot := NewCurrentEval[T]()

	err := ceRoot.shift(source)
	if err != nil {
		return err
	}

	sols := evaluate(p.dt, source, ceRoot)

	results, err := extractResults(sols)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		return errors.New("no parse trees were found")
	}

	p.evals = results

	return nil
}

// GetParseTree returns the parse tree that the parser has generated.
//
// Parse() must be called before calling this method. If it is not, an error will
// be returned.
//
// Returns:
//   - []*gr.TokenTree: A slice of parse trees.
//   - error: An error if the parse tree could not be retrieved.
func (p *Parser[T]) GetParseTree() ([]*gr.TokenTree[T], error) {
	if len(p.evals) == 0 {
		return nil, errors.New("nothing was parsed. Use Parse() to parse the input stream")
	}

	var forest []*gr.TokenTree[T]

	for _, eval := range p.evals {
		tmp, err := eval.GetParseTree()
		if err == nil {
			forest = append(forest, tmp...)
		}
	}

	if len(forest) == 0 {
		return nil, errors.New("no parse trees were found")
	}

	return forest, nil
}
