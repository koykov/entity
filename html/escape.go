package html

import (
	"github.com/koykov/bytealg"
	"github.com/koykov/byteseq"
)

var (
	esc = []byte("&'<>\"")
)

func Escape[T byteseq.Byteseq](x T) T {
	r := AppendEscape(nil, x)
	return T(r)
}

func AppendEscape[T byteseq.Byteseq](dst []byte, x T) []byte {
	p := byteseq.Q2B(x)
	off := 0
	for {
		i := bytealg.IndexAnyAt(p, esc, off)
		if i < 0 {
			break
		}
		dst = append(dst, p[off:i]...)
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
		}
		off = i + 1
	}
	dst = append(dst, p[off:]...)

	return dst
}
