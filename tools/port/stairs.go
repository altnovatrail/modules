package main

import (
	"fmt"
	"math"
)

// stairsPicture is a staircase of 3 to 16 steps up from the floor, step 0, as the schema of
// standalone's tasks has it, with labels on some steps and jumps drawn from one step to another.
// The steps are not numbered: counting them is the child's.
type stairsPicture struct {
	Kind  string      `json:"kind"`
	Steps *int        `json:"steps"`
	Marks []stairMark `json:"marks,omitempty"`
	Jumps []stairJump `json:"jumps,omitempty"`
}

type stairMark struct {
	Step  *int    `json:"step"`
	Label *string `json:"label"`
}

type stairJump struct {
	From *int `json:"from"`
	To   *int `json:"to"`
}

func (s *stairsPicture) valid() error {
	if s.Steps == nil || *s.Steps < 3 || *s.Steps > 16 {
		return fmt.Errorf("a staircase has 3 to 16 steps")
	}
	onStep := func(step *int) bool { return step != nil && *step >= 0 && *step <= *s.Steps }
	if len(s.Marks) > 4 || len(s.Jumps) > 4 {
		return fmt.Errorf("a staircase has at most 4 marks and 4 jumps")
	}
	marked := map[int]bool{}
	for _, mark := range s.Marks {
		switch {
		case !onStep(mark.Step):
			return fmt.Errorf("a mark is on no step of the staircase: a whole number from 0, the floor, to its top step")
		case marked[*mark.Step]:
			return fmt.Errorf("two marks stand on one step")
		case mark.Label == nil || !isLabel(*mark.Label):
			return fmt.Errorf("a mark of a staircase has a label")
		}
		marked[*mark.Step] = true
	}
	jumped := map[[2]int]bool{}
	for _, jump := range s.Jumps {
		switch {
		case !onStep(jump.From) || !onStep(jump.To):
			return fmt.Errorf("a jump goes between steps of the staircase: whole numbers from 0, the floor, to its top step")
		case *jump.From == *jump.To:
			return fmt.Errorf("a jump goes from one step to another")
		case jumped[[2]int{*jump.From, *jump.To}]:
			return fmt.Errorf("a jump is drawn twice")
		}
		jumped[[2]int{*jump.From, *jump.To}] = true
	}
	return nil
}

// The geometry of stairs, in CSS px: how wide and how tall a step is, how far a label stands above
// its step, and how high a jump rises over the higher of its ends.
const (
	stepWide  = 30
	stepRise  = 16
	labelOver = 12
	jumpRise  = 30
)

func (s *stairsPicture) size() (int, int) {
	return 2*margin + (*s.Steps+1)*stepWide, margin + *s.Steps*stepRise + labelOver + 12 + jumpRise + 8
}

// stepTop is the middle of the top of a step, from 0, the floor.
func (s *stairsPicture) stepTop(step int, floor float64) (float64, float64) {
	return float64(margin+step*stepWide) + stepWide/2, floor - float64(step*stepRise)
}

func (s *stairsPicture) paint(c canvas) {
	_, height := c.size()
	floor := height - margin/2
	left, right := float64(margin), float64(margin+(*s.Steps+1)*stepWide)
	c.line(left-8, floor, right, floor, stroke)
	x, y := left+stepWide, floor
	for step := 1; step <= *s.Steps; step++ {
		c.line(x, y, x, y-stepRise, stroke)
		c.line(x, y-stepRise, x+stepWide, y-stepRise, stroke)
		x, y = x+stepWide, y-stepRise
	}
	c.line(right, y, right, floor, stroke)
	marked := map[int]bool{}
	for _, mark := range s.Marks {
		mx, my := s.stepTop(*mark.Step, floor)
		c.text(*mark.Label, mx, my-labelOver)
		marked[*mark.Step] = true
	}
	// end is where a jump leaves or lands on a step: just over it, or over its label.
	end := func(step int) (float64, float64) {
		x, y := s.stepTop(step, floor)
		if marked[step] {
			return x, y - labelOver - 12
		}
		return x, y - 6
	}
	for _, jump := range s.Jumps {
		fx, fy := end(*jump.From)
		tx, ty := end(*jump.To)
		lift := min(fy, ty) - jumpRise
		// A parabola from one end to the other, rising jumpRise over the higher, and an arrowhead
		// where it lands.
		const parts = 24
		point := func(t float64) (float64, float64) {
			base := fy + (ty-fy)*t
			return fx + (tx-fx)*t, base - 4*t*(1-t)*(base-lift)
		}
		px, py := point(0)
		for i := 1; i <= parts; i++ {
			nx, ny := point(float64(i) / parts)
			c.line(px, py, nx, ny, 1.5)
			px, py = nx, ny
		}
		bx, by := point(1 - 1.0/parts)
		angle := math.Atan2(py-by, px-bx)
		for _, turn := range []float64{2.6, -2.6} {
			c.line(px, py, px+7*math.Cos(angle+turn), py+7*math.Sin(angle+turn), 1.5)
		}
	}
}
