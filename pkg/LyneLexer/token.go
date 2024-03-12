package LyneLexer

import (
	"fmt"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

type Token struct {
	id           string
	data         []byte
	matchingRule int
}

func (t *Token) String() string {
	if t == nil {
		return "Token[nil]"
	}

	return fmt.Sprintf("Token[id=%s, Value=%s]", t.id, t.data)
}

func NewToken(id string, data []byte, matchingRule int) (*Token, error) {
	if !IsLexToken(id) {
		return nil, ers.NewErrInvalidParameter("id").
			Wrap(NewErrInvalidTokenName(id))
	}

	return &Token{
		id:           id,
		data:         data,
		matchingRule: matchingRule,
	}, nil
}

func (t *Token) GetID() string {
	return t.id
}

func (t *Token) GetData() string {
	return string(t.data)
}

func (t *Token) MatchingRuleIndex() int {
	return t.matchingRule
}
