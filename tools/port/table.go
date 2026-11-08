package main

import "fmt"

// tablePicture is a table, as the schema of standalone's tasks has it: a header row, if it has one,
// and 1 to 8 rows of 1 to 6 cells, every row as long as the header, or as the first row where there
// is none. A cell is a label, a time, … for cells left out, or empty.
type tablePicture struct {
	Kind   string     `json:"kind"`
	Header []string   `json:"header,omitempty"`
	Rows   [][]string `json:"rows"`
}

// ellipsis is what a table writes in a cell for the cells it leaves out.
const ellipsis = "…"

func (t *tablePicture) valid() error {
	if len(t.Rows) < 1 || len(t.Rows) > 8 {
		return fmt.Errorf("a table has 1 to 8 rows, not %d", len(t.Rows))
	}
	rows := t.all()
	width := len(rows[0])
	for i, row := range rows {
		if len(row) < 1 || len(row) > 6 {
			return fmt.Errorf("a row of a table has 1 to 6 cells, not %d", len(row))
		}
		if len(row) != width {
			if t.Header != nil {
				i--
			}
			return fmt.Errorf("the row %d of rows is not as long as the header, or as the first row where there is none", i+1)
		}
		for _, cell := range row {
			if cell != "" && cell != ellipsis && !isLabel(cell) && !timeOfDay.MatchString(cell) {
				return fmt.Errorf("the cell %q is no label, time, %s or empty cell", cell, ellipsis)
			}
		}
	}
	return nil
}

// The geometry of a table, in CSS px.
const (
	cellTall   = 28
	cellNarrow = 36
	tableEdge  = 12
)

// all are the rows of the table, the header first where it has one.
func (t *tablePicture) all() [][]string {
	if t.Header == nil {
		return t.Rows
	}
	return append([][]string{t.Header}, t.Rows...)
}

// columns are how wide each column is: its widest cell with room around it, or a narrow column.
func (t *tablePicture) columns() []float64 {
	widths := make([]float64, len(t.Rows[0]))
	for _, row := range t.all() {
		for i, cell := range row {
			widths[i] = max(widths[i], cellNarrow, measure(cell, fontSize)+16)
		}
	}
	return widths
}

func (t *tablePicture) size() (int, int) {
	width := 0.0
	for _, column := range t.columns() {
		width += column
	}
	return 2*tableEdge + int(width+0.5), 2*tableEdge + cellTall*len(t.all())
}

func (t *tablePicture) paint(c canvas) {
	columns := t.columns()
	width, height := c.size()
	left, top := float64(tableEdge), float64(tableEdge)
	right, bottom := width-tableEdge, height-tableEdge
	if t.Header != nil {
		c.rect(left, top, right, top+cellTall, paper)
	}
	for i, row := range t.all() {
		y := top + float64(i*cellTall)
		if i > 0 {
			c.rect(left, y-0.5, right, y+0.5, ink)
		}
		x := left
		for j, cell := range row {
			c.text(cell, x+columns[j]/2, y+cellTall/2)
			x += columns[j]
		}
	}
	x := left
	for _, column := range columns[:len(columns)-1] {
		x += column
		c.rect(x-0.5, top, x+0.5, bottom, ink)
	}
	if t.Header != nil {
		c.rect(left, top+cellTall-1, right, top+cellTall+1, ink)
	}
	c.frame(left, top, right, bottom, 1.5)
}
