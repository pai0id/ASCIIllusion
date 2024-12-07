package asciiser

import (
	"fmt"
	"math"

	"github.com/pai0id/CgCourseProject/internal/asciiser/mapping"
	"github.com/pai0id/CgCourseProject/internal/fontparser"
)

type cellType [][]bool

func (c cellType) clean() {
	for i := range c {
		for j := range c[i] {
			c[i][j] = false
		}
	}
}

type CellInfo struct {
	isLine     bool
	isPolygon  bool
	cell       cellType
	brightness float64
}

type Canvas [][]CellInfo

func (c CellInfo) GetData() [][]bool {
	return c.cell
}

func fontMapToCellMap(cms map[fontparser.Char]fontparser.CharMatrix) map[fontparser.Char]mapping.Cell {
	cells := make(map[fontparser.Char]mapping.Cell, len(cms))
	for ch, cm := range cms {
		cells[ch] = cm
	}
	return cells
}

func SplitToCells(img Image, cellWidth, cellHeight int) (Canvas, error) {
	if img == nil {
		return nil, nil
	}
	n := len(img)    // высота
	m := len(img[0]) // ширина
	if n%cellHeight != 0 || m%cellWidth != 0 {
		return nil, fmt.Errorf("error: cant split image %dx%d to cells %dx%d", n, m, cellHeight, cellWidth)
	}

	rowBlocks := m / cellWidth
	colBlocks := n / cellHeight

	cells := make([][]CellInfo, colBlocks)
	for i := range cells {
		cells[i] = make([]CellInfo, rowBlocks)
		for j := range cells[i] {
			cells[i][j] = CellInfo{cell: make([][]bool, cellHeight)}
			for k := range cells[i][j].cell {
				cells[i][j].cell[k] = make([]bool, cellWidth)
				for l := range cells[i][j].cell[k] {
					p := img[i*cellHeight+k][j*cellWidth+l]
					if cells[i][j].isLine {
						cells[i][j].cell[k][l] = p.IsLine
					} else if p.IsLine {
						cells[i][j].isLine = true
						cells[i][j].cell.clean()
						cells[i][j].cell[k][l] = true
					} else if p.IsPolygon {
						cells[i][j].isPolygon = true
						cells[i][j].cell[k][l] = p.Brightness > 0
						cells[i][j].brightness += p.Brightness
					}
				}
			}
			if cells[i][j].isLine {
				cells[i][j].brightness = 1
			} else {
				cells[i][j].cell = nil
				cells[i][j].brightness = cells[i][j].brightness / float64(cellWidth*cellHeight)
			}
		}
	}
	return cells, nil
}

type ASCIImtx [][]fontparser.Char

func (c Canvas) ConvertToASCII(ctx *DrawContext) (ASCIImtx, error) {
	res := make([][]fontparser.Char, len(c))
	for i := range c {
		res[i] = make([]fontparser.Char, len(c[i]))
	}
	for i := range c {
		for j := range c[i] {
			var minch fontparser.Char
			if c[i][j].isLine {
				var mindelt int = math.MaxInt
				dv := mapping.GetDescriptionVector(ctx.shapeContext, c[i][j])
				for ch, dvf := range ctx.shapeMap {
					d, err := mapping.GetVectorDelt(dv, dvf)
					if err != nil {
						return nil, fmt.Errorf("error: could not draw Canvas: %v", err)
					}
					if d < mindelt {
						mindelt = d
						minch = ch
					}
				}
			} else if c[i][j].isPolygon {
				var mindelt float64 = math.MaxFloat64
				for ch, b := range ctx.brightnessMap {
					d := math.Abs(c[i][j].brightness - b)
					if d < mindelt {
						mindelt = d
						minch = ch
					}
				}
			} else {
				minch = ctx.bg
			}
			res[i][j] = minch
		}
	}
	return res, nil
}
