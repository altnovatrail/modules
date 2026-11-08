package main

import (
	"fmt"
	"math"
	"strconv"
)

// ringPicture is places in a ring, as the schema of standalone's tasks has it: 3 to 24 of them,
// numbered from 1 in the order of travel and drawn clockwise from the top, and the place to start
// at, where the question names one, with a label or none.
type ringPicture struct {
	Kind  string     `json:"kind"`
	Count *int       `json:"count"`
	Start *ringStart `json:"start,omitempty"`
}

type ringStart struct {
	At    *int   `json:"at"`
	Label string `json:"label,omitempty"`
}

func (r *ringPicture) valid() error {
	switch {
	case r.Count == nil || *r.Count < 3 || *r.Count > 24:
		return fmt.Errorf("a ring has 3 to 24 places")
	case r.Start == nil:
		return nil
	case r.Start.At == nil || *r.Start.At < 1 || *r.Start.At > *r.Count:
		return fmt.Errorf("the start of a ring is a place of it, from 1 to its count")
	case r.Start.Label != "" && !isLabel(r.Start.Label):
		return fmt.Errorf("the label %q of the start is none", r.Start.Label)
	}
	return nil
}

// The geometry of a ring, in CSS px: a place's radius, the share of the ring each place takes at
// least, and the least radius of the ring.
const (
	placeRadius = 12
	placeShare  = 40
	ringLeast   = 64
)

func (r *ringPicture) radius() float64 {
	return max(ringLeast, float64(*r.Count*placeShare)/(2*math.Pi))
}

func (r *ringPicture) size() (int, int) {
	side := 2*margin + int(2*(r.radius()+placeRadius)+0.5)
	return side, side
}

// placeAt is the middle of a place, from 1 at the top on round clockwise.
func (r *ringPicture) placeAt(place int, middle float64) (float64, float64) {
	angle := 2 * math.Pi * float64(place-1) / float64(*r.Count)
	return middle + r.radius()*math.Sin(angle), middle - r.radius()*math.Cos(angle)
}

func (r *ringPicture) paint(c canvas) {
	side, _ := c.size()
	middle, radius := side/2, r.radius()
	c.ring(middle, middle, radius, 1.5)
	for place := 1; place <= *r.Count; place++ {
		// An arrowhead on the ring halfway to the next place, pointing on, clockwise.
		half := 2 * math.Pi * (float64(place) - 0.5) / float64(*r.Count)
		ax, ay := middle+radius*math.Sin(half), middle-radius*math.Cos(half)
		tx, ty := math.Cos(half), math.Sin(half)
		nx, ny := math.Sin(half), -math.Cos(half)
		c.line(ax+3*tx, ay+3*ty, ax-3*tx-5*nx, ay-3*ty-5*ny, 1.5)
		c.line(ax+3*tx, ay+3*ty, ax-3*tx+5*nx, ay-3*ty+5*ny, 1.5)
	}
	for place := 1; place <= *r.Count; place++ {
		x, y := r.placeAt(place, middle)
		if r.Start != nil && *r.Start.At == place {
			c.discColor(x, y, placeRadius, shade)
			c.outline(x, y, placeRadius, 2.5)
		} else {
			c.discColor(x, y, placeRadius, white)
			c.outline(x, y, placeRadius, 1.5)
		}
		c.text(strconv.Itoa(place), x, y)
	}
	if r.Start != nil && r.Start.Label != "" {
		x, y := r.placeAt(*r.Start.At, middle)
		toward := math.Hypot(x-middle, y-middle)
		ux, uy := (x-middle)/toward, (y-middle)/toward
		ex, ey := x-ux*(placeRadius+4), y-uy*(placeRadius+4)
		sx, sy := middle+ux*14, middle+uy*14
		c.line(sx, sy, ex, ey, 1.5)
		c.line(ex, ey, ex-ux*7-uy*4, ey-uy*7+ux*4, 1.5)
		c.line(ex, ey, ex-ux*7+uy*4, ey-uy*7-ux*4, 1.5)
		c.plate(r.Start.Label, middle, middle)
	}
}
