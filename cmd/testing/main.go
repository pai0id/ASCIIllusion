package main

import (
	"fmt"

	"github.com/pai0id/CgCourseProject/internal/drawer"
	"github.com/pai0id/CgCourseProject/internal/drawer/mapping"
	"github.com/pai0id/CgCourseProject/internal/fontparser"
)

func main() {
	mctx := mapping.NewContext(11, 11, 4, 4, 44)
	chars := make([]fontparser.Char, 0, 128)
	for i := 32; i < 127; i++ {
		chars = append(chars, fontparser.Char(i))
	}
	f, err := fontparser.GetFontMap("fonts/IBM.ttf", 44, 44, 20, 144, chars)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	dctx := drawer.NewDrawContext()
	dctx.SetBrightnessMap(f)
	delete(f, ' ')
	dctx.SetShapeMap(mctx, f)

	canvas := drawer.NewImage(44*40, 44*100)
	for x := range canvas {
		y := (x * x) / 600
		if y >= 0 && y < len(canvas[x]) {
			canvas[x][y].IsLine = true
		}
	}

	cells, err := drawer.SplitToCells(canvas, 44, 44)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	cells.Draw(dctx)
}
