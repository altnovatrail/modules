package main

import (
	"fmt"
	"math"
	"slices"
)

// shapesPicture is rows of shapes with a total beside each, as the schema of standalone's tasks
// has it: 1 to 5 rows of 1 to 6 shapes — triangle, square, circle or diamond — the same shape
// standing for the same number wherever it is drawn.
type shapesPicture struct {
	Kind string     `json:"kind"`
	Rows []shapeRow `json:"rows"`
}

type shapeRow struct {
	Shapes []string `json:"shapes"`
	Total  *string  `json:"total"`
}

var shapeNames = []string{"triangle", "square", "circle", "diamond"}

func (s *shapesPicture) valid() error {
	if len(s.Rows) < 1 || len(s.Rows) > 5 {
		return fmt.Errorf("shapes stand in 1 to 5 rows, not %d", len(s.Rows))
	}
	for i, row := range s.Rows {
		if len(row.Shapes) < 1 || len(row.Shapes) > 6 {
			return fmt.Errorf("the row %d has 1 to 6 shapes, not %d", i+1, len(row.Shapes))
		}
		for _, shape := range row.Shapes {
			if !slices.Contains(shapeNames, shape) {
				return fmt.Errorf("the shape %q is no triangle, square, circle or diamond", shape)
			}
		}
		if row.Total == nil || !isLabel(*row.Total) {
			return fmt.Errorf("the row %d has a total, a label", i+1)
		}
	}
	return nil
}

// The geometry of rows of shapes, in CSS px: how far apart two shapes stand, half a shape's side,
// how tall a row is, and the room before the total.
const (
	shapePitch = 40
	shapeHalf  = 13
	shapeRowH  = 40
	totalGap   = 24
)

func (s *shapesPicture) size() (int, int) {
	most, total := 0, 0.0
	for _, row := range s.Rows {
		most = max(most, len(row.Shapes))
		total = max(total, measure(*row.Total, fontSize))
	}
	return 2*margin + most*shapePitch + totalGap + int(total+0.5), margin + len(s.Rows)*shapeRowH
}

func (s *shapesPicture) paint(c canvas) {
	most := 0
	for _, row := range s.Rows {
		most = max(most, len(row.Shapes))
	}
	for r, row := range s.Rows {
		y := float64(margin/2+r*shapeRowH) + shapeRowH/2
		for i, shape := range row.Shapes {
			drawShape(c, shape, float64(margin+i*shapePitch)+shapePitch/2, y)
		}
		x := float64(margin+most*shapePitch) + totalGap
		c.text(*row.Total, x+measure(*row.Total, fontSize)/2, y)
	}
}

// drawShape draws one shape, its outline, centred on (x, y).
func drawShape(c canvas, shape string, x, y float64) {
	h := float64(shapeHalf)
	switch shape {
	case "circle":
		c.outline(x, y, h, stroke)
	case "square":
		c.polygon([][2]float64{{x - h, y - h}, {x + h, y - h}, {x + h, y + h}, {x - h, y + h}}, stroke)
	case "diamond":
		c.polygon([][2]float64{{x, y - h - 2}, {x + h + 2, y}, {x, y + h + 2}, {x - h - 2, y}}, stroke)
	default:
		side := 2 * h * 1.15
		tall := side * math.Sqrt(3) / 2
		c.polygon([][2]float64{{x, y - tall*2/3}, {x + side/2, y + tall/3}, {x - side/2, y + tall/3}}, stroke)
	}
}
