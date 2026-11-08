package main

import (
	"fmt"
	"regexp"
)

// barsPicture is quantities as bars, as the schema of standalone's tasks has it: 1 to 4 bars drawn
// to one scale, and at most 3 notes under them all. A bar is cut into equal parts, some of them
// shaded, or into segments of their own sizes, or into neither; it has a label before it, a value
// at its end, braces under it and a span of the whole.
type barsPicture struct {
	Kind  string    `json:"kind"`
	Bars  []barItem `json:"bars"`
	Notes []string  `json:"notes,omitempty"`
}

type barItem struct {
	Label    string       `json:"label,omitempty"`
	Parts    *int         `json:"parts,omitempty"`
	Shaded   *int         `json:"shaded,omitempty"`
	Segments []barSegment `json:"segments,omitempty"`
	Length   *int         `json:"length,omitempty"`
	Value    string       `json:"value,omitempty"`
	Braces   []barBrace   `json:"braces,omitempty"`
	Span     string       `json:"span,omitempty"`
}

type barSegment struct {
	Size  *int   `json:"size"`
	Label string `json:"label,omitempty"`
}

type barBrace struct {
	From  *int    `json:"from"`
	To    *int    `json:"to"`
	Label *string `json:"label"`
}

// note is a note under bars: an equality or an inequality of labels, numbers and ?, joined by
// + − × ÷ = < > ( ) and spaces, at most 24 characters long.
var note = regexp.MustCompile(`^[A-Z0-9.,+−×÷()? -]*[=<>][A-Z0-9.,+−×÷=<>()? -]*$`)

func (b *barsPicture) valid() error {
	if len(b.Bars) < 1 || len(b.Bars) > 4 {
		return fmt.Errorf("bars are 1 to 4, not %d", len(b.Bars))
	}
	if len(b.Notes) > 3 {
		return fmt.Errorf("bars have at most 3 notes, not %d", len(b.Notes))
	}
	for _, text := range b.Notes {
		if !note.MatchString(text) || len([]rune(text)) > 24 {
			return fmt.Errorf("the note %q is no equality or inequality of at most 24 characters", text)
		}
	}
	for i, bar := range b.Bars {
		if err := bar.valid(); err != nil {
			return fmt.Errorf("the bar %d: %w", i+1, err)
		}
	}
	return nil
}

func (bar barItem) valid() error {
	labels := []string{bar.Label, bar.Value, bar.Span}
	switch {
	case bar.Parts != nil && bar.Segments != nil:
		return fmt.Errorf("it is cut into parts or into segments, not both")
	case bar.Parts != nil && (*bar.Parts < 1 || *bar.Parts > 12):
		return fmt.Errorf("it has 1 to 12 parts, not %d", *bar.Parts)
	case bar.Segments != nil && (len(bar.Segments) < 1 || len(bar.Segments) > 12):
		return fmt.Errorf("it has 1 to 12 segments, not %d", len(bar.Segments))
	case bar.Shaded != nil && bar.Parts == nil:
		return fmt.Errorf("only a bar cut into parts shades some of them")
	case bar.Shaded != nil && (*bar.Shaded < 0 || *bar.Shaded > *bar.Parts):
		return fmt.Errorf("it shades %d of its %d parts", *bar.Shaded, *bar.Parts)
	case bar.Length != nil && (*bar.Length < 1 || *bar.Length > 60):
		return fmt.Errorf("its length is from 1 to 60, not %d", *bar.Length)
	case len(bar.Braces) > 4:
		return fmt.Errorf("it has at most 4 braces, not %d", len(bar.Braces))
	}
	for _, segment := range bar.Segments {
		if segment.Size == nil || *segment.Size < 1 || *segment.Size > 60 {
			return fmt.Errorf("the size of a segment is from 1 to 60")
		}
		labels = append(labels, segment.Label)
	}
	for _, brace := range bar.Braces {
		if brace.From == nil || brace.To == nil || brace.Label == nil {
			return fmt.Errorf("a brace has from, to and a label")
		}
		if *brace.From < 0 || *brace.From >= *brace.To || *brace.To > bar.divisions() {
			return fmt.Errorf("the brace from %d to %d lies outside its bar of %d edges after its start", *brace.From, *brace.To, bar.divisions())
		}
		if !isLabel(*brace.Label) {
			return fmt.Errorf("the label %q of a brace is none", *brace.Label)
		}
	}
	for _, text := range labels {
		if text != "" && !isLabel(text) {
			return fmt.Errorf("the label %q is none", text)
		}
	}
	return nil
}

// divisions is how many parts or segments the bar is cut into: 1 for a bar cut into neither.
func (bar barItem) divisions() int {
	switch {
	case bar.Parts != nil:
		return *bar.Parts
	case bar.Segments != nil:
		return len(bar.Segments)
	}
	return 1
}

// units is how long a bar is in the units all the bars share: its length where it gives one, else
// the sizes of its segments or the count of its parts, and one unit for a bar cut into neither.
func (bar barItem) units() int {
	switch {
	case bar.Length != nil:
		return *bar.Length
	case bar.Segments != nil:
		sum := 0
		for _, segment := range bar.Segments {
			sum += *segment.Size
		}
		return sum
	case bar.Parts != nil:
		return *bar.Parts
	}
	return 1
}

