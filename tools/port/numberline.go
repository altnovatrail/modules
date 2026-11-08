package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// numberLinePicture is a number line, as the schema of standalone's tasks has it: 2 to 16 evenly
// spaced ticks, from from to to by step, which the line numbers itself, each from -9999 to 99999,
// and the marks the question names, each on a tick, one to a tick, with a label or none.
type numberLinePicture struct {
	Kind  string     `json:"kind"`
	From  *int       `json:"from"`
	To    *int       `json:"to"`
	Step  *int       `json:"step,omitempty"`
	Marks []lineMark `json:"marks,omitempty"`
}

type lineMark struct {
	At    *int   `json:"at"`
	Label string `json:"label,omitempty"`
}

const lowestTick, highestTick = -9999, 99999

func (n *numberLinePicture) valid() error {
	if n.From == nil || n.To == nil {
		return fmt.Errorf("a number line has from and to")
	}
	from, to, step := *n.From, *n.To, n.step()
	switch {
	case from < lowestTick || from > highestTick || to < lowestTick || to > highestTick:
		return fmt.Errorf("a tick of a number line is from %d to %d", lowestTick, highestTick)
	case n.Step != nil && (step < 1 || step > highestTick-lowestTick):
		return fmt.Errorf("the step of a number line is from 1 to %d, not %d", highestTick-lowestTick, step)
	case to <= from:
		return fmt.Errorf("to, %d, is not greater than from, %d", to, from)
	case (to-from)%step != 0:
		return fmt.Errorf("the step %d does not reach %d from %d in whole steps", step, to, from)
	case (to-from)/step+1 > 16:
		return fmt.Errorf("a number line has 2 to 16 ticks, not %d", (to-from)/step+1)
	case len(n.Marks) > 16:
		return fmt.Errorf("a number line has at most 16 marks, not %d", len(n.Marks))
	}
	var taken []int
	for _, mark := range n.Marks {
		switch {
		case mark.At == nil:
			return fmt.Errorf("a mark of a number line has no at")
		case *mark.At < from || *mark.At > to || (*mark.At-from)%step != 0:
			return fmt.Errorf("the mark at %d stands on no tick of the line", *mark.At)
		case slices.Contains(taken, *mark.At):
			return fmt.Errorf("two marks stand on the tick %d", *mark.At)
		case mark.Label != "" && !isLabel(mark.Label):
			return fmt.Errorf("the label %q of a mark is none", mark.Label)
		}
		taken = append(taken, *mark.At)
	}
	return nil
}

// step is the line's step, 1 unless it gives one.
func (n *numberLinePicture) step() int {
	if n.Step == nil {
		return 1
	}
	return *n.Step
}

// The geometry of a number line, in CSS px.
const (
	tickNarrow = 36
	lineEnd    = 24
	tickHalf   = 6
	lineBand   = 24
)

func (n *numberLinePicture) ticks() []int {
	var ticks []int
	for at := *n.From; at <= *n.To; at += n.step() {
		ticks = append(ticks, at)
	}
	return ticks
}

// spacing is how far apart the ticks are: far enough for the widest number under them.
func (n *numberLinePicture) spacing() float64 {
	spacing := float64(tickNarrow)
	for _, tick := range n.ticks() {
		spacing = max(spacing, measure(tickNumber(tick), fontSize)+10)
	}
	return spacing
}

// tickNumber is a tick's number as the line writes it, a minus sign before a negative one.
func tickNumber(tick int) string {
	return strings.Replace(strconv.Itoa(tick), "-", "−", 1)
}

func (n *numberLinePicture) labelled() bool {
	return slices.ContainsFunc(n.Marks, func(m lineMark) bool { return m.Label != "" })
}

func (n *numberLinePicture) size() (int, int) {
	height := margin/2 + lineBand + band + margin/2
	if n.labelled() {
		height += band
	}
	return 2*margin + int(float64(len(n.ticks())-1)*n.spacing()+0.5), height
}

func (n *numberLinePicture) paint(c canvas) {
	ticks, spacing := n.ticks(), n.spacing()
	at := func(tick int) float64 { return margin + float64((tick-*n.From)/n.step())*spacing }
	y := float64(margin/2) + lineBand/2
	if n.labelled() {
		y += band
	}
	left, right := at(ticks[0])-lineEnd, at(ticks[len(ticks)-1])+lineEnd
	c.line(left, y, right, y, stroke)
	c.line(right, y, right-7, y-5, stroke)
	c.line(right, y, right-7, y+5, stroke)
	for _, tick := range ticks {
		c.line(at(tick), y-tickHalf, at(tick), y+tickHalf, 1.5)
		c.text(tickNumber(tick), at(tick), y+lineBand/2+band/2)
	}
	for _, mark := range n.Marks {
		c.disc(at(*mark.At), y, 5)
		if mark.Label != "" {
			c.text(mark.Label, at(*mark.At), y-lineBand/2-band/2+4)
		}
	}
}
