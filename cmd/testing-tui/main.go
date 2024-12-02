package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pai0id/CgCourseProject/internal/asciiser"
	"github.com/pai0id/CgCourseProject/internal/asciiser/mapping"
	"github.com/pai0id/CgCourseProject/internal/fontparser"
	"github.com/pai0id/CgCourseProject/internal/reader"
	"github.com/pai0id/CgCourseProject/internal/renderer"
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

	resizeCmd := func(m tui.TUIModel, w, h int) (tea.Model, tea.Cmd) {
		options := renderer.RenderOptions{
			Width:  48 * w,
			Height: 66 * h,
			Fov:    60,
		}
		options.CameraDist = renderer.OptimalCameraDist(obj, options)

		canvas := renderer.RenderModel(obj, options)

		cells, err := asciiser.SplitToCells(canvas, 48, 66)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return nil, tea.Quit
		}

		mtx, err := cells.ConvertToASCII(dctx)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return nil, tea.Quit
		}

		m.SetMatrix(mtx)
		return m, nil
	}
	initialModel := tui.NewModel(mtx, make(map[string]tui.Command, 10), resizeCmd)

	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	initialModel.AppendCommand("q", func(m tui.TUIModel) (tea.Model, tea.Cmd) { return m, tea.Quit })

	if _, err := p.Run(); err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
}
