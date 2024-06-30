package Parser

import (
	"fmt"

	uc "github.com/PlayerR9/MyGoLib/Units/common"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
)

type Parser struct {
	input []gr.Token
	stack *lls.ArrayStack[gr.Token]
}

func NewParser() *Parser {
	p := &Parser{}
	return p
}

func (p *Parser) getDecision(top gr.Token) (Actioner, error) {
	var act Actioner

	id := top.GetID()

	switch id {
	case "EOF":
		act = NewReduceAct(true)
	case "register", "unary_operator", "binary_operator", "cl_paren":
		act = NewReduceAct(false)
	case "newline", "right_arrow", "UnaryInstruction", "LoadImmediate", "Operand", "immediate", "op_paren":
		act = NewShiftAct()
	case "Source1":
		la := top.GetLookahead()

		if la == nil {
			// [Source1] newline Statement -> Source1 : Reduce
			act = NewReduceAct(false)

			return act, nil
		}

		laID := la.GetID()

		if laID == "EOF" {
			// EOF [Source1] -> Source : Shift
			act = NewShiftAct()
		} else {
			// [Source1] newline Statement -> Source1 : Reduce
			act = NewReduceAct(false)
		}
	case "Statement":
		la := top.GetLookahead()

		if la == nil {
			// [Statement] -> Source1 : Reduce
			act = NewReduceAct(false)

			return act, nil
		}

		laID := la.GetID()

		if laID == "newline" {
			// Source1 newline [Statement] -> Source1 : Shift
			act = NewShiftAct()
		} else {
			// [Statement] -> Source1 : Reduce
			act = NewReduceAct(false)
		}

	case "BinaryInstruction":
		la := top.GetLookahead()

		if la == nil {
			// [BinaryInstruction] -> Operand : Reduce
			act = NewReduceAct(false)

			return act, nil
		}

		laID := la.GetID()

		if laID == "right_arrow" {
			// register right_arrow [BinaryInstruction] -> Statement : Shift
			act = NewShiftAct()
		} else {
			// [BinaryInstruction] -> Operand : Reduce
			act = NewReduceAct(false)
		}
	default:
		return nil, fmt.Errorf("unexpected token %s", id)
	}

	return act, nil
}

func (p *Parser) shift() {
	uc.Assert(len(p.input) > 0, "input is empty")

	first := p.input[0]
	p.input = p.input[1:]

	p.stack.Push(first)
}

