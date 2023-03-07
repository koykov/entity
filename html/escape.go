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

func WriteEscape[T byteseq.Byteseq](w Writer, x T) (n int, err error) {
	p := byteseq.Q2B(x)
	l := len(p)
	if l == 0 {
		return
	}

	if i := bytes.IndexAny(p, "&'<>\""); i == -1 {
		n, err = w.Write(p)
		return
	}

	_ = p[l-1]
	for i := 0; i < l; i++ {
		var n1 int
		switch p[i] {
		case '&':
			if n1, err = w.WriteString("&amp;"); err != nil {
				return
			}
		case '\'':
			if n1, err = w.WriteString("&#39;"); err != nil {
				return
			}
		case '<':
			if n1, err = w.WriteString("&lt;"); err != nil {
				return
			}
		case '>':
			if n1, err = w.WriteString("&gt;"); err != nil {
				return
			}
		case '"':
			if n1, err = w.WriteString("&#34;"); err != nil {
				return
			}
		default:
			if err = w.WriteByte(p[i]); err != nil {
				return
			}
			n1 = 1
		}
		n += n1
	}

	return
}
