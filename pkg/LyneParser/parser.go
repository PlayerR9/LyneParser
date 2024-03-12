package LyneParser

import (
	ll "LyneParser/pkg/LyneLexer"
	"errors"
	"fmt"

	gr "LyneParser/pkg/Grammar"

	Stack "github.com/PlayerR9/MyGoLib/CustomData/ListLike/Stack"

	itf "github.com/PlayerR9/MyGoLib/Interfaces"
)

type ParserDecision int

const (
	ActShift ParserDecision = iota
	ActReduce
	ActAccept
	ActError
)

type ASTer interface {
	Root() *Node
	Nodes() []*Node
	Print() string
}

type ParserOption func(*LyneParser) error

func WithGrammar(grammar *gr.Grammar) ParserOption {
	return func(lp *LyneParser) error {
		if grammar == nil {
			return nil
		}

		lp.grammar = grammar

		return nil
	}
}

type LyneParser struct {
	grammar *gr.Grammar
	states  Stack.Stacker[*State]

	inputStream Stack.Stacker[*ll.Token]
	trees       Stack.Stacker[*Node]
}

func (lp *LyneParser) Copy() itf.Copier {
	pCopy := LyneParser{
		grammar:     lp.grammar,
		states:      lp.states.Copy().(Stack.Stacker[*State]),
		inputStream: lp.inputStream.Copy().(Stack.Stacker[*ll.Token]),
		trees:       lp.trees.Copy().(Stack.Stacker[*Node]),
	}

	return &pCopy
}

func NewParser(options ...ParserOption) (*LyneParser, error) {
	lp := &LyneParser{}

	for _, opt := range options {
		err := opt(lp)
		if err != nil {
			return nil, err
		}
	}

	return lp, nil
}

func (lp *LyneParser) Parse(inputStream Stack.Stacker[*ll.Token]) (*Node, error) {
	lp.inputStream = inputStream
	lp.states = Stack.NewLinkedStack(NewInitialState())
	lp.trees = Stack.NewLinkedStack[*Node]()

	// var lookahead wppRawTokener
	var err error

	// While there are still input symbols
	for !lp.inputStream.IsEmpty() {
		// Peek at the current state.
		currentState, err := lp.states.Peek()
		if err != nil {
			return nil, errors.New("no state to get")
		}

		// Peek at the next symbol in the queue.
		next, _ := lp.inputStream.Peek()

		// Based on the current state and the next symbol, look up the action in the parsing table.
		decision := lp.decision(next, currentState)
		switch decision {
		case ActReduce:
			err := lp.Reduce()
			if err != nil {
				return nil, err
			}
		case ActShift:
			err := lp.Shift()
			if err != nil {
				return nil, err
			}
		case ActAccept:
		default:
			return nil, errors.New("could not decide whether to reduce or shift")
		}
	}

	// - If the action is to shift:
	// 	- Push the next symbol and the new state onto the stack.
	// 	- Remove the next symbol from the queue.
	// - Else if the action is to reduce:
	// 	- Pop the appropriate number of states off the stack.
	// 	- Push the left-hand side of the reduced rule and the new state onto the stack.
	// - Else if the action is to accept:
	// 	- Stop the parsing process.
	// - Else if the action is an error:
	// 	- Report the error and stop the parsing process.
	// If the parsing process was stopped due to an error, report the error.
	// Else, return the parse tree.

	if size := lp.trees.Size(); size == 0 {
		return nil, errors.New("no trees to get")
	} else if size > 1 {
		return nil, fmt.Errorf("expected 1 tree to get. Got %d instead", size)
	}

	top, err := lp.trees.Pop()
	if err != nil {
		return nil, errors.New("no trees to get")
	}

	return top, nil
}

func (lp *LyneParser) SetLexer(lexer *ll.Lexer) error {
	return nil
}

func (lp *LyneParser) Shift() error {
	token, err := lp.inputStream.Pop()
	if err != nil {
		return errors.New("no more tokens to shift")
	}

	lp.trees.Push(NewNode(token.GetID(), token.GetData()))

	return nil
}

