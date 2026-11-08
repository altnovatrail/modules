package main

import (
	"fmt"
	"math"
	"regexp"
)

// tilesPicture is identical tiles in a row or a ring, as the schema of standalone's tasks has it:
// 2 to 12 tiles, each a regular polygon of 3 to 6 sides cut into triangles from its middle, with a
// number of at most 2 digits, ? or nothing in each triangle, clockwise from the side it shares with
// the tile before it. A tile shares its side joint, half its sides unless given, with the tile after
// it, and a row runs straight on while a ring goes round once and closes.
type tilesPicture struct {
	Kind   string     `json:"kind"`
	Sides  *int       `json:"sides"`
	Layout string     `json:"layout"`
	Joint  *int       `json:"joint,omitempty"`
	Tiles  [][]string `json:"tiles"`
}

// tileText is what a triangle of a tile holds: a number of at most 2 digits, ?, or nothing.
var tileText = regexp.MustCompile(`^(|\?|[0-9]{1,2})$`)

func (t *tilesPicture) valid() error {
	switch {
	case t.Sides == nil || *t.Sides < 3 || *t.Sides > 6:
		return fmt.Errorf("a tile has 3 to 6 sides")
	case t.Layout != "row" && t.Layout != "ring":
		return fmt.Errorf("tiles are laid in a row or a ring, not %q", t.Layout)
	case t.Joint != nil && (*t.Joint < 1 || *t.Joint > *t.Sides-1):
		return fmt.Errorf("the joint of a tile is a side of it from 1 to %d", *t.Sides-1)
	case len(t.Tiles) < 2 || len(t.Tiles) > 12:
		return fmt.Errorf("tiles are 2 to 12, not %d", len(t.Tiles))
	}
	for i, tile := range t.Tiles {
		if len(tile) != *t.Sides {
			return fmt.Errorf("the tile %d has as many triangles as a tile has sides, not %d", i+1, len(tile))
		}
		for _, text := range tile {
			if !tileText.MatchString(text) {
				return fmt.Errorf("the triangle %q holds no number of at most 2 digits, ? or nothing", text)
			}
		}
	}
	if !tilesLaidOut(*t.Sides, t.joint(), len(t.Tiles), t.Layout == "ring") {
		if t.Layout == "row" {
			return fmt.Errorf("a row of tiles runs straight: tiles of 4 or 6 sides, each against the side opposite its first, at most 6 tiles and 4 hexagons")
		}
		return fmt.Errorf("a ring of tiles goes round once and closes, its last tile against its first")
	}
	return nil
}

func (t *tilesPicture) joint() int {
	if t.Joint == nil {
		return *t.Sides / 2
	}
	return *t.Joint
}

// turnOf is how far each tile is turned from the one before it, in halves of a side, from less
// than half round back to half round on, as standalone's service counts it: a tile laid against
// side joint of the one before faces it with its own first side.
func turnOf(sides, joint int) int {
	turn := (2*joint + sides) % (2 * sides)
	if turn > sides {
		turn -= 2 * sides
	}
	return turn
}

// tilesLaidOut says whether tiles lie as their layout says: a row runs straight on and holds at
// most 6 tiles, or 4 hexagons; a ring goes round once and closes.
func tilesLaidOut(sides, joint, count int, ring bool) bool {
	turn := turnOf(sides, joint)
	if ring {
		return count*max(turn, -turn) == 2*sides
	}
	most := 6
	if sides == 6 {
		most = 4
	}
	return turn == 0 && count <= most
}

// placedTile is a tile laid out, in units of the distance from a tile's middle to its corners: its
// middle, and the direction its first side faces, in radians clockwise from the right.
type placedTile struct {
	x, y, turn float64
}

// laidOut are the tiles laid out: the first facing left with its first side, each after it laid
// against side joint of the one before and turned as turnOf says.
func (t *tilesPicture) laidOut() []placedTile {
	sides := float64(*t.Sides)
	across := 2 * math.Cos(math.Pi/sides)
	step := float64(turnOf(*t.Sides, t.joint())) * math.Pi / sides
	placed := []placedTile{{turn: math.Pi}}
	for range len(t.Tiles) - 1 {
		before := placed[len(placed)-1]
		angle := before.turn + float64(t.joint())*2*math.Pi/sides
		placed = append(placed, placedTile{
			x: before.x + across*math.Cos(angle), y: before.y + across*math.Sin(angle), turn: before.turn + step,
		})
	}
	return placed
}

// corner is a tile's corner where its side k ends, in units.
func (p placedTile) corner(sides, k int) (float64, float64) {
	angle := p.turn + (float64(k)+0.5)*2*math.Pi/float64(sides)
	return p.x + math.Cos(angle), p.y + math.Sin(angle)
}

// The geometry of tiles, in CSS px: the largest a tile's radius is drawn, and the widest and
// tallest the tiles may be.
const (
	tileRadius = 44
	tilesWide  = 520
	tilesTall  = 300
)

// frame is the tiles' unit in CSS px and the corner of their bounds, in units.
func (t *tilesPicture) frame() (unit, minX, minY, spanX, spanY float64) {
	minX, minY = math.Inf(1), math.Inf(1)
	maxX, maxY := math.Inf(-1), math.Inf(-1)
	for _, p := range t.laidOut() {
		for k := range *t.Sides {
			x, y := p.corner(*t.Sides, k)
			minX, maxX, minY, maxY = min(minX, x), max(maxX, x), min(minY, y), max(maxY, y)
		}
	}
	spanX, spanY = maxX-minX, maxY-minY
	return min(tileRadius, tilesWide/spanX, tilesTall/spanY), minX, minY, spanX, spanY
}

func (t *tilesPicture) size() (int, int) {
	unit, _, _, spanX, spanY := t.frame()
	return 2*margin + int(math.Ceil(spanX*unit)), 2*margin + int(math.Ceil(spanY*unit))
}

func (t *tilesPicture) paint(c canvas) {
	unit, minX, minY, _, _ := t.frame()
	at := func(x, y float64) (float64, float64) { return margin + (x-minX)*unit, margin + (y-minY)*unit }
	sides := *t.Sides
	size := min(fontSize, unit*0.36)
	for i, p := range t.laidOut() {
		mx, my := at(p.x, p.y)
		var corners [][2]float64
		for k := range sides {
			x, y := at(p.corner(sides, k))
			corners = append(corners, [2]float64{x, y})
			c.line(mx, my, x, y, 1)
		}
		c.polygon(corners, stroke)
		apothem := math.Cos(math.Pi/float64(sides)) * unit
		for k, text := range t.Tiles[i] {
			if text == "" {
				continue
			}
			angle := p.turn + float64(k)*2*math.Pi/float64(sides)
			c.textAt(text, mx+0.62*apothem*math.Cos(angle), my+0.62*apothem*math.Sin(angle), size)
		}
	}
}
