package Parser

import (
	"errors"
	"fmt"

	gr "github.com/PlayerR9/LyneParser/Grammar"

	ds "github.com/PlayerR9/MyGoLib/ListLike/DoubleLL"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// DecisionFunc is a function that is used to determine the next action to take
// in the parser.
//
// Parameters:
//
//   - stack: The stack that the parser is using.
//   - lookahead: The lookahead token.
//
// Returns:
//
//   - Action: The action to take.
type DecisionFunc func(stack *ds.DoubleStack[gr.Tokener]) Action

// Parser is a parser that uses a stack to parse a stream of tokens.
type Parser struct {
	// productions represents the productions that the parser will use.
	productions []gr.Production

	// inputStream represents the stream of tokens that the parser will parse.
	inputStream gr.TokenStream

	// stack represents the stack that the parser will use.
	stack *ds.DoubleStack[gr.Tokener]

	// decisionFunc represents the function that the parser will use to determine
	// the next action to take.
	decisionFunc DecisionFunc
}

// NewParser creates a new parser with the given grammar.
//
// If the grammar is nil or has no productions that are of type *gr.Production,
// an error of type *ers.ErrInvalidParameter will be returned.
//
// Parameters:
//
//   - grammar: The grammar that the parser will use.
//
// Returns:
//
//   - *Parser: A pointer to the new parser.
//   - error: An error if the parser could not be created.
func NewParser(grammar *gr.Grammar) (Parser, error) {
	p := Parser{productions: make([]gr.Production, 0)}

	if grammar == nil {
		return p, ers.NewErrNilParameter("grammar")
	}

	for _, production := range grammar.Productions {
		prod, ok := production.(*gr.Production)
		if !ok {
			continue
		}

		p.productions = append(p.productions, *prod)
	}

	if len(p.productions) == 0 {
		return p, ers.NewErrInvalidParameter(
			"grammar",
			errors.New("no productions found"),
		)
	}

	return p, nil
}

// SetDecisionFunc sets the decision function that the parser will use to
// determine the next action to take.
//
// If the decision function is nil, an error of type *ers.ErrInvalidParameter
// will be returned.
//
// Parameters:
//
//   - decisionFunc: The decision function that the parser will use.
//
// Returns:
//
//   - error: An error if the decision function could not be set.
func (p *Parser) SetDecisionFunc(decisionFunc DecisionFunc) error {
	if decisionFunc == nil {
		return ers.NewErrNilParameter("decisionFunc")
	}

	p.decisionFunc = decisionFunc

	return nil
}

// SetInputStream sets the input stream that the parser will parse. It also adds
// an EOF token to the end of the input stream if it is not already present.
//
// If the input stream is nil or empty, an error of type *ers.ErrInvalidParameter
// will be returned.
//
// Parameters:
//
//   - inputStream: The input stream that the parser will parse.
//
// Returns:
//
//   - error: An error if the input stream could not be set.
func (p *Parser) SetInputStream(inputStream gr.TokenStream) error {
	if inputStream.IsEmpty() {
		return ers.NewErrInvalidParameter(
			"inputStream",
			errors.New("value is empty"),
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
//
//   - error: An error if the input stream could not be parsed.
func (p *Parser) Parse() error {
	if p.inputStream.Size() == 0 || p.inputStream.IsDone() {
		return errors.New("call SetInputStream() first")
	} else if p.decisionFunc == nil {
		return errors.New("call SetDecisionFunc() first")
	}

	p.stack = ds.NewDoubleLinkedStack[gr.Tokener]()

	// Initial shift
	decision := NewShiftAction()

	err := p.shift()
	if err != nil {
		return err
	}

	for !p.stack.IsEmpty() && decision.Type != ActAccept {
		decision = p.decisionFunc(p.stack)
		p.stack.Refuse()

		switch decision.Type {
		case ActShift:
			err := p.shift()
			if err != nil {
				return err
			}
		case ActReduce:
			err := p.reduce(decision.Data.(int))
			if err != nil {
				p.stack.Refuse()
				return err
			}

			p.stack.Accept()
		case ActAccept:
			err := p.reduce(decision.Data.(int))
			if err != nil {
				p.stack.Refuse()
				return err
			}

			p.stack.Accept()
		case ActError:
			return decision.Data.(error)
		default:
			return fmt.Errorf("unknown action type: %v", decision.Type)
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
//
//   - []gr.NonLeafToken: The parse tree.
//   - error: An error if the parse tree could not be retrieved.
func (p *Parser) GetParseTree() ([]gr.NonLeafToken, error) {
	if p.stack.IsEmpty() {
		return nil, errors.New("call Parse() first")
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
//
//   - error: An error if the token could not be shifted.
func (p *Parser) shift() error {
	if p.inputStream.IsDone() {
		return NewErrNoAccept()
	}

	tok := p.inputStream.Consume()

	p.stack.Push(tok)

	return nil
}

// reduce is a helper method that reduces the stack by a rule.
//
// Parameters:
//
//   - rule: The index of the rule to reduce by.
//
// Returns:
//
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
			return fmt.Errorf("after %s: %v", lhs, ers.NewErrUnexpected(nil, value))
		}

		top := p.stack.Pop()

		if lookahead == nil {
			lookahead = top.GetLookahead()
		}

		if top.GetID() != value {
			return fmt.Errorf("after %s: %v", lhs, ers.NewErrUnexpected(top, value))
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
func FullParse(grammar *gr.Grammar, inputStream gr.TokenStream, decisionFunc DecisionFunc) ([]gr.NonLeafToken, error) {
	parser, err := NewParser(grammar)
	if err != nil {
		return nil, fmt.Errorf("could not create parser: %v", err)
	}

	err = parser.SetInputStream(inputStream)
	if err != nil {
		return nil, fmt.Errorf("could not set input stream: %v", err)
	}

	err = parser.SetDecisionFunc(decisionFunc)
	if err != nil {
		return nil, fmt.Errorf("could not set decision function: %v", err)
	}

	err = parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("parse error: %v", err)
	}

	roots, err := parser.GetParseTree()
	if err != nil {
		return nil, fmt.Errorf("could not get parse tree: %v", err)
	}

	return roots, nil
}
