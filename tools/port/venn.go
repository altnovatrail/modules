package main

import "fmt"

// vennPicture is two groups that overlap, inside a box of everyone the question counts, as the
// schema of standalone's tasks has it: each group's label and how many it holds, how many are in
// both, in neither and in all, every count a label.
type vennPicture struct {
	Kind    string    `json:"kind"`
	Sets    []vennSet `json:"sets"`
	Both    string    `json:"both,omitempty"`
	Neither string    `json:"neither,omitempty"`
	Total   string    `json:"total,omitempty"`
}

type vennSet struct {
	Label *string `json:"label"`
	Count string  `json:"count,omitempty"`
}

func (v *vennPicture) valid() error {
	if len(v.Sets) != 2 {
		return fmt.Errorf("a picture of groups has exactly 2 sets, not %d", len(v.Sets))
	}
	texts := []string{v.Both, v.Neither, v.Total}
	for _, set := range v.Sets {
		if set.Label == nil || !isLabel(*set.Label) {
			return fmt.Errorf("a set has a label")
		}
		texts = append(texts, set.Count)
	}
	for _, text := range texts {
		if text != "" && !isLabel(text) {
			return fmt.Errorf("the count %q is no label", text)
		}
	}
	return nil
}

// The geometry of two groups, in CSS px: a circle's radius, how far apart the middles stand, and
// the room inside the box.
const (
	vennRadius = 64
	vennApart  = 76
	vennRoom   = 14
)

func (v *vennPicture) size() (int, int) {
	return 2*margin/2 + 2*vennRoom + 2*vennRadius + vennApart, 2*margin/2 + 2*vennRoom + 2*vennRadius + band + 10
}

// caption is what is written over a group: its label, and how many it holds where the picture says.
func (set vennSet) caption() string {
	if set.Count == "" {
		return *set.Label
	}
	return *set.Label + ": " + set.Count
}

func (v *vennPicture) paint(c canvas) {
	width, height := c.size()
	left, top := float64(margin/2), float64(margin/2)+10
	right, bottom := width-margin/2, height-margin/2
	c.frame(left, top, right, bottom, 1.5)
	cy := top + vennRoom + band + vennRadius
	one, two := left+vennRoom+vennRadius, left+vennRoom+vennRadius+vennApart
	c.outline(one, cy, vennRadius, stroke)
	c.outline(two, cy, vennRadius, stroke)
	c.text(v.Sets[0].caption(), one-vennRadius/2, cy-vennRadius-band/2+4)
	c.text(v.Sets[1].caption(), two+vennRadius/2, cy-vennRadius-band/2+4)
	if v.Both != "" {
		c.plate(v.Both, (one+two)/2, cy)
	}
	if v.Neither != "" {
		c.text(v.Neither, left+vennRoom+measure(v.Neither, fontSize)/2, bottom-band/2)
	}
	if v.Total != "" {
		c.plate(v.Total, (left+right)/2, top)
	}
}
