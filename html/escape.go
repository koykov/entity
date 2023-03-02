package html

import (
	"bytes"

	"github.com/koykov/byteseq"
)

func Escape[T byteseq.Byteseq](x T) T {
	r := AppendEscape(nil, x)
	return T(r)
}

func AppendEscape[T byteseq.Byteseq](dst []byte, x T) []byte {
	p := byteseq.Q2B(x)
	l := len(p)
	if l == 0 {
		return dst
	}

	if i := bytes.IndexAny(p, "&'<>\""); i == -1 {
		dst = append(dst, p...)
		return dst
	}

	_ = p[l-1]
	for i := 0; i < l; i++ {
		switch p[i] {
		case '&':
			dst = append(dst, "&amp;"...)
		case '\'':
			dst = append(dst, "&#39;"...)
		case '<':
			dst = append(dst, "&lt;"...)
		case '>':
			dst = append(dst, "&gt;"...)
		case '"':
			dst = append(dst, "&#34;"...)
		default:
			dst = append(dst, p[i])
		}
	}

	return dst
}
