package asciiser

import (
	"fmt"
	"math"

	"github.com/pai0id/CgCourseProject/internal/asciiser/mapping"
	"github.com/pai0id/CgCourseProject/internal/fontparser"
)

type DrawContext struct {
	brightnessMap map[fontparser.Char]int
	shapeMap      map[fontparser.Char]mapping.DescriptionVector
	shapeContext  *mapping.ApproximationContext
	bg            fontparser.Char
}

func NewDrawContext(ctx *mapping.ApproximationContext, fm map[fontparser.Char]fontparser.CharMatrix) *DrawContext {
	dc := DrawContext{}
	dc.SetBrightnessMap(fm)
	dc.SetShapeMap(ctx, fm)
	return &dc
}

func (c *DrawContext) SetBrightnessMap(fm map[fontparser.Char]fontparser.CharMatrix) {
	c.brightnessMap = mapping.GetBrightnessMap(fontMapToCellMap(fm))
	var minch fontparser.Char
	var minb int = mapping.MaxBrigtness + 1
	for i, v := range c.brightnessMap {
		if v < minb {
			minb = v
			minch = i
		}
	}
	c.bg = minch
}

func (c *DrawContext) SetShapeMap(ctx *mapping.ApproximationContext, fm map[fontparser.Char]fontparser.CharMatrix) {
	c.shapeContext = ctx
	bgMtx, ok := fm[c.bg]
	delete(fm, c.bg)
	c.shapeMap = mapping.GetShapeMap(ctx, fontMapToCellMap(fm))
	if ok {
		fm[c.bg] = bgMtx
	}
}

type Pixel struct {
	Brightness int // [0, 100]
	IsLine     bool
}

type Image [][]Pixel

func NewImage(width, height int) Image {
	res := make([][]Pixel, height)
	for i := range res {
		res[i] = make([]Pixel, width)
	}
	return res
}

type CellInfo struct {
	isLine     bool
	cell       [][]bool
	brightness int // [0, 100]
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
			cells[i][j] = CellInfo{cell: make([][]bool, cellHeight), isLine: false, brightness: 0}
			for k := range cells[i][j].cell {
				cells[i][j].cell[k] = make([]bool, cellWidth)
				for l := range cells[i][j].cell[k] {
					p := img[i*cellHeight+k][j*cellWidth+l]
					if p.IsLine {
						cells[i][j].isLine = true
						cells[i][j].cell[k][l] = true
					} else {
						cells[i][j].cell[k][l] = p.Brightness > 0
						cells[i][j].brightness += p.Brightness
					}
				}
			}
			if cells[i][j].isLine {
				cells[i][j].brightness = mapping.MaxBrigtness
			} else {
				cells[i][j].cell = nil
				cells[i][j].brightness = cells[i][j].brightness / (cellWidth * cellHeight)
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
			var mindelt int = math.MaxInt
			if c[i][j].isLine {
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
			} else {
				for ch, b := range ctx.brightnessMap {
					d := absInt(c[i][j].brightness - b)
					if d < mindelt {
						mindelt = b
						minch = ch
					}
				}
			}
			res[i][j] = minch
		}
	}
	return res, nil
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
