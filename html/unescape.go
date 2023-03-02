package html

import (
	"bytes"

	"github.com/koykov/byteseq"
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
	for i := 0; i < l; i++ {
		if p[i] == '&' {
			//
		} else {
			dst = append(dst, p[i])
		}
	}

	return dst
}
