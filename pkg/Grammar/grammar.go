package Grammar

import (
	"fmt"
	"slices"
	"strings"
)

type GrammarBuilder struct {
	productions []*Production
	direction   ProductionDirection
}

func (b *GrammarBuilder) String() string {
	if b.productions == nil {
		return "GrammarBuilder[nil]"
	}

	if len(b.productions) == 0 {
		return "GrammarBuilder[total=0, productions=[]]"
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "GrammarBuilder[total=%d, productions=[%v", len(b.productions), b.productions[0])

	for _, production := range b.productions[1:] {
		fmt.Fprintf(&builder, ", %v", production)
	}

	builder.WriteString("]]")

	return builder.String()
}

func (b *GrammarBuilder) AddProduction(p *Production) {
	if p == nil {
		return
	}

	if b.productions == nil {
		b.productions = []*Production{p}
	} else {
		b.productions = append(b.productions, p)
	}
}

func (b *GrammarBuilder) SetDirection(direction ProductionDirection) {
	b.direction = direction
}

func (b *GrammarBuilder) Build() *Grammar {
	if b.productions == nil {
		b.direction = LeftToRight

		return &Grammar{
			productions: make([]*Production, 0),
			symbols:     make([]string, 0),
		}
	}

	grammar := Grammar{
		symbols: make([]string, 0),
	}

	// 1. Remove duplicates
	for i := 0; i < len(b.productions); {
		index := slices.IndexFunc(b.productions[i+1:], func(p *Production) bool {
			return p.IsEqual(b.productions[i])
		})

		if index != -1 {
			b.productions = append(b.productions[:index], b.productions[index+1:]...)
		} else {
			i++
		}
	}

	// 2. Make sure all productions follow the same direction
	for _, p := range b.productions {
		p.SetDirection(b.direction)
	}

	// 3. Add productions to grammar
	grammar.productions = make([]*Production, len(b.productions))
	copy(grammar.productions, b.productions)

	// 3. Add symbols to grammar
	for _, p := range b.productions {
		for _, symbol := range p.GetSymbols() {
			if !slices.Contains(grammar.symbols, symbol) {
				grammar.symbols = append(grammar.symbols, symbol)
			}
		}
	}

	// 4. Clear builder
	for i := range b.productions {
		b.productions[i] = nil
	}

	b.productions = nil
	b.direction = LeftToRight

	return &grammar
}

func (b *GrammarBuilder) Clear() {
	for i := range b.productions {
		b.productions[i] = nil
	}

	b.productions = nil

	b.direction = LeftToRight
}

type Grammar struct {
	productions []*Production
	symbols     []string
}

func (g *Grammar) String() string {
	if g == nil {
		return "Grammar[nil]"
	}

	if len(g.productions) == 0 {
		return "Grammar[prouctions=[], symbols=[]]"
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "Grammar[productions=[%v", g.productions[0])

	for _, production := range g.productions[1:] {
		fmt.Fprintf(&builder, ", %v", production)
	}

	fmt.Fprintf(&builder, "], symbols=[%v", g.symbols[0])

	for _, symbol := range g.symbols[1:] {
		fmt.Fprintf(&builder, ", %v", symbol)
	}

	builder.WriteString("]]")

	return builder.String()
}

func (g *Grammar) GetProductions() []*Production {
	productions := make([]*Production, len(g.productions))
	copy(productions, g.productions)

	return productions
}

func (g *Grammar) GetSymbols() []string {
	symbols := make([]string, len(g.symbols))
	copy(symbols, g.symbols)

	return symbols
}
