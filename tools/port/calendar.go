package main

import (
	"fmt"
	"regexp"
	"strconv"
)

// calendarPicture is the page of a month, as the schema of standalone's tasks has it: the weekday
// of its 1st, 1 for Monday to 7 for Sunday, its 28 to 31 days, the day its weeks start on, and a
// label on some of its days, keyed by the day. The page writes the weekdays' names in the module's
// language, English.
type calendarPicture struct {
	Kind       string            `json:"kind"`
	First      *int              `json:"first"`
	Days       *int              `json:"days"`
	WeekStarts string            `json:"week_starts,omitempty"`
	Marks      map[string]string `json:"marks,omitempty"`
}

// dayOfMonth is a day as a mark of a calendar keys it: a whole number from 1 to 31, with no nought
// before it, so that no day takes two marks.
var dayOfMonth = regexp.MustCompile(`^([1-9]|[12][0-9]|3[01])$`)

func (k *calendarPicture) valid() error {
	switch {
	case k.First == nil || *k.First < 1 || *k.First > 7:
		return fmt.Errorf("a calendar's first is the weekday of the 1st, from 1 for Monday to 7 for Sunday")
	case k.Days == nil || *k.Days < 28 || *k.Days > 31:
		return fmt.Errorf("a calendar's month has 28 to 31 days")
	case k.WeekStarts != "" && k.WeekStarts != "monday" && k.WeekStarts != "sunday":
		return fmt.Errorf("a calendar's weeks start on monday or sunday, not %q", k.WeekStarts)
	}
	for day, text := range k.Marks {
		number, _ := strconv.Atoi(day)
		if !dayOfMonth.MatchString(day) || number > *k.Days {
			return fmt.Errorf("the mark %q is on no day of the month", day)
		}
		if !isLabel(text) {
			return fmt.Errorf("the label %q of the day %s is none", text, day)
		}
	}
	return nil
}

// The geometry of a calendar, in CSS px.
const (
	dayWide    = 40
	weekTall   = 24
	markedWide = 26
	markedTall = 20
	markBand   = 18
	pageEdge   = 12
)

// weekdays are the names of the weekdays, from Monday.
var weekdays = []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

// offset is the column of the 1st, from 0.
func (k *calendarPicture) offset() int {
	if k.WeekStarts == "sunday" {
		return *k.First % 7
	}
	return *k.First - 1
}

// weekTalls are how tall each week is drawn: taller where a day of it has a mark.
func (k *calendarPicture) weekTalls() []int {
	weeks := (k.offset() + *k.Days + 6) / 7
	talls := make([]int, weeks)
	for week := range talls {
		talls[week] = weekTall
	}
	for day := range k.Marks {
		number, _ := strconv.Atoi(day)
		talls[(k.offset()+number-1)/7] = weekTall + markBand
	}
	return talls
}

func (k *calendarPicture) size() (int, int) {
	height := 2*pageEdge + weekTall + 4
	for _, tall := range k.weekTalls() {
		height += tall
	}
	return 2*pageEdge + 7*dayWide, height
}

func (k *calendarPicture) paint(c canvas) {
	width, _ := c.size()
	c.rect(pageEdge, pageEdge, width-pageEdge, pageEdge+weekTall, paper)
	for column := range 7 {
		name := weekdays[column]
		if k.WeekStarts == "sunday" {
			name = weekdays[(column+6)%7]
		}
		c.text(name, pageEdge+float64(column)*dayWide+dayWide/2, pageEdge+weekTall/2)
	}
	talls := k.weekTalls()
	tops := make([]float64, len(talls))
	top := float64(pageEdge + weekTall + 4)
	for week, tall := range talls {
		tops[week] = top
		top += float64(tall)
	}
	for day := 1; day <= *k.Days; day++ {
		cell := k.offset() + day - 1
		x := pageEdge + float64(cell%7)*dayWide + dayWide/2
		y := tops[cell/7] + weekTall/2
		c.text(strconv.Itoa(day), x, y)
		if text, marked := k.Marks[strconv.Itoa(day)]; marked {
			c.frame(x-markedWide/2, y-markedTall/2, x+markedWide/2, y+markedTall/2, 1.5)
			c.text(text, x, y+markedTall/2+markBand/2)
		}
	}
}
