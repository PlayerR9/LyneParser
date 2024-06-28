package Stream

import "fmt"

type Tokener interface {
	GetID() string
	GetLookahead() *LeafToken

	fmt.GoStringer
}

func CheckID(tok Tokener, id string) bool {
	if tok == nil {
		return false
	}

	switch tok := tok.(type) {
	case *LeafToken:
		return tok.ID == id
	case *NonLeafToken:
		return tok.ID == id
	default:
		return false
	}
}

type LeafToken struct {
	ID   string
	Data string
	At   int

	lookahead *LeafToken
}

func (lt *LeafToken) GetID() string {
	return lt.ID
}

func (lt *LeafToken) GetLookahead() *LeafToken {
	return lt.lookahead
}

func (lt *LeafToken) GoString() string {
	return fmt.Sprintf("%+v", lt)
}

func NewLeafToken(id string, data string, at int) *LeafToken {
	lt := &LeafToken{
		ID:   id,
		Data: data,
		At:   at,
	}
	return lt
}

type NonLeafToken struct {
	ID   string
	Data []Tokener
	Pos  int

	lookahead *LeafToken
}

func (lt *NonLeafToken) GetID() string {
	return lt.ID
}

func (lt *NonLeafToken) GoString() string {
	return fmt.Sprintf("%+v", lt)
}

func (lt *NonLeafToken) GetLookahead() *LeafToken {
	return lt.lookahead
}

func NewNonLeafToken(id string, data []Tokener, pos int, lookahead *LeafToken) *NonLeafToken {
	lt := &NonLeafToken{
		ID:        id,
		Data:      data,
		Pos:       pos,
		lookahead: lookahead,
	}
	return lt
}

func SetLookaheads(tokens []*LeafToken) {
	for i := 0; i < len(tokens)-1; i++ {
		tokens[i].lookahead = tokens[i+1]
	}

	tokens[len(tokens)-1].lookahead = nil
}
