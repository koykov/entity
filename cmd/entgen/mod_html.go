package main

import (
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/koykov/entry"
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
		target = "html_repo.go"
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

	_, _ = w.WriteString("__bufH = []HTML{\n")
	for i := 0; i < len(names); i++ {
		var e entry.Entry32
		lo := len(buf)
		buf = append(buf, raw[names[i]].Characters...)
		hi := len(buf)
		e.Encode(uint16(lo), uint16(hi))

		_, _ = w.WriteString("HTML{name:")
		_, _ = w.WriteString(fmt.Sprintf("0x%08x", e))
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

	_, _ = w.WriteString("__bufHN = map[string]int{\n")
	for i := 0; i < len(names); i++ {
		_, _ = w.WriteString(`"`)
		_, _ = w.WriteString(names[i])
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
