package visualiser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/pai0id/CgCourseProject/internal/asciiser"
	"github.com/pai0id/CgCourseProject/internal/asciiser/mapping"
	"github.com/pai0id/CgCourseProject/internal/fontparser"
	"github.com/pai0id/CgCourseProject/internal/reader"
	"github.com/pai0id/CgCourseProject/internal/renderer"
	"github.com/pai0id/CgCourseProject/internal/transformer"
)

const fov = 60

type visualiserConfig struct {
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

type Visualiser struct {
	objs          []*reader.Model
	dctx          *asciiser.DrawContext
	renderOptions *renderer.RenderOptions
	cfg           *visualiserConfig
}

func NewVisualiser(cfgFileName, sliceFileName, fontFileName string) (*Visualiser, error) {
	v := &Visualiser{}
	err := v.readConfig(cfgFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	chars, err := reader.ReadCharsJson(sliceFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read slice: %w", err)
	}

	f, err := fontparser.GetFontMap(
		fontFileName,
		v.cfg.ImgWidth,
		v.cfg.ImgHeight,
		v.cfg.FontSize,
		v.cfg.DPI,
		chars)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}
	mctx := mapping.NewContext(
		v.cfg.HorParts,
		v.cfg.VertParts,
		v.cfg.PhiParts,
		v.cfg.RParts,
		v.cfg.RMax)
	dctx := asciiser.NewDrawContext(mctx, f)

	v.dctx = dctx
	v.objs = make([]*reader.Model, 0, 10)
	v.renderOptions = &renderer.RenderOptions{Fov: fov, LightSources: make([]reader.Vertex, 0, 10)}

	return v, nil
}

func (v *Visualiser) readConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	err = json.Unmarshal(data, &v.cfg)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	return nil
}

func (v *Visualiser) AddObj(obj *reader.Model) {
	v.objs = append(v.objs, obj)
}

func (v *Visualiser) TranslateObj(id int, tx, ty, tz float64) {
	transformer.Translate(v.objs[id-1], tx, ty, tz)
}

func (v *Visualiser) ScaleObj(id int, sx, sy, sz float64) {
	transformer.Scale(v.objs[id-1], sx, sy, sz)
}

func (v *Visualiser) RotateObj(id int, angle float64, axis int) {
	transformer.Rotate(v.objs[id-1], angle, axis)
}

func (v *Visualiser) Reconvert() (asciiser.ASCIImtx, error) {
	canvas := renderer.RenderModels(v.objs, v.renderOptions)

	cells, err := asciiser.SplitToCells(canvas, v.cfg.ImgWidth, v.cfg.ImgHeight)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	if cells == nil {
		return nil, nil
	} else {
		mtx, err := cells.ConvertToASCII(v.dctx)
		if err != nil {
			return nil, fmt.Errorf("error: %w", err)
		}
		return mtx, nil
	}
}

func (v *Visualiser) AddLightSource(x, y, z float64) {
	v.renderOptions.LightSources = append(v.renderOptions.LightSources, reader.Vertex{X: x, Y: y, Z: z})
}

func (v *Visualiser) Resize(w, h int) {
	v.renderOptions.Width = v.cfg.ImgWidth * w
	v.renderOptions.Height = v.cfg.ImgHeight * h

	v.renderOptions.CameraDist = renderer.OptimalCameraDist(v.objs, v.renderOptions)
}

func (v *Visualiser) OptimizeCamera() {
	v.renderOptions.CameraDist = renderer.OptimalCameraDist(v.objs, v.renderOptions)
}
