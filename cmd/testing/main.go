package main

import (
	"fmt"

	"github.com/pai0id/CgCourseProject/internal/drawer"
	"github.com/pai0id/CgCourseProject/internal/drawer/mapping"
	"github.com/pai0id/CgCourseProject/internal/fontparser"
)

func main() {
	mctx := mapping.NewContext(11, 11, 4, 4, 44)
	f, err := fontparser.GetFontMap("fonts/IBM.ttf", 44, 44, 20, 144)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	dctx := drawer.NewDrawContext()
	dctx.SetBrightnessMap(f)
	dctx.SetShapeMap(mctx, f)

	canvas := drawer.NewImage(44*40, 44*100)
	for x := range canvas {
		if (x-880)*(x-880)/100 < len(canvas[0]) {
			canvas[x][(x-880)*(x-880)/100] = drawer.Pixel{Brightness: mapping.MaxBrigtness, IsLine: true}
		}
	}

	cells, err := drawer.SplitToCells(canvas, 44, 44)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	cells.Draw(dctx)
}
