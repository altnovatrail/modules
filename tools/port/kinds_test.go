package main

import (
	"bytes"
	"encoding/json"
	"image/png"
	"slices"
	"strings"
	"testing"

	"golang.org/x/image/font"
)

// Every kind the port draws, in every form the tests give it, draws as a row does: an opaque PNG
// of the size its part names, twice as many pixels each way, with IHDR, IDAT and IEND alone, the
// same bytes again, and something inked (T19).
func TestKindsDraw(t *testing.T) {
	for _, raw := range []string{
		`{"kind": "clock", "time": "3:00"}`,
		`{"kind": "clock", "time": "23:59"}`,
		`{"kind": "clock", "time": "0:05", "hour_label": "A", "minute_label": "B"}`,
		`{"kind": "clock", "time": "9:15", "minute_label": "M"}`,
		`{"kind": "clock", "time": "12:00", "hour_label": "A", "minute_label": "B"}`,
		`{"kind": "table", "rows": [["1"]]}`,
		`{"kind": "table", "header": ["A", "B", "C"], "rows": [["8:30", "…", ""], ["10000", "?", "-1,5"]]}`,
		`{"kind": "table", "rows": [["1", "2", "3", "4", "5", "6"], ["1", "2", "3", "4", "5", "6"], ["1", "2", "3", "4", "5", "6"],
			["1", "2", "3", "4", "5", "6"], ["1", "2", "3", "4", "5", "6"], ["1", "2", "3", "4", "5", "6"], ["1", "2", "3", "4", "5", "6"],
			["1", "2", "3", "4", "5", "6"]]}`,
		`{"kind": "bars", "bars": [{}]}`,
		`{"kind": "bars", "bars": [{"label": "N", "parts": 3, "shaded": 2, "length": 12, "braces": [{"from": 0, "to": 2, "label": "24"}]},
			{"label": "M", "parts": 4, "shaded": 0, "length": 12, "value": "?", "span": "48"}], "notes": ["N + M = 48", "N < M"]}`,
		`{"kind": "bars", "bars": [{"segments": [{"size": 60, "label": "60"}, {"size": 1}], "braces": [{"from": 0, "to": 1, "label": "A"},
			{"from": 0, "to": 2, "label": "B"}, {"from": 1, "to": 2, "label": "C"}]}, {"length": 1}, {"parts": 12}, {"parts": 1, "shaded": 1}]}`,
		`{"kind": "bars", "bars": [{"length": 2, "segments": [{"size": 1}, {"size": 3}]}], "notes": ["ABCDEFGHIJ + KLMNOP = 10"]}`,
		`{"kind": "number_line", "from": 0, "to": 6}`,
		`{"kind": "number_line", "from": -9999, "to": 99999, "step": 109998, "marks": [{"at": -9999, "label": "P"}, {"at": 99999}]}`,
		`{"kind": "number_line", "from": 0, "to": 15, "marks": [{"at": 0}, {"at": 15, "label": "?"}]}`,
		`{"kind": "calendar", "first": 7, "days": 30, "marks": {"4": "1", "11": "2", "18": "3"}}`,
		`{"kind": "calendar", "first": 1, "days": 28, "week_starts": "sunday"}`,
		`{"kind": "calendar", "first": 6, "days": 31, "week_starts": "monday", "marks": {"31": "?"}}`,
		`{"kind": "calendar", "first": 7, "days": 31, "week_starts": "sunday", "marks": {"1": "A"}}`,
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

// The hands of a clock stand where the time puts them: the hour hand moved on by the minutes, and
// a time after 12:59 as a twelve-hour face shows it.
func TestClockHands(t *testing.T) {
	for time, want := range map[string][2]float64{
		"3:00": {90, 0}, "15:00": {90, 0}, "0:00": {0, 0}, "5:45": {172.5, 270}, "12:30": {15, 180}, "23:59": {359.5, 354},
	} {
		hour, minute := (&clockPicture{Time: time}).hands()
		if hour != want[0] || minute != want[1] {
			t.Errorf("%s: hands %v, %v, want %v", time, hour, minute, want)
		}
	}
}

// The weeks of a calendar: the 1st in its column, by the day the weeks start on.
func TestCalendarWeeks(t *testing.T) {
	seven, one := 7, 1
	thirty := 30
	if got := (&calendarPicture{First: &seven, Days: &thirty}).offset(); got != 6 {
		t.Errorf("a Sunday, weeks from Monday: column %d", got)
	}
	if got := (&calendarPicture{First: &seven, Days: &thirty, WeekStarts: "sunday"}).offset(); got != 0 {
		t.Errorf("a Sunday, weeks from Sunday: column %d", got)
	}
	if got := len((&calendarPicture{First: &seven, Days: &thirty}).weekTalls()); got != 6 {
		t.Errorf("30 days from a Sunday: %d weeks, want 6", got)
	}
	if got := len((&calendarPicture{First: &one, Days: &[]int{28}[0]}).weekTalls()); got != 4 {
		t.Errorf("28 days from a Monday: %d weeks, want 4", got)
	}
}

// Every character a picture writes is in Go Regular: no label falls back to a missing glyph.
func TestGlyphs(t *testing.T) {
	face := faceAt(fontSize)
	for _, r := range "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ?.,−-+×÷=<>()…MonTuesWdhFriSat" {
		if _, ok := face.GlyphAdvance(r); !ok {
			t.Errorf("Go Regular has no glyph %q", r)
		}
	}
	if measure("12", fontSize) <= 0 || font.MeasureString(face, "12") <= 0 {
		t.Error("a label measures nothing")
	}
}

func TestKindsRefuse(t *testing.T) {
	for _, test := range []struct{ raw, message string }{
		{`{"kind": "clock"}`, "a time from 0:00 to 23:59"},
		{`{"kind": "clock", "time": "24:00"}`, "a time from 0:00 to 23:59"},
		{`{"kind": "clock", "time": "7:5"}`, "a time from 0:00 to 23:59"},
		{`{"kind": "clock", "time": "7:05", "hour_label": "ab"}`, `label "ab"`},
		{`{"kind": "clock", "time": "7:05", "hands": 2}`, "unknown field"},
		{`{"kind": "table", "rows": []}`, "1 to 8 rows"},
		{`{"kind": "table", "rows": [["1"], ["2"], ["3"], ["4"], ["5"], ["6"], ["7"], ["8"], ["9"]]}`, "1 to 8 rows"},
		{`{"kind": "table", "rows": [["1", "2", "3", "4", "5", "6", "7"]]}`, "1 to 6 cells"},
		{`{"kind": "table", "header": [], "rows": [["1"]]}`, "1 to 6 cells"},
		{`{"kind": "table", "header": ["A", "B"], "rows": [["1"]]}`, "row 1 of rows is not as long as the header"},
		{`{"kind": "table", "rows": [["1", "2"], ["1"]]}`, "row 2 of rows is not as long"},
		{`{"kind": "table", "rows": [["one"]]}`, `cell "one"`},
		{`{"kind": "table", "rows": [["123456"]]}`, `cell "123456"`},
		{`{"kind": "bars", "bars": []}`, "1 to 4"},
		{`{"kind": "bars", "bars": [{}, {}, {}, {}, {}]}`, "1 to 4"},
		{`{"kind": "bars", "bars": [{}], "notes": ["A = 1", "A = 1", "A = 1", "A = 1"]}`, "at most 3 notes"},
		{`{"kind": "bars", "bars": [{}], "notes": ["A plus 1"]}`, `note "A plus 1"`},
		{`{"kind": "bars", "bars": [{}], "notes": ["A + B + C + D + E + F = 10"]}`, "at most 24"},
		{`{"kind": "bars", "bars": [{"parts": 2, "segments": [{"size": 1}]}]}`, "not both"},
		{`{"kind": "bars", "bars": [{"parts": 13}]}`, "1 to 12 parts"},
		{`{"kind": "bars", "bars": [{"parts": 0}]}`, "1 to 12 parts"},
		{`{"kind": "bars", "bars": [{"segments": []}]}`, "1 to 12 segments"},
		{`{"kind": "bars", "bars": [{"shaded": 1}]}`, "only a bar cut into parts"},
		{`{"kind": "bars", "bars": [{"parts": 2, "shaded": 3}]}`, "shades 3 of its 2 parts"},
		{`{"kind": "bars", "bars": [{"parts": 2, "shaded": -1}]}`, "shades -1"},
		{`{"kind": "bars", "bars": [{"length": 61}]}`, "from 1 to 60"},
		{`{"kind": "bars", "bars": [{"segments": [{"size": 0}]}]}`, "size of a segment"},
		{`{"kind": "bars", "bars": [{"segments": [{"label": "A"}]}]}`, "size of a segment"},
		{`{"kind": "bars", "bars": [{"parts": 2, "braces": [{"from": 0, "to": 3, "label": "A"}]}]}`, "outside its bar"},
		{`{"kind": "bars", "bars": [{"parts": 2, "braces": [{"from": 1, "to": 1, "label": "A"}]}]}`, "outside its bar"},
		{`{"kind": "bars", "bars": [{"parts": 2, "braces": [{"from": 0, "to": 1}]}]}`, "from, to and a label"},
		{`{"kind": "bars", "bars": [{"parts": 2, "braces": [{"from": 0, "to": 1, "label": "ab"}]}]}`, `brace is none`},
		{`{"kind": "bars", "bars": [{"parts": 4, "braces": [{"from": 0, "to": 1, "label": "A"}, {"from": 0, "to": 1, "label": "A"},
			{"from": 0, "to": 1, "label": "A"}, {"from": 0, "to": 1, "label": "A"}, {"from": 0, "to": 1, "label": "A"}]}]}`, "at most 4 braces"},
		{`{"kind": "bars", "bars": [{"value": "ab"}]}`, `label "ab"`},
		{`{"kind": "bars", "bars": [{"segments": [{"size": 1, "label": "ab"}]}]}`, `label "ab"`},
		{`{"kind": "bars", "bars": [{"colour": "red"}]}`, "unknown field"},
		{`{"kind": "number_line", "from": 0}`, "has from and to"},
		{`{"kind": "number_line", "from": -10000, "to": 0}`, "from -9999 to 99999"},
		{`{"kind": "number_line", "from": 0, "to": 100000}`, "from -9999 to 99999"},
		{`{"kind": "number_line", "from": 0, "to": 6, "step": 0}`, "step of a number line"},
		{`{"kind": "number_line", "from": 6, "to": 6}`, "is not greater than from"},
		{`{"kind": "number_line", "from": 0, "to": 7, "step": 2}`, "whole steps"},
		{`{"kind": "number_line", "from": 0, "to": 16}`, "2 to 16 ticks"},
		{`{"kind": "number_line", "from": 0, "to": 6, "marks": [{"label": "A"}]}`, "has no at"},
		{`{"kind": "number_line", "from": 0, "to": 6, "step": 2, "marks": [{"at": 3}]}`, "stands on no tick"},
		{`{"kind": "number_line", "from": 0, "to": 6, "marks": [{"at": 7}]}`, "stands on no tick"},
		{`{"kind": "number_line", "from": 0, "to": 6, "marks": [{"at": 1}, {"at": 1}]}`, "two marks"},
		{`{"kind": "number_line", "from": 0, "to": 6, "marks": [{"at": 1, "label": "ab"}]}`, `label "ab"`},
		{`{"kind": "number_line", "from": 0, "to": 15, "marks": [{"at": 0}, {"at": 1}, {"at": 2}, {"at": 3}, {"at": 4}, {"at": 5}, {"at": 6},
			{"at": 7}, {"at": 8}, {"at": 9}, {"at": 10}, {"at": 11}, {"at": 12}, {"at": 13}, {"at": 14}, {"at": 15}, {"at": 15}]}`, "at most 16 marks"},
		{`{"kind": "number_line", "from": 0.5, "to": 6}`, "cannot unmarshal"},
		{`{"kind": "calendar", "days": 30}`, "weekday of the 1st"},
		{`{"kind": "calendar", "first": 8, "days": 30}`, "weekday of the 1st"},
		{`{"kind": "calendar", "first": 1}`, "28 to 31 days"},
		{`{"kind": "calendar", "first": 1, "days": 32}`, "28 to 31 days"},
		{`{"kind": "calendar", "first": 1, "days": 30, "week_starts": "friday"}`, "monday or sunday"},
		{`{"kind": "calendar", "first": 1, "days": 30, "marks": {"01": "A"}}`, `mark "01"`},
		{`{"kind": "calendar", "first": 1, "days": 30, "marks": {"31": "A"}}`, `mark "31"`},
		{`{"kind": "calendar", "first": 1, "days": 30, "marks": {"5": "ab"}}`, `label "ab"`},
		{`{"kind": "clock", "TIME": "3:15"}`, `member "TIME"`},
		{`{"kind": "clock", "time": "3:15", "time": "9:00"}`, `member "time"`},
		{`{"kind": "bars", "bars": [{"Parts": 2}]}`, `member "Parts"`},
		{`{"kind": "clock", "time": "3:15"`, "unexpected end"},
	} {
		_, _, err := pictureOf("pic-1", json.RawMessage(test.raw), "A picture.", "")
		if err == nil || !strings.Contains(err.Error(), test.message) {
			t.Errorf("%s: error = %v, want one that says %q", test.raw, err, test.message)
		}
	}
}

// Hands that lie together have their labels on two sides of them, and hands apart each on the side
// away from the other.
func TestClockLabelSides(t *testing.T) {
	if awayFrom(0, 0) != -90 || awayFrom(0, 90) != -90 || awayFrom(90, 0) != 90 {
		t.Error("a label stands on the side of the other hand")
	}
}
