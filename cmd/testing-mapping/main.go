package main

import (
	"fmt"

	"github.com/pai0id/CgCourseProject/internal/asciiser"
	"github.com/pai0id/CgCourseProject/internal/asciiser/mapping"
	"github.com/pai0id/CgCourseProject/internal/fontparser"
	"github.com/pai0id/CgCourseProject/internal/reader"
	"github.com/pai0id/CgCourseProject/internal/tui"
)

func main() {
	mctx := mapping.NewContext(8, 11, 4, 4, 66)
	chars, err := reader.ReadCharsTxt("fonts/slice.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	f, err := fontparser.GetFontMap("fonts/IBM.ttf", 48, 66, 25, 144, chars)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	dctx := asciiser.NewDrawContext(mctx, f)

	canvas := asciiser.NewImage(48*150, 66*30)
	for x := range canvas {
		y := (x * x) / 300
		if y >= 0 && y < len(canvas[x]) {
			canvas[x][y].IsLine = true
		}
	}

	cells, err := asciiser.SplitToCells(canvas, 48, 66)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	mtx, err := cells.ConvertToASCII(dctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	tui.DrawASCIImtr(mtx)
}
