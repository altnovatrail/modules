package main

import (
	"bytes"
	"encoding/json"
	"image/png"
	"slices"
	"strings"
	"testing"
)

// The kinds of T20 draw as every kind does: opaque, deterministic, of the size their part names.
func TestKinds20Draw(t *testing.T) {
	for _, raw := range []string{
		`{"kind": "containers", "items": [{"capacity": 3, "amount": 0}, {"capacity": 5, "amount": 5, "label": "A"}]}`,
		`{"kind": "containers", "items": [{"capacity": 20, "amount": 1}, {"capacity": 1, "amount": 1}, {"capacity": 1, "amount": 0},
			{"capacity": 10, "amount": 9}]}`,
		`{"kind": "grid", "rows": ["A"], "cols": ["1"]}`,
		`{"kind": "grid", "rows": ["A", "B", "C", "D"], "cols": ["1", "2", "3", "4"], "filled": ["B2", "D4"],
			"marks": {"A1": "1", "C4": "?", "D4": "8"}}`,
		`{"kind": "grid", "rows": ["10", "20"], "cols": ["A", "B", "C", "D", "E", "F", "G", "H"], "filled": ["10H"]}`,
		`{"kind": "piles", "piles": [{"count": 3}, {"count": 2, "shape": "square"}]}`,
		`{"kind": "piles", "piles": [{"label": "10", "count": 40, "group": 5, "value": "?"}, {"skip": true}, {"count": 0},
			{"label": "A", "boxed": true, "fill": "light"}, {"count": 30, "shown": 7, "shape": "square", "fill": "light"}], "box": true}`,
		`{"kind": "piles", "piles": [{"label": "A", "count": 7, "value": "B"}, {"skip": true}, {"count": 1, "boxed": true},
			{"count": 0}, {"count": 12, "shown": 4}], "across": true, "box": true}`,
		`{"kind": "venn", "sets": [{"label": "F"}, {"label": "C", "count": "12"}]}`,
		`{"kind": "venn", "total": "25", "sets": [{"label": "F", "count": "15"}, {"label": "C", "count": "12"}], "both": "5", "neither": "?"}`,
		`{"kind": "balance", "left": [], "right": ["1"]}`,
		`{"kind": "balance", "left": ["?", "AB", "100", "2"], "right": ["1", "3", "9"]}`,
		`{"kind": "ring", "count": 3}`,
		`{"kind": "ring", "count": 24, "start": {"at": 24, "label": "S"}}`,
		`{"kind": "ring", "count": 7, "start": {"at": 3}}`,
	} {
		part, data, err := pictureOf("pic-1", json.RawMessage(raw), "A picture.", "")
		if err != nil {
			t.Errorf("%s: pictureOf() error = %v", raw, err)
			continue
		}
		img, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			t.Fatalf("%s: png.Decode() error = %v", raw, err)
		}
		width, height := part[2].value.(int), part[3].value.(int)
		if b := img.Bounds(); b.Dx() != width*scale || b.Dy() != height*scale {
			t.Errorf("%s: the PNG is %v, the part says %d × %d CSS px", raw, b, width, height)
		}
		if !opaque(img) || !inked(img) {
			t.Errorf("%s: a pixel is transparent, or nothing is drawn", raw)
		}
		if got := strings.Join(slices.Compact(chunks(t, data)), " "); got != "IHDR IDAT IEND" {
			t.Errorf("%s: the chunks are %s", raw, got)
		}
		if _, again, _ := pictureOf("pic-1", json.RawMessage(raw), "A picture.", ""); !bytes.Equal(again, data) {
			t.Errorf("%s: a second run draws other bytes", raw)
		}
	}
}

// A pile's counters: cut short to half before an ellipsis and half after, a heap of no stated size
// as a counter, an ellipsis and a counter, groups apart, and lines that break before the width.
func TestPileTokens(t *testing.T) {
	ten, seven, five, three := 10, 7, 5, 3
	shown := func(p pileItem) string {
		tokens, _, _ := p.tokens(widestPile)
		var out []byte
		for _, tk := range tokens {
			out = append(out, map[bool]byte{true: '.', false: 'o'}[tk.ellipsis])
		}
		return string(out)
	}
	if got := shown(pileItem{Count: &ten, Shown: &seven}); got != "oooo.ooo" {
		t.Errorf("cut short: %s", got)
	}
	if got := shown(pileItem{}); got != "o.o" {
		t.Errorf("a heap: %s", got)
	}
	tokens, lines, reach := pileItem{Count: &ten, Group: &five}.tokens(widestPile)
	if lines != 1 || tokens[5].x != 5*counterPitch+groupGap || reach != 10*counterPitch+groupGap {
		t.Errorf("groups: %v, %d, %v", tokens, lines, reach)
	}
	if _, lines, _ := (pileItem{Count: &ten, Group: &three}).tokens(4 * counterPitch); lines != 4 {
		t.Errorf("lines: %d, want 4", lines)
	}
	// A group that would pass the line's end starts the next line; one wider than a line breaks in it.
	if tokens, _, _ := (pileItem{Count: &ten, Group: &five}).tokens(7 * counterPitch); tokens[5].x != 0 || tokens[5].line != 1 {
		t.Errorf("a group split: %v", tokens[5])
	}
	if _, lines, _ := (pileItem{Count: &ten, Group: &ten}).tokens(4 * counterPitch); lines != 3 {
		t.Errorf("a group wider than a line: %d lines, want 3", lines)
	}
}

