package main

import (
	"fmt"
	"slices"
)

// latticePicture is a figure on squared paper, as the schema of standalone's tasks has it: width
// cells across and height down, 1 to 8 each, a cell [x, y] from [0, 0] at the top left; cells
// filled, lines from corner to corner — walls, sides, matches, sticks — dashed lines — a fold, a
// cut, an axis — points at corners, each with a label or none, and labels in cells. Its facts may
// stand on the picture alone (docs/architecture.md, section 3.5).
type latticePicture struct {
	Kind   string         `json:"kind"`
	Width  *int           `json:"width"`
	Height *int           `json:"height"`
	Cells  [][]int        `json:"cells,omitempty"`
	Lines  [][]int        `json:"lines,omitempty"`
	Dashed [][]int        `json:"dashed,omitempty"`
	Points []latticeLabel `json:"points,omitempty"`
	Marks  []latticeLabel `json:"marks,omitempty"`
}

type latticeLabel struct {
	At    []int   `json:"at"`
	Label *string `json:"label,omitempty"`
}

// isShortLabel tells a label of a lattice, at most 2 characters long, from any other text.
func isShortLabel(text string) bool { return isLabel(text) && len([]rune(text)) <= 2 }

func (l *latticePicture) valid() error {
	if l.Width == nil || l.Height == nil || *l.Width < 1 || *l.Width > 8 || *l.Height < 1 || *l.Height > 8 {
		return fmt.Errorf("a lattice is 1 to 8 cells across and 1 to 8 down")
	}
	w, h := *l.Width, *l.Height
	within := func(pair []int, x, y int) bool {
		return len(pair) == 2 && pair[0] >= 0 && pair[0] <= x && pair[1] >= 0 && pair[1] <= y
	}
	switch {
	case len(l.Cells) > 64:
		return fmt.Errorf("a lattice fills at most 64 cells")
	case len(l.Lines) > 48 || len(l.Dashed) > 8:
		return fmt.Errorf("a lattice has at most 48 lines and 8 dashed lines")
	case len(l.Points) > 16 || len(l.Marks) > 36:
		return fmt.Errorf("a lattice has at most 16 points and 36 marks")
	}
	var filled [][2]int
	for _, cell := range l.Cells {
		if !within(cell, w-1, h-1) {
			return fmt.Errorf("the cell %v is no cell of the lattice: [x, y], from [0, 0] at the top left", cell)
		}
		if slices.Contains(filled, [2]int{cell[0], cell[1]}) {
			return fmt.Errorf("the cell %v is filled twice", cell)
		}
		filled = append(filled, [2]int{cell[0], cell[1]})
	}
	var drawn [][4]int
	for _, lines := range [][][]int{l.Lines, l.Dashed} {
		for _, line := range lines {
			if len(line) != 4 || !within(line[:2], w, h) || !within(line[2:], w, h) {
				return fmt.Errorf("the line %v does not join two corners of the lattice: [x1, y1, x2, y2], from [0, 0] at the top left to [width, height]", line)
			}
			one, other := [4]int(line), [4]int{line[2], line[3], line[0], line[1]}
			switch {
			case one[0] == one[2] && one[1] == one[3]:
				return fmt.Errorf("the line %v joins a corner to itself", line)
			case slices.Contains(drawn, one) || slices.Contains(drawn, other):
				return fmt.Errorf("the line %v is drawn twice", line)
			}
			drawn = append(drawn, one)
		}
	}
	for i, line := range l.Lines {
		for _, earlier := range l.Lines[:i] {
			if meetsBetweenCorners([4]int(line), [4]int(earlier)) {
				return fmt.Errorf("the line %v meets another between corners of the lattice: lines may meet only at its corners, and never run along each other", line)
			}
		}
	}
	for _, set := range []struct {
		labels []latticeLabel
		x, y   int
		needed bool
	}{{l.Points, w, h, false}, {l.Marks, w - 1, h - 1, true}} {
		var taken [][2]int
		for _, item := range set.labels {
			switch {
			case !within(item.At, set.x, set.y):
				return fmt.Errorf("the label at %v stands outside the lattice", item.At)
			case slices.Contains(taken, [2]int{item.At[0], item.At[1]}):
				return fmt.Errorf("two labels stand at %v", item.At)
			case set.needed && item.Label == nil:
				return fmt.Errorf("a mark of a lattice has a label")
			case item.Label != nil && !isShortLabel(*item.Label):
				return fmt.Errorf("the label %q of a lattice is no label of at most 2 characters", *item.Label)
			}
			taken = append(taken, [2]int{item.At[0], item.At[1]})
		}
	}
	if l.pointLabels() == nil && slices.ContainsFunc(l.Points, func(p latticeLabel) bool { return p.Label != nil && *p.Label != "" }) {
		return fmt.Errorf("a point's label has no cell of its own among the four round its corner that no mark holds")
	}
	return nil
}

