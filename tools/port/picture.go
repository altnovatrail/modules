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
// same bytes on every run, since a picture's bytes are in the digest of the module's package.
//
// Each kind is read strictly, by the $defs of its kind in content/schemas/task.json of standalone
// at 42b4408, and held to what standalone's service checks and the schema cannot say
// (internal/domain/picture/ of standalone): a field the port does not know stops it. It is drawn in
// CSS px, twice as many pixels each way, with labels in Go Regular. The kinds not drawn yet are
// those of T20 and T21 (E3).

// drawing is a picture of one kind, read: what the service checks, the size it is drawn at in CSS
// px, and how it is drawn there.
type drawing interface {
	valid() error
	size() (width, height int)
	paint(c canvas)
}

// kinds are the kinds the port draws, each a new empty description to read one into.
var kinds = map[string]func() drawing{
	"row":         func() drawing { return &rowPicture{} },
	"clock":       func() drawing { return &clockPicture{} },
	"table":       func() drawing { return &tablePicture{} },
	"bars":        func() drawing { return &barsPicture{} },
	"number_line": func() drawing { return &numberLinePicture{} },
	"calendar":    func() drawing { return &calendarPicture{} },
}

var (
	// label is a label of a picture of standalone: a Latin capital or a run of them, a number, or ?,
	// at most 5 characters long.
	label = regexp.MustCompile(`^(\?|[A-Z]+|[-−]?[0-9]+([.,][0-9]+)?)$`)
	// timeOfDay is a time a picture of standalone shows: H:MM, from 0:00 to 23:59.
	timeOfDay = regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]$`)
)

// isLabel tells a label of a picture, at most 5 characters long, from any other text.
func isLabel(text string) bool {
	return label.MatchString(text) && len([]rune(text)) <= 5
}

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
	newDrawing, drawn := kinds[kind.Kind]
	if !drawn {
		return nil, nil, notDrawn{id: id, kind: kind.Kind}
	}
	if err := exactKeys(raw); err != nil {
		return nil, nil, fmt.Errorf("%s: its picture: %w", id, err)
	}
	picture := newDrawing()
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(picture); err != nil {
		return nil, nil, fmt.Errorf("%s: its picture: %w", id, err)
	}
	if err := picture.valid(); err != nil {
		return nil, nil, fmt.Errorf("%s: its picture: %w", id, err)
	}
	data, width, height, err := render(picture)
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

// key is a member's name as standalone's schema writes every one: lowercase, or a day of a
// calendar's marks.
var key = regexp.MustCompile(`^[a-z0-9_]+$`)

// exactKeys refuses what encoding/json lets through: a member named in other letters than the
// schema's, such as TIME for time, which it matches without regard to case, and a member named
// twice in one object, whose last value it takes. Every name of the schema is lowercase, so a
// lowercase name matches only itself.
func exactKeys(raw json.RawMessage) error {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	var walk func() error
	walk = func() error {
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		switch token {
		case json.Delim('{'):
			seen := map[string]bool{}
			for decoder.More() {
				name, err := decoder.Token()
				if err != nil {
					return err
				}
				member := name.(string)
				if !key.MatchString(member) || seen[member] {
					return fmt.Errorf("the member %q is named in other letters than the schema's, or twice", member)
				}
				seen[member] = true
				if err := walk(); err != nil {
					return err
				}
			}
			_, err = decoder.Token()
			return err
		case json.Delim('['):
			for decoder.More() {
				if err := walk(); err != nil {
					return err
				}
			}
			_, err = decoder.Token()
			return err
		}
		return nil
	}
	return walk()
}

// pathOfPicture is the path of a picture part's file from the module's folder.
func pathOfPicture(part mapping) string {
	return "media/" + part[0].value.(string)
}

// render draws a picture on a white opaque image of its size and gives the PNG and the size in CSS
// px.
func render(picture drawing) ([]byte, int, int, error) {
	width, height := picture.size()
	img := image.NewRGBA(image.Rect(0, 0, width*scale, height*scale))
	draw.Draw(img, img.Bounds(), image.NewUniform(white), image.Point{}, draw.Src)
	picture.paint(canvas{img: img, raster: vector.NewRasterizer(width*scale, height*scale)})
	var out bytes.Buffer
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	if err := encoder.Encode(&out, img); err != nil {
		return nil, 0, 0, err
	}
	return out.Bytes(), width, height, nil
}

// notDrawn is a picture of a kind the port does not draw yet, which a survey counts by kind.
type notDrawn struct{ id, kind string }

func (e notDrawn) Error() string {
	return fmt.Sprintf("%s: its picture is of the kind %q, which the port does not draw yet", e.id, e.kind)
}

// The scale of every picture, its stroke and the size of its labels, in CSS px, and what a row,
// bars and a number line share: the margin, a band of labels and the band of a span.
const (
	scale    = 2
	stroke   = 2
	fontSize = 13
	margin   = 24
	band     = 22
	spanBand = 30
)

var (
	white = color.RGBA{0xff, 0xff, 0xff, 0xff}
	ink   = color.RGBA{0x1f, 0x29, 0x37, 0xff}
	// shade fills what a picture shades, such as the shaded parts of a bar; paper fills a header.
	shade = color.RGBA{0xbf, 0xc5, 0xce, 0xff}
	paper = color.RGBA{0xee, 0xf0, 0xf3, 0xff}
)

// goRegular is the font of the labels, Go Regular, BSD-3-Clause, parsed once; faces holds a face
// of it at each size, in CSS px.
var (
	goRegular = sync.OnceValue(func() *opentype.Font {
		parsed, err := opentype.Parse(goregular.TTF)
		if err != nil {
			panic(err) // the font is in the binary: it always parses
		}
		return parsed
	})
	facesMu sync.Mutex
	faces   = map[float64]font.Face{}
)

// faceAt is the face of the labels at a size in CSS px.
func faceAt(size float64) font.Face {
	facesMu.Lock()
	defer facesMu.Unlock()
	if face, made := faces[size]; made {
		return face
	}
	face, err := opentype.NewFace(goRegular(), &opentype.FaceOptions{Size: size * scale, DPI: 72, Hinting: font.HintingNone})
	if err != nil {
		panic(err) // the options are constant: a face is always made
	}
	faces[size] = face
	return face
}

// measure is how wide a text is drawn at a size, in CSS px.
func measure(text string, size float64) float64 {
	return float64(font.MeasureString(faceAt(size), text)) / 64 / scale
}

// canvas draws on the image in CSS px, anti-aliased, always over the opaque background, so that no
// pixel of the picture is transparent. One rasterizer serves every shape of a picture.
type canvas struct {
	img    *image.RGBA
	raster *vector.Rasterizer
}

// size is the size of the picture the canvas draws, in CSS px.
func (c canvas) size() (float64, float64) {
	b := c.img.Bounds()
	return float64(b.Dx() / scale), float64(b.Dy() / scale)
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

// ring draws a circle's outline of a width, centred on its radius.
func (c canvas) ring(cx, cy, r, width float64) {
	c.disc(cx, cy, r+width/2)
	c.discColor(cx, cy, r-width/2, white)
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

// frame draws a rectangle's outline of a width, inside its edges.
func (c canvas) frame(x0, y0, x1, y1, width float64) {
	c.rect(x0, y0, x1, y0+width, ink)
	c.rect(x0, y1-width, x1, y1, ink)
	c.rect(x0, y0, x0+width, y1, ink)
	c.rect(x1-width, y0, x1, y1, ink)
}

// line draws a straight line of a width from one point to another, with round ends.
func (c canvas) line(x0, y0, x1, y1, width float64) {
	length := math.Hypot(x1-x0, y1-y0)
	if length > 0 {
		nx, ny := -(y1-y0)/length*width/2, (x1-x0)/length*width/2
		path := c.rasterizer()
		path.MoveTo(float32((x0+nx)*scale), float32((y0+ny)*scale))
		path.LineTo(float32((x1+nx)*scale), float32((y1+ny)*scale))
		path.LineTo(float32((x1-nx)*scale), float32((y1-ny)*scale))
		path.LineTo(float32((x0-nx)*scale), float32((y0-ny)*scale))
		path.ClosePath()
		path.Draw(c.img, c.img.Bounds(), image.NewUniform(ink), image.Point{})
	}
	c.disc(x0, y0, width/2)
	c.disc(x1, y1, width/2)
}

func (c canvas) rasterizer() *vector.Rasterizer {
	b := c.img.Bounds()
	c.raster.Reset(b.Dx(), b.Dy())
	c.raster.DrawOp = draw.Over
	return c.raster
}

// text draws a label centred on (cx, cy), at the size of the labels.
func (c canvas) text(s string, cx, cy float64) { c.textAt(s, cx, cy, fontSize) }

// textAt draws a text centred on (cx, cy), at a size in CSS px.
func (c canvas) textAt(s string, cx, cy, size float64) {
	face := faceAt(size)
	metrics := face.Metrics()
	advance := font.MeasureString(face, s)
	d := font.Drawer{Dst: c.img, Src: image.NewUniform(ink), Face: face}
	d.Dot = fixed.Point26_6{
		X: fixed.I(int(math.Round(cx*scale))) - advance/2,
		Y: fixed.I(int(math.Round(cy*scale))) + (metrics.Ascent-metrics.Descent)/2,
	}
	d.DrawString(s)
}

// plate draws a text centred on (cx, cy) on a white plate a little wider than it, so that it reads
// over the lines under it.
func (c canvas) plate(s string, cx, cy float64) {
	c.plateOf(s, cx, cy, measure(s, fontSize)/2+3, 9)
}

// plateOf draws a text centred on (cx, cy) on a white plate of a half width and a half height.
func (c canvas) plateOf(s string, cx, cy, halfWide, halfTall float64) {
	c.rect(cx-halfWide, cy-halfTall, cx+halfWide, cy+halfTall, white)
	c.text(s, cx, cy)
}

// span draws the span of a whole length, from left to right at y: a line with a tick at each end,
// and its label on a plate in the middle, of the half width given.
func (c canvas) span(left, right, y float64, s string, halfWide, halfTall float64) {
	c.rect(left, y-stroke/2, right, y+stroke/2, ink)
	c.rect(left-stroke/2, y-6, left+stroke/2, y+6, ink)
	c.rect(right-stroke/2, y-6, right+stroke/2, y+6, ink)
	c.plateOf(s, (left+right)/2, y, halfWide, halfTall)
}
