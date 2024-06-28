package Parser

type Actioner interface {
}

type ShiftAct struct{}

func NewShiftAct() *ShiftAct {
	sa := &ShiftAct{}
	return sa
}

type ReduceAct struct {
	isAccept bool
}

func NewReduceAct(isAccept bool) *ReduceAct {
	ra := &ReduceAct{
		isAccept: isAccept,
	}
	return ra
}
