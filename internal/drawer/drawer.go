package drawer

import (
	"fmt"
	"math"

	"github.com/pai0id/CgCourseProject/internal/drawer/mapping"
	"github.com/pai0id/CgCourseProject/internal/fontparser"
)

type DrawContext struct {
	brightnessMap []int
	shapeMap      []mapping.DescriptionVector
	shapeContext  *mapping.ApproximationContext
}

func NewDrawContext() *DrawContext {
	return &DrawContext{}
}

func (c *DrawContext) SetBrightnessMap(fm []fontparser.CharMatrix) {
	c.brightnessMap = mapping.GetBrightnessMap(fontMapToCellSlice(fm))
}

func (c *DrawContext) SetShapeMap(ctx *mapping.ApproximationContext, fm []fontparser.CharMatrix) {
	c.shapeContext = ctx
	c.shapeMap = mapping.GetShapeMap(ctx, fontMapToCellSlice(fm))
}

type Pixel struct {
	Brightness int // [0, 100]
	IsLine     bool
}

type Image [][]Pixel

func NewImage(width, height int) Image {
	res := make([][]Pixel, width)
	for i := 0; i < width; i++ {
		res[i] = make([]Pixel, height)
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

func fontMapToCellSlice(cms []fontparser.CharMatrix) []mapping.Cell {
	cells := make([]mapping.Cell, len(cms))
	for i, cm := range cms {
		cells[i] = cm
	}
	return cells
}

func SplitToCells(img Image, cellWidth, cellHeight int) (Canvas, error) {
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
					}
					cells[i][j].cell[k][l] = p.Brightness > 0
					cells[i][j].brightness += p.Brightness
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

func (c Canvas) Draw(ctx *DrawContext) error {
	for i := range c {
		for j := range c[i] {
			var minid int
			var mindelt int = math.MaxInt
			if c[i][j].isLine {
				dv := mapping.GetDescriptionVector(ctx.shapeContext, c[i][j])
				for id, dvf := range ctx.shapeMap {
					d, err := mapping.GetVectorDelt(dv, dvf)
					if err != nil {
						return fmt.Errorf("error: could not draw Canvas: %w", err)

					}
					if d < mindelt {
						mindelt = d
						minid = id
					}
				}
			} else {
				for id, b := range ctx.brightnessMap {
					d := absInt(c[i][j].brightness - b)
					if d < mindelt {
						mindelt = b
						minid = id
					}
				}
			}
			c, err := fontparser.GetChar(minid)
			if err != nil {
				return fmt.Errorf("error: could not draw Canvas: %w", err)
			}
			fmt.Printf("%c", c)
		}
		fmt.Println()
	}
	return nil
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
