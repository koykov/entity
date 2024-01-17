package html

import (
	"github.com/koykov/byteconv"
	"github.com/koykov/entry"
)

type Entity struct {
	name entry.Entry32
	val  entry.Entry32
	cp   int64
}

func (e Entity) Name() string {
	lo, hi := e.name.Decode()
	return byteconv.B2S(__buf[lo:hi])
}

func (e Entity) Value() string {
	lo, hi := e.val.Decode()
	return byteconv.B2S(__buf[lo:hi])
}

func (e Entity) Codepoint() int64 {
	return e.cp
}
