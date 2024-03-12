package LyneParser

import (
	"errors"
	"fmt"
	"strings"

	gr "LyneParser/pkg/Grammar"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

type Item struct {
	// The production that is being used
	production *gr.Production

	// The position of the last rhs symbol that was matched.
	// len(production.RHS) == position means that the production is fully matched.
	position int
}

func (i *Item) String() string {
	if i == nil {
		return "Item[nil]"
	}

	// LHS -> RHS # RHS

	var builder strings.Builder

	if i.production.IsLeftToRight() {
		fmt.Fprintf(&builder, "Item[%v ->", i.production.GetLHS())
	} else {
		fmt.Fprintf(&builder, "Item[%v <-", i.production.GetLHS())
	}

	iter := i.production.Iterator()

	if i.position == -1 {
		// Print everything before the position
		for j := 0; iter.Next(); j++ {
			val, _ := iter.Value()

			fmt.Fprintf(&builder, " %s", val)

			if j == i.position {
				break
			}
		}
	}

	builder.WriteString(" #")

	// Print everything after the position
	for j := i.position + 1; iter.Next(); j++ {
		val, _ := iter.Value()

		fmt.Fprintf(&builder, " %s", val)
	}

	builder.WriteString("]")

	return builder.String()
}

func NewItem(p *gr.Production) (*Item, error) {
	if p == nil {
		return nil, ers.NewErrInvalidParameter("p").
			Wrap(errors.New("production is nil"))
	}

	i := &Item{
		production: p,
		position:   -1,
	}

	return i, nil
}

func (i *Item) Next() bool {
	if i.production.Size() == i.position {
		return false
	}

	i.position++

	return true
}

func (i *Item) PeekNext() (string, error) {
	if i.production.Size() == i.position+1 {
		return "", errors.New("position is at the end")
	}

	return i.production.GetRhsAt(i.position + 1)
}

func (i *Item) Value() (string, error) {
	if i.position == -1 {
		return "", errors.New("item is not started yet")
	} else if i.position == i.production.Size() {
		return "", errors.New("item is finished")
	}

	return i.production.GetRhsAt(i.position)
}

func (i *Item) Restart() {
	i.position = -1
}

func (i *Item) GetLHS() string {
	return i.production.GetLHS()
}

func NewLockedItem(p *gr.Production, position int) (*Item, error) {
	if p == nil {
		return nil, ers.NewErrInvalidParameter("p").
			Wrap(errors.New("production is nil"))
	}

	if position < -1 || position > p.Size() {
		return nil, ers.NewErrInvalidParameter("position").
			Wrap(errors.New("position is out of range"))
	}

	i := &Item{
		production: p,
		position:   position,
	}

	return i, nil
}

type ActionType int

const (
	Shift ActionType = iota
	Reduce
	Accept
	Error
)

type Action struct {
	actionType ActionType
	production *gr.Production
}

func NewAction(actionType ActionType, production *gr.Production) *Action {
	return &Action{actionType: actionType, production: production}
}

type ActionTable struct {
	actions map[*Item]map[string]*Action
}

func NewActionTable(g *gr.Grammar) *ActionTable {
	table := ActionTable{
		actions: make(map[*Item]map[string]*Action),
	}

	for _, production := range g.GetProductions() {
		for i := 0; i <= production.Size(); i++ {
			item, _ := NewLockedItem(production, i)

			table.actions[item] = make(map[string]*Action)

			for _, symbol := range g.GetSymbols() {
				table.actions[item][symbol] = nil
			}
		}
	}

	return &table
}

func (t *ActionTable) SetAction(item *Item, symbol string, action *Action) {
	if item == nil || action == nil {
		return
	}

	t.actions[item][symbol] = action
}

func (t *ActionTable) GetAction(item *Item, symbol string) *Action {
	if item == nil {
		return nil
	}

	return t.actions[item][symbol]
}

func (t *ActionTable) Fill(table *GotoTable) {
	for item, actions := range t.actions {
		for symbol, action := range actions {
			if action != nil {
				continue
			}

			nextItem := table.GetGoto(item, symbol)

			if nextItem != nil {
				if nextItem.position == nextItem.production.Size() {
					action = NewAction(Reduce, nextItem.production)
				} else {
					action = NewAction(Shift, nil)
				}

				t.SetAction(item, symbol, action)
			}
		}
	}
}

type GotoTable struct {
	gotos map[*Item]map[string]*Item
}

func NewGotoTable(g *gr.Grammar) *GotoTable {
	table := GotoTable{
		gotos: make(map[*Item]map[string]*Item),
	}

	for _, production := range g.GetProductions() {
		for i := 0; i <= production.Size(); i++ {
			item, _ := NewLockedItem(production, i)

			table.gotos[item] = make(map[string]*Item)

			for _, symbol := range g.GetSymbols() {
				table.gotos[item][symbol] = nil
			}
		}
	}

	return &table
}

func (t *GotoTable) SetGoto(item *Item, symbol string, nextItem *Item) {
	if item == nil || nextItem == nil {
		return
	}

	t.gotos[item][symbol] = nextItem
}

func (t *GotoTable) GetGoto(item *Item, symbol string) *Item {
	if item == nil {
		return nil
	}

	return t.gotos[item][symbol]
}

func (t *GotoTable) Fill(table *ActionTable) {
	for item, gotos := range t.gotos {
		for symbol, nextItem := range gotos {
			if nextItem != nil {
				continue
			}

			nextItem, _ = NewLockedItem(item.production, item.position+1)

			if item.position < item.production.Size() {
				nextItem, _ = NewLockedItem(item.production, item.position+1)
			}

			t.SetGoto(item, symbol, nextItem)
		}
	}
}
