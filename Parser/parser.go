package Parser

import (
	"errors"
	"fmt"

	cs "github.com/PlayerR9/LyneParser/ConflictSolver"
	gr "github.com/PlayerR9/LyneParser/Grammar"

	ds "github.com/PlayerR9/MyGoLib/ListLike/DoubleLL"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

/////////////////////////////////////////////////////////////

// Parser is a parser that uses a stack to parse a stream of tokens.
type Parser struct {
	// stack represents the stack that the parser will use.
	stack *ds.DoubleStack[gr.Tokener]

	// decisionFunc represents the function that the parser will use to determine
	// the next action to take.
	dt *DecisionTable
}

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
		dt: &DecisionTable{
			table: table,
		},
	}

	return p, nil
}

// Parse parses the input stream using the parser's decision function.
//
// SetInputStream() and SetDecisionFunc() must be called before calling this
// method. If they are not, an error will be returned.
//
// Returns:
//   - error: An error if the input stream could not be parsed.
func (p *Parser) Parse(source *gr.TokenStream) error {
	if source == nil || source.IsEmpty() {
		return errors.New("source is empty")
	}

	if p.dt == nil {
		return errors.New("no grammar was set")
	}

	p.stack = ds.NewDoubleLinkedStack[gr.Tokener]()

	// Initial shift
	var decision cs.Actioner

	decision = cs.NewActShift()

	err := p.shift(source)
	if err != nil {
		return err
	}

	for !p.stack.IsEmpty() {
		if _, ok := decision.(*cs.ActAccept); ok {
			break
		}

		decision, err = p.dt.Match(p.stack)
		p.stack.Refuse()

		if err != nil {
			return err
		}

		switch decision := decision.(type) {
		case *cs.ActShift:
			err := p.shift(source)
			if err != nil {
				return err
			}
		case *cs.ActReduce:
			err := p.reduce(decision.GetRule())
			if err != nil {
				p.stack.Refuse()
				return err
			}

			p.stack.Accept()
		case *cs.ActAccept:
			err := p.reduce(decision.GetRule())
			if err != nil {
				p.stack.Refuse()
				return err
			}

			p.stack.Accept()
		default:
			return NewErrUnknownAction(decision)
		}
	}

	return nil
}

// GetParseTree returns the parse tree that the parser has generated.
//
// Parse() must be called before calling this method. If it is not, an error will
// be returned.
//
// Returns:
//   - []gr.NonLeafToken: The parse tree.
//   - error: An error if the parse tree could not be retrieved.
func (p *Parser) GetParseTree() ([]gr.NonLeafToken, error) {
	if p.stack.IsEmpty() {
		return nil, errors.New("nothing was parsed. Use Parse() to parse the input stream")
	}

	roots := make([]gr.NonLeafToken, 0)

	for {
		top, err := p.stack.Pop()
		if err != nil {
			break
		}

		root, ok := top.(*gr.NonLeafToken)
		if !ok {
			continue
		}

		roots = append(roots, *root)
	}

	return roots, nil
}

// shift is a helper method that shifts the current token onto the stack.
//
// Returns:
//   - error: An error of type *ErrNoAccept if the input stream is done.
func (p *Parser) shift(source *gr.TokenStream) error {
	tok, err := source.Consume()
	if err != nil {
		return NewErrNoAccept()
	}

	p.stack.Push(tok)

	return nil
}

// reduce is a helper method that reduces the stack by a rule.
//
// Parameters:
//   - rule: The rule to reduce by.
//
// Returns:
//   - error: An error if the stack could not be reduced.
func (p *Parser) reduce(rule *gr.Production) error {
	lhs := rule.GetLhs()
	rhss := rule.ReverseIterator()

	var lookahead *gr.LeafToken = nil

	for {
		value, err := rhss.Consume()
		if err != nil {
			break
		}

		top, err := p.stack.Pop()
		if err != nil {
			return NewErrAfter(lhs, ers.NewErrUnexpected(nil, value))
		}

		if lookahead == nil {
			lookahead = top.GetLookahead()
		}

		if top.GetID() != value {
			return NewErrAfter(lhs, ers.NewErrUnexpected(top, value))
		}
	}

	data := p.stack.GetExtracted()
	tok := gr.NewNonLeafToken(lhs, 0, data...)
	tok.Lookahead = lookahead
	p.stack.Push(tok)

	return nil
}

// FullParse parses the input stream using the given grammar and decision
// function. It is a convenience function intended for simple parsing tasks.
//
// Parameters:
//
//   - grammar: The grammar that the parser will use.
//   - inputStream: The input stream that the parser will parse.
//   - decisionFunc: The decision function that the parser will use.
//
// Returns:
//
//   - []gr.NonLeafToken: The parse tree.
//   - error: An error if the input stream could not be parsed.
func FullParse(grammar *gr.Grammar, source *gr.TokenStream, dt *DecisionTable) ([]gr.NonLeafToken, error) {
	parser, err := NewParser(grammar)
	if err != nil {
		return nil, fmt.Errorf("could not create parser: %s", err.Error())
	}

	err = parser.Parse(source)
	if err != nil {
		return nil, fmt.Errorf("parse error: %s", err.Error())
	}

	roots, err := parser.GetParseTree()
	if err != nil {
		return nil, fmt.Errorf("could not get parse tree: %s", err.Error())
	}

	return roots, nil
}
