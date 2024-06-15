package Parser

import (
	"errors"

	cs "github.com/PlayerR9/LyneParser/ConflictSolver"
	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	ds "github.com/PlayerR9/MyGoLib/ListLike/DoubleLL"
	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
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
//   - uc.Copier: A copy of the current evaluation.
func (ce *CurrentEval) Copy() uc.Copier {
	return &CurrentEval{
		stack:        ce.stack.Copy().(*ds.DoubleStack[gr.Tokener]),
		currentIndex: ce.currentIndex,
		isDone:       ce.isDone,
	}
}

// Accept returns true if the current evaluation has accepted the input stream.
//
// Returns:
//   - bool: True if the current evaluation has accepted the input stream.
func (ce *CurrentEval) Accept() bool {
	return ce.isDone
}

// NewCurrentEval creates a new current evaluation.
//
// Returns:
//   - *CurrentEval: A new current evaluation.
func NewCurrentEval() *CurrentEval {
	ce := &CurrentEval{
		currentIndex: 0,
		isDone:       false,
	}

	stack, err := ds.NewDoubleStack(Stacker.NewLinkedStack[gr.Tokener]())
	if err != nil {
		panic(err)
	}

	ce.stack = stack

	return ce
}

// GetParseTree returns the parse tree that the parser has generated.
//
// Parse() must be called before calling this method. If it is not, an error will
// be returned.
//
// Returns:
//   - []*gr.TokenTree: A slice of parse trees.
//   - error: An error if the parse tree could not be retrieved.
//
// Errors:
//   - *ue.ErrInvalidUsage: If Parse() has not been called.
//   - *gr.ErrCycleDetected: A cycle is detected in the token tree.
//   - *ue.ErrInvalidParameter: The top of the stack is nil.
//   - *gr.ErrUnknowToken: The root is not a known token.
func (ce *CurrentEval) GetParseTree() ([]*gr.TokenTree, error) {
	if ce.stack.IsEmpty() {
		return nil, ue.NewErrInvalidUsage(
			NewErrNothingWasParsed(),
			"Use Parse() to parse the input stream",
		)
	}

	var forest []*gr.TokenTree

	for {
		top, ok := ce.stack.Pop()
		if !ok {
			break
		}

		tree, err := gr.NewTokenTree(top)
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
func (ce *CurrentEval) shift(source *cds.Stream[*gr.LeafToken]) error {
	toks, err := source.Get(ce.currentIndex, 1)
	if err != nil || len(toks) == 0 {
		return NewErrNoAccept()
	}

	ok := ce.stack.Push(toks[0])
	if !ok {
		return ue.NewErrUnexpectedError(errors.New("could not push token onto stack"))
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

		top, ok := ce.stack.Pop()
		if !ok {
			ce.stack.Refuse()
			return ue.NewErrAfter(lhs, ue.NewErrUnexpected("", value))
		}

		if lookahead == nil {
			lookahead = top.GetLookahead()
		}

		id := top.GetID()
		if id != value {
			ce.stack.Refuse()
			return ue.NewErrAfter(lhs, ue.NewErrUnexpected(top.GoString(), value))
		}
	}

	data := ce.stack.GetExtracted()
	ce.stack.Accept()

	tok := gr.NewNonLeafToken(lhs, 0, data...)
	tok.Lookahead = lookahead

	ok := ce.stack.Push(tok)
	if !ok {
		return ue.NewErrUnexpectedError(errors.New("could not push token onto stack"))
	}

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
func (ce *CurrentEval) ActOnDecision(decision cs.HelperElem, source *cds.Stream[*gr.LeafToken]) error {
	var err error

	switch decision := decision.(type) {
	case *cs.ActShift:
		err = ce.shift(source)
	case *cs.ActReduce:
		rule := decision.GetOriginal()

		err = ce.reduce(rule)
		if err == nil && decision.ShouldAccept() {
			ce.isDone = true
		}
	default:
		err = NewErrUnknownAction(decision)
	}

	if err != nil {
		return err
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
func (ce *CurrentEval) Parse(source *cds.Stream[*gr.LeafToken], dt *cs.ConflictSolver) ([]*CurrentEval, error) {
	decisions, err := dt.Match(ce.stack)
	ce.stack.Refuse()

	if err != nil {
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
