package asciiser

import (
	"sort"

	"github.com/pai0id/CgCourseProject/internal/asciiser/mapping"
	"github.com/pai0id/CgCourseProject/internal/fontparser"
)

func halfMap(input map[fontparser.Char]float64) []fontparser.Char {
	type kv struct {
		Key   fontparser.Char
		Value float64
	}
	var kvSlice []kv
	for k, v := range input {
		kvSlice = append(kvSlice, kv{k, v})
	}

	sort.Slice(kvSlice, func(i, j int) bool {
		return kvSlice[i].Value < kvSlice[j].Value
	})

	res := make([]fontparser.Char, 0, len(kvSlice)/2)
	for i, kv := range kvSlice {
		if i%2 == 0 {
			res = append(res, kv.Key)
		}
	}

	return res
}

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
	delete(c.brightnessMap, c.bg)
	// half := halfMap(c.brightnessMap)
	// for _, ch := range half {
	// 	delete(fm, ch)
	// 	delete(c.brightnessMap, ch)
	// }
	// half = halfMap(c.brightnessMap)
	// for _, ch := range half {
	// 	delete(fm, ch)
	// 	delete(c.brightnessMap, ch)
	// }
	// c.brightnessMap[' '] = 0
}

func (c *DrawContext) SetShapeMap(ctx *mapping.ApproximationContext, fm map[fontparser.Char]fontparser.CharMatrix) {
	c.shapeContext = ctx
	c.shapeMap = mapping.GetShapeMap(ctx, fontMapToCellMap(fm))
	delete(c.shapeMap, c.bg)
}
