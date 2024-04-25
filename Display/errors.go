package Display

type ErrESCPressed struct{}

func (e *ErrESCPressed) Error() string {
	return "ESC key pressed"
}

func NewErrESCPressed() *ErrESCPressed {
	return &ErrESCPressed{}
}