func (p *Parser) reduce() error {
	top1, ok := p.stack.Pop()
	uc.Assert(ok, "stack is empty")

	id := top1.GetID()

	switch id {
	case "EOF":
		// [EOF] Source1 -> Source : Reduce
		top2, ok := p.stack.Pop()
		if !ok {
			return fmt.Errorf("expected source1, got EOF instead")
		}

		id2 := top2.GetID()

		if id2 != "Source1" {
			return fmt.Errorf("expected Source1, got %s instead", id2)
		}

		la := top2.GetLookahead()

		tok := gr.NewToken("Source", []gr.Token{top2, top1}, 0, la)

		p.stack.Push(tok)
	case "Source1":
		// [Source1] newline Statement -> Source1 : Reduce
		top2, ok := p.stack.Pop()
		if !ok {
			return fmt.Errorf("expected newline, got EOF instead")
		}

		id2 := top2.GetID()

		if id2 != "newline" {
			return fmt.Errorf("expected newline, got %s instead", id2)
		}

		top3, ok := p.stack.Pop()
		if !ok {
			return fmt.Errorf("expected Statement, got nothing instead")
		}

		id3 := top3.GetID()

		if id3 != "Statement" {
			return fmt.Errorf("expected Statement, got %s instead", id3)
		}

		la := top3.GetLookahead()

		tok := gr.NewToken("Source1", []gr.Token{top3, top1}, 0, la)

		p.stack.Push(tok)
	case "Statement":
		// [Statement] -> Source1 : Reduce
		la := top1.GetLookahead()

		tok := gr.NewToken("Source1", []gr.Token{top1}, 0, la)

		p.stack.Push(tok)
	case "register":
		top2, ok := p.stack.Pop()
		if !ok {
			// [register] -> Operand : Reduce
			la := top1.GetLookahead()

			tok := gr.NewToken("Operand", []gr.Token{top1}, 0, la)

			p.stack.Push(tok)

			return nil
		}

		id2 := top2.GetID()

		if id2 != "right_arrow" {
			p.stack.Push(top2)

			// [register] -> Operand : Reduce
			la := top1.GetLookahead()

			tok := gr.NewToken("Operand", []gr.Token{top1}, 0, la)

			p.stack.Push(tok)

			return nil
		}

		top3, ok := p.stack.Pop()
		if !ok {
			return fmt.Errorf("expected UnaryInstruction, BinaryInstruction, or LoadImmediate, got nothing instead")
		}

		id3 := top3.GetID()

		switch id3 {
		case "UnaryInstruction":
			// [register] right_arrow UnaryInstruction -> Statement : Reduce

			la := top3.GetLookahead()

			tok := gr.NewToken("Statement", []gr.Token{top3, top1}, 0, la)

			p.stack.Push(tok)
		case "BinaryInstruction":
			// [register] right_arrow BinaryInstruction -> Statement : Reduce

			la := top3.GetLookahead()

			tok := gr.NewToken("Statement", []gr.Token{top3, top1}, 0, la)

			p.stack.Push(tok)
		case "LoadImmediate":
			// [register] right_arrow LoadImmediate -> Statement : Reduce

			la := top3.GetLookahead()

			tok := gr.NewToken("Statement", []gr.Token{top3, top1}, 0, la)

			p.stack.Push(tok)
		default:
			return fmt.Errorf("expected UnaryInstruction, BinaryInstruction, or LoadImmediate, got %s instead", id3)
		}
	case "BinaryInstruction":
		// [BinaryInstruction] -> Operand : Reduce

		la := top1.GetLookahead()

		tok := gr.NewToken("Operand", []gr.Token{top1}, 0, la)

		p.stack.Push(tok)
	case "unary_operator":
		// [unary_operator] Operand -> UnaryInstruction : Reduce

		top2, ok := p.stack.Pop()
		if !ok {
			return fmt.Errorf("expected operand, got nothing instead")
		}

		id2 := top2.GetID()

		if id2 != "Operand" {
			return fmt.Errorf("expected Operand, got %s instead", id2)
		}

		la := top2.GetLookahead()

		tok := gr.NewToken("UnaryInstruction", []gr.Token{top2, top1}, 0, la)

		p.stack.Push(tok)
	case "binary_operator":
		// [binary_operator] Operand Operand -> BinaryInstruction : Reduce

		top2, ok := p.stack.Pop()
		if !ok {
			return fmt.Errorf("expected operand, got nothing instead")
		}

		id2 := top2.GetID()

		if id2 != "Operand" {
			return fmt.Errorf("expected Operand, got %s instead", id2)
		}

		top3, ok := p.stack.Pop()
		if !ok {
			return fmt.Errorf("expected operand, got nothing instead")
		}

		id3 := top3.GetID()

		if id3 != "Operand" {
			return fmt.Errorf("expected Operand, got %s instead", id3)
		}

		la := top3.GetLookahead()

		tok := gr.NewToken("BinaryInstruction", []gr.Token{top3, top2, top1}, 0, la)

		p.stack.Push(tok)
	case "cl_paren":
		// [cl_paren] immediate op_paren -> LoadImmediate : Reduce

		top2, ok := p.stack.Pop()
		if !ok {
			return fmt.Errorf("expected immediate, got nothing instead")
		}

		id2 := top2.GetID()

		if id2 != "immediate" {
			return fmt.Errorf("expected immediate, got %s instead", id2)
		}

		top3, ok := p.stack.Pop()
		if !ok {
			return fmt.Errorf("expected op_paren, got nothing instead")
		}

		id3 := top3.GetID()

		if id3 != "op_paren" {
			return fmt.Errorf("expected op_paren, got %s instead", id3)
		}

		la := top3.GetLookahead()

		tok := gr.NewToken("LoadImmediate", []gr.Token{top2}, 0, la)

		p.stack.Push(tok)
	default:
		return fmt.Errorf("cannot reduce token %s", id)
	}

	return nil
}

func (p *Parser) Parse(tokens []gr.Token) (gr.Token, error) {
	if len(tokens) == 0 {
		return gr.Token{}, fmt.Errorf("no tokens to parse")
	}

	p.input = tokens

	p.stack = lls.NewArrayStack[gr.Token]()

	p.shift()

	err := p.parse()
	if err != nil {
		return gr.Token{}, fmt.Errorf("failed to parse: %w", err)
	}

	top, ok := p.stack.Pop()
	if !ok {
		return gr.Token{}, fmt.Errorf("no result found")
	}

	ok = top.IsNonLeaf()
	if !ok {
		return gr.Token{}, fmt.Errorf("expected non-leaf token, got %T instead", top)
	}

	return top, nil
}

func (p *Parser) parse() error {
	for {
		top, ok := p.stack.Peek()
		if !ok {
			break
		}

		act, err := p.getDecision(top)
		if err != nil {
			return fmt.Errorf("failed to get decision: %w", err)
		}

		switch act := act.(type) {
		case *ShiftAct:
			p.shift()
		case *ReduceAct:
			err := p.reduce()
			if err != nil {
				return fmt.Errorf("failed to reduce: %w", err)
			}

			if act.isAccept {
				return nil
			}
		}
	}

	return fmt.Errorf("done parsing but no accept state found")
}
