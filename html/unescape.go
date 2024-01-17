package html

import (
	"bytes"
	"strconv"
	"unicode"

	"github.com/koykov/byteconv"
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
			dst = unescB(dst, p[lo:hi])
		case tag && !unicode.IsLetter(rune(c)) && !unicode.IsDigit(rune(c)) && c != '#' && c != 'x' && c != 'X':
			tag = false
			hi = i
			dst = unescB(dst, p[lo:hi])
			dst = append(dst, c)
		default:
			if !tag {
				dst = append(dst, c)
			}
		}
	}
	if tag {
		dst = unescB(dst, p[lo:l])
	}

	return dst
}

func unescB(dst, ent []byte) []byte {
	if len(ent) < 3 {
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
		if pent[0] == 'x' || pent[0] == 'X' {
			base = 16
			pent = pent[1:]
		}
		var rest []byte
		if base == 10 {
			for i := 0; i < len(pent); i++ {
				if '0' > pent[i] || pent[i] > '9' {
					rest = pent[i:]
					pent = pent[:i]
					break
				}
			}
		}
		if len(pent) == 0 {
			dst = append(dst, ent...)
			return dst
		}
		i, err := strconv.ParseInt(byteconv.B2S(pent), base, 64)
		if err != nil {
			dst = append(dst, ent...)
			return dst
		}
		if i1, ok := __bufHCP[i]; ok {
			h := __bufH[i1]
			dst = append(dst, h.Value()...)
		} else {
			if 0x80 <= i && i <= 0x9F {
				i = int64(cp1252[i-0x80])
			} else if i == 0 || (0xD800 <= i && i <= 0xDFFF) || i > 0x10FFFF {
				i = '\uFFFD'
			}
			if i1, ok := __bufHCP[i]; ok {
				h := __bufH[i1]
				dst = append(dst, h.Value()...)
			} else {
				dst = append(dst, ent...)
			}
		}
		dst = append(dst, rest...)
	default:
		s := byteconv.B2S(ent)
		if i1, ok := __bufHN[s]; ok {
			h := __bufH[i1]
			dst = append(dst, h.Value()...)
		} else {
			dst = append(dst, s...)
		}
	}
	return dst
}

func WriteUnescape[T byteseq.Byteseq](w Writer, x T) (n int, err error) {
	p := byteseq.Q2B(x)
	l := len(p)
	if l == 0 {
		return
	}

	if i := bytes.IndexByte(p, '&'); i == -1 {
		n, err = w.Write(p)
		return
	}

	_ = p[l-1]
	lo, hi, n1 := 0, 0, 0
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
			if n1, err = unescW(w, p[lo:hi]); err != nil {
				return
			}
			n += n1
		case tag && !unicode.IsLetter(rune(c)) && !unicode.IsDigit(rune(c)) && c != '#' && c != 'x' && c != 'X':
			tag = false
			hi = i
			if n1, err = unescW(w, p[lo:hi]); err != nil {
				return
			}
			n += n1
			if err = w.WriteByte(c); err != nil {
				return
			}
			n++
		default:
			if !tag {
				if err = w.WriteByte(c); err != nil {
					return
				}
				n++
			}
		}
	}
	if tag {
		if n1, err = unescW(w, p[lo:l]); err != nil {
			return
		}
		n += n1
	}

	return
}

func unescW(w Writer, ent []byte) (n int, err error) {
	var n1 int
	if len(ent) < 3 {
		if n1, err = w.Write(ent); err != nil {
			return
		}
		n += n1
		return
	}
	switch {
	case ent[1] == '#':
		lo, hi := 2, len(ent)
		if ent[hi-1] == ';' {
			hi--
		}
		base := 10
		pent := ent[lo:hi]
		if pent[0] == 'x' || pent[0] == 'X' {
			base = 16
			pent = pent[1:]
		}
		var rest []byte
		if base == 10 {
			for i := 0; i < len(pent); i++ {
				if '0' > pent[i] || pent[i] > '9' {
					rest = pent[i:]
					pent = pent[:i]
					break
				}
			}
		}
		if len(pent) == 0 {
			if n1, err = w.Write(ent); err != nil {
				return
			}
			n += n1
			return
		}
		i, err1 := strconv.ParseInt(byteconv.B2S(pent), base, 64)
		if err1 != nil {
			if n1, err = w.Write(ent); err != nil {
				return
			}
			n += n1
			return
		}
		if i1, ok := __bufHCP[i]; ok {
			h := __bufH[i1]
			if n1, err = w.WriteString(h.Value()); err != nil {
				return
			}
			n += n1
		} else {
			if 0x80 <= i && i <= 0x9F {
				i = int64(cp1252[i-0x80])
			} else if i == 0 || (0xD800 <= i && i <= 0xDFFF) || i > 0x10FFFF {
				i = '\uFFFD'
			}
			n1 = 0
			if i1, ok := __bufHCP[i]; ok {
				h := __bufH[i1]
				if n1, err = w.WriteString(h.Value()); err != nil {
					return
				}
			} else {
				if n1, err = w.Write(ent); err != nil {
					return
				}
			}
			n += n1
		}
		if n1, err = w.Write(rest); err != nil {
			return
		}
		n += n1
	default:
		s := byteconv.B2S(ent)
		if i1, ok := __bufHN[s]; ok {
			h := __bufH[i1]
			if n1, err = w.WriteString(h.Value()); err != nil {
				return
			}
		} else {
			if n1, err = w.WriteString(s); err != nil {
				return
			}
		}
		n += n1
	}
	return
}
