package html

import (
	"bytes"
	"strconv"
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
		c := p[i]
		switch {
		case c == '&':
			tag = true
			lo = i
			hi = lo
		case tag && c == ';':
			tag = false
			hi = i + 1
			dst = unesc(dst, p[lo:hi])
		case tag && !unicode.IsLetter(rune(c)) && !unicode.IsDigit(rune(c)) && c != '#' && c != 'x' && c != 'X':
			tag = false
			hi = i
			dst = unesc(dst, p[lo:hi])
			dst = append(dst, c)
		default:
			if !tag {
				dst = append(dst, c)
			}
		}
	}
	if tag {
		dst = unesc(dst, p[lo:l])
	}

	return dst
}

func unesc(dst, ent []byte) []byte {
	if len(ent) < 2 {
		dst = append(dst, ent...)
		return dst
	}
	switch {
	case ent[1] == '#':
		lo, hi := 2, len(ent)
		if ent[hi-1] == ';' {
			hi--
		}
		base := 10
		pent := ent[lo:hi]
		if pent[0] == 'x' || pent[1] == 'X' {
			base = 16
			pent = pent[1:]
		}
		i, err := strconv.ParseInt(fastconv.B2S(pent), base, 64)
		if err != nil {
			dst = append(dst, ent...)
			return dst
		}
		if i1, ok := __bufHCP[i]; ok {
			h := __bufH[i1]
			dst = append(dst, h.Value()...)
		} else {
			dst = append(dst, ent...)
		}
	default:
		s := fastconv.B2S(ent)
		if i1, ok := __bufHN[s]; ok {
			h := __bufH[i1]
			dst = append(dst, h.Value()...)
		} else {
			dst = append(dst, s...)
		}
	}
	return dst
}
