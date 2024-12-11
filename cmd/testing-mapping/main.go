package main

import (
	"fmt"

	"github.com/pai0id/CgCourseProject/internal/asciiser"
	"github.com/pai0id/CgCourseProject/internal/asciiser/mapping"
	"github.com/pai0id/CgCourseProject/internal/fontparser"
	"github.com/pai0id/CgCourseProject/internal/reader"
)

func main() {
	mctx := mapping.NewContext(8, 8, 8, 10, 10)
	chars, err := reader.ReadCharsTxt("fonts/slice.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	f, err := fontparser.GetFontMap("fonts/IBM.ttf", 8, 8, 6, 72, chars)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, ch := range f {
		for i := range ch {
			for j := range ch[i] {
				if ch[i][j] {
					fmt.Print("@")
				} else {
					fmt.Print(" ")
				}
			}
			fmt.Println()
		}
		fmt.Println("--------------------------------")
	}

	dctx := asciiser.NewDrawContext(mctx, f)

	canvas := asciiser.NewImage(8*150, 8*30)
	for x := range canvas {
		y := (x * x) / 10
		if y >= 0 && y < len(canvas[x]) {
			for i := 0; i < 12; i++ {
				canvas[x][y+i].IsLine = true
			}
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
