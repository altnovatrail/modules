package main

import "fmt"

// balancePicture is a pan balance, as the schema of standalone's tasks has it: the labels of what
// stands on each pan, at most 4 on a pan, each drawn as a weight.
type balancePicture struct {
	Kind  string    `json:"kind"`
	Left  *[]string `json:"left"`
	Right *[]string `json:"right"`
}

func (b *balancePicture) valid() error {
	for _, pan := range []*[]string{b.Left, b.Right} {
		if pan == nil {
			return fmt.Errorf("a balance has a left pan and a right pan")
		}
		if len(*pan) > 4 {
			return fmt.Errorf("a pan holds at most 4 weights, not %d", len(*pan))
		}
		for _, text := range *pan {
			if !isLabel(text) {
				return fmt.Errorf("the label %q of a weight is none", text)
			}
		}
	}
	return nil
}

// The geometry of a balance, in CSS px: a weight's height and least width, the room between
// weights, how many stand in a row, a pan's least width, the room between the pans, and the post,
// the stand and the foot under them.
const (
	weightTall   = 26
	weightNarrow = 28
	weightGap    = 6
	weightsInRow = 2
	panNarrow    = 96
	pansApart    = 40
	postTall     = 22
	standTall    = 56
	footWide     = 90
)

// weightWide is how wide a weight is: its label with room around it.
func weightWide(text string) float64 { return max(weightNarrow, measure(text, fontSize)+12) }

// rows are a pan's weights in rows from the pan up.
func rowsOf(pan []string) [][]string {
	var rows [][]string
	for start := 0; start < len(pan); start += weightsInRow {
		rows = append(rows, pan[start:min(start+weightsInRow, len(pan))])
	}
	return rows
}

// panWide is how wide a pan is: its widest row of weights with room around it.
func panWide(pan []string) float64 {
	wide := float64(panNarrow)
	for _, row := range rowsOf(pan) {
		sum := weightGap * float64(len(row)+1)
		for _, text := range row {
			sum += weightWide(text)
		}
		wide = max(wide, sum)
	}
	return wide
}

func (b *balancePicture) weightRows() int {
	return max(len(rowsOf(*b.Left)), len(rowsOf(*b.Right)), 1)
}

func (b *balancePicture) size() (int, int) {
	width := 2*margin + panWide(*b.Left) + pansApart + panWide(*b.Right)
	height := 2*margin/2 + b.weightRows()*(weightTall+weightGap) + postTall + standTall + 8
	return int(width + 0.5), height
}

func (b *balancePicture) paint(c canvas) {
	left, right := panWide(*b.Left), panWide(*b.Right)
	_, height := c.size()
	floor := height - margin/2
	beam := floor - standTall
	pan := beam - postTall
	middle := margin + left + pansApart/2
	leftMiddle, rightMiddle := margin+left/2, margin+left+pansApart+right/2
	c.rect(middle-footWide/2, floor-4, middle+footWide/2, floor, ink)
	c.line(middle, floor-4, middle, beam, 4)
	c.line(leftMiddle, beam, rightMiddle, beam, 4)
	c.disc(middle, beam, 6)
	for _, side := range []struct {
		middle, wide float64
		weights      []string
	}{{leftMiddle, left, *b.Left}, {rightMiddle, right, *b.Right}} {
		c.line(side.middle, beam, side.middle, pan, 3)
		c.line(side.middle-side.wide/2, pan, side.middle+side.wide/2, pan, 3)
		c.line(side.middle-side.wide/2, pan, side.middle-side.wide/2-4, pan-6, 3)
		c.line(side.middle+side.wide/2, pan, side.middle+side.wide/2+4, pan-6, 3)
		for r, row := range rowsOf(side.weights) {
			sum := weightGap * float64(len(row)-1)
			for _, text := range row {
				sum += weightWide(text)
			}
			x := side.middle - sum/2
			bottom := pan - 2 - float64(r*(weightTall+weightGap))
			for _, text := range row {
				wide := weightWide(text)
				c.frame(x, bottom-weightTall, x+wide, bottom, 1.5)
				c.text(text, x+wide/2, bottom-weightTall/2)
				x += wide + weightGap
			}
		}
	}
}
