package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"regexp"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/vector"
)

// A picture of standalone is a description that its card draws (content/pictures/ of standalone).
// The format takes a picture as a file (R51 spec), so the port draws it: a PNG on an opaque white
// background, as the format asks of every picture of a module (format/README.md, section 14; the
// second amendment of 2026-11-01 to R51 spec), with only the chunks IHDR, IDAT and IEND, and the
// same bytes on every run, since a picture's bytes are in the digest of the module's package. Only
// the kind row is drawn here; the other kinds are the rest of the port's (E3).

// rowPicture is a picture of the kind row, as the schema of standalone's tasks has it: 2 to 12
// items, each a mark with a label over it and one under it, or a skip that stands for the middle
// of a long row cut short; a line through the marks, unless it is false; a label under every gap
// drawn; a label of the whole length; and the row once or twice, one under the other.
type rowPicture struct {
	Kind   string    `json:"kind"`
	Items  []rowItem `json:"items"`
	Line   *bool     `json:"line,omitempty"`
	Gaps   string    `json:"gaps,omitempty"`
	Span   string    `json:"span,omitempty"`
	Copies int       `json:"copies,omitempty"`
}

type rowItem struct {
	Skip  bool   `json:"skip,omitempty"`
	Label string `json:"label,omitempty"`
	Below string `json:"below,omitempty"`
	Mark  string `json:"mark,omitempty"`
}

// label is a label of a picture of standalone: a Latin capital or a run of them, a number, or ?,
// at most 5 characters long.
var label = regexp.MustCompile(`^(\?|[A-Z]+|[-−]?[0-9]+([.,][0-9]+)?)$`)

// pictureOf draws a picture of standalone and gives the part of the statement or the solution that
// shows it: its path from media/, its alt and its size in CSS px. The alt is the overlay's: what
// the picture shows, in words, which a person writes, so that no answer is in it. A picture without
// one is refused, since alt is empty only for a picture that only decorates.
func pictureOf(id string, raw json.RawMessage, alt, suffix string) (mapping, []byte, error) {
	if alt == "" {
		return nil, nil, fmt.Errorf("%s: its picture has no alt in its overlay", id)
	}
	var kind struct {
		Kind string `json:"kind"`
	}
	if err := json.Unmarshal(raw, &kind); err != nil {
		return nil, nil, fmt.Errorf("%s: its picture: %w", id, err)
	}
	if kind.Kind != "row" {
		return nil, nil, fmt.Errorf("%s: its picture is of the kind %q, which the port does not draw yet: only row", id, kind.Kind)
	}
	var row rowPicture
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&row); err != nil {
		return nil, nil, fmt.Errorf("%s: its picture: %w", id, err)
	}
	if err := row.valid(); err != nil {
		return nil, nil, fmt.Errorf("%s: its picture: %w", id, err)
	}
	data, width, height, err := row.draw()
	if err != nil {
		return nil, nil, fmt.Errorf("%s: its picture: %w", id, err)
	}
	part := mapping{
		{"picture", "pictures/" + id + suffix + ".png"},
		{"alt", alt},
		{"width", width},
		{"height", height},
	}
	return part, data, nil
}

// pathOfPicture is the path of a picture part's file from the module's folder.
func pathOfPicture(part mapping) string {
	return "media/" + part[0].value.(string)
}

