package ConflictSolver

import "strings"

type Helper struct {
	Item   *Item
	Action Actioner
}

func (h *Helper) String() string {
	if h == nil {
		return ""
	}

	var builder strings.Builder

	builder.WriteString(h.Item.String())
	builder.WriteRune(' ')
	builder.WriteRune('(')

	if h.Action != nil {
		builder.WriteString(h.Action.String())
	}

	builder.WriteRune(')')

	return builder.String()
}

func NewHelper(item *Item, action Actioner) *Helper {
	return &Helper{
		Item:   item,
		Action: action,
	}
}
