package ConflictSolver

type Helper struct {
	Item   *Item
	Action Actioner
}

func NewHelper(item *Item, action Actioner) *Helper {
	return &Helper{
		Item:   item,
		Action: action,
	}
}
