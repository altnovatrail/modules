package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// clockPicture is a clock face, as the schema of standalone's tasks has it: the time its hands
// show, from 0:00 to 23:59, with a time after 12:59 drawn as a twelve-hour face shows it, and a
// label at the tip of a hand where the question calls the hand so. The numbers 1 to 12 and the
// minute marks are the face's own.
type clockPicture struct {
	Kind        string `json:"kind"`
	Time        string `json:"time"`
	HourLabel   string `json:"hour_label,omitempty"`
	MinuteLabel string `json:"minute_label,omitempty"`
}

func (c *clockPicture) valid() error {
	if !timeOfDay.MatchString(c.Time) {
		return fmt.Errorf("a clock shows a time from 0:00 to 23:59, written H:MM, not %q", c.Time)
	}
	for _, text := range []string{c.HourLabel, c.MinuteLabel} {
		if text != "" && !isLabel(text) {
			return fmt.Errorf("the label %q of a hand is none", text)
		}
	}
	return nil
}

// The face, in CSS px, as standalone's card draws it: its middle, its rim, where its numbers
// stand, its hands and how wide they are, and how far a hand's label stands from the hand.
const (
	clockSide     = 200
	clockMiddle   = clockSide / 2
	clockRim      = 94
	clockNumerals = 72
	hourHand      = 46
	minuteHand    = 56
	numeralSize   = 16
	handAside     = 16
)

func (c *clockPicture) size() (int, int) { return clockSide, clockSide }

// hands are where the hands point, in degrees clockwise from twelve: the hour hand moved on by the
// minutes.
func (c *clockPicture) hands() (hour, minute float64) {
	h, m, _ := strings.Cut(c.Time, ":")
	hours, _ := strconv.Atoi(h)
	minutes, _ := strconv.Atoi(m)
	return float64(30*(hours%12)) + float64(minutes)/2, float64(6 * minutes)
}

// onFace is the point at a distance from the middle, in a direction in degrees clockwise from
// twelve.
func onFace(degrees, distance float64) (float64, float64) {
	angle := degrees * math.Pi / 180
	return clockMiddle + distance*math.Sin(angle), clockMiddle - distance*math.Cos(angle)
}

func (c *clockPicture) paint(cv canvas) {
	cv.ring(clockMiddle, clockMiddle, clockRim, 2.5)
	for minute := range 60 {
		inner, width := 87.0, 1.0
		if minute%5 == 0 {
			inner, width = 84, 2
		}
		x0, y0 := onFace(float64(minute*6), inner)
		x1, y1 := onFace(float64(minute*6), 91)
		cv.line(x0, y0, x1, y1, width)
	}
	for numeral := 1; numeral <= 12; numeral++ {
		x, y := onFace(float64(numeral*30), clockNumerals)
		cv.textAt(strconv.Itoa(numeral), x, y, numeralSize)
	}
	hour, minute := c.hands()
	hx, hy := onFace(hour, hourHand)
	mx, my := onFace(minute, minuteHand)
	cv.line(clockMiddle, clockMiddle, hx, hy, 6)
	cv.line(clockMiddle, clockMiddle, mx, my, 3.5)
	cv.disc(clockMiddle, clockMiddle, 4.5)
	hourSide, minuteSide := awayFrom(hour, minute), awayFrom(minute, hour)
	if hourSide == minuteSide {
		// The hands lie together: their labels stand on the two sides of them.
		minuteSide = -hourSide
	}
	if c.HourLabel != "" {
		handLabel(cv, c.HourLabel, hour, hourHand, hourSide)
	}
	if c.MinuteLabel != "" {
		handLabel(cv, c.MinuteLabel, minute, minuteHand, minuteSide)
	}
}

// awayFrom is the side of a hand away from the other, in degrees from it: 90 clockwise, -90 not.
func awayFrom(direction, other float64) float64 {
	if math.Mod(other-direction+360, 360) < 180 {
		return -90
	}
	return 90
}

// handLabel writes a hand's label beside it, on a side.
func handLabel(cv canvas, text string, direction, length, side float64) {
	x, y := onFace(direction, length*0.62)
	ax, ay := onFace(direction+side, handAside)
	cv.plate(text, x+ax-clockMiddle, y+ay-clockMiddle)
}
