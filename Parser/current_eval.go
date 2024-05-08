package Parser

import (
	"errors"

	com "github.com/PlayerR9/LyneParser/Common"
	cs "github.com/PlayerR9/LyneParser/ConflictSolver"
	gr "github.com/PlayerR9/LyneParser/Grammar"
	ds "github.com/PlayerR9/MyGoLib/ListLike/DoubleLL"

	intf "github.com/PlayerR9/MyGoLib/Units/Common"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// CurrentEval is a struct that represents the current evaluation of the parser.
type CurrentEval struct {
	// stack represents the stack that the parser will use.
	stack *ds.DoubleStack[gr.Tokener]

	// currentIndex is the current index of the input stream.
	currentIndex int

	// isDone is a flag that represents if the parser has finished parsing.
	isDone bool
}

// Copy creates a copy of the current evaluation.
//
// Returns:
//   - intf.Copier: A copy of the current evaluation.
func (ce *CurrentEval) Copy() intf.Copier {
	return &CurrentEval{
		stack:        ce.stack.Copy().(*ds.DoubleStack[gr.Tokener]),
		currentIndex: ce.currentIndex,
		isDone:       ce.isDone,
	}
}

// NewCurrentEval creates a new current evaluation.
//
// Returns:
//   - *CurrentEval: A new current evaluation.
func NewCurrentEval() *CurrentEval {
	return &CurrentEval{
		stack:        ds.NewDoubleLinkedStack[gr.Tokener](),
		currentIndex: 0,
		isDone:       false,
	}
}

// GetParseTree returns the parse tree that the parser has generated.
//
// Parse() must be called before calling this method. If it is not, an error will
// be returned.
//
// Returns:
//   - []*com.TokenTree: A slice of parse trees.
//   - error: An error if the parse tree could not be retrieved.
func (ce *CurrentEval) GetParseTree() ([]*com.TokenTree, error) {
	if ce.stack.IsEmpty() {
		return nil, errors.New("nothing was parsed. Use Parse() to parse the input stream")
	}

	forest := make([]*com.TokenTree, 0)

	for {
		top, err := ce.stack.Pop()
		if err != nil {
			break
		}

		tree, err := com.NewTokenTree(top)
		if err != nil {
			return nil, err
		}

		forest = append(forest, tree)
	}

	return forest, nil
}

// shift is a helper method that shifts the current token onto the stack.
//
// Returns:
//   - error: An error of type *ErrNoAccept if the input stream is done.
func (ce *CurrentEval) shift(source *com.TokenStream) error {
	toks, err := source.Get(ce.currentIndex, 1)
	if err != nil || len(toks) == 0 {
		return NewErrNoAccept()
	}

	err = ce.stack.Push(toks[0])
	if err != nil {
		return err
	}

	ce.currentIndex++

	return nil
}

// reduce is a helper method that reduces the stack by a rule.
//
// Parameters:
//   - rule: The rule to reduce by.
//
// Returns:
//   - error: An error if the stack could not be reduced.
func (ce *CurrentEval) reduce(rule *gr.Production) error {
	lhs := rule.GetLhs()
	rhss := rule.ReverseIterator()

	var lookahead *gr.LeafToken = nil

	for {
		value, err := rhss.Consume()
		if err != nil {
			break
		}

		top, err := ce.stack.Pop()
		if err != nil {
			return ers.NewErrAfter(lhs, ers.NewErrUnexpected(nil, value))
		}

		if lookahead == nil {
			lookahead = top.GetLookahead()
		}

		if top.GetID() != value {
			return ers.NewErrAfter(lhs, ers.NewErrUnexpected(top, value))
		}
	}

	data := ce.stack.GetExtracted()
	tok := gr.NewNonLeafToken(lhs, 0, data...)
	tok.Lookahead = lookahead
	ce.stack.Push(tok)

	return nil
}

// ActOnDecision acts on a decision that the parser has made.
//
// Parameters:
//   - decision: The decision that the parser has made.
//   - source: The source of the input stream.
//
// Returns:
//   - bool: True if the parser has accepted the input stream.
//   - error: An error if the parser could not act on the decision.
func (ce *CurrentEval) ActOnDecision(decision cs.Actioner, source *com.TokenStream) error {
	switch decision := decision.(type) {
	case *cs.ActShift:
		err := ce.shift(source)
		if err != nil {
			return err
		}
	case *cs.ActReduce:
		err := ce.reduce(decision.GetRule())
		if err != nil {
			ce.stack.Refuse()
			return err
		}

		ce.stack.Accept()
	case *cs.ActAccept:
		err := ce.reduce(decision.GetRule())
		if err != nil {
			ce.stack.Refuse()
			return err
		}

		ce.stack.Accept()

		ce.isDone = true
	default:
		return NewErrUnknownAction(decision)
	}

	return nil
}

// Parse parses the input stream using the parser's decision table.
//
// Parameters:
//   - source: The source of the input stream.
//   - dt: The decision table to use.
//
// Returns:
//   - []*CurrentEval: A slice of current evaluations.
//   - error: An error if the input stream could not be parsed.
func (ce *CurrentEval) Parse(source *com.TokenStream, dt *cs.ConflictSolver) ([]*CurrentEval, error) {
	if ce.stack.IsEmpty() {
		ce.isDone = true

		return []*CurrentEval{ce}, nil
	}

	decisions, err := dt.Match(ce.stack)
	ce.stack.Refuse()

	if err != nil && !ers.As[*cs.ErrAmbiguousGrammar](err) {
		return nil, err
	}

	switch len(decisions) {
	case 0:
		return nil, NewErrNoAccept()
	case 1:
		err := ce.ActOnDecision(decisions[0], source)
		if err != nil {
			return nil, err
		}

		return []*CurrentEval{ce}, nil
	default:
		ceCopies := make([]*CurrentEval, 0, len(decisions))

		for _, decision := range decisions {
			ceCopy := ce.Copy().(*CurrentEval)

			err := ceCopy.ActOnDecision(decision, source)
			if err != nil {
				continue
			}

			ceCopies = append(ceCopies, ceCopy)
		}

		return ceCopies, nil
	}
}
