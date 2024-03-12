package Grammar

import (
	"fmt"
	"slices"
	"strings"

	itf "github.com/PlayerR9/MyGoLib/Interfaces"
)

const (
	RightArrowSymbol string = "->"
	LeftArrowSymbol  string = "<-"
)

type ProductionDirection bool

const (
	LeftToRight ProductionDirection = false
	RightToLeft ProductionDirection = true
)

type Production struct {
	lhs       string
	rhs       []string
	direction ProductionDirection
}

func (p *Production) String() string {
	if p == nil {
		return "Production[nil]"
	}

	if len(p.rhs) == 0 {
		return fmt.Sprintf("Production[%s %s]", p.lhs, RightArrowSymbol)
	}

	if p.direction == LeftToRight {
		return fmt.Sprintf("Production[%s %s %s]", p.lhs, RightArrowSymbol, strings.Join(p.rhs, " "))
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "Production[%s %s %s", p.lhs, LeftArrowSymbol, p.rhs[len(p.rhs)-1])

	for i := len(p.rhs) - 2; i >= 0; i-- {
		fmt.Fprintf(&builder, " %s", p.rhs[i])
	}

	builder.WriteString("]")

	return builder.String()
}

func (p *Production) Iterator() itf.Iterater[string] {
	if p.direction == LeftToRight {
		return itf.IteratorFromSlice(p.rhs)
	}

	var builder itf.Builder[string]

	for i := len(p.rhs) - 1; i >= 0; i-- {
		builder.Append(p.rhs[i])
	}

	return builder.Build()
}

func NewProduction(lhs string, rhs ...string) *Production {
	return &Production{lhs: lhs, rhs: rhs}
}

func (p *Production) IsEqual(other *Production) bool {
	if other == nil {
		return false
	} else if p.lhs != other.lhs {
		return false
	} else if len(p.rhs) != len(other.rhs) {
		return false
	}

	for i, symbol := range p.rhs {
		if symbol != other.rhs[i] {
			return false
		}
	}

	return true
}

func (p *Production) GetSymbols() []string {
	symbols := make([]string, len(p.rhs)+1)
	copy(symbols, p.rhs)

	symbols[len(symbols)-1] = p.lhs

	// Remove duplicates
	for i := 0; i < len(symbols); {
		index := slices.Index(symbols[i+1:], symbols[i])

		if index != -1 {
			symbols = append(symbols[:index], symbols[index+1:]...)
		} else {
			i++
		}
	}

	return symbols
}

func (p *Production) SetDirection(direction ProductionDirection) {
	p.direction = direction
}

func (p *Production) Size() int {
	return len(p.rhs)
}

func (p *Production) GetRhsAt(index int) (string, error) {
	if index < 0 || index >= len(p.rhs) {
		return "", fmt.Errorf("index %d out of range", index)
	}

	return p.rhs[index], nil
}

func (p *Production) GetLHS() string {
	return p.lhs
}

func (p *Production) IsLeftToRight() bool {
	return p.direction == LeftToRight
}
