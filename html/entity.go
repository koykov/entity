package html

import (
	"github.com/koykov/entry"
	"github.com/koykov/fastconv"
)

type Entity struct {
	name entry.Entry32
	val  entry.Entry32
	cp   int64
}

func (e Entity) Name() string {
	lo, hi := e.name.Decode()
	return fastconv.B2S(__buf[lo:hi])
}

func (e Entity) Value() string {
	lo, hi := e.val.Decode()
	return fastconv.B2S(__buf[lo:hi])
}

func (e Entity) Codepoint() int64 {
	return e.cp
}
