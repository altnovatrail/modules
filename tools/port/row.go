package main

import (
	"fmt"

	"golang.org/x/image/font"
)

// rowPicture is a picture of the kind row, as the schema of standalone's tasks has it: 2 to 12
// items, each a mark with a label over it and one under it, or a skip that stands for the middle
// of a long row cut short; a line through the marks, unless it is false; a label under every gap
// drawn; a label of the whole length; and the row once or twice, one under the other.
type rowPicture struct {
	Kind   string    `json:"kind"`
	Items  []rowItem `json:"items"`
	Line   *bool     `json:"line,omitempty"`
	Gaps   string    `json:"gaps,omitempty"`
	Span   string    `json:"span,omitempty"`
	Copies int       `json:"copies,omitempty"`
}

type rowItem struct {
	Skip  bool   `json:"skip,omitempty"`
	Label string `json:"label,omitempty"`
	Below string `json:"below,omitempty"`
	Mark  string `json:"mark,omitempty"`
}

func (r *rowPicture) valid() error {
	if len(r.Items) < 2 || len(r.Items) > 12 {
		return fmt.Errorf("a row has 2 to 12 items, not %d", len(r.Items))
	}
	if r.Copies < 0 || r.Copies > 2 {
		return fmt.Errorf("a row is drawn once or twice, not %d times", r.Copies)
	}
	for i, item := range r.Items {
		if item.Skip {
			if item.Label != "" || item.Below != "" || item.Mark != "" {
				return fmt.Errorf("the skip %d has a label or a mark", i+1)
			}
			if i == 0 || i == len(r.Items)-1 || r.Items[i-1].Skip {
				return fmt.Errorf("the skip %d is first, last or next to another", i+1)
			}
			continue
		}
		switch item.Mark {
		case "", "dot", "ring", "square":
		default:
			return fmt.Errorf("the item %d has the mark %q: dot, ring or square", i+1, item.Mark)
		}
		for _, text := range []string{item.Label, item.Below} {
			if text != "" && !isLabel(text) {
				return fmt.Errorf("the item %d has the label %q, which is none", i+1, text)
			}
		}
	}
	for _, text := range []string{r.Gaps, r.Span} {
		if text != "" && !isLabel(text) {
			return fmt.Errorf("the label %q is none", text)
		}
	}
	return nil
}

// The geometry of a row, in CSS px, beside what picture.go gives every kind.
const (
	step    = 56
	marks   = 28
	between = 14
	radius  = 6
)

// rowLayout is how a row is laid out: which bands it has, how tall one copy is, and how many
// copies there are.
type rowLayout struct {
	hasLabel, hasBelow bool
	copies, rowHeight  int
}

func (r *rowPicture) layout() rowLayout {
	l := rowLayout{copies: max(r.Copies, 1), rowHeight: marks}
	for _, item := range r.Items {
		l.hasLabel = l.hasLabel || item.Label != ""
		l.hasBelow = l.hasBelow || item.Below != ""
	}
	if l.hasLabel {
		l.rowHeight += band
	}
	if l.hasBelow {
		l.rowHeight += band
	}
	if r.Gaps != "" {
		l.rowHeight += band
	}
	return l
}

func (r *rowPicture) size() (int, int) {
	l := r.layout()
	width := 2*margin + (len(r.Items)-1)*step
	height := 2*margin/2 + l.copies*l.rowHeight + (l.copies-1)*between
	if r.Span != "" {
		height += spanBand
	}
	return width, height
}

func (r *rowPicture) paint(c canvas) {
	l := r.layout()
	_, height := c.size()
	top := float64(margin / 2)
	for n := range l.copies {
		y := top + float64(n*(l.rowHeight+between))
		if l.hasLabel {
			for i, item := range r.Items {
				if item.Label != "" {
					c.text(item.Label, x(i), y+band/2)
				}
			}
			y += band
		}
		line := y + marks/2
		r.drawLine(c, line)
		for i, item := range r.Items {
			if item.Skip {
				for _, dx := range []float64{-8, 0, 8} {
					c.disc(x(i)+dx, line, 1.5)
				}
				continue
			}
			switch item.Mark {
			case "ring":
				c.ring(x(i), line, radius-stroke/2, stroke)
			case "square":
				c.rect(x(i)-radius, line-radius, x(i)+radius, line+radius, ink)
			default:
				c.disc(x(i), line, radius)
			}
		}
		y += marks
		if l.hasBelow {
			for i, item := range r.Items {
				if item.Below != "" {
					c.text(item.Below, x(i), y+band/2)
				}
			}
			y += band
		}
		if r.Gaps != "" {
			for i := 0; i+1 < len(r.Items); i++ {
				if !r.Items[i].Skip && !r.Items[i+1].Skip {
					c.text(r.Gaps, (x(i)+x(i+1))/2, y+band/2)
				}
			}
		}
	}
	if r.Span != "" {
		y := height - margin/2 - spanBand/2
		advance := float64(font.MeasureString(faceAt(fontSize), r.Span).Round()) / scale
		c.span(x(0), x(len(r.Items)-1), y, r.Span, advance/2+4, band/2)
	}
}

// drawLine draws the line through the marks: whole between two marks, and broken around a skip.
func (r *rowPicture) drawLine(c canvas, y float64) {
	if r.Line != nil && !*r.Line {
		return
	}
	for i := 0; i+1 < len(r.Items); i++ {
		left, right := x(i), x(i+1)
		if r.Items[i].Skip {
			left += 14
		}
		if r.Items[i+1].Skip {
			right -= 14
		}
		c.rect(left, y-stroke/2, right, y+stroke/2, ink)
	}
}

// x is the centre of the item i, in CSS px.
func x(i int) float64 { return float64(margin + i*step) }
