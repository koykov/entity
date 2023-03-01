package html

import (
	"strings"
	"testing"
)

type stage struct {
	key    string
	raw    string
	expect string
}

var (
	stagesEsc = []stage{
		{"copy", "Lorem ipsum dolor sit amet, consectetur adipiscing elit.", "Lorem ipsum dolor sit amet, consectetur adipiscing elit."},
		{"simple", "foo & > < bar", "foo &amp; &gt; &lt; bar"},
		{"stringEnd", "foobar '", "foobar &#39;"},
		{
			"multiple",
			strings.Repeat("foo < bar > asd & fgh ' zzz \" ", 100),
			strings.Repeat("foo &lt; bar &gt; asd &amp; fgh &#39; zzz &#34; ", 100),
		},
	}
	stagesUnesc = []stage{
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
)

func TestEscape(t *testing.T) {
	for _, stage := range stagesEsc {
		t.Run(stage.key, func(t *testing.T) {
			r := Escape(stage.raw)
			if r != stage.expect {
				t.FailNow()
			}
		})
	}
}

func BenchmarkUnescape(b *testing.B) {
	for _, stage := range stagesEsc {
		b.Run(stage.key, func(b *testing.B) {
			b.ReportAllocs()
			var buf []byte
			for i := 0; i < b.N; i++ {
				buf = AppendEscape(buf[:0], stage.raw)
			}
		})
	}
}
