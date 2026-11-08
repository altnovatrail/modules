package main

import (
	"bytes"
	"encoding/json"
	"image/png"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// The kinds of T21 draw as every kind does: opaque, deterministic, of the size their part names.
func TestKinds21Draw(t *testing.T) {
	for _, raw := range []string{
		`{"kind": "lattice", "width": 1, "height": 1, "lines": [[0, 0, 1, 0]]}`,
		`{"kind": "lattice", "width": 4, "height": 4, "lines": [[1, 0, 1, 3], [2, 1, 2, 4], [3, 0, 3, 3]],
			"marks": [{"at": [0, 0], "label": "S"}, {"at": [3, 3], "label": "E"}]}`,
		`{"kind": "lattice", "width": 8, "height": 8, "cells": [[0, 0], [7, 7]], "lines": [[0, 0, 8, 8]], "dashed": [[0, 4, 8, 4]],
			"points": [{"at": [0, 0], "label": "A"}, {"at": [8, 8], "label": "B"}, {"at": [4, 4]}, {"at": [0, 8], "label": "12"}],
			"marks": [{"at": [1, 0], "label": "?"}]}`,
		`{"kind": "tiles", "sides": 4, "layout": "row", "tiles": [["1", "2", "", ""], ["", "6", "", ""]]}`,
		`{"kind": "tiles", "sides": 6, "layout": "row", "joint": 3, "tiles": [["1", "", "", "", "", ""], ["", "", "", "", "", ""],
			["", "", "", "", "", ""], ["", "", "", "", "", "99"]]}`,
		`{"kind": "tiles", "sides": 5, "layout": "ring", "joint": 2, "tiles": [["1", "2", "3", "4", "5"], ["", "", "", "", ""],
			["", "", "", "", ""], ["", "", "", "", ""], ["", "", "", "", ""], ["", "", "", "", ""], ["", "", "", "", ""],
			["", "", "", "?", ""], ["", "", "", "", ""], ["", "", "", "", ""]]}`,
		`{"kind": "tiles", "sides": 3, "layout": "ring", "joint": 1, "tiles": [["1", "2", "3"], ["", "", ""], ["", "", ""],
			["", "", ""], ["", "", ""], ["", "", "?"]]}`,
		`{"kind": "tiles", "sides": 4, "layout": "ring", "joint": 1, "tiles": [["1", "", "", ""], ["", "", "", ""], ["", "", "", ""], ["", "", "", ""]]}`,
		`{"kind": "shapes", "rows": [{"shapes": ["triangle"], "total": "?"}]}`,
		`{"kind": "shapes", "rows": [{"shapes": ["triangle", "square", "circle", "diamond", "circle", "circle"], "total": "100"},
			{"shapes": ["diamond"], "total": "4"}, {"shapes": ["square"], "total": "3"}, {"shapes": ["circle"], "total": "2"},
			{"shapes": ["triangle"], "total": "1"}]}`,
		`{"kind": "stairs", "steps": 3}`,
		`{"kind": "stairs", "steps": 16, "marks": [{"step": 0, "label": "A"}, {"step": 16, "label": "B"}, {"step": 8, "label": "?"}],
			"jumps": [{"from": 0, "to": 16}, {"from": 16, "to": 0}, {"from": 3, "to": 5}, {"from": 8, "to": 7}]}`,
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

// A row of tiles runs straight on, and a ring goes round once and closes, its last tile against its
// first: every tile's middle is one tile's width from the next.
func TestTilesLaidOut(t *testing.T) {
	if !tilesLaidOut(4, 2, 6, false) || tilesLaidOut(4, 2, 7, false) || !tilesLaidOut(6, 3, 4, false) || tilesLaidOut(6, 3, 5, false) {
		t.Error("a row's length is not held")
	}
	if tilesLaidOut(4, 1, 4, false) || !tilesLaidOut(4, 1, 4, true) || !tilesLaidOut(5, 2, 10, true) || tilesLaidOut(5, 2, 9, true) {
		t.Error("a ring does not close as it should")
	}
	five, two := 5, 2
	tiles := &tilesPicture{Sides: &five, Joint: &two, Tiles: make([][]string, 10)}
	placed := tiles.laidOut()
	first, last := placed[0], placed[len(placed)-1]
	across := 2 * 0.8090169943749475 // 2 cos(π/5)
	if d := (first.x-last.x)*(first.x-last.x) + (first.y-last.y)*(first.y-last.y); d < across*across-1e-9 || d > across*across+1e-9 {
		t.Errorf("the ring of ten pentagons does not close: the last tile is %v from the first", d)
	}
}

// A point's label takes a cell of its own round its corner, cells of the paper first, and a label
// already seated moves where one beside it has none left.
func TestPointLabels(t *testing.T) {
	label := func(s string) *string { return &s }
	two := 2
	l := &latticePicture{Width: &two, Height: &two,
		Points: []latticeLabel{{At: []int{1, 1}, Label: label("A")}, {At: []int{1, 0}, Label: label("B")}, {At: []int{0, 0}}},
		Marks:  []latticeLabel{{At: []int{1, 1}, Label: label("M")}}}
	seats := l.pointLabels()
	if len(seats) != 2 || seats[0] == seats[1] || seats[0] == [2]int{1, 1} {
		t.Errorf("seats %v", seats)
	}
	if meetsBetweenCorners([4]int{0, 0, 2, 2}, [4]int{0, 2, 2, 0}) {
		t.Error("crossing at the corner [1, 1] is between corners")
	}
	if !meetsBetweenCorners([4]int{0, 0, 1, 1}, [4]int{0, 1, 1, 0}) || !meetsBetweenCorners([4]int{0, 0, 2, 0}, [4]int{1, 0, 3, 0}) {
		t.Error("lines that cross between corners, or run along each other, meet only at corners")
	}
	if meetsBetweenCorners([4]int{0, 0, 1, 0}, [4]int{1, 0, 2, 0}) || meetsBetweenCorners([4]int{0, 0, 1, 0}, [4]int{0, 1, 1, 1}) {
		t.Error("lines that touch at a corner, or never meet, meet between corners")
	}
}

func TestKinds21Refuse(t *testing.T) {
	for _, test := range []struct{ raw, message string }{
		{`{"kind": "lattice", "width": 0, "height": 1}`, "1 to 8 cells across"},
		{`{"kind": "lattice", "width": 1}`, "1 to 8 cells across"},
		{`{"kind": "lattice", "width": 2, "height": 2, "cells": [[2, 0]]}`, "no cell of the lattice"},
		{`{"kind": "lattice", "width": 2, "height": 2, "cells": [[1]]}`, "no cell of the lattice"},
		{`{"kind": "lattice", "width": 2, "height": 2, "cells": [[1, 1], [1, 1]]}`, "filled twice"},
		{`{"kind": "lattice", "width": 2, "height": 2, "lines": [[0, 0, 3, 0]]}`, "does not join two corners"},
		{`{"kind": "lattice", "width": 2, "height": 2, "lines": [[0, 0, 0]]}`, "does not join two corners"},
		{`{"kind": "lattice", "width": 2, "height": 2, "lines": [[1, 1, 1, 1]]}`, "joins a corner to itself"},
		{`{"kind": "lattice", "width": 2, "height": 2, "lines": [[0, 0, 1, 0]], "dashed": [[1, 0, 0, 0]]}`, "drawn twice"},
		{`{"kind": "lattice", "width": 2, "height": 2, "lines": [[0, 0, 1, 1], [0, 1, 1, 0]]}`, "meets another between corners"},
		{`{"kind": "lattice", "width": 2, "height": 2, "points": [{"at": [3, 0]}]}`, "stands outside"},
		{`{"kind": "lattice", "width": 2, "height": 2, "points": [{"at": [1, 1]}, {"at": [1, 1]}]}`, "two labels stand"},
		{`{"kind": "lattice", "width": 2, "height": 2, "marks": [{"at": [0, 0]}]}`, "has a label"},
		{`{"kind": "lattice", "width": 2, "height": 2, "marks": [{"at": [2, 0], "label": "A"}]}`, "stands outside"},
		{`{"kind": "lattice", "width": 2, "height": 2, "marks": [{"at": [0, 0], "label": "ABC"}]}`, "at most 2 characters"},
		{`{"kind": "lattice", "width": 2, "height": 2, "points": [{"at": [1, 1], "label": "A"}], "marks": [{"at": [0, 0], "label": "M"},
			{"at": [1, 0], "label": "M"}, {"at": [0, 1], "label": "M"}, {"at": [1, 1], "label": "M"}]}`, "has no cell of its own"},
		{`{"kind": "lattice", "width": 2, "height": 2, "cells": [` + strings.Repeat(`[0, 0], `, 64) + `[0, 0]]}`, "at most 64 cells"},
		{`{"kind": "lattice", "width": 2, "height": 2, "dashed": [` + strings.Repeat(`[0, 0, 1, 0], `, 8) + `[0, 0, 1, 0]]}`, "at most 48 lines"},
		{`{"kind": "lattice", "width": 2, "height": 2, "points": [` + strings.Repeat(`{"at": [0, 0]}, `, 16) + `{"at": [0, 0]}]}`, "at most 16 points"},
		{`{"kind": "tiles", "sides": 2, "layout": "row", "tiles": [[], []]}`, "3 to 6 sides"},
		{`{"kind": "tiles", "sides": 4, "layout": "fan", "tiles": [[], []]}`, "a row or a ring"},
		{`{"kind": "tiles", "sides": 4, "layout": "row", "joint": 4, "tiles": [[], []]}`, "from 1 to 3"},
		{`{"kind": "tiles", "sides": 4, "layout": "row", "tiles": [["", "", "", ""]]}`, "2 to 12"},
		{`{"kind": "tiles", "sides": 4, "layout": "row", "tiles": [["", "", ""], ["", "", "", ""]]}`, "as many triangles"},
		{`{"kind": "tiles", "sides": 4, "layout": "row", "tiles": [["100", "", "", ""], ["", "", "", ""]]}`, `triangle "100"`},
		{`{"kind": "tiles", "sides": 4, "layout": "row", "joint": 1, "tiles": [["", "", "", ""], ["", "", "", ""]]}`, "runs straight"},
		{`{"kind": "tiles", "sides": 4, "layout": "ring", "joint": 1, "tiles": [["", "", "", ""], ["", "", "", ""]]}`, "closes"},
		{`{"kind": "shapes", "rows": []}`, "1 to 5 rows"},
		{`{"kind": "shapes", "rows": [{"shapes": [], "total": "1"}]}`, "1 to 6 shapes"},
		{`{"kind": "shapes", "rows": [{"shapes": ["star"], "total": "1"}]}`, `shape "star"`},
		{`{"kind": "shapes", "rows": [{"shapes": ["circle"]}]}`, "has a total"},
		{`{"kind": "stairs", "steps": 2}`, "3 to 16 steps"},
		{`{"kind": "stairs", "steps": 3, "marks": [{"step": 4, "label": "A"}]}`, "on no step"},
		{`{"kind": "stairs", "steps": 3, "marks": [{"step": 1, "label": "A"}, {"step": 1, "label": "B"}]}`, "two marks"},
		{`{"kind": "stairs", "steps": 3, "marks": [{"step": 1}]}`, "has a label"},
		{`{"kind": "stairs", "steps": 3, "jumps": [{"from": 0, "to": 4}]}`, "between steps"},
		{`{"kind": "stairs", "steps": 3, "jumps": [{"from": 1, "to": 1}]}`, "from one step to another"},
		{`{"kind": "stairs", "steps": 3, "jumps": [{"from": 0, "to": 1}, {"from": 0, "to": 1}]}`, "drawn twice"},
		{`{"kind": "stairs", "steps": 3, "jumps": [{"from": 0, "to": 1}, {"from": 0, "to": 2}, {"from": 0, "to": 3}, {"from": 1, "to": 2},
			{"from": 1, "to": 3}]}`, "at most 4 marks and 4 jumps"},
	} {
		_, _, err := pictureOf("pic-1", json.RawMessage(test.raw), "A picture.", "")
		if err == nil || !strings.Contains(err.Error(), test.message) {
			t.Errorf("%s: error = %v, want one that says %q", test.raw, err, test.message)
		}
	}
}

// The texts each kind writes, which the alt of a task solved from the picture names.
func TestWrittenLabels(t *testing.T) {
	for raw, want := range map[string]string{
		`{"kind": "row", "items": [{"label": "A"}, {"below": "1"}], "gaps": "2", "span": "?"}`: "2 A 1",
		`{"kind": "clock", "time": "7:30", "hour_label": "H"}`:                                 "H",
		`{"kind": "table", "header": ["A"], "rows": [["…"], [""], ["1"]]}`:                     "A … 1",
		`{"kind": "bars", "bars": [{"label": "N", "segments": [{"size": 1, "label": "6"}], "braces": [{"from": 0, "to": 1, "label": "B"}]}],
			"notes": ["N = 6"]}`: "N = 6 N 6 B",
		`{"kind": "number_line", "from": 0, "to": 3, "marks": [{"at": 1, "label": "P"}, {"at": 2}]}`: "P",
		`{"kind": "calendar", "first": 1, "days": 30, "marks": {"4": "1", "11": "2"}}`:               "2 1",
		`{"kind": "containers", "items": [{"capacity": 5, "amount": 2, "label": "A"}]}`:              "A 5 2",
		`{"kind": "grid", "rows": ["A"], "cols": ["1", "2"], "marks": {"A2": "9", "A1": "8"}}`:       "8 9",
		`{"kind": "piles", "piles": [{"label": "1", "value": "2"}, {"count": 3}]}`:                   "1 2",
		`{"kind": "venn", "sets": [{"label": "F", "count": "15"}, {"label": "C"}], "both": "5"}`:     "5 F 15 C",
		`{"kind": "balance", "left": ["?"], "right": ["1", "3"]}`:                                    "1 3",
		`{"kind": "ring", "count": 5, "start": {"at": 1, "label": "S"}}`:                             "S",
		`{"kind": "ring", "count": 5}`: "",
		`{"kind": "lattice", "width": 2, "height": 2, "points": [{"at": [0, 0], "label": "A"}, {"at": [2, 2]}], "marks": [{"at": [1, 1], "label": "7"}]}`: "A 7",
		`{"kind": "tiles", "sides": 4, "layout": "row", "tiles": [["1", "", "?", ""], ["", "6", "1", ""]]}`:                                               "1 6",
		`{"kind": "shapes", "rows": [{"shapes": ["circle"], "total": "9"}, {"shapes": ["circle"], "total": "?"}]}`:                                        "9",
		`{"kind": "stairs", "steps": 9, "marks": [{"step": 0, "label": "F"}]}`:                                                                            "F",
	} {
		said, err := writtenLabels("pic-1", json.RawMessage(raw))
		if got := strings.Join(said, " "); got != want || err != nil {
			t.Errorf("%s: writes %q, %v, want %q", raw, got, err, want)
		}
	}
}

// The alt of a task solved from the picture names every text the picture writes, and no alt names
// the key unless the question does (R20; T21).
func TestConvertAltOfPicture(t *testing.T) {
	ex := example("fig-1")
	ex.Picture = json.RawMessage(`{"kind": "lattice", "width": 4, "height": 4, "lines": [[1, 0, 1, 3]],
		"marks": [{"at": [0, 0], "label": "S"}, {"at": [3, 3], "label": "E"}]}`)
	ex.fromPicture = true
	good := Overlay{PictureAlt: "A maze of 4 by 4 cells, S in the top left cell and E in the bottom right one; a wall runs down from the top along the right side of the first column for 3 cells."}
	if _, err := convert(ex, testBatch(ex.ID), metres, good, "content/examples/x.json"); err != nil {
		t.Errorf("an alt that names every label: %v", err)
	}
	short := Overlay{PictureAlt: "A maze of 4 by 4 cells with S in the top left cell."}
	if _, err := convert(ex, testBatch(ex.ID), metres, short, "content/examples/x.json"); err == nil ||
		!strings.Contains(err.Error(), `does not name "E", which the picture writes`) {
		t.Errorf("an alt that leaves a label out: error = %v", err)
	}
	ex.fromPicture = false
	if _, err := convert(ex, testBatch(ex.ID), metres, short, "content/examples/x.json"); err != nil {
		t.Errorf("a task not solved from the picture: %v", err)
	}
	keyed := Overlay{PictureAlt: "A maze; the way from S to E takes 14 moves."}
	if _, err := convert(ex, testBatch(ex.ID), metres, keyed, "content/examples/x.json"); err == nil ||
		!strings.Contains(err.Error(), `picture_alt names its key "14"`) {
		t.Errorf("an alt that names the key: error = %v", err)
	}
	ex.Question = "Is the way from S to E 14 moves long, or how many?"
	if _, err := convert(ex, testBatch(ex.ID), metres, keyed, "content/examples/x.json"); err != nil {
		t.Errorf("an alt that names a key the question names: %v", err)
	}
	batch := testBatch(ex.ID)
	batch.survey = true
	ex.fromPicture = true
	if _, err := convert(ex, batch, metres, Overlay{}, "content/examples/x.json"); err != nil {
		t.Errorf("a survey: %v", err)
	}
}

// A task is solved from the picture when its solver, content/examples/solvers/<id>.star, reads the
// picture; a task with no such file reads none.
func TestSolvedFromPicture(t *testing.T) {
	clone := t.TempDir()
	dir := filepath.Join(clone, "content", "examples", "solvers")
	writeFile(t, filepath.Join(dir, "fig-1.star"), "# From the picture.\ndef solve(options, picture):\n    return [\"A\"]\n")
	writeFile(t, filepath.Join(dir, "lamps-1.star"), "def solve(options):\n    return [\"C\"]\n")
	for id, want := range map[string]bool{"fig-1": true, "lamps-1": false, "none-1": false} {
		if got, err := solvedFromPicture(clone, id); err != nil || got != want {
			t.Errorf("%s: %v, %v, want %v", id, got, err, want)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "dir-1.star"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := solvedFromPicture(clone, "dir-1"); err == nil {
		t.Error("a solver that cannot be read: no error")
	}
}

// The review of T21: a picture that writes its key lets the alt name it, with a note; a label of
// capitals is named as written; a solver reads the picture whatever its head's spacing.
func TestAltAfterReview(t *testing.T) {
	ex := example("tiles-1")
	ex.Picture = json.RawMessage(`{"kind": "tiles", "sides": 4, "layout": "row", "tiles": [["14", "", "", ""], ["", "", "", ""]]}`)
	ex.fromPicture = true
	c, err := convert(ex, testBatch(ex.ID), metres, Overlay{PictureAlt: "Two square tiles in a row; the first has 14 in its left triangle."}, "content/examples/x.json")
	if err != nil || !strings.Contains(strings.Join(c.Notes, "\n"), `its picture writes its key "14"`) {
		t.Errorf("a picture that writes its key: %v, %v", err, c.Notes)
	}

	ex = example("maze-1")
	ex.Picture = json.RawMessage(`{"kind": "lattice", "width": 2, "height": 2, "lines": [[1, 0, 1, 1]], "marks": [{"at": [0, 0], "label": "A"}]}`)
	ex.fromPicture = true
	if _, err := convert(ex, testBatch(ex.ID), metres, Overlay{PictureAlt: "The maze has 2 by 2 cells and a wall."}, "content/examples/x.json"); err == nil ||
		!strings.Contains(err.Error(), `does not name "A"`) {
		t.Errorf("the word a for the label A: error = %v", err)
	}
	if _, err := convert(ex, testBatch(ex.ID), metres, Overlay{PictureAlt: "A maze of 2 by 2 cells, A in the top left cell."}, "content/examples/x.json"); err != nil {
		t.Errorf("the label A named: %v", err)
	}
	if namesCased("a maze", "A") || !namesCased("cell A here", "A") || !names("a maze", "A") {
		t.Error("namesCased tells no capitals from small letters")
	}

	ex.Picture = json.RawMessage(`{"kind": "lattice", "width": 0}`)
	if _, err := altOf(ex, Overlay{PictureAlt: "x"}, "1", false); err == nil {
		t.Error("an invalid picture: no error")
	}

	clone := t.TempDir()
	writeFile(t, filepath.Join(clone, "content", "examples", "solvers", "a-1.star"), "def  solve( options ,pic ):\n    pass\n")
	if got, err := solvedFromPicture(clone, "a-1"); !got || err != nil {
		t.Errorf("a solver head spaced otherwise: %v, %v", got, err)
	}
}
