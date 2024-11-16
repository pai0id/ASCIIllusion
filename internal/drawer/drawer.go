package drawer

import (
	"fmt"
	"math"

	"github.com/pai0id/CgCourseProject/internal/drawer/fontparser"
	"github.com/pai0id/CgCourseProject/internal/drawer/mapping"
)

type DrawContext struct {
	brightnessMap []int
	shapeMap      []mapping.DescriptionVector
	shapeContext  *mapping.ApproximationContext
}

type Pixel struct {
	brightness int // [0, 100]
	isLine     bool
}

type Image [][]Pixel

type CellInfo struct {
	isLine     bool
	cell       [][]bool
	brightness int // [0, 100]
}

type Canvas [][]CellInfo

func (c CellInfo) GetData() [][]bool {
	return c.cell
}

func FontMapToCellSlice(cms []fontparser.CharMatrix) []mapping.Cell {
	cells := make([]mapping.Cell, len(cms))
	for i, cm := range cms {
		cells[i] = cm
	}
	return cells
}

func SplitToCells(img Image, cellWidth, cellHeight int) (Canvas, error) {
	n := len(img)
	m := len(img[0])
	if n%cellWidth != 0 || m%cellHeight != 0 {
		return nil, fmt.Errorf("error: cant split image %dx%d to cells %dx%d", n, m, cellWidth, cellHeight)
	}

	rowBlocks := n / cellWidth
	colBlocks := m / cellHeight

	cells := make([][]CellInfo, rowBlocks)
	for i := 0; i < rowBlocks; i++ {
		cells[i] = make([]CellInfo, colBlocks)
		for j := 0; j < colBlocks; j++ {
			cells[i][j] = CellInfo{cell: make([][]bool, cellWidth), isLine: false, brightness: 0}
			for k := 0; k < m; k++ {
				cells[i][j].cell[k] = make([]bool, cellHeight)
				for l := 0; l < n; l++ {
					p := img[i*m+k][j*n+l]
					if p.isLine {
						cells[i][j].isLine = true
					}
					cells[i][j].cell[k][l] = p.brightness > 0
					cells[i][j].brightness += p.brightness
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
