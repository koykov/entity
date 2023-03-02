package html

import (
	"bytes"

	"github.com/koykov/byteseq"
	"github.com/koykov/fastconv"
)

func Unescape[T byteseq.Byteseq](x T) T {
	r := AppendUnescape(nil, x)
	return T(r)
}

func AppendUnescape[T byteseq.Byteseq](dst []byte, x T) []byte {
	p := byteseq.Q2B(x)
	l := len(p)
	if l == 0 {
		return dst
	}

	if i := bytes.IndexByte(p, '&'); i == -1 {
		dst = append(dst, p...)
		return dst
	}

	_ = p[l-1]
	lo, hi := 0, 0
	var tag bool
	for i := 0; i < l; i++ {
		switch {
		case p[i] == '&':
			tag = true
			lo = i
			hi = lo
		case tag && p[i] == ';':
			tag = false
			hi = i + 1
			t := fastconv.B2S(p[lo:hi])
			if i1, ok := __bufHN[t]; ok {
				h := __bufH[i1]
				dst = append(dst, h.Value()...)
			} else {
				dst = append(dst, t...)
			}
		default:
			if !tag {
				dst = append(dst, p[i])
			}
		}
	}

	return dst
}
