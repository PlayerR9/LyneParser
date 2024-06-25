package Parser

import (
	cs "github.com/PlayerR9/LyneParser/ConflictSolver"
	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	ud "github.com/PlayerR9/MyGoLib/Units/Debugging"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// CurrentEval is a struct that represents the current evaluation of the parser.
type CurrentEval struct {
	// stack represents the stack that the parser will use.
	stack *ud.History[lls.Stacker[gr.Tokener]]

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
		stack:        ce.stack.Copy().(*ud.History[lls.Stacker[gr.Tokener]]),
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

	ce.stack = lls.NewStackWithHistory(lls.NewLinkedStack[gr.Tokener]())

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
//   - *uc.ErrInvalidUsage: If Parse() has not been called.
//   - *gr.ErrCycleDetected: A cycle is detected in the token tree.
//   - *uc.ErrInvalidParameter: The top of the stack is nil.
//   - *gr.ErrUnknowToken: The root is not a known token.
func (ce *CurrentEval) GetParseTree() ([]*gr.TokenTree, error) {
	var forest []*gr.TokenTree

	for {
		cmd := lls.NewPop[gr.Tokener]()
		err := ce.stack.ExecuteCommand(cmd)
		if err != nil {
			break
		}
		top := cmd.Value()

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

	cmd := lls.NewPush[gr.Tokener](toks[0])
	ce.stack.ExecuteCommand(cmd)

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
	var popped []gr.Tokener

	for {
		value, err := rhss.Consume()
		if err != nil {
			break
		}

		cmd := lls.NewPop[gr.Tokener]()
		err = ce.stack.ExecuteCommand(cmd)
		if err != nil {
			ce.stack.Reject()
			return uc.NewErrAfter(lhs, uc.NewErrUnexpected("", value))
		}
		top := cmd.Value()

		popped = append(popped, top)

		if lookahead == nil {
			lookahead = top.GetLookahead()
		}

		id := top.GetID()
		if id != value {
			ce.stack.Reject()
			return uc.NewErrAfter(lhs, uc.NewErrUnexpected(top.GoString(), value))
		}
	}

	ce.stack.Accept()

	tok := gr.NewNonLeafToken(lhs, 0, popped...)
	tok.Lookahead = lookahead

	cmd := lls.NewPush[gr.Tokener](tok)
	ce.stack.ExecuteCommand(cmd)

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
	ce.stack.Reject()

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
