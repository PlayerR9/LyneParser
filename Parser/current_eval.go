package Parser

import (
	"fmt"

	cs "github.com/PlayerR9/LyneParser/ConflictSolver"
	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	ud "github.com/PlayerR9/MyGoLib/Units/Debugging"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// CurrentEval is a struct that represents the current evaluation of the parser.
type CurrentEval[T gr.TokenTyper] struct {
	// stack represents the stack that the parser will use.
	stack *ud.History[lls.Stacker[*gr.Token[T]]]

	// current_index is the current index of the input stream.
	current_index int

	// is_done is a flag that represents if the parser has finished parsing.
	is_done bool
}

// Copy creates a copy of the current evaluation.
//
// Returns:
//   - uc.Copier: A copy of the current evaluation.
func (ce *CurrentEval[T]) Copy() uc.Copier {
	ce_copy := &CurrentEval[T]{
		stack:         ce.stack.Copy().(*ud.History[lls.Stacker[*gr.Token[T]]]),
		current_index: ce.current_index,
		is_done:       ce.is_done,
	}
	return ce_copy
}

// Accept returns true if the current evaluation has accepted the input stream.
//
// Returns:
//   - bool: True if the current evaluation has accepted the input stream.
func (ce *CurrentEval[T]) Accept() bool {
	return ce.is_done
}

// NewCurrentEval creates a new current evaluation.
//
// Returns:
//   - *CurrentEval: A new current evaluation.
func NewCurrentEval[T gr.TokenTyper]() *CurrentEval[T] {
	ce := &CurrentEval[T]{
		current_index: 0,
		is_done:       false,
	}

	ce.stack = lls.NewStackWithHistory(lls.NewLinkedStack[*gr.Token[T]]())

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
func (ce *CurrentEval[T]) GetParseTree() ([]*gr.TokenTree[T], error) {
	var forest []*gr.TokenTree[T]

	for {
		cmd := lls.NewPop[*gr.Token[T]]()
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
func (ce *CurrentEval[T]) shift(source *cds.Stream[*gr.Token[T]]) error {
	toks, err := source.Get(ce.current_index, 1)
	if err != nil || len(toks) == 0 {
		return NewErrNoAccept()
	}

	cmd := lls.NewPush(toks[0])
	err = ce.stack.ExecuteCommand(cmd)
	if err != nil {
		return fmt.Errorf("could not push token: %s", err.Error())
	}

	ce.current_index++

	return nil
}

// reduce is a helper method that reduces the stack by a rule.
//
// Parameters:
//   - rule: The rule to reduce by.
//
// Returns:
//   - error: An error if the stack could not be reduced.
func (ce *CurrentEval[T]) reduce(rule *gr.Production[T]) error {
	lhs := rule.GetLhs()
	rhss := rule.ReverseIterator()

	var lookahead *gr.Token[T]
	var popped []*gr.Token[T]

	for {
		value, err := rhss.Consume()
		if err != nil {
			break
		}

		cmd := lls.NewPop[*gr.Token[T]]()
		err = ce.stack.ExecuteCommand(cmd)
		if err != nil {
			ce.stack.Reject()
			return uc.NewErrAfter(lhs.String(), uc.NewErrUnexpected("", value.String()))
		}
		top := cmd.Value()

		popped = append(popped, top)

		if lookahead == nil {
			lookahead = top.GetLookahead()
		}

		id := top.GetID()
		if id != value {
			ce.stack.Reject()
			return uc.NewErrAfter(lhs.String(), uc.NewErrUnexpected(top.GoString(), value.String()))
		}
	}

	ce.stack.Accept()

	tok := gr.NewToken(lhs, popped, 0, lookahead)

	cmd := lls.NewPush(tok)
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
func (ce *CurrentEval[T]) ActOnDecision(decision cs.HelperElem[T], source *cds.Stream[*gr.Token[T]]) error {
	var err error

	switch decision := decision.(type) {
	case *cs.ActShift[T]:
		err = ce.shift(source)
	case *cs.ActReduce[T]:
		err = ce.reduce(decision.Original)
	case *cs.ActAccept[T]:
		err = ce.reduce(decision.Original)
		if err == nil {
			ce.is_done = true
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
func (ce *CurrentEval[T]) Parse(source *cds.Stream[*gr.Token[T]], dt *cs.ConflictSolver[T]) ([]*CurrentEval[T], error) {
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

		return []*CurrentEval[T]{ce}, nil
	default:
		ce_copies := make([]*CurrentEval[T], 0, len(decisions))

		for _, decision := range decisions {
			ce_copy := ce.Copy().(*CurrentEval[T])

			err := ce_copy.ActOnDecision(decision, source)
			if err != nil {
				continue
			}

			ce_copies = append(ce_copies, ce_copy)
		}

		return ce_copies, nil
	}
}
