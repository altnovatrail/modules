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

// A picture of standalone is a description that its card draws, and the format takes a picture as
// a file, so the port draws it: a PNG on an opaque white background, as the format asks of every
// picture, with only the chunks IHDR, IDAT and IEND, and the same bytes on every run, since a
// picture's bytes are in the digest of the module's package.
//
// Each kind is read strictly, by its $defs in standalone's content/schemas/task.json, and held to
// what standalone's service checks and the schema cannot say (internal/domain/picture/ of
// standalone): a field the port does not know stops it.

// drawing is a picture of one kind, read: what the service checks, its size in CSS px, how it is
// drawn, and the texts it writes.
type drawing interface {
	valid() error
	size() (width, height int)
	paint(c canvas)
	labels() []string
}

// kinds give a new empty description of each kind the port draws, to read one into.
var kinds = map[string]func() drawing{
	"row":         func() drawing { return &rowPicture{} },
	"clock":       func() drawing { return &clockPicture{} },
	"table":       func() drawing { return &tablePicture{} },
	"bars":        func() drawing { return &barsPicture{} },
	"number_line": func() drawing { return &numberLinePicture{} },
	"calendar":    func() drawing { return &calendarPicture{} },
	"containers":  func() drawing { return &containersPicture{} },
	"grid":        func() drawing { return &gridPicture{} },
	"piles":       func() drawing { return &pilesPicture{} },
	"venn":        func() drawing { return &vennPicture{} },
	"balance":     func() drawing { return &balancePicture{} },
	"ring":        func() drawing { return &ringPicture{} },
	"lattice":     func() drawing { return &latticePicture{} },
	"tiles":       func() drawing { return &tilesPicture{} },
	"shapes":      func() drawing { return &shapesPicture{} },
	"stairs":      func() drawing { return &stairsPicture{} },
}

