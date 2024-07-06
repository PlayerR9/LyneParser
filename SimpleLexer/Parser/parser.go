package Parser

import (
	"fmt"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	com "github.com/PlayerR9/LyneParser/SimpleLexer/Common"
	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

type Parser struct {
	input []*gr.Token[com.TestTkType]
	stack *lls.ArrayStack[*gr.Token[com.TestTkType]]
}

func NewParser() *Parser {
	p := &Parser{}
	return p
}

func (p *Parser) getDecision(top *gr.Token[com.TestTkType]) (Actioner, error) {
	var act Actioner

	id := top.GetID()

	switch id {
	case com.TkEof:
		act = NewReduceAct(true)
	case com.TkRegister, com.TkUnaryOp, com.TkBinOp, com.TkClParen:
		act = NewReduceAct(false)
	case com.TkNewline, com.TkRightArrow, com.TkUnaryInstruction,
		com.TkLoadImmediate, com.TkOperand, com.TkImmediate,
		com.TkOpParen:
		act = NewShiftAct()
	case com.TkSource1:
		la := top.GetLookahead()

		if la == nil {
			// [Source1] newline Statement -> Source1 : Reduce
			act = NewReduceAct(false)

			return act, nil
		}

		laID := la.GetID()

		if laID == com.TkEof {
			// EOF [Source1] -> Source : Shift
			act = NewShiftAct()
		} else {
			// [Source1] newline Statement -> Source1 : Reduce
			act = NewReduceAct(false)
		}
	case com.TkStatement:
		la := top.GetLookahead()

		if la == nil {
			// [Statement] -> Source1 : Reduce
			act = NewReduceAct(false)

			return act, nil
		}

		laID := la.GetID()

		if laID == com.TkNewline {
			// Source1 newline [Statement] -> Source1 : Shift
			act = NewShiftAct()
		} else {
			// [Statement] -> Source1 : Reduce
			act = NewReduceAct(false)
		}

	case com.TkBinInstruction:
		la := top.GetLookahead()

		if la == nil {
			// [BinaryInstruction] -> Operand : Reduce
			act = NewReduceAct(false)

			return act, nil
		}

		laID := la.GetID()

		if laID == com.TkRightArrow {
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
	case com.TkEof:
		// [EOF] Source1 -> Source : Reduce
		top2, ok := p.stack.Pop()
		if !ok {
			return fmt.Errorf("expected source1, got EOF instead")
		}

		id2 := top2.GetID()

		if id2 != com.TkSource1 {
			return fmt.Errorf("expected Source1, got %s instead", id2)
		}

		la := top2.GetLookahead()

		tok := gr.NewToken(com.TkSource, []*gr.Token[com.TestTkType]{top2, top1}, 0, la)

		p.stack.Push(tok)
	case com.TkSource1:
		// [Source1] newline Statement -> Source1 : Reduce
		top2, ok := p.stack.Pop()
		if !ok {
			return fmt.Errorf("expected newline, got EOF instead")
		}

		id2 := top2.GetID()

		if id2 != com.TkNewline {
			return fmt.Errorf("expected newline, got %s instead", id2)
		}

		top3, ok := p.stack.Pop()
		if !ok {
			return fmt.Errorf("expected Statement, got nothing instead")
		}

		id3 := top3.GetID()

		if id3 != com.TkStatement {
			return fmt.Errorf("expected Statement, got %s instead", id3)
		}

		la := top3.GetLookahead()

		tok := gr.NewToken(com.TkSource1, []*gr.Token[com.TestTkType]{top3, top1}, 0, la)

		p.stack.Push(tok)
	case com.TkStatement:
		// [Statement] -> Source1 : Reduce
		la := top1.GetLookahead()

		tok := gr.NewToken(com.TkSource1, []*gr.Token[com.TestTkType]{top1}, 0, la)

		p.stack.Push(tok)
	case com.TkRegister:
		top2, ok := p.stack.Pop()
		if !ok {
			// [register] -> Operand : Reduce
			la := top1.GetLookahead()

			tok := gr.NewToken(com.TkOperand, []*gr.Token[com.TestTkType]{top1}, 0, la)

			p.stack.Push(tok)

			return nil
		}

		id2 := top2.GetID()

		if id2 != com.TkRightArrow {
			p.stack.Push(top2)

			// [register] -> Operand : Reduce
			la := top1.GetLookahead()

			tok := gr.NewToken(com.TkOperand, []*gr.Token[com.TestTkType]{top1}, 0, la)

			p.stack.Push(tok)

			return nil
		}

		top3, ok := p.stack.Pop()
		if !ok {
			return fmt.Errorf("expected UnaryInstruction, BinaryInstruction, or LoadImmediate, got nothing instead")
		}

		id3 := top3.GetID()

		switch id3 {
		case com.TkUnaryInstruction:
			// [register] right_arrow UnaryInstruction -> Statement : Reduce

			la := top3.GetLookahead()

			tok := gr.NewToken(com.TkStatement, []*gr.Token[com.TestTkType]{top3, top1}, 0, la)

			p.stack.Push(tok)
		case com.TkBinInstruction:
			// [register] right_arrow BinaryInstruction -> Statement : Reduce

			la := top3.GetLookahead()

			tok := gr.NewToken(com.TkStatement, []*gr.Token[com.TestTkType]{top3, top1}, 0, la)

			p.stack.Push(tok)
		case com.TkLoadImmediate:
			// [register] right_arrow LoadImmediate -> Statement : Reduce

			la := top3.GetLookahead()

			tok := gr.NewToken(com.TkStatement, []*gr.Token[com.TestTkType]{top3, top1}, 0, la)

			p.stack.Push(tok)
		default:
			return fmt.Errorf("expected UnaryInstruction, BinaryInstruction, or LoadImmediate, got %s instead", id3)
		}
	case com.TkBinInstruction:
		// [BinaryInstruction] -> Operand : Reduce

		la := top1.GetLookahead()

		tok := gr.NewToken(com.TkOperand, []*gr.Token[com.TestTkType]{top1}, 0, la)

		p.stack.Push(tok)
	case com.TkUnaryInstruction:
		// [unary_operator] Operand -> UnaryInstruction : Reduce

		top2, ok := p.stack.Pop()
		if !ok {
			return fmt.Errorf("expected operand, got nothing instead")
		}

		id2 := top2.GetID()

		if id2 != com.TkOperand {
			return fmt.Errorf("expected Operand, got %s instead", id2)
		}

		la := top2.GetLookahead()

		tok := gr.NewToken(com.TkUnaryInstruction, []*gr.Token[com.TestTkType]{top2, top1}, 0, la)

		p.stack.Push(tok)
	case com.TkBinOp:
		// [binary_operator] Operand Operand -> BinaryInstruction : Reduce

		top2, ok := p.stack.Pop()
		if !ok {
			return fmt.Errorf("expected operand, got nothing instead")
		}

		id2 := top2.GetID()

		if id2 != com.TkOperand {
			return fmt.Errorf("expected Operand, got %s instead", id2)
		}

		top3, ok := p.stack.Pop()
		if !ok {
			return fmt.Errorf("expected operand, got nothing instead")
		}

		id3 := top3.GetID()

		if id3 != com.TkOperand {
			return fmt.Errorf("expected Operand, got %s instead", id3)
		}

		la := top3.GetLookahead()

		tok := gr.NewToken(com.TkBinInstruction, []*gr.Token[com.TestTkType]{top3, top2, top1}, 0, la)

		p.stack.Push(tok)
	case com.TkClParen:
		// [cl_paren] immediate op_paren -> LoadImmediate : Reduce

		top2, ok := p.stack.Pop()
		if !ok {
			return fmt.Errorf("expected immediate, got nothing instead")
		}

		id2 := top2.GetID()

		if id2 != com.TkImmediate {
			return fmt.Errorf("expected immediate, got %s instead", id2)
		}

		top3, ok := p.stack.Pop()
		if !ok {
			return fmt.Errorf("expected op_paren, got nothing instead")
		}

		id3 := top3.GetID()

		if id3 != com.TkOpParen {
			return fmt.Errorf("expected op_paren, got %s instead", id3)
		}

		la := top3.GetLookahead()

		tok := gr.NewToken(com.TkLoadImmediate, []*gr.Token[com.TestTkType]{top2}, 0, la)

		p.stack.Push(tok)
	default:
		return fmt.Errorf("cannot reduce token %s", id.String())
	}

	return nil
}

func (p *Parser) Parse(tokens []*gr.Token[com.TestTkType]) (*gr.Token[com.TestTkType], error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no tokens to parse")
	}

	p.input = tokens

	p.stack = lls.NewArrayStack[*gr.Token[com.TestTkType]]()

	p.shift()

	err := p.parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}

	top, ok := p.stack.Pop()
	if !ok {
		return nil, fmt.Errorf("no result found")
	}

	ok = top.IsNonLeaf()
	if !ok {
		return nil, fmt.Errorf("expected non-leaf token, got %T instead", top)
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
