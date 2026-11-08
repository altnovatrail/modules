package main

import "fmt"

// pilesPicture is the piles of a game, as the schema of standalone's tasks has it: 1 to 6 piles,
// one under another or, across, side by side, all in one box or not, and a skip among them as in a
// row. A pile is a row of counters of one shape and fill, with a gap after every group of them; it
// is cut short to the counters it shows, half before an ellipsis and half after; a pile with no
// count is a heap of no stated size; it has a label at its start and a value at its end, and a box
// of its own where it says so. How many counters a pile has is never written: the counters show it.
type pilesPicture struct {
	Kind   string     `json:"kind"`
	Piles  []pileItem `json:"piles"`
	Box    bool       `json:"box,omitempty"`
	Across bool       `json:"across,omitempty"`
}

type pileItem struct {
	Skip  bool   `json:"skip,omitempty"`
	Label string `json:"label,omitempty"`
	Value string `json:"value,omitempty"`
	Count *int   `json:"count,omitempty"`
	Shown *int   `json:"shown,omitempty"`
	Group *int   `json:"group,omitempty"`
	Fill  string `json:"fill,omitempty"`
	Shape string `json:"shape,omitempty"`
	Boxed bool   `json:"boxed,omitempty"`
}

func (p *pilesPicture) valid() error {
	if len(p.Piles) < 1 || len(p.Piles) > 6 {
		return fmt.Errorf("piles are 1 to 6, not %d", len(p.Piles))
	}
	skips := make([]bool, len(p.Piles))
	for i, pile := range p.Piles {
		skips[i] = pile.Skip
		if err := pile.valid(); err != nil {
			return fmt.Errorf("the pile %d: %w", i+1, err)
		}
	}
	return skipsPlaced(skips)
}

func (pile pileItem) valid() error {
	if pile.Skip {
		if pile != (pileItem{Skip: true}) {
			return fmt.Errorf("it is a skip, which holds nothing else")
		}
		return nil
	}
	switch {
	case pile.Count != nil && (*pile.Count < 0 || *pile.Count > 40):
		return fmt.Errorf("its count is from 0 to 40, not %d", *pile.Count)
	case pile.Shown != nil && pile.Count == nil:
		return fmt.Errorf("only a pile of a stated size is cut short")
	case pile.Shown != nil && (*pile.Shown < 2 || *pile.Shown > *pile.Count-1):
		return fmt.Errorf("it shows from 2 to one less than its count, not %d", *pile.Shown)
	case pile.Group != nil && (*pile.Group < 1 || *pile.Group > 40):
		return fmt.Errorf("its group is from 1 to 40, not %d", *pile.Group)
	case pile.Fill != "" && pile.Fill != "dark" && pile.Fill != "light":
		return fmt.Errorf("its fill is dark or light, not %q", pile.Fill)
	case pile.Shape != "" && pile.Shape != "dot" && pile.Shape != "square":
		return fmt.Errorf("its shape is dot or square, not %q", pile.Shape)
	}
	for _, text := range []string{pile.Label, pile.Value} {
		if text != "" && !isLabel(text) {
			return fmt.Errorf("the label %q is none", text)
		}
	}
	return nil
}

// skipsPlaced holds the skips among the items of a row or of piles: never first or last, and never
// next to another.
func skipsPlaced(skips []bool) error {
	for i, skip := range skips {
		if skip && (i == 0 || i == len(skips)-1 || skips[i-1]) {
			return fmt.Errorf("the skip %d is first, last or next to another", i+1)
		}
	}
	return nil
}

// The geometry of piles, in CSS px: a counter's radius, how far apart two counters stand, the gap
// between groups, how wide an ellipsis is, how tall a line of counters is, the room between piles,
// the room around what stands beside a pile, the room inside a box, and how wide a pile across is.
const (
	counterRadius = 7
	counterPitch  = 18
	groupGap      = 8
	ellipsisWide  = 24
	pileLine      = 24
	pileGap       = 12
	pileAside     = 10
	boxRoom       = 8
	columnWide    = 6 * counterPitch
	widestPile    = 20 * counterPitch
)

// token is one thing of a pile laid out in a line: a counter, or an ellipsis where it is cut short.
type token struct {
	ellipsis bool
	x        float64
	line     int
}

