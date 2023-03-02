package html

import (
	"bytes"
	"unicode"

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
			dst = unesc(dst, p[lo:hi])
		case tag && !unicode.IsLetter(rune(p[i])):
			tag = false
			hi = i
			dst = unesc(dst, p[lo:hi])
			dst = append(dst, p[i])
		default:
			if !tag {
				dst = append(dst, p[i])
			}
		}
	}
	if tag {
		dst = unesc(dst, p[lo:l])
	}

	return dst
}

func unesc(dst, ent []byte) []byte {
	s := fastconv.B2S(ent)
	if i1, ok := __bufHN[s]; ok {
		h := __bufH[i1]
		dst = append(dst, h.Value()...)
	} else {
		dst = append(dst, s...)
	}
	return dst
}
