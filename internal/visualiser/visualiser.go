package visualiser

import (
	"encoding/json"
	"fmt"
	"io"
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

const fov = 60

var visualiserConfig struct {
	ImgWidth  int     `json:"imgWidth"`
	ImgHeight int     `json:"imgHeight"`
	FontSize  float64 `json:"fontSize"`
	DPI       float64 `json:"dpi"`
	HorParts  int     `json:"horParts"`
	VertParts int     `json:"vertParts"`
	PhiParts  int     `json:"phiParts"`
	RParts    int     `json:"rParts"`
	RMax      float64 `json:"rMax"`
}

func readConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	err = json.Unmarshal(data, &visualiserConfig)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	return nil
}

var visualiser struct {
	objs          []*reader.Model
	dctx          *asciiser.DrawContext
	renderOptions *renderer.RenderOptions
	TUIModel      *tui.TUIModel
}

func NewVisualiser(cfgFileName, sliceFileName, fontFileName string) error {
	err := readConfig(cfgFileName)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	chars, err := reader.ReadCharsJson(sliceFileName)
	if err != nil {
		return fmt.Errorf("failed to read slice: %w", err)
	}

	f, err := fontparser.GetFontMap(
		fontFileName,
		visualiserConfig.ImgWidth,
		visualiserConfig.ImgHeight,
		visualiserConfig.FontSize,
		visualiserConfig.DPI,
		chars)
	if err != nil {
		return fmt.Errorf("failed to parse font: %w", err)
	}
	mctx := mapping.NewContext(
		visualiserConfig.HorParts,
		visualiserConfig.VertParts,
		visualiserConfig.PhiParts,
		visualiserConfig.RParts,
		visualiserConfig.RMax)
	dctx := asciiser.NewDrawContext(mctx, f)

	visualiser.dctx = dctx
	visualiser.objs = make([]*reader.Model, 0, 10)
	visualiser.renderOptions = &renderer.RenderOptions{}

	initialModel := tui.NewTUIModel(nil, make(map[string]tui.Command, 10), resize)
	p := tea.NewProgram(initialModel, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
	return nil
}

func AddObj(obj *reader.Model) {
	visualiser.objs = append(visualiser.objs, obj)
}

func Reconvert(m tui.TUIModel) (tui.TUIModel, error) {
	visualiser.renderOptions.CameraDist = renderer.OptimalCameraDist(visualiser.objs, visualiser.renderOptions)

	canvas := renderer.RenderModels(visualiser.objs, visualiser.renderOptions)

	cells, err := asciiser.SplitToCells(canvas, visualiserConfig.ImgWidth, visualiserConfig.ImgHeight)
	if err != nil {
		return m, fmt.Errorf("error: %w", err)
	}

	if cells == nil {
		m.SetMatrix(nil)
	} else {
		mtx, err := cells.ConvertToASCII(visualiser.dctx)
		if err != nil {
			return m, fmt.Errorf("error: %w", err)
		}
		m.SetMatrix(mtx)
	}

	return m, nil
}

func resize(m tui.TUIModel, w, h int) (tui.TUIModel, tea.Cmd) {
	options := renderer.RenderOptions{
		Width:  visualiserConfig.ImgWidth * w,
		Height: visualiserConfig.ImgHeight * h,
		Fov:    fov,
	}
	options.CameraDist = renderer.OptimalCameraDist(visualiser.objs, &options)

	canvas := renderer.RenderModels(visualiser.objs, visualiser.renderOptions)

	cells, err := asciiser.SplitToCells(canvas, visualiserConfig.ImgWidth, visualiserConfig.ImgHeight)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return m, tea.Quit
	}

	if cells == nil {
		m.SetMatrix(nil)
	} else {
		mtx, err := cells.ConvertToASCII(visualiser.dctx)
		if err != nil {
			return m, tea.Quit
		}
		m.SetMatrix(mtx)
	}

	return m, nil
}
