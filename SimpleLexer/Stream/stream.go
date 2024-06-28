package Stream

import (
	"fmt"
	"unicode/utf8"
)

type lookahead struct {
	char rune
	size int
}

type InputStream struct {
	data []byte
	la   *lookahead
	pos  int
}

func NewInputStream(data []byte) *InputStream {
	is := &InputStream{
		data: data,
		pos:  0,
		la:   nil,
	}
	return is
}

func (is *InputStream) Next() (rune, error) {
	if is.la != nil {
		r := is.la.char
		is.pos += is.la.size
		is.la = nil

		return r, nil
	}

	if len(is.data) == 0 {
		return 0, NewErrStreamExhausted()
	}

	r, size := utf8.DecodeRune(is.data)
	if r == utf8.RuneError {
		return 0, fmt.Errorf("invalid utf8 character")
	}

	is.data = is.data[size:]
	is.pos += size

	return r, nil
}

func (is *InputStream) Pos() int {
	if is.la == nil {
		return is.pos
	} else {
		return is.pos - is.la.size
	}
}

func (is *InputStream) Peek() (rune, error) {
	if is.la != nil {
		return is.la.char, nil
	}

	if len(is.data) == 0 {
		return 0, NewErrStreamExhausted()
	}

	r, size := utf8.DecodeRune(is.data)
	if r == utf8.RuneError {
		return 0, fmt.Errorf("invalid utf8 character")
	}

	is.data = is.data[size:]

	is.la = &lookahead{
		char: r,
		size: size,
	}

	return r, nil
}

func (is *InputStream) Accept() {
	is.la = nil
}