func TestKinds20Refuse(t *testing.T) {
	for _, test := range []struct{ raw, message string }{
		{`{"kind": "containers", "items": []}`, "1 to 4"},
		{`{"kind": "containers", "items": [{"capacity": 1, "amount": 0}, {"capacity": 1, "amount": 0}, {"capacity": 1, "amount": 0},
			{"capacity": 1, "amount": 0}, {"capacity": 1, "amount": 0}]}`, "1 to 4"},
		{`{"kind": "containers", "items": [{"amount": 0}]}`, "from 1 to 20"},
		{`{"kind": "containers", "items": [{"capacity": 21, "amount": 0}]}`, "from 1 to 20"},
		{`{"kind": "containers", "items": [{"capacity": 3}]}`, "from 0 up to its capacity"},
		{`{"kind": "containers", "items": [{"capacity": 3, "amount": 4}]}`, "from 0 up to its capacity"},
		{`{"kind": "containers", "items": [{"capacity": 3, "amount": 1, "label": "ab"}]}`, `label "ab"`},
		{`{"kind": "grid", "rows": [], "cols": ["1"]}`, "1 to 8 rows"},
		{`{"kind": "grid", "rows": ["A"], "cols": ["1", "2", "3", "4", "5", "6", "7", "8", "9"]}`, "1 to 8 rows"},
		{`{"kind": "grid", "rows": ["a"], "cols": ["1"]}`, `name "a"`},
		{`{"kind": "grid", "rows": ["A", "A"], "cols": ["1"]}`, `names two lines`},
		{`{"kind": "grid", "rows": ["1", "12"], "cols": ["22", "2"]}`, "name two cells alike"},
		{`{"kind": "grid", "rows": ["A"], "cols": ["1"], "filled": ["B1"]}`, `cell "B1" is no cell`},
		{`{"kind": "grid", "rows": ["A"], "cols": ["1"], "filled": ["A1", "A1"]}`, "filled twice"},
		{`{"kind": "grid", "rows": ["A"], "cols": ["1"], "marks": {"B1": "1"}}`, `mark "B1"`},
		{`{"kind": "grid", "rows": ["A"], "cols": ["1"], "marks": {"A1": "one"}}`, `label "one"`},
		{`{"kind": "grid", "rows": ["A"], "cols": ["1"], "marks": {"A1": "1", "A1": "2"}}`, `member "A1"`},
		{`{"kind": "number_line", "from": 0, "to": 6, "marks": {"AT": 1}}`, `member "AT"`},
		{`{"kind": "grid", "rows": ["A"], "cols": ["1"], "marks": {"A1": {"Deep": 1}}}`, `member "Deep"`},
		{`{"kind": "piles", "piles": []}`, "1 to 6"},
		{`{"kind": "piles", "piles": [{}, {}, {}, {}, {}, {}, {}]}`, "1 to 6"},
		{`{"kind": "piles", "piles": [{"skip": true}, {}]}`, "first, last or next to another"},
		{`{"kind": "piles", "piles": [{}, {"skip": true}, {"skip": true}, {}]}`, "first, last or next to another"},
		{`{"kind": "piles", "piles": [{}, {"skip": true, "count": 1}, {}]}`, "a skip, which holds nothing else"},
		{`{"kind": "piles", "piles": [{"count": 41}]}`, "from 0 to 40"},
		{`{"kind": "piles", "piles": [{"shown": 2}]}`, "of a stated size"},
		{`{"kind": "piles", "piles": [{"count": 5, "shown": 5}]}`, "from 2 to one less"},
		{`{"kind": "piles", "piles": [{"count": 5, "shown": 1}]}`, "from 2 to one less"},
		{`{"kind": "piles", "piles": [{"group": 0}]}`, "group is from 1 to 40"},
		{`{"kind": "piles", "piles": [{"fill": "grey"}]}`, "dark or light"},
		{`{"kind": "piles", "piles": [{"shape": "star"}]}`, "dot or square"},
		{`{"kind": "piles", "piles": [{"value": "ab"}]}`, `label "ab"`},
		{`{"kind": "venn", "sets": [{"label": "A"}]}`, "exactly 2 sets"},
		{`{"kind": "venn", "sets": [{"label": "A"}, {"count": "3"}]}`, "has a label"},
		{`{"kind": "venn", "sets": [{"label": "A"}, {"label": "ab"}]}`, "has a label"},
		{`{"kind": "venn", "sets": [{"label": "A"}, {"label": "B"}], "both": "ab"}`, `count "ab"`},
		{`{"kind": "balance", "left": ["1"]}`, "a left pan and a right pan"},
		{`{"kind": "balance", "left": ["1", "2", "3", "4", "5"], "right": []}`, "at most 4 weights"},
		{`{"kind": "balance", "left": ["one"], "right": []}`, `label "one"`},
		{`{"kind": "ring", "count": 2}`, "3 to 24 places"},
		{`{"kind": "ring", "count": 25}`, "3 to 24 places"},
		{`{"kind": "ring", "count": 5, "start": {"label": "S"}}`, "from 1 to its count"},
		{`{"kind": "ring", "count": 5, "start": {"at": 6}}`, "from 1 to its count"},
		{`{"kind": "ring", "count": 5, "start": {"at": 1, "label": "ab"}}`, `label "ab"`},
	} {
		_, _, err := pictureOf("pic-1", json.RawMessage(test.raw), "A picture.", "")
		if err == nil || !strings.Contains(err.Error(), test.message) {
			t.Errorf("%s: error = %v, want one that says %q", test.raw, err, test.message)
		}
	}
}