var (
	// label is a label of a picture of standalone: a Latin capital or a run of them, a number, or ?,
	// at most 5 characters long.
	label = regexp.MustCompile(`^(\?|[A-Z]+|[-−]?[0-9]+([.,][0-9]+)?)$`)
	// timeOfDay is a time a picture of standalone shows: H:MM, from 0:00 to 23:59.
	timeOfDay = regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]$`)
)

func isLabel(text string) bool {
	return label.MatchString(text) && len([]rune(text)) <= 5
}

// pictureOf draws a picture and gives the part of the statement or the solution that shows it. The
// alt is the overlay's, written by a person so that no answer is in it; a picture without one is
// refused, since the format leaves alt empty only for a decoration.
func pictureOf(id string, raw json.RawMessage, alt, suffix string) (mapping, []byte, error) {
	if alt == "" {
		return nil, nil, fmt.Errorf("%s: its picture has no alt in its overlay", id)
	}
	picture, err := readPicture(id, raw)
	if err != nil {
		return nil, nil, err
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

// key is a member's name as standalone's schema writes every one: lowercase.
var key = regexp.MustCompile(`^[a-z0-9_]+$`)

// keyedByData are the members of a kind that are maps keyed by data, not by the schema's names:
// the marks of a grid, keyed by a cell, B2, and of a calendar, keyed by a day.
var keyedByData = map[string]string{"grid": "marks", "calendar": "marks"}

// exactKeys refuses what encoding/json lets through: a member named in other letters than the
// schema's, such as TIME for time, which it matches without regard to case, and a member named
// twice, whose last value it takes. Every name of the schema is lowercase, as key holds, so a
// lowercase name matches only itself. The member that keyed names, a map keyed by data, has its
// keys held only to being named once.
func exactKeys(raw json.RawMessage, keyed string) error {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	depth := 0
	var walk func(byData bool) error
	walk = func(byData bool) error {
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		switch token {
		case json.Delim('{'):
			depth++
			defer func() { depth-- }()
			seen := map[string]bool{}
			for decoder.More() {
				name, err := decoder.Token()
				if err != nil {
					return err
				}
				member := name.(string)
				if (!byData && !key.MatchString(member)) || seen[member] {
					return fmt.Errorf("the member %q is named in other letters than the schema's, or twice", member)
				}
				seen[member] = true
				if err := walk(depth == 1 && keyed != "" && member == keyed); err != nil {
					return err
				}
			}
			_, err = decoder.Token()
			return err
		case json.Delim('['):
			for decoder.More() {
				if err := walk(false); err != nil {
					return err
				}
			}
			_, err = decoder.Token()
			return err
		}
		return nil
	}
	return walk(false)
}

// readPicture reads a picture strictly, by its kind: a kind the port does not draw is notDrawn.
func readPicture(id string, raw json.RawMessage) (drawing, error) {
	var kind struct {
		Kind string `json:"kind"`
	}
	if err := json.Unmarshal(raw, &kind); err != nil {
		return nil, fmt.Errorf("%s: its picture: %w", id, err)
	}
	newDrawing, drawn := kinds[kind.Kind]
	if !drawn {
		return nil, notDrawn{id: id, kind: kind.Kind}
	}
	if err := exactKeys(raw, keyedByData[kind.Kind]); err != nil {
		return nil, fmt.Errorf("%s: its picture: %w", id, err)
	}
	picture := newDrawing()
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(picture); err != nil {
		return nil, fmt.Errorf("%s: its picture: %w", id, err)
	}
	if err := picture.valid(); err != nil {
		return nil, fmt.Errorf("%s: its picture: %w", id, err)
	}
	return picture, nil
}

// pathOfPicture is the path of a picture part's file from the module's folder.
func pathOfPicture(part mapping) string {
	return "media/" + part[0].value.(string)
}

// render gives the PNG of a picture and its size in CSS px.
func render(picture drawing) ([]byte, int, int, error) {
	width, height := picture.size()
	data, err := encodePNG(paintImage(picture))
	return data, width, height, err
}

func paintImage(picture drawing) *image.RGBA {
	width, height := picture.size()
	img := whiteImage(width, height)
	picture.paint(canvas{img: img, raster: vector.NewRasterizer(width*scale, height*scale)})
	return img
}

// whiteImage is a white opaque image of a size in CSS px, at the scale the port draws.
func whiteImage(width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width*scale, height*scale))
	draw.Draw(img, img.Bounds(), image.NewUniform(white), image.Point{}, draw.Src)
	return img
}

func encodePNG(img image.Image) ([]byte, error) {
	var out bytes.Buffer
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	if err := encoder.Encode(&out, img); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// notDrawn is the error of a picture of a kind the port does not draw, which a survey counts by
// kind.
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
	// rule is the lines of squared paper: light, and still there to count the cells by.
	rule = color.RGBA{0xa8, 0xaf, 0xb8, 0xff}
)

// goRegular is the font of the labels, BSD-3-Clause; faces holds a face of it at each size, in CSS
// px. A face is not safe for concurrent use, and template-pictures draws on several goroutines, so
// facesMu is held for each use of one.
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

// withFace calls use with the face of the labels at a size in CSS px, holding facesMu.
func withFace(size float64, use func(font.Face)) {
	facesMu.Lock()
	defer facesMu.Unlock()
	face, made := faces[size]
	if !made {
		var err error
		face, err = opentype.NewFace(goRegular(), &opentype.FaceOptions{Size: size * scale, DPI: 72, Hinting: font.HintingNone})
		if err != nil {
			panic(err) // the options are constant: a face is always made
		}
		faces[size] = face
	}
	use(face)
}

// advance is how wide a text is drawn at a size, in 26.6 fixed point at the scale the port draws.
func advance(text string, size float64) fixed.Int26_6 {
	var wide fixed.Int26_6
	withFace(size, func(face font.Face) { wide = font.MeasureString(face, text) })
	return wide
}

// measure is how wide a text is drawn at a size, in CSS px.
func measure(text string, size float64) float64 {
	return float64(advance(text, size)) / 64 / scale
}

// canvas draws on the image in CSS px, anti-aliased, always over the opaque background, so that no
// pixel of the picture is transparent. One rasterizer serves every shape of a picture.
type canvas struct {
	img    *image.RGBA
	raster *vector.Rasterizer
}

func (c canvas) size() (float64, float64) {
	b := c.img.Bounds()
	return float64(b.Dx() / scale), float64(b.Dy() / scale)
}

func (c canvas) disc(cx, cy, r float64) { c.discColor(cx, cy, r, ink) }

func (c canvas) discColor(cx, cy, r float64, col color.RGBA) {
	path := c.rasterizer()
	circle(path, cx, cy, r, 48, 1)
	path.Draw(c.img, c.img.Bounds(), image.NewUniform(col), image.Point{})
}

// circle adds a closed circle of a number of sides to a path, one way round or, with turn -1, the
// other.
func circle(path *vector.Rasterizer, cx, cy, r float64, sides int, turn float64) {
	for i := range sides + 1 {
		angle := turn * 2 * math.Pi * float64(i) / float64(sides)
		px, py := float32((cx+r*math.Cos(angle))*scale), float32((cy+r*math.Sin(angle))*scale)
		if i == 0 {
			path.MoveTo(px, py)
		} else {
			path.LineTo(px, py)
		}
	}
	path.ClosePath()
}

// ring draws a circle's outline centred on its radius.
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

// frame draws a rectangle's outline inside its edges.
func (c canvas) frame(x0, y0, x1, y1, width float64) {
	c.rect(x0, y0, x1, y0+width, ink)
	c.rect(x0, y1-width, x1, y1, ink)
	c.rect(x0, y0, x0+width, y1, ink)
	c.rect(x1-width, y0, x1, y1, ink)
}

// line draws a line with round ends.
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

func (c canvas) text(s string, cx, cy float64) { c.textAt(s, cx, cy, fontSize) }

// textAt draws a text centred on (cx, cy).
func (c canvas) textAt(s string, cx, cy, size float64) {
	withFace(size, func(face font.Face) {
		metrics := face.Metrics()
		wide := font.MeasureString(face, s)
		d := font.Drawer{Dst: c.img, Src: image.NewUniform(ink), Face: face}
		d.Dot = fixed.Point26_6{
			X: fixed.I(int(math.Round(cx*scale))) - wide/2,
			Y: fixed.I(int(math.Round(cy*scale))) + (metrics.Ascent-metrics.Descent)/2,
		}
		d.DrawString(s)
	})
}

// plate draws a text centred on (cx, cy) on a white plate a little wider than it, so that it reads
// over the lines under it.
func (c canvas) plate(s string, cx, cy float64) {
	c.plateOf(s, cx, cy, measure(s, fontSize)/2+3, 9)
}

func (c canvas) plateOf(s string, cx, cy, halfWide, halfTall float64) {
	c.rect(cx-halfWide, cy-halfTall, cx+halfWide, cy+halfTall, white)
	c.text(s, cx, cy)
}

// span draws the span of a whole length at y: a line with a tick at each end, and its label on a
// plate in the middle.
func (c canvas) span(left, right, y float64, s string, halfWide, halfTall float64) {
	c.rect(left, y-stroke/2, right, y+stroke/2, ink)
	c.rect(left-stroke/2, y-6, left+stroke/2, y+6, ink)
	c.rect(right-stroke/2, y-6, right+stroke/2, y+6, ink)
	c.plateOf(s, (left+right)/2, y, halfWide, halfTall)
}

// outline draws a circle's outline centred on its radius and leaves the inside as it is: the outer
// edge one way round and the inner the other, in one path.
func (c canvas) outline(cx, cy, r, width float64) {
	path := c.rasterizer()
	circle(path, cx, cy, r+width/2, 96, 1)
	circle(path, cx, cy, r-width/2, 96, -1)
	path.Draw(c.img, c.img.Bounds(), image.NewUniform(ink), image.Point{})
}

func (c canvas) polygon(corners [][2]float64, width float64) {
	for i, from := range corners {
		to := corners[(i+1)%len(corners)]
		c.line(from[0], from[1], to[0], to[1], width)
	}
}

// dashed draws dashes of 6 CSS px, 4 apart.
func (c canvas) dashed(x0, y0, x1, y1, width float64) {
	length := math.Hypot(x1-x0, y1-y0)
	for from := 0.0; from < length; from += 10 {
		to := min(from+6, length)
		c.line(x0+(x1-x0)*from/length, y0+(y1-y0)*from/length, x0+(x1-x0)*to/length, y0+(y1-y0)*to/length, width)
	}
}
