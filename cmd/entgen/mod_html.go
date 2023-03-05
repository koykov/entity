package main

import (
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/koykov/entry"
)

type entities map[string]struct {
	Codepoints []int  `json:"codepoints"`
	Characters string `json:"characters"`
}

var (
	raw   entities
	names []string
)

type htmlModule struct{}

func (m htmlModule) Validate(input, _ string) error {
	if len(input) == 0 {
		return fmt.Errorf("param -input is required")
	}
	return nil
}

func (m htmlModule) Compile(w moduleWriter, input, target string) (err error) {
	if len(target) == 0 {
		target = "html/repo.go"
	}

	resp, err := http.Get(input)
	if err != nil {
		return
	}
	defer func() { _ = resp.Body.Close() }()

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if err = json.Unmarshal(contents, &raw); err != nil {
		return
	}

	for name := range raw {
		names = append(names, name)
	}
	sort.Strings(names)

	var buf []byte

	_, _ = w.WriteString("var (\n")

	// Write list of HTML objects.
	_, _ = w.WriteString("__bufH = []Entity{\n")
	for i := 0; i < len(names); i++ {
		var en, ev entry.Entry32
		lo := len(buf)
		buf = append(buf, names[i]...)
		hi := len(buf)
		en.Encode(uint16(lo), uint16(hi))

		lo = len(buf)
		buf = append(buf, raw[names[i]].Characters...)
		hi = len(buf)
		ev.Encode(uint16(lo), uint16(hi))

		_, _ = w.WriteString("Entity{name:")
		_, _ = w.WriteString(fmt.Sprintf("0x%08x", en))
		_, _ = w.WriteString(",val:")
		_, _ = w.WriteString(fmt.Sprintf("0x%08x", ev))
		_, _ = w.WriteString(",cp:")
		var cp int64
		cp = int64(raw[names[i]].Codepoints[0]) << 32
		if len(raw[names[i]].Codepoints) > 1 {
			cp = cp | int64(raw[names[i]].Codepoints[1])
		}
		_, _ = w.WriteString(fmt.Sprintf("0x%08x", cp))
		_, _ = w.WriteString("},\n")
	}
	_, _ = w.WriteString("}\n")

	// Write name/HTML registry.
	_, _ = w.WriteString("__bufHN = map[string]int{\n")
	for i := 0; i < len(names); i++ {
		_, _ = w.WriteString(`"`)
		_, _ = w.WriteString(names[i])
		_, _ = w.WriteString(`":`)
		_, _ = w.WriteString(fmt.Sprintf("0x%04x", i))
		_, _ = w.WriteString(",\n")
	}
	_, _ = w.WriteString("}\n")

	// Write codepoint/HTML registry.
	var hcp = map[int64]struct{}{}
	_, _ = w.WriteString("__bufHCP = map[int64]int{\n")
	for i := 0; i < len(names); i++ {
		var cp int64
		cp = int64(raw[names[i]].Codepoints[0])
		// todo implement HCP2 registry.
		// cp = int64(raw[names[i]].Codepoints[0]) << 32
		// if len(raw[names[i]].Codepoints) > 1 {
		// 	cp = cp | int64(raw[names[i]].Codepoints[1])
		// }
		if _, ok := hcp[cp]; ok {
			continue
		}
		hcp[cp] = struct{}{}
		_, _ = w.WriteString(fmt.Sprintf("0x%08x", cp))
		_, _ = w.WriteString(":")
		_, _ = w.WriteString(fmt.Sprintf("0x%04x", i))
		_, _ = w.WriteString(",\n")
	}
	_, _ = w.WriteString("}\n")

	// Write char/HTML registry.
	var hc = map[string]struct{}{}
	_, _ = w.WriteString("__bufHC = map[string]int{\n")
	for i := 0; i < len(names); i++ {
		c := raw[names[i]].Characters
		c = strings.ReplaceAll(c, "\\", "\\\\")
		c = strings.ReplaceAll(c, "\n", "\\n")
		c = strings.ReplaceAll(c, "\"", `\"`)
		if _, ok := hc[c]; ok {
			continue
		}
		hc[c] = struct{}{}
		_, _ = w.WriteString(`"`)
		_, _ = w.WriteString(c)
		_, _ = w.WriteString(`":`)
		_, _ = w.WriteString(fmt.Sprintf("0x%04x", i))
		_, _ = w.WriteString(",\n")
	}
	_, _ = w.WriteString("}\n")

	_, _ = w.WriteString("__buf = []byte{\n")
	for i := 0; i < len(buf); i++ {
		if i > 0 && i%16 == 0 {
			_ = w.WriteByte('\n')
		}
		_, _ = w.WriteString(fmt.Sprintf("0x%02x, ", buf[i]))
	}
	_, _ = w.WriteString("\n}\n")

	_, _ = w.WriteString(")\n")

	source := w.Bytes()
	var fmtSource []byte
	if fmtSource, err = format.Source(source); err != nil {
		return
	}

	if err = ioutil.WriteFile(target, fmtSource, 0644); err != nil {
		return
	}

	return
}
