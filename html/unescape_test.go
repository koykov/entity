package html

import (
	"bytes"
	"testing"
)

var stagesUnesc = []stage{
	{"copy", "A\ttext\nstring", "A\ttext\nstring"},
	{"simple", "&amp; &gt; &lt;", "& > <"},
	{"stringEnd", "&amp &amp", "& &"},
	{"multiCodepoint", "text &gesl; blah", "text \u22db\ufe00 blah"},
	{"decimalEntity", "Delta = &#916; ", "Delta = Δ "},
	{"hexadecimalEntity", "Lambda = &#x3bb; = &#X3Bb ", "Lambda = λ = λ "},
	{"numericEnds", "&# &#x &#128;43 &copy = &#169f = &#xa9", "&# &#x €43 © = ©f = ©"},
	{"numericReplacements", "Footnote&#x87;", "Footnote‡"},
	{"copySingleAmpersand", "&", "&"},
	{"copyAmpersandNonEntity", "text &test", "text &test"},
	{"copyAmpersandHash", "text &#", "text &#"},
}

func TestUnescape(t *testing.T) {
	for _, stage := range stagesUnesc {
		t.Run(stage.key, func(t *testing.T) {
			r := Unescape(stage.raw)
			if r != stage.expect {
				t.FailNow()
			}
		})
	}
}

func BenchmarkUnescape(b *testing.B) {
	for _, stage := range stagesUnesc {
		b.Run(stage.key, func(b *testing.B) {
			b.ReportAllocs()
			var buf bytes.Buffer
			for i := 0; i < b.N; i++ {
				buf.Reset()
				_, _ = WriteUnescape(&buf, stage.raw)
			}
			_ = buf
		})
	}
}