func (lp *LyneParser) Reduce() error {
	top, err := lp.trees.Pop()
	if err != nil {
		return errors.New("no more trees to reduce")
	}

	G := lp.grammar.GetProductions()

	rules := make([]*Item, 0)
	for _, production := range G {
		rhs, _ := production.GetRhsAt(0)
		if rhs == top.GetID() {
			rule, err := NewItem(production)
			if err != nil {
				return err
			}

			rules = append(rules, rule)
		}
	}

	if len(rules) == 0 {
		return fmt.Errorf("no rules found for %v", top.GetID())
	}

	if len(rules) == 1 {
		ev := newRuleEvaluator(lp)

		for {
			isDone, err := ev.EvaluateRuleOnce(rules[0], top)
			if err != nil {
				return err
			}

			if isDone {
				break
			}
		}

		node := NewNode(rules[0].GetLHS(), top.GetData())
		node.AddChildren(ev.childrenToAdd...)
		lp.trees.Push(node)

		return nil
	}

	for _, rule := range rules {
		ev := newRuleEvaluator(lp.Copy().(*LyneParser))

		for {
			isDone, err := ev.EvaluateRuleOnce(rule, top)
			if err != nil {
				break
			}

			if isDone {
				node := NewNode(rules[0].GetLHS(), top.GetData())
				node.AddChildren(ev.childrenToAdd...)
				lp.trees.Push(node)

				return nil
			}
		}
	}

	return fmt.Errorf("no rules found for %v", top.GetID())
}

type ruleEvaluator struct {
	childrenToAdd []*Node
	p             *LyneParser
}

func newRuleEvaluator(lp *LyneParser) *ruleEvaluator {
	return &ruleEvaluator{
		childrenToAdd: make([]*Node, 0),
		p:             lp,
	}
}

func (ev *ruleEvaluator) EvaluateRuleOnce(rule *Item, prev *Node) (bool, error) {
	next, err := rule.PeekNext()
	if err != nil {
		return true, nil
	}

	top, err := ev.p.trees.Pop()
	if err != nil {
		return false, NewErrUnexpectedToken(next, prev.GetID(), nil)
	}

	if top.GetID() != next {
		return false, NewErrUnexpectedToken(next, prev.GetID(), top)
	}

	ev.childrenToAdd = append(ev.childrenToAdd, top)

	return false, nil
}

func (lp *LyneParser) decision(nextInput *ll.Token, lookahead any) ParserDecision {
	if lookahead != nil {
		switch x := lookahead.(type) {
		case *Node:
			switch x.GetID() {
			case "tkField1":
				if nextInput != nil && nextInput.GetID() == "tkClParen" {
					return ActReduce
				} else {
					return ActShift
				}
			case "tkField":
				return ActReduce
			}
		case *ll.Token:
			switch x.GetID() {
			case "tkEOF", "tkClCurly", "tkClParen":
				return ActReduce
			case "tkOpSquare", "tkClSquare", "tkOpCurly", "tkOpParen", "tkPlus", "tkKey", "tkNewline":
				return ActShift
			case "tkValue":
				if nextInput != nil && nextInput.GetID() == "tkPlus" {
					return ActShift
				} else {
					return ActReduce
				}
			}
		default:
			return ActShift
		}
	} else {
		return ActShift
	}

	// EOF	 		CL_SQUARE	obj			OP_SQUARE 	-> arrayObj
	// CL_CURLY	 	OP_CURLY 									-> obj
	// CL_CURLY 	obj'	 		OP_CURLY 					-> obj
	// field 		NEWLINE 		obj' 							-> obj'
	// CL_PAREN 	OP_PAREN 					KEY 			-> field
	// CL_PAREN 	field' 		OP_PAREN 	KEY 			-> field
	// field' 		PLUS	 		VALUE 						-> field'
	// VALUE 														-> field'

	return ActError
}