// tokens are a pile's counters and ellipses laid out in lines no wider than given, from 0: a pile
// cut short shows half its counters, an ellipsis and the other half; a heap of no stated size shows
// a counter, an ellipsis and a counter; any other pile every counter, with a gap after every group.
// A line breaks before a group that would pass the width, and within a group wider than a line
// before a counter that would. It gives how many lines they take and how far the widest reaches.
func (pile pileItem) tokens(width float64) ([]token, int, float64) {
	var kinds []bool // true for an ellipsis
	var gaps []bool  // a gap before the token
	switch {
	case pile.Count == nil:
		kinds = []bool{false, true, false}
		gaps = []bool{false, false, false}
	case pile.Shown != nil:
		before, after := (*pile.Shown+1)/2, *pile.Shown/2
		for i := range before + 1 + after {
			kinds = append(kinds, i == before)
			gaps = append(gaps, false)
		}
	default:
		for i := range *pile.Count {
			kinds = append(kinds, false)
			gaps = append(gaps, i > 0 && pile.Group != nil && i%*pile.Group == 0)
		}
	}
	var laid []token
	x, line, reach := 0.0, 0, 0.0
	for i, ellipsis := range kinds {
		wide := float64(counterPitch)
		if ellipsis {
			wide = ellipsisWide
		}
		if gaps[i] {
			group := float64(min(*pile.Group, *pile.Count-i) * counterPitch)
			if x+groupGap+group > width && group <= width {
				x, line = 0, line+1
			} else {
				x += groupGap
			}
		}
		if x > 0 && x+wide > width {
			x, line = 0, line+1
		}
		laid = append(laid, token{ellipsis: ellipsis, x: x, line: line})
		x += wide
		reach = max(reach, x)
	}
	return laid, line + 1, reach
}

// emptyWide is how wide the place of a pile of no counters is.
const emptyWide = 40

// pileBox is how a pile is laid out: its counters, how many lines they take, how wide they reach.
type pileBox struct {
	tokens []token
	lines  int
	reach  float64
}

func (pile pileItem) box(width float64) pileBox {
	if pile.Count != nil && *pile.Count == 0 {
		return pileBox{lines: 1, reach: emptyWide}
	}
	tokens, lines, reach := pile.tokens(width)
	return pileBox{tokens: tokens, lines: lines, reach: reach}
}

// beside is how wide the widest label, or the widest value, of the piles is.
func (p *pilesPicture) beside() (labels, values float64) {
	for _, pile := range p.Piles {
		if pile.Label != "" {
			labels = max(labels, measure(pile.Label, fontSize)+pileAside)
		}
		if pile.Value != "" {
			values = max(values, measure(pile.Value, fontSize)+pileAside)
		}
	}
	return labels, values
}

func (p *pilesPicture) room() float64 {
	if p.Box {
		return boxRoom
	}
	return 0
}

func (p *pilesPicture) size() (int, int) {
	room, own := p.room(), p.ownRoom()
	if p.Across {
		width, height := 0.0, 0.0
		for i, pile := range p.Piles {
			if i > 0 {
				width += pileGap
			}
			if pile.Skip {
				width += ellipsisWide
				continue
			}
			b := pile.box(columnWide)
			width += max(b.reach, measure(pile.Label, fontSize), measure(pile.Value, fontSize)) + 2*own
			height = max(height, float64(b.lines*pileLine)+2*own)
		}
		return 2*margin + int(width+2*room+0.5), 2*margin/2 + int(2*band+height+2*room+0.5)
	}
	labels, values := p.beside()
	width, height := 0.0, 0.0
	for i, pile := range p.Piles {
		if i > 0 {
			height += pileGap
		}
		if pile.Skip {
			height += ellipsisWide
			continue
		}
		b := pile.box(widestPile)
		width = max(width, b.reach+2*own)
		height += float64(b.lines*pileLine) + 2*own
	}
	return 2*margin/2 + int(labels+width+values+2*room+0.5), 2*margin/2 + int(height+2*room+0.5)
}

// ownRoom is the room inside a pile's own box, where any pile has one.
func (p *pilesPicture) ownRoom() float64 {
	for _, pile := range p.Piles {
		if pile.Boxed {
			return boxRoom
		}
	}
	return 0
}

