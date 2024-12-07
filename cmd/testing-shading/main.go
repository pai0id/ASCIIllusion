package main

import (
	"fmt"

	"github.com/pai0id/CgCourseProject/internal/asciiser"
	"github.com/pai0id/CgCourseProject/internal/asciiser/mapping"
	"github.com/pai0id/CgCourseProject/internal/fontparser"
	"github.com/pai0id/CgCourseProject/internal/reader"
)

func main() {
	mctx := mapping.NewContext(8, 10, 4, 5, 5)
	chars, err := reader.ReadCharsTxt("fonts/slice.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	f, err := fontparser.GetFontMap("fonts/IBM.ttf", 8, 10, 8, 72, chars)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	dctx := asciiser.NewDrawContext(mctx, f)

	canvas := asciiser.NewImage(8*100, 10*200)
	for x := range canvas {
		for y := range canvas[x] {
			canvas[x][y].Brightness = float64(x) / float64(len(canvas))
		}
	}

	cells, err := asciiser.SplitToCells(canvas, 8, 10)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	mtx, err := cells.ConvertToASCII(dctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for i := range mtx {
		for j := range mtx[i] {
			fmt.Printf("%c", mtx[i][j])
		}
		fmt.Println()
	}
	fmt.Println()
}
