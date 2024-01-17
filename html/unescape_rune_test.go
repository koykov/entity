package html

import (
	"testing"

	"github.com/koykov/byteconv"
)

func TestUnescapeRune(t *testing.T) {
	for _, stage := range stagesUnesc {
		t.Run(stage.key, func(t *testing.T) {
			r := UnescapeRune(stage.raw)
			e := byteconv.AppendS2R(nil, stage.expect)
			if !assertR(r, e) {
				t.FailNow()
			}
		})
	}
}

func BenchmarkUnescapeRune(b *testing.B) {
	for _, stage := range stagesUnesc {
		b.Run(stage.key, func(b *testing.B) {
			b.ReportAllocs()
			var buf []rune
			for i := 0; i < b.N; i++ {
				buf = AppendUnescapeRune(buf[:0], stage.raw)
			}
			_ = buf
		})
	}
}