// meetsBetweenCorners says whether two lines share a point that is no corner of the lattice, or
// share a stretch of their length, as standalone's service judges it.
func meetsBetweenCorners(one, other [4]int) bool {
	px, py := one[0], one[1]
	rx, ry := one[2]-one[0], one[3]-one[1]
	qx, qy := other[0], other[1]
	sx, sy := other[2]-other[0], other[3]-other[1]
	d := rx*sy - ry*sx
	wx, wy := qx-px, qy-py
	if d == 0 {
		if wx*ry-wy*rx != 0 {
			return false
		}
		t0, t1 := wx*rx+wy*ry, (wx+sx)*rx+(wy+sy)*ry
		return max(min(t0, t1), 0) < min(max(t0, t1), rx*rx+ry*ry)
	}
	t, u := wx*sy-wy*sx, wx*ry-wy*rx
	if d < 0 {
		d, t, u = -d, -t, -u
	}
	if t < 0 || t > d || u < 0 || u > d {
		return false
	}
	return (px*d+t*rx)%d != 0 || (py*d+t*ry)%d != 0
}

// quarters are the four cells round a corner, as steps from the cell below and to the right of it,
// in the order a point's label is given them: above right, above left, below right, below left.
var quarters = [4][2]int{{0, -1}, {-1, -1}, {0, 0}, {-1, 0}}

// pointLabels gives each labelled point a cell of its own among the four round its corner, one no
// mark holds, cells of the paper first, moving a label already seated where one beside it has none
// left: standalone's own way of setting them. It is nil when some label has no cell, and empty when
// no point has a label.
func (l *latticePicture) pointLabels() [][2]int {
	marked := map[[2]int]bool{}
	for _, mark := range l.Marks {
		marked[[2]int{mark.At[0], mark.At[1]}] = true
	}
	var choices [][][2]int
	for _, point := range l.Points {
		if point.Label == nil || *point.Label == "" {
			continue
		}
		var on, off [][2]int
		for _, q := range quarters {
			cell := [2]int{point.At[0] + q[0], point.At[1] + q[1]}
			switch {
			case marked[cell]:
			case cell[0] >= 0 && cell[1] >= 0 && cell[0] < *l.Width && cell[1] < *l.Height:
				on = append(on, cell)
			default:
				off = append(off, cell)
			}
		}
		choices = append(choices, slices.Concat(on, off))
	}
	holder, chosen := map[[2]int]int{}, make([][2]int, len(choices))
	var seat func(i int, tried map[[2]int]bool) bool
	seat = func(i int, tried map[[2]int]bool) bool {
		for _, cell := range choices[i] {
			if tried[cell] {
				continue
			}
			tried[cell] = true
			held, taken := holder[cell]
			if !taken || seat(held, tried) {
				holder[cell] = i
				chosen[i] = cell
				return true
			}
		}
		return false
	}
	for i := range choices {
		if !seat(i, map[[2]int]bool{}) {
			return nil
		}
	}
	return chosen
}

// The geometry of a lattice, in CSS px: a cell's side at most, the widest the paper may be, and how
// far a point's label stands into its cell from the corner.
const (
	latticeCell = 40
	latticeWide = 320
	pointAside  = 11
)

func (l *latticePicture) cell() float64 {
	return min(latticeCell, float64(latticeWide)/float64(max(*l.Width, *l.Height)))
}

func (l *latticePicture) size() (int, int) {
	return 2*margin + int(float64(*l.Width)*l.cell()+0.5), 2*margin + int(float64(*l.Height)*l.cell()+0.5)
}

func (l *latticePicture) paint(c canvas) {
	side := l.cell()
	at := func(x, y int) (float64, float64) { return margin + float64(x)*side, margin + float64(y)*side }
	for _, cell := range l.Cells {
		x, y := at(cell[0], cell[1])
		c.rect(x, y, x+side, y+side, shade)
	}
	right, bottom := at(*l.Width, *l.Height)
	for x := 0; x <= *l.Width; x++ {
		px, _ := at(x, 0)
		c.rect(px-0.5, margin, px+0.5, bottom, rule)
	}
	for y := 0; y <= *l.Height; y++ {
		_, py := at(0, y)
		c.rect(margin, py-0.5, right, py+0.5, rule)
	}
	for _, line := range l.Lines {
		x0, y0 := at(line[0], line[1])
		x1, y1 := at(line[2], line[3])
		c.line(x0, y0, x1, y1, 3)
	}
	for _, line := range l.Dashed {
		x0, y0 := at(line[0], line[1])
		x1, y1 := at(line[2], line[3])
		c.dashed(x0, y0, x1, y1, stroke)
	}
	for _, mark := range l.Marks {
		x, y := at(mark.At[0], mark.At[1])
		c.text(*mark.Label, x+side/2, y+side/2)
	}
	seats, n := l.pointLabels(), 0
	for _, point := range l.Points {
		x, y := at(point.At[0], point.At[1])
		c.disc(x, y, 4)
		if point.Label == nil || *point.Label == "" {
			continue
		}
		cell := seats[n]
		n++
		dx, dy := float64(cell[0]-point.At[0])*2+1, float64(cell[1]-point.At[1])*2+1
		c.text(*point.Label, x+dx*pointAside, y+dy*pointAside)
	}
}
