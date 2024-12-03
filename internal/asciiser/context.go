package asciiser

import (
	"github.com/pai0id/CgCourseProject/internal/asciiser/mapping"
	"github.com/pai0id/CgCourseProject/internal/fontparser"
)

type DrawContext struct {
	brightnessMap map[fontparser.Char]float64
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
	var minb float64 = 1
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
