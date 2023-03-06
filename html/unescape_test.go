package html

import "testing"

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