func (r rowPicture) valid() error {
	if len(r.Items) < 2 || len(r.Items) > 12 {
		return fmt.Errorf("a row has 2 to 12 items, not %d", len(r.Items))
	}
	if r.Copies < 0 || r.Copies > 2 {
		return fmt.Errorf("a row is drawn once or twice, not %d times", r.Copies)
	}
	for i, item := range r.Items {
		if item.Skip {
			if item.Label != "" || item.Below != "" || item.Mark != "" {
				return fmt.Errorf("the skip %d has a label or a mark", i+1)
			}
			if i == 0 || i == len(r.Items)-1 || r.Items[i-1].Skip {
				return fmt.Errorf("the skip %d is first, last or next to another", i+1)
			}
			continue
		}
		switch item.Mark {
		case "", "dot", "ring", "square":
		default:
			return fmt.Errorf("the item %d has the mark %q: dot, ring or square", i+1, item.Mark)
		}
		for _, text := range []string{item.Label, item.Below} {
			if text != "" && (!label.MatchString(text) || len([]rune(text)) > 5) {
				return fmt.Errorf("the item %d has the label %q, which is none", i+1, text)
			}
		}
	}
	for _, text := range []string{r.Gaps, r.Span} {
		if text != "" && (!label.MatchString(text) || len([]rune(text)) > 5) {
			return fmt.Errorf("the label %q is none", text)
		}
	}
	return nil
}

// The geometry of a row, in CSS px; the PNG holds twice as many pixels each way, for screens that
// show two to a CSS px.
const (
	scale    = 2
	margin   = 24
	step     = 56
	band     = 22
	marks    = 28
	spanBand = 30
	between  = 14
	stroke   = 2
	radius   = 6
	fontSize = 13
)

var (
	white = color.RGBA{0xff, 0xff, 0xff, 0xff}
	ink   = color.RGBA{0x1f, 0x29, 0x37, 0xff}
)

// draw lays the row out and draws it, and gives the PNG and its size in CSS px.
func (r rowPicture) draw() ([]byte, int, int, error) {
	face, err := faceOf()
	if err != nil {
		return nil, 0, 0, err
	}
	hasLabel, hasBelow := false, false
	for _, item := range r.Items {
		hasLabel = hasLabel || item.Label != ""
		hasBelow = hasBelow || item.Below != ""
	}
	copies := max(r.Copies, 1)
	rowHeight := marks
	if hasLabel {
		rowHeight += band
	}
	if hasBelow {
		rowHeight += band
	}
	if r.Gaps != "" {
		rowHeight += band
	}
	width := 2*margin + (len(r.Items)-1)*step
	height := 2*margin/2 + copies*rowHeight + (copies-1)*between
	if r.Span != "" {
		height += spanBand
	}
	img := image.NewRGBA(image.Rect(0, 0, width*scale, height*scale))
	draw.Draw(img, img.Bounds(), image.NewUniform(white), image.Point{}, draw.Src)
	c := canvas{img: img, face: face, raster: vector.NewRasterizer(width*scale, height*scale)}

	top := float64(margin / 2)
	for n := range copies {
		y := top + float64(n*(rowHeight+between))
		if hasLabel {
			for i, item := range r.Items {
				if item.Label != "" {
					c.text(item.Label, x(i), y+band/2)
				}
			}
			y += band
		}
		line := y + marks/2
		r.drawLine(c, line)
		for i, item := range r.Items {
			if item.Skip {
				for _, dx := range []float64{-8, 0, 8} {
					c.disc(x(i)+dx, line, 1.5)
				}
				continue
			}
			switch item.Mark {
			case "ring":
				c.disc(x(i), line, radius)
				c.discColor(x(i), line, radius-stroke, white)
			case "square":
				c.rect(x(i)-radius, line-radius, x(i)+radius, line+radius, ink)
			default:
				c.disc(x(i), line, radius)
			}
		}
		y += marks
		if hasBelow {
			for i, item := range r.Items {
				if item.Below != "" {
					c.text(item.Below, x(i), y+band/2)
				}
			}
			y += band
		}
		if r.Gaps != "" {
			for i := 0; i+1 < len(r.Items); i++ {
				if !r.Items[i].Skip && !r.Items[i+1].Skip {
					c.text(r.Gaps, (x(i)+x(i+1))/2, y+band/2)
				}
			}
		}
	}
	if r.Span != "" {
		y := float64(height - margin/2 - spanBand/2)
		left, right := x(0), x(len(r.Items)-1)
		c.rect(left, y-stroke/2, right, y+stroke/2, ink)
		c.rect(left-stroke/2, y-6, left+stroke/2, y+6, ink)
		c.rect(right-stroke/2, y-6, right+stroke/2, y+6, ink)
		advance := float64(font.MeasureString(face, r.Span).Round()) / scale
		middle := (left + right) / 2
		c.rect(middle-advance/2-4, y-band/2, middle+advance/2+4, y+band/2, white)
		c.text(r.Span, middle, y)
	}

	var out bytes.Buffer
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	if err := encoder.Encode(&out, img); err != nil {
		return nil, 0, 0, err
	}
	return out.Bytes(), width, height, nil
}