func (p *pilesPicture) paint(c canvas) {
	width, height := c.size()
	room, own := p.room(), p.ownRoom()
	if p.Box {
		c.frame(margin/2, margin/2, width-margin/2, height-margin/2, 1.5)
	}
	if p.Across {
		x := float64(margin) + room
		for i, pile := range p.Piles {
			if i > 0 {
				x += pileGap
			}
			if pile.Skip {
				dots(c, x+ellipsisWide/2, height/2, false)
				x += ellipsisWide
				continue
			}
			b := pile.box(columnWide)
			wide := max(b.reach, measure(pile.Label, fontSize), measure(pile.Value, fontSize)) + 2*own
			top := float64(margin/2) + room + band
			if pile.Label != "" {
				c.text(pile.Label, x+wide/2, top-band/2)
			}
			pile.draw(c, b, x+(wide-b.reach-2*own)/2, top, own)
			if pile.Value != "" {
				c.text(pile.Value, x+wide/2, top+float64(b.lines*pileLine)+2*own+band/2)
			}
			x += wide
		}
		return
	}
	labels, _ := p.beside()
	y := float64(margin/2) + room
	for i, pile := range p.Piles {
		if i > 0 {
			y += pileGap
		}
		if pile.Skip {
			dots(c, float64(margin/2)+room+labels+own+counterPitch/2, y+ellipsisWide/2, true)
			y += ellipsisWide
			continue
		}
		b := pile.box(widestPile)
		left := float64(margin/2) + room + labels
		if pile.Label != "" {
			c.text(pile.Label, left-pileAside/2-measure(pile.Label, fontSize)/2, y+own+pileLine/2)
		}
		pile.draw(c, b, left, y, own)
		if pile.Value != "" {
			last := y + own + float64((b.lines-1)*pileLine) + pileLine/2
			c.text(pile.Value, left+b.reach+2*own+pileAside/2+measure(pile.Value, fontSize)/2, last)
		}
		y += float64(b.lines*pileLine) + 2*own
	}
}

// draw draws a pile's counters from its left edge and top, in its own box where it has one.
func (pile pileItem) draw(c canvas, b pileBox, left, top, own float64) {
	if pile.Boxed {
		c.frame(left, top, left+b.reach+2*own, top+float64(b.lines*pileLine)+2*own, 1.5)
	}
	left, top = left+own, top+own
	if pile.Count != nil && *pile.Count == 0 {
		for x := left; x < left+emptyWide; x += 8 {
			c.rect(x, top+3, min(x+4, left+emptyWide), top+4.5, ink)
			c.rect(x, top+pileLine-4.5, min(x+4, left+emptyWide), top+pileLine-3, ink)
		}
		c.rect(left, top+3, left+1.5, top+pileLine-3, ink)
		c.rect(left+emptyWide-1.5, top+3, left+emptyWide, top+pileLine-3, ink)
		return
	}
	for _, t := range b.tokens {
		y := top + float64(t.line*pileLine) + pileLine/2
		if t.ellipsis {
			dots(c, left+t.x+ellipsisWide/2, y, false)
			continue
		}
		pile.counter(c, left+t.x+counterPitch/2, y)
	}
}

// counter draws one counter centred on (x, y): a dot or a square, dark or, light, an outline.
func (pile pileItem) counter(c canvas, x, y float64) {
	light := pile.Fill == "light"
	if pile.Shape == "square" {
		if light {
			c.frame(x-counterRadius, y-counterRadius, x+counterRadius, y+counterRadius, 1.5)
			return
		}
		c.rect(x-counterRadius, y-counterRadius, x+counterRadius, y+counterRadius, ink)
		return
	}
	if light {
		c.ring(x, y, counterRadius-0.75, 1.5)
		return
	}
	c.disc(x, y, counterRadius)
}

// dots are three small dots in a row, or down a column, centred on (x, y).
func dots(c canvas, x, y float64, down bool) {
	for _, d := range []float64{-6, 0, 6} {
		if down {
			c.disc(x, y+d, 1.8)
		} else {
			c.disc(x+d, y, 1.8)
		}
	}
}