// edges are where the edges between the bar's parts or segments stand, from 0 at its start to its
// length: segments in proportion to their sizes, parts evenly.
func (bar barItem) edges(length float64) []float64 {
	edges := []float64{0}
	if bar.Segments != nil {
		whole, reached := 0.0, 0.0
		for _, segment := range bar.Segments {
			whole += float64(*segment.Size)
		}
		for _, segment := range bar.Segments {
			reached += float64(*segment.Size)
			edges = append(edges, reached/whole*length)
		}
		return edges
	}
	parts := bar.divisions()
	for i := 1; i <= parts; i++ {
		edges = append(edges, float64(i)*length/float64(parts))
	}
	return edges
}

// braceRows lays a bar's braces in rows under it: each in the first row where it overlaps no brace
// already there.
func (bar barItem) braceRows() []int {
	rows := make([]int, len(bar.Braces))
	for i, brace := range bar.Braces {
		for row := 0; ; row++ {
			free := true
			for j := range i {
				other := bar.Braces[j]
				if rows[j] == row && *brace.From < *other.To && *other.From < *brace.To {
					free = false
				}
			}
			if free {
				rows[i] = row
				break
			}
		}
	}
	return rows
}

// The geometry of bars, in CSS px.
const (
	barTall     = 28
	widestUnit  = 28
	longestBar  = 400
	barEdge     = 16
	barBetween  = 14
	braceTall   = 30
	noteTall    = 20
	besideGap   = 8
	shortestBar = 3
)

// layoutBars is how bars are laid out: how long a unit is, and the columns before and after them.
type layoutBars struct {
	unit, before, after float64
	longest             int
}

func (b *barsPicture) layout() layoutBars {
	var l layoutBars
	for _, bar := range b.Bars {
		l.longest = max(l.longest, bar.units())
		if bar.Label != "" {
			l.before = max(l.before, measure(bar.Label, fontSize)+besideGap)
		}
		if bar.Value != "" {
			l.after = max(l.after, measure(bar.Value, fontSize)+besideGap)
		}
	}
	l.unit = min(widestUnit, float64(longestBar)/float64(l.longest))
	return l
}

// tall is how tall a bar is drawn, with its braces and its span.
func (bar barItem) tall() float64 {
	height := float64(barTall)
	rows := 0
	for _, row := range bar.braceRows() {
		rows = max(rows, row+1)
	}
	height += float64(rows * braceTall)
	if bar.Span != "" {
		height += spanBand
	}
	return height
}

func (b *barsPicture) size() (int, int) {
	l := b.layout()
	width := 2*barEdge + l.before + float64(l.longest)*l.unit + l.after
	for _, text := range b.Notes {
		width = max(width, 2*barEdge+measure(text, fontSize))
	}
	height := float64(2 * margin / 2)
	for i, bar := range b.Bars {
		if i > 0 {
			height += barBetween
		}
		height += bar.tall()
	}
	if len(b.Notes) > 0 {
		height += 6 + float64(len(b.Notes)*noteTall)
	}
	return int(width + 0.5), int(height + 0.5)
}

func (b *barsPicture) paint(c canvas) {
	l := b.layout()
	width, _ := c.size()
	start := barEdge + l.before
	y := float64(margin / 2)
	for _, bar := range b.Bars {
		length := max(float64(bar.units())*l.unit, shortestBar)
		edges := bar.edges(length)
		if bar.Shaded != nil {
			c.rect(start, y, start+edges[*bar.Shaded], y+barTall, shade)
		}
		for _, edge := range edges[1 : len(edges)-1] {
			c.rect(start+edge-0.75, y, start+edge+0.75, y+barTall, ink)
		}
		c.frame(start, y, start+length, y+barTall, stroke)
		for i, segment := range bar.Segments {
			if segment.Label != "" {
				c.text(segment.Label, start+(edges[i]+edges[i+1])/2, y+barTall/2)
			}
		}
		if bar.Label != "" {
			c.text(bar.Label, barEdge+(l.before-besideGap)/2, y+barTall/2)
		}
		if bar.Value != "" {
			c.text(bar.Value, start+length+besideGap/2+measure(bar.Value, fontSize)/2, y+barTall/2)
		}
		below := y + barTall
		rows := bar.braceRows()
		for i, brace := range bar.Braces {
			top := below + float64(rows[i]*braceTall) + 4
			from, to := start+edges[*brace.From]+2, start+edges[*brace.To]-2
			c.line(from, top, from, top+6, 1.5)
			c.line(to, top, to, top+6, 1.5)
			c.line(from, top+6, to, top+6, 1.5)
			c.line((from+to)/2, top+6, (from+to)/2, top+10, 1.5)
			c.text(*brace.Label, (from+to)/2, top+20)
		}
		for _, row := range rows {
			below = max(below, y+barTall+float64((row+1)*braceTall))
		}
		if bar.Span != "" {
			c.span(start, start+length, below+spanBand/2, bar.Span, measure(bar.Span, fontSize)/2+3, 9)
		}
		y += bar.tall() + barBetween
	}
	y += 6 - barBetween
	for _, text := range b.Notes {
		c.text(text, width/2, y+noteTall/2)
		y += noteTall
	}
}
