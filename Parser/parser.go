package Parser

import (
	"errors"

	com "github.com/PlayerR9/LyneParser/Common"
	cs "github.com/PlayerR9/LyneParser/ConflictSolver"
	gr "github.com/PlayerR9/LyneParser/Grammar"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"

	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

// Parser is a parser that uses a stack to parse a stream of tokens.
type Parser struct {

	// evals is a list of evaluations that the parser will use.
	evals []*CurrentEval

	// decisionFunc represents the function that the parser will use to determine
	// the next action to take.
	dt *cs.ConflictSolver
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
//   - *ers.ErrInvalidParameter: The grammar is nil.
//   - *gr.ErrNoProductionRulesFound: No production rules are found in the grammar.
func NewParser(grammar *gr.Grammar) (*Parser, error) {
	if grammar == nil {
		return nil, ers.NewErrNilParameter("grammar")
	}

	productions := grammar.GetProductions()
	if len(productions) == 0 {
		return nil, gr.NewErrNoProductionRulesFound()
	}

	table, err := cs.SolveConflicts(grammar.Symbols, productions)
	if err != nil {
		return nil, err
	}

	p := &Parser{
		dt: table,
	}

	return p, nil
}

// Parse parses the input stream using the parser's decision function.
//
// Returns:
//   - error: An error if the input stream could not be parsed.
func (p *Parser) Parse(source *com.TokenStream) error {
	if source == nil || source.IsEmpty() {
		return errors.New("source is empty")
	}

	if p.dt == nil {
		return errors.New("no grammar was set")
	}

	todo := []*CurrentEval{NewCurrentEval()}

	err := todo[0].shift(source)
	if err != nil {
		return err
	}

	done := slext.DoWhile(
		todo,
		func(eval *CurrentEval) bool { return eval.isDone },
		func(eval *CurrentEval) ([]*CurrentEval, error) {
			sol, err := eval.Parse(source, p.dt)
			if err != nil {
				panic(err)
			}

			return sol, nil
		},
	)

	if len(done) == 0 {
		return errors.New("no parse trees were found")
	}

	p.evals = done

	return nil
}

// GetParseTree returns the parse tree that the parser has generated.
//
// Parse() must be called before calling this method. If it is not, an error will
// be returned.
//
// Returns:
//   - []*com.TokenTree: A slice of parse trees.
//   - error: An error if the parse tree could not be retrieved.
func (p *Parser) GetParseTree() ([]*com.TokenTree, error) {
	if len(p.evals) == 0 {
		return nil, errors.New("nothing was parsed. Use Parse() to parse the input stream")
	}

	forest := make([]*com.TokenTree, 0)

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
