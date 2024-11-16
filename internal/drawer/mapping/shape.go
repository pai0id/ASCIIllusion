package mapping

import (
	"fmt"
	"math"

	"github.com/pai0id/CgCourseProject/internal/fontparser"
)

type Hyst [][]bool
type DescriptionVector []Hyst

type polarSector struct {
	rMin, rMax, thetaMin, thetaMax float64
	ox, oy                         int
}

type ApproximationContext struct {
	horParts  int
	vertParts int
	phiParts  int
	rParts    int
	rMax      float64
}

func NewContext(horParts, vertParts, phiParts, rParts int, rMax float64) *ApproximationContext {
	return &ApproximationContext{
		horParts:  horParts,
		vertParts: vertParts,
		phiParts:  phiParts,
		rParts:    rParts,
		rMax:      rMax,
	}
}

func getVal(c Cell, ox, oy, x, y int) bool {
	cellData := c.GetData()
	i := ox + x
	j := oy + y

	if i < 0 || i >= len(cellData) || j < 0 || j >= len(cellData[0]) {
		return false
	}
	return cellData[i][j]
}

func isPointInSector(x, y int, s polarSector) bool {
	r := math.Sqrt(float64(x*x + y*y))
	if r < s.rMin || r > s.rMax {
		return false
	}

	theta := math.Atan2(float64(y), float64(x))
	if theta < 0 {
		theta += 2 * math.Pi
	}

	return theta >= s.thetaMin && theta <= s.thetaMax
}

func arePointsInSector(s polarSector, c Cell) bool {
	xMin := int(math.Floor(s.rMin * math.Cos(s.thetaMax)))
	xMax := int(math.Ceil(s.rMax * math.Cos(s.thetaMin)))
	yMin := int(math.Floor(s.rMin * math.Sin(s.thetaMin)))
	yMax := int(math.Ceil(s.rMax * math.Sin(s.thetaMax)))

	for x := xMin; x <= xMax; x++ {
		for y := yMin; y <= yMax; y++ {
			if getVal(c, s.ox, s.oy, x, y) && isPointInSector(x, y, s) {
				return true
			}
		}
	}

	return false
}

func GetDescriptionVector(ctx *ApproximationContext, c Cell) DescriptionVector {
	cellData := c.GetData()
	resVector := make(DescriptionVector, ctx.horParts*ctx.vertParts)
	for i := range resVector {
		resVector[i] = make(Hyst, ctx.phiParts)
		for j := range resVector[i] {
			resVector[i][j] = make([]bool, ctx.rParts)
		}
	}

	n := float64(len(cellData[0]))
	m := float64(len(cellData))

	nStep := n / float64(ctx.horParts)
	mStep := m / float64(ctx.vertParts)

	rStep := ctx.rMax / float64(ctx.rParts)

	phiStep := 2 * math.Pi / float64(ctx.phiParts)

	currHystId := 0
	for i := 0.0; i < n; i += nStep {
		for j := 0.0; j < m; j += mStep {
			phi := 0.0
			for k := 0; k < ctx.phiParts; k++ {
				phiNext := phi + phiStep
				r := 0.0
				for q := 0; q < ctx.rParts; q++ {
					rNext := r + rStep

					s := polarSector{
						rMin:     r,
						rMax:     rNext,
						thetaMin: phi,
						thetaMax: phiNext,
						ox:       int(i),
						oy:       int(j),
					}
					resVector[currHystId][k][q] = arePointsInSector(s, c)

					r = rNext
				}
				phi = phiNext
			}
			currHystId++
		}
	}

	return resVector
}

func GetShapeMap(ctx *ApproximationContext, cArr map[fontparser.Char]Cell) map[fontparser.Char]DescriptionVector {
	res := make(map[fontparser.Char]DescriptionVector, len(cArr))
	for ch, c := range cArr {
		res[ch] = GetDescriptionVector(ctx, c)
	}
	return res
}

func subHyst(a, b Hyst) (int, error) {
	if len(a) != len(b) || len(a[0]) != len(b[0]) {
		return 0, fmt.Errorf("error: hysts have different dimensions")
	}
	result := 0
	rows := len(a)
	cols := len(a[0])

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			result += absInt(boolToInt(a[i][j]) - boolToInt(b[i][j]))
		}
	}

	return result, nil
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func GetVectorDelt(dv1, dv2 DescriptionVector) (int, error) {
	if len(dv1) != len(dv2) {
		return 0, fmt.Errorf("error: description vectors have different dimensions")
	}

	result := 0

	for i := range dv1 {
		subRes, err := subHyst(dv1[i], dv2[i])
		if err != nil {
			return 0, fmt.Errorf("error: %w", err)
		}
		result += subRes
	}

	return result, nil
}
