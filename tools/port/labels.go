package main

import (
	"encoding/json"
	"maps"
	"slices"
	"strconv"
)

// The texts a picture writes, as a description gives them: what the alt of a task solved from the
// picture names, so that a learner who cannot see the picture reads its facts (R20; T21). The
// numbers a picture numbers itself — a clock's 1 to 12, a line's ticks, a calendar's days, a ring's
// places — are not given by the description, and are left out, and so are ? and the empty text,
// which say nothing.

// writtenLabels are the texts a picture writes, in order, once each: the picture read and held as
// pictureOf reads it.
func writtenLabels(id string, raw json.RawMessage) ([]string, error) {
	picture, err := readPicture(id, raw)
	if err != nil {
		return nil, err
	}
	var said []string
	for _, text := range picture.labels() {
		if text != "" && text != "?" && !slices.Contains(said, text) {
			said = append(said, text)
		}
	}
	return said, nil
}

func (r *rowPicture) labels() []string {
	texts := []string{r.Gaps, r.Span}
	for _, item := range r.Items {
		texts = append(texts, item.Label, item.Below)
	}
	return texts
}

// A clock writes no time: its hands show it, and its alt says it as any fact.
func (c *clockPicture) labels() []string { return []string{c.HourLabel, c.MinuteLabel} }

func (t *tablePicture) labels() []string { return slices.Concat(t.all()...) }

func (b *barsPicture) labels() []string {
	texts := slices.Clone(b.Notes)
	for _, bar := range b.Bars {
		texts = append(texts, bar.Label, bar.Value, bar.Span)
		for _, segment := range bar.Segments {
			texts = append(texts, segment.Label)
		}
		for _, brace := range bar.Braces {
			texts = append(texts, *brace.Label)
		}
	}
	return texts
}

func (n *numberLinePicture) labels() []string {
	var texts []string
	for _, mark := range n.Marks {
		texts = append(texts, mark.Label)
	}
	return texts
}

func (k *calendarPicture) labels() []string {
	var texts []string
	for _, day := range slices.Sorted(maps.Keys(k.Marks)) {
		texts = append(texts, k.Marks[day])
	}
	return texts
}

func (p *containersPicture) labels() []string {
	var texts []string
	for _, item := range p.Items {
		texts = append(texts, item.Label, strconv.Itoa(*item.Capacity), strconv.Itoa(*item.Amount))
	}
	return texts
}

func (g *gridPicture) labels() []string {
	var texts []string
	for _, cell := range slices.Sorted(maps.Keys(g.Marks)) {
		texts = append(texts, g.Marks[cell])
	}
	return texts
}

func (p *pilesPicture) labels() []string {
	var texts []string
	for _, pile := range p.Piles {
		texts = append(texts, pile.Label, pile.Value)
	}
	return texts
}

func (v *vennPicture) labels() []string {
	texts := []string{v.Both, v.Neither, v.Total}
	for _, set := range v.Sets {
		texts = append(texts, *set.Label, set.Count)
	}
	return texts
}

func (b *balancePicture) labels() []string { return slices.Concat(*b.Left, *b.Right) }

func (r *ringPicture) labels() []string {
	if r.Start == nil {
		return nil
	}
	return []string{r.Start.Label}
}

func (l *latticePicture) labels() []string {
	var texts []string
	for _, item := range slices.Concat(l.Points, l.Marks) {
		if item.Label != nil {
			texts = append(texts, *item.Label)
		}
	}
	return texts
}

func (t *tilesPicture) labels() []string { return slices.Concat(t.Tiles...) }

func (s *shapesPicture) labels() []string {
	var texts []string
	for _, row := range s.Rows {
		texts = append(texts, *row.Total)
	}
	return texts
}

func (s *stairsPicture) labels() []string {
	var texts []string
	for _, mark := range s.Marks {
		texts = append(texts, *mark.Label)
	}
	return texts
}
