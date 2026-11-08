package main

import (
	"fmt"
	"maps"
	"slices"
)

// gridPicture is a grid, as the schema of standalone's tasks has it: 1 to 8 rows and 1 to 8
// columns, each named by a label once, the cells it fills, each once, and a label in some of its
// cells, keyed by the cell. A cell is named by its row's name run into its column's, as B2 is.
type gridPicture struct {
	Kind   string            `json:"kind"`
	Rows   []string          `json:"rows"`
	Cols   []string          `json:"cols"`
	Filled []string          `json:"filled,omitempty"`
	Marks  map[string]string `json:"marks,omitempty"`
}

func (g *gridPicture) valid() error {
	for _, names := range [][]string{g.Rows, g.Cols} {
		if len(names) < 1 || len(names) > 8 {
			return fmt.Errorf("a grid has 1 to 8 rows and 1 to 8 columns")
		}
		for i, name := range names {
			if !isLabel(name) {
				return fmt.Errorf("the name %q of a line of the grid is no label", name)
			}
			if slices.Contains(names[:i], name) {
				return fmt.Errorf("the name %q names two lines of the grid", name)
			}
		}
	}
	cells := g.cells()
	if cells == nil {
		return fmt.Errorf("the names of the rows run into those of the columns name two cells alike")
	}
	for i, cell := range g.Filled {
		if _, of := cells[cell]; !of {
			return fmt.Errorf("the cell %q is no cell of the grid: its row and its column, as B2 is", cell)
		}
		if slices.Contains(g.Filled[:i], cell) {
			return fmt.Errorf("the cell %q is filled twice", cell)
		}
	}
	for _, cell := range slices.Sorted(maps.Keys(g.Marks)) {
		text := g.Marks[cell]
		if _, of := cells[cell]; !of {
			return fmt.Errorf("the mark %q is in no cell of the grid", cell)
		}
		if !isLabel(text) {
			return fmt.Errorf("the label %q of the cell %s is none", text, cell)
		}
	}
	return nil
}

// cells are where each cell of the grid stands, by its name, its row and column from 0; nil where
// two cells would have one name, as rows 1 and 12 with columns 22 and 2 would.
func (g *gridPicture) cells() map[string][2]int {
	cells := map[string][2]int{}
	for r, row := range g.Rows {
		for k, col := range g.Cols {
			if _, twice := cells[row+col]; twice {
				return nil
			}
			cells[row+col] = [2]int{r, k}
		}
	}
	return cells
}

// The geometry of a grid, in CSS px.
const (
	cellSide  = 36
	namesGap  = 8
	gridFrame = 12
)

// namesWide is how wide the column of the rows' names is.
func (g *gridPicture) namesWide() float64 {
	wide := 0.0
	for _, name := range g.Rows {
		wide = max(wide, measure(name, fontSize))
	}
	return wide + namesGap
}

func (g *gridPicture) size() (int, int) {
	return 2*gridFrame + int(g.namesWide()+0.5) + len(g.Cols)*cellSide, 2*gridFrame + band + len(g.Rows)*cellSide
}

func (g *gridPicture) paint(c canvas) {
	left, top := gridFrame+float64(int(g.namesWide()+0.5)), float64(gridFrame+band)
	right, bottom := left+float64(len(g.Cols)*cellSide), top+float64(len(g.Rows)*cellSide)
	cells := g.cells()
	for _, cell := range g.Filled {
		at := cells[cell]
		x, y := left+float64(at[1]*cellSide), top+float64(at[0]*cellSide)
		c.rect(x, y, x+cellSide, y+cellSide, shade)
	}
	for i := 1; i < len(g.Rows); i++ {
		y := top + float64(i*cellSide)
		c.rect(left, y-0.75, right, y+0.75, ink)
	}
	for k := 1; k < len(g.Cols); k++ {
		x := left + float64(k*cellSide)
		c.rect(x-0.75, top, x+0.75, bottom, ink)
	}
	c.frame(left, top, right, bottom, stroke)
	for r, name := range g.Rows {
		c.text(name, left-namesGap/2-measure(name, fontSize)/2, top+float64(r*cellSide)+cellSide/2)
	}
	for k, name := range g.Cols {
		c.text(name, left+float64(k*cellSide)+cellSide/2, top-band/2)
	}
	for _, cell := range slices.Sorted(maps.Keys(g.Marks)) {
		at, text := cells[cell], g.Marks[cell]
		c.text(text, left+float64(at[1]*cellSide)+cellSide/2, top+float64(at[0]*cellSide)+cellSide/2)
	}
}