// drawLine draws the line through the marks: whole between two marks, and broken around a skip.
func (r rowPicture) drawLine(c canvas, y float64) {
	if r.Line != nil && !*r.Line {
		return
	}
	for i := 0; i+1 < len(r.Items); i++ {
		left, right := x(i), x(i+1)
		if r.Items[i].Skip {
			left += 14
		}
		if r.Items[i+1].Skip {
			right -= 14
		}
		c.rect(left, y-stroke/2, right, y+stroke/2, ink)
	}
}

// x is the centre of the item i, in CSS px.
func x(i int) float64 { return float64(margin + i*step) }

// faceOf is the face of the labels: Go Regular, BSD-3-Clause, at fontSize CSS px, parsed once.
var faceOf = sync.OnceValues(func() (font.Face, error) {
	parsed, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}
	return opentype.NewFace(parsed, &opentype.FaceOptions{Size: fontSize * scale, DPI: 72, Hinting: font.HintingNone})
})

// canvas draws on the image in CSS px, anti-aliased, always over the opaque background, so that no
// pixel of the picture is transparent. One rasterizer serves every shape of a picture.
type canvas struct {
	img    *image.RGBA
	face   font.Face
	raster *vector.Rasterizer
}

func (c canvas) disc(cx, cy, r float64) { c.discColor(cx, cy, r, ink) }

func (c canvas) discColor(cx, cy, r float64, col color.RGBA) {
	const sides = 48
	path := c.rasterizer()
	for i := range sides + 1 {
		angle := 2 * math.Pi * float64(i) / sides
		px, py := float32((cx+r*math.Cos(angle))*scale), float32((cy+r*math.Sin(angle))*scale)
		if i == 0 {
			path.MoveTo(px, py)
		} else {
			path.LineTo(px, py)
		}
	}
	path.ClosePath()
	path.Draw(c.img, c.img.Bounds(), image.NewUniform(col), image.Point{})
}

func (c canvas) rect(x0, y0, x1, y1 float64, col color.RGBA) {
	path := c.rasterizer()
	path.MoveTo(float32(x0*scale), float32(y0*scale))
	path.LineTo(float32(x1*scale), float32(y0*scale))
	path.LineTo(float32(x1*scale), float32(y1*scale))
	path.LineTo(float32(x0*scale), float32(y1*scale))
	path.ClosePath()
	path.Draw(c.img, c.img.Bounds(), image.NewUniform(col), image.Point{})
}

func (c canvas) rasterizer() *vector.Rasterizer {
	b := c.img.Bounds()
	c.raster.Reset(b.Dx(), b.Dy())
	c.raster.DrawOp = draw.Over
	return c.raster
}

// text draws a label centred on (cx, cy).
func (c canvas) text(s string, cx, cy float64) {
	metrics := c.face.Metrics()
	advance := font.MeasureString(c.face, s)
	d := font.Drawer{Dst: c.img, Src: image.NewUniform(ink), Face: c.face}
	d.Dot = fixed.Point26_6{
		X: fixed.I(int(math.Round(cx*scale))) - advance/2,
		Y: fixed.I(int(math.Round(cy*scale))) + (metrics.Ascent-metrics.Descent)/2,
	}
	d.DrawString(s)
}
