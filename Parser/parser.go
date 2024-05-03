package Parser

import (
	"errors"
	"fmt"

	cs "github.com/PlayerR9/LyneParser/ConflictSolver"
	gr "github.com/PlayerR9/LyneParser/Grammar"

	ds "github.com/PlayerR9/MyGoLib/ListLike/DoubleLL"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// Parser is a parser that uses a stack to parse a stream of tokens.
type Parser struct {
	// productions represents the productions that the parser will use.
	productions []*gr.Production

	// inputStream represents the stream of tokens that the parser will parse.
	inputStream *gr.TokenStream

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
func NewParser(grammar *gr.Grammar) (Parser, error) {
	p := Parser{productions: make([]*gr.Production, 0)}

	if grammar == nil {
		return p, ers.NewErrNilParameter("grammar")
	}

	for _, production := range grammar.Productions {
		prod, ok := production.(*gr.Production)
		if !ok {
			continue
		}

		p.productions = append(p.productions, prod)
	}

	if len(p.productions) == 0 {
		return p, gr.NewErrNoProductionRulesFound()
	}

	p.dt = NewDecisionTable()

	err := p.dt.GenerateItems(p.productions)
	if err != nil {
		return p, err
	}

	err = p.dt.FixConflicts()
	if err != nil {
		return p, err
	}

	return p, nil
}

// SetInputStream sets the input stream that the parser will parse. It also adds
// an EOF token to the end of the input stream if it is not already present.
//
// Parameters:
//   - inputStream: The input stream that the parser will parse.
//
// Returns:
//   - error: An error of type *ers.ErrInvalidParameter if the input stream is nil
//     or empty.
func (p *Parser) SetInputStream(inputStream *gr.TokenStream) error {
	if inputStream == nil || inputStream.IsEmpty() {
		return ers.NewErrInvalidParameter(
			"inputStream",
			ers.NewErrEmptySlice(),
		)
	}

	// Reset the input stream to the beginning.
	inputStream.Reset()

	// Add EOF token to the end of the input stream (if it is not already present).
	inputStream.SetEOFToken()

	// Add lookahead to all tokens
	inputStream.SetLookahead()

	p.inputStream = inputStream

	return nil
}

// Parse parses the input stream using the parser's decision function.
//
// SetInputStream() and SetDecisionFunc() must be called before calling this
// method. If they are not, an error will be returned.
//
// Returns:
//   - error: An error if the input stream could not be parsed.
func (p *Parser) Parse() error {
	if p.inputStream == nil {
		return errors.New("call SetInputStream() first")
	} else if p.inputStream.IsEmpty() || p.inputStream.IsDone() {
		return errors.New("input stream is empty or done. Use SetInputStream() to set a new stream")
	}

	if p.dt == nil {
		return errors.New("no grammar was set")
	}

	p.stack = ds.NewDoubleLinkedStack[gr.Tokener]()

	// Initial shift
	var decision cs.Actioner

	decision = cs.NewActShift()

	err := p.shift()
	if err != nil {
		return err
	}

	for !p.stack.IsEmpty() {
		if _, ok := decision.(*cs.ActAccept); ok {
			break
		}

		decision = p.dt.Match(p.stack)
		p.stack.Refuse()

		switch decision := decision.(type) {
		case *cs.ActShift:
			err := p.shift()
			if err != nil {
				return err
			}
		case *cs.ActReduce:
			err := p.reduce(decision.RuleIndex)
			if err != nil {
				p.stack.Refuse()
				return err
			}

			p.stack.Accept()
		case *cs.ActAccept:
			err := p.reduce(decision.RuleIndex)
			if err != nil {
				p.stack.Refuse()
				return err
			}

			p.stack.Accept()
		case *cs.ActError:
			return decision.Reason
		default:
			return fmt.Errorf("unknown action type: %T", decision)
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

	for !p.stack.IsEmpty() {
		top := p.stack.Pop()

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
func (p *Parser) shift() error {
	tok, err := p.inputStream.Consume()
	if err != nil {
		return NewErrNoAccept()
	}

	p.stack.Push(tok)

	return nil
}

// reduce is a helper method that reduces the stack by a rule.
//
// Parameters:
//   - rule: The index of the rule to reduce by.
//
// Returns:
//   - error: An error if the stack could not be reduced.
func (p *Parser) reduce(rule int) error {
	lhs := p.productions[rule].GetLhs()
	rhss := p.productions[rule].ReverseIterator()

	var lookahead *gr.LeafToken = nil

	for {
		value, err := rhss.Consume()
		if err != nil {
			break
		}

		if p.stack.IsEmpty() {
			return NewErrAfter(lhs, ers.NewErrUnexpected(nil, value))
		}

		top := p.stack.Pop()

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
func FullParse(grammar *gr.Grammar, inputStream *gr.TokenStream, dt *DecisionTable) ([]gr.NonLeafToken, error) {
	parser, err := NewParser(grammar)
	if err != nil {
		return nil, fmt.Errorf("could not create parser: %s", err.Error())
	}

	err = parser.SetInputStream(inputStream)
	if err != nil {
		return nil, fmt.Errorf("could not set input stream: %s", err.Error())
	}

	err = parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("parse error: %s", err.Error())
	}

	roots, err := parser.GetParseTree()
	if err != nil {
		return nil, fmt.Errorf("could not get parse tree: %s", err.Error())
	}

	return roots, nil
}
