package main

import (
	"fmt"
	"strconv"
)

// containersPicture is jugs, buckets and the like to pour between, as the schema of standalone's
// tasks has it: 1 to 4 of them, each with how much it can hold, from 1 to 20, how much it holds
// now, from 0 up to that, and a label. The picture writes both numbers on it.
type containersPicture struct {
	Kind  string      `json:"kind"`
	Items []container `json:"items"`
}

type container struct {
	Capacity *int   `json:"capacity"`
	Amount   *int   `json:"amount"`
	Label    string `json:"label,omitempty"`
}

func (p *containersPicture) valid() error {
	if len(p.Items) < 1 || len(p.Items) > 4 {
		return fmt.Errorf("containers are 1 to 4, not %d", len(p.Items))
	}
	for i, item := range p.Items {
		switch {
		case item.Capacity == nil || *item.Capacity < 1 || *item.Capacity > 20:
			return fmt.Errorf("the container %d can hold from 1 to 20", i+1)
		case item.Amount == nil || *item.Amount < 0 || *item.Amount > *item.Capacity:
			return fmt.Errorf("the container %d holds from 0 up to its capacity", i+1)
		case item.Label != "" && !isLabel(item.Label):
			return fmt.Errorf("the label %q of the container %d is none", item.Label, i+1)
		}
	}
	return nil
}

// The geometry of containers, in CSS px: how wide one is, the room between two, how tall the
// largest may be and a unit at most, how thick a wall is, and the room a number needs.
const (
	jugWide    = 56
	jugBetween = 28
	jugTallest = 180
	jugUnit    = 14
	jugWall    = 3
	numberRoom = 18
)

func (p *containersPicture) unit() float64 {
	most := 0
	for _, item := range p.Items {
		most = max(most, *item.Capacity)
	}
	return min(jugUnit, float64(jugTallest)/float64(most))
}

// overTheTop says whether the number a container holds is written over it, at the unit given,
// having room neither in its water nor above the water.
func overTheTop(item container, u float64) bool {
	return float64(*item.Amount)*u < numberRoom && float64(*item.Capacity-*item.Amount)*u < numberRoom
}

// bands are the bands over the containers: one for their labels, one for a number written over a
// container.
func (p *containersPicture) bands() (labels, over bool) {
	u := p.unit()
	for _, item := range p.Items {
		labels = labels || item.Label != ""
		over = over || overTheTop(item, u)
	}
	return labels, over
}

func (p *containersPicture) size() (int, int) {
	most, u := 0.0, p.unit()
	for _, item := range p.Items {
		most = max(most, float64(*item.Capacity)*u)
	}
	labels, over := p.bands()
	height := margin + most + jugWall + band
	if labels {
		height += band
	}
	if over {
		height += band
	}
	return 2*margin + len(p.Items)*jugWide + (len(p.Items)-1)*jugBetween, int(height + 0.5)
}

func (p *containersPicture) paint(c canvas) {
	_, height := c.size()
	u := p.unit()
	floor := height - margin/2 - band - jugWall
	for i, item := range p.Items {
		left := float64(margin + i*(jugWide+jugBetween))
		right := left + jugWide
		top := floor - float64(*item.Capacity)*u
		water := floor - float64(*item.Amount)*u
		c.rect(left, water, right, floor, shade)
		c.rect(left-jugWall, top, left, floor+jugWall, ink)
		c.rect(right, top, right+jugWall, floor+jugWall, ink)
		c.rect(left-jugWall, floor, right+jugWall, floor+jugWall, ink)
		for unit := 1; unit < *item.Capacity; unit++ {
			y := floor - float64(unit)*u
			c.rect(left, y-0.5, left+6, y+0.5, ink)
		}
		amount := strconv.Itoa(*item.Amount)
		switch {
		case float64(*item.Amount)*u >= numberRoom:
			c.text(amount, (left+right)/2, (water+floor)/2)
		case !overTheTop(item, u):
			c.text(amount, (left+right)/2, water-numberRoom/2)
		default:
			c.text(amount, (left+right)/2, top-band/2)
		}
		c.text(strconv.Itoa(*item.Capacity), (left+right)/2, floor+jugWall+band/2)
		if item.Label != "" {
			labelAt := float64(margin/2) + band/2
			c.text(item.Label, (left+right)/2, labelAt)
		}
	}
}
