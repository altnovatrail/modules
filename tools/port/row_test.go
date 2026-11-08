package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"image"
	"image/png"
	"strings"
	"testing"
)

// chunks lists the chunk types of a PNG file.
func chunks(t *testing.T, data []byte) []string {
	t.Helper()
	signature := []byte("\x89PNG\r\n\x1a\n")
	if !bytes.HasPrefix(data, signature) {
		t.Fatal("no PNG signature")
	}
	var types []string
	for rest := data[len(signature):]; len(rest) >= 12; {
		length := binary.BigEndian.Uint32(rest)
		types = append(types, string(rest[4:8]))
		rest = rest[12+length:]
	}
	return types
}

func drawRow(t *testing.T, raw string) ([]byte, mapping) {
	t.Helper()
	part, data, err := pictureOf("row-1", json.RawMessage(raw), "A row.", "")
	if err != nil {
		t.Fatalf("pictureOf() error = %v", err)
	}
	return data, part
}

// Every form of a row draws: an opaque PNG of the size the part names, twice as many pixels each
// way, with the chunks IHDR, IDAT and IEND alone, and the same bytes again.
func TestRowDraws(t *testing.T) {
	for _, raw := range []string{
		rowJSON,
		`{"kind": "row", "items": [{"label": "1"}, {"label": "2", "mark": "ring"}, {"skip": true}, {"label": "50", "mark": "square", "below": "B"}]}`,
		`{"kind": "row", "items": [{}, {}, {}, {}], "gaps": "10", "span": "90", "copies": 2}`,
		`{"kind": "row", "items": [{"label": "A"}, {"label": "B"}], "line": false, "span": "?"}`,
	} {
		data, part := drawRow(t, raw)
		img, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			t.Fatalf("%s: png.Decode() error = %v", raw, err)
		}
		width, height := part[2].value.(int), part[3].value.(int)
		if b := img.Bounds(); b.Dx() != width*scale || b.Dy() != height*scale {
			t.Errorf("%s: the PNG is %v, the part says %d × %d CSS px", raw, b, width, height)
		}
		if !opaque(img) {
			t.Errorf("%s: a pixel is transparent", raw)
		}
		if got := strings.Join(chunks(t, data), " "); got != "IHDR IDAT IEND" {
			t.Errorf("%s: the chunks are %s", raw, got)
		}
		if again, _ := drawRow(t, raw); !bytes.Equal(again, data) {
			t.Errorf("%s: a second run draws other bytes", raw)
		}
		if !inked(img) {
			t.Errorf("%s: nothing is drawn", raw)
		}
	}
}

func opaque(img image.Image) bool {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if _, _, _, a := img.At(x, y).RGBA(); a != 0xffff {
				return false
			}
		}
	}
	return true
}

func inked(img image.Image) bool {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if r, _, _, _ := img.At(x, y).RGBA(); r < 0x8000 {
				return true
			}
		}
	}
	return false
}

func TestRowRefuses(t *testing.T) {
	for _, test := range []struct{ raw, message string }{
		{`{"kind": "lattice", "width": 2}`, `"lattice", which the port does not draw yet`},
		{`not json`, "its picture"},
		{`{"kind": "row", "items": [{}, {}], "colour": "red"}`, "unknown field"},
		{`{"kind": "row", "items": [{}]}`, "2 to 12 items"},
		{`{"kind": "row", "items": [{}, {}], "copies": 3}`, "once or twice"},
		{`{"kind": "row", "items": [{"skip": true}, {}, {}]}`, "first, last or next to another"},
		{`{"kind": "row", "items": [{}, {"skip": true}, {"skip": true}, {}]}`, "first, last or next to another"},
		{`{"kind": "row", "items": [{}, {"skip": true, "label": "1"}, {}]}`, "has a label or a mark"},
		{`{"kind": "row", "items": [{"mark": "star"}, {}]}`, `mark "star"`},
		{`{"kind": "row", "items": [{"label": "abc"}, {}]}`, `label "abc"`},
		{`{"kind": "row", "items": [{}, {}], "gaps": "123456"}`, `label "123456"`},
	} {
		_, _, err := pictureOf("row-1", json.RawMessage(test.raw), "A row.", "")
		if err == nil || !strings.Contains(err.Error(), test.message) {
			t.Errorf("%s: error = %v, want one that says %q", test.raw, err, test.message)
		}
	}
}
