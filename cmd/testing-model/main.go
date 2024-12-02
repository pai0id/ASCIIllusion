package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/pai0id/CgCourseProject/internal/asciiser"
	"github.com/pai0id/CgCourseProject/internal/asciiser/mapping"
	"github.com/pai0id/CgCourseProject/internal/fontparser"
	"github.com/pai0id/CgCourseProject/internal/reader"
	"github.com/pai0id/CgCourseProject/internal/renderer"
	"github.com/pai0id/CgCourseProject/internal/transformer"
	"github.com/pai0id/CgCourseProject/internal/tui"
)

func ClearTerminal() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

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
	// for _, m := range f {
	// 	for i := range m {
	// 		for j := range m[i] {
	// 			if m[i][j] {
	// 				fmt.Print("@")
	// 			} else {
	// 				fmt.Print(" ")
	// 			}
	// 		}
	// 		fmt.Println()
	// 	}
	// 	fmt.Println("-----------------------------")
	// }

	dctx := asciiser.NewDrawContext(mctx, f)

	obj, err := reader.LoadOBJ("data/tetra.obj")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	options := renderer.RenderOptions{
		Width:  48 * 400,
		Height: 66 * 180,
		Fov:    60,
	}
	options.CameraDist = renderer.OptimalCameraDist(obj, options)

	canvas := renderer.RenderModel(obj, options)

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

	for {
		ClearTerminal()
		transformer.Rotate(obj, 10, transformer.YAxis)
		transformer.Rotate(obj, 10, transformer.ZAxis)

		canvas := renderer.RenderModel(obj, options)

		cells, err := asciiser.SplitToCells(canvas, 48, 66)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		mtx, err = cells.ConvertToASCII(dctx)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		tui.DrawASCIImtr(mtx)
		time.Sleep(time.Second)
	}
}
