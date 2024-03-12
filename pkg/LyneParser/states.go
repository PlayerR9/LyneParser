package LyneParser

type ParserState int

const (
	InitialState ParserState = iota
)

type State struct {
	code ParserState
}

func NewInitialState() *State {
	return &State{code: InitialState}
}
