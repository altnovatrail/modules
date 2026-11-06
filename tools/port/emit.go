package main

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
)

// The YAML that the port writes keeps to the rules of the format (format/README.md, section 4) by
// writing little: block mappings with plain keys, which are all snake_case words; block sequences
// of mappings; flow sequences of strings; integers; and every string in double quotes, as JSON
// writes it. JSON's escapes are YAML's, and Go never writes \/, which the format refuses. The order
// of keys is the order the port gives them, so a run writes the same bytes again.

type (
	// mapping is a block mapping, its keys in order.
	mapping []pair
	pair    struct {
		key   string
		value any
	}
	// sequence is a sequence: of strings, written in flow, or of mappings, written in block.
	sequence []any
)

// emit writes a document: the comment, a line each, then the mapping.
func emit(comment []string, root mapping) []byte {
	var out strings.Builder
	for _, line := range comment {
		out.WriteString("# " + line + "\n")
	}
	writeMapping(&out, root, 0, false)
	return []byte(out.String())
}

// writeMapping writes a mapping, each key at indent; inline says that the first key follows a "- "
// already written, as an item of a block sequence.
func writeMapping(out *strings.Builder, m mapping, indent int, inline bool) {
	for i, p := range m {
		if i > 0 || !inline {
			out.WriteString(strings.Repeat(" ", indent))
		}
		writePair(out, p, indent)
	}
}

func writePair(out *strings.Builder, p pair, indent int) {
	out.WriteString(p.key + ":")
	switch value := p.value.(type) {
	case mapping:
		out.WriteString("\n")
		writeMapping(out, value, indent+2, false)
	case sequence:
		if flow(value) {
			out.WriteString(" [")
			for i, item := range value {
				if i > 0 {
					out.WriteString(", ")
				}
				out.WriteString(scalar(item))
			}
			out.WriteString("]\n")
			return
		}
		out.WriteString("\n")
		for _, item := range value {
			out.WriteString(strings.Repeat(" ", indent+2) + "- ")
			writeMapping(out, item.(mapping), indent+4, true)
		}
	default:
		out.WriteString(" " + scalar(value) + "\n")
	}
}

// flow tells a sequence of scalars, written on one line, from one of mappings.
func flow(s sequence) bool {
	for _, item := range s {
		if _, isMapping := item.(mapping); isMapping {
			return false
		}
	}
	return true
}

// scalar writes a string in double quotes as JSON does, and an integer as it is.
func scalar(value any) string {
	switch v := value.(type) {
	case string:
		var buffer bytes.Buffer
		encoder := json.NewEncoder(&buffer)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(v); err != nil {
			panic(err) // a Go string always encodes
		}
		return strings.TrimSuffix(buffer.String(), "\n")
	case int:
		return strconv.Itoa(v)
	default:
		panic("port: no scalar of this type is written")
	}
}
