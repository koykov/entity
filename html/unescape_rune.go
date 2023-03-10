package html

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/koykov/byteseq"
	"github.com/koykov/fastconv"
)

func AppendUnescapeRune[T byteseq.Byteseq](dst []rune, x T) []rune {
	s := byteseq.Q2S(x)
	l := len(s)
	if l == 0 {
		return dst
	}

	if i := strings.Index(s, "&"); i == -1 {
		dst = fastconv.AppendS2R(dst, s)
		return dst
	}

	lo, hi := 0, 0
	var tag bool
	for i, r := range s {
		switch {
		case r == '&':
			tag = true
			lo = i
			hi = lo
		case tag && r == ';':
			tag = false
			hi = i + 1
			dst = unescR(dst, s[lo:hi])
		case tag && !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '#' && r != 'x' && r != 'X':
			tag = false
			hi = i
			dst = fastconv.AppendS2R(dst, s[lo:hi])
			dst = append(dst, r)
		default:
			if !tag {
				dst = append(dst, r)
			}
		}
	}
	if tag {
		dst = unescR(dst, s[lo:l])
	}

	return dst
}

func unescR(dst []rune, ent string) []rune {
	if len(ent) < 3 {
		dst = fastconv.AppendS2R(dst, ent)
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
		var rest string
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
			dst = fastconv.AppendS2R(dst, ent)
			return dst
		}
		i, err := strconv.ParseInt(pent, base, 64)
		if err != nil {
			dst = fastconv.AppendS2R(dst, ent)
			return dst
		}
		if i1, ok := __bufHCP[i]; ok {
			h := __bufH[i1]
			dst = fastconv.AppendS2R(dst, h.Value())
		} else {
			if 0x80 <= i && i <= 0x9F {
				i = int64(cp1252[i-0x80])
			} else if i == 0 || (0xD800 <= i && i <= 0xDFFF) || i > 0x10FFFF {
				i = '\uFFFD'
			}
			if i1, ok := __bufHCP[i]; ok {
				h := __bufH[i1]
				dst = fastconv.AppendS2R(dst, h.Value())
			} else {
				dst = fastconv.AppendS2R(dst, ent)
			}
		}
		dst = fastconv.AppendS2R(dst, rest)
	default:
		if i1, ok := __bufHN[ent]; ok {
			h := __bufH[i1]
			dst = fastconv.AppendS2R(dst, h.Value())
		} else {
			dst = fastconv.AppendS2R(dst, ent)
		}
	}
	return dst
}
