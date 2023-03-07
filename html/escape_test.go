package html

import (
	"bytes"
	"strings"
	"testing"
)

var stagesEsc = []stage{
	{"copy", "Lorem ipsum dolor sit amet, consectetur adipiscing elit.", "Lorem ipsum dolor sit amet, consectetur adipiscing elit."},
	{"simple", "foo & > < bar", "foo &amp; &gt; &lt; bar"},
	{"stringEnd", "foobar '", "foobar &#39;"},
	{
		"long",
		strings.Repeat("foo < bar > asd & fgh ' zzz \" ", 100),
		strings.Repeat("foo &lt; bar &gt; asd &amp; fgh &#39; zzz &#34; ", 100),
	},
}

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

func BenchmarkEscape(b *testing.B) {
	for _, stage := range stagesEsc {
		b.Run(stage.key, func(b *testing.B) {
			b.ReportAllocs()
			var buf bytes.Buffer
			for i := 0; i < b.N; i++ {
				buf.Reset()
				_, _ = WriteEscape(&buf, stage.raw)
			}
			_ = buf
		})
	}
}
