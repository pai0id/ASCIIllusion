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

const fov = 10
const zNear = 0.01
const zFar = 100.0

type visualiserConfig struct {
	CharWidth  int     `json:"charWidth"`
	CharHeight int     `json:"charHeight"`
	FontSize   float64 `json:"fontSize"`
	DPI        float64 `json:"dpi"`
	HorParts   int     `json:"horParts"`
	VertParts  int     `json:"vertParts"`
	PhiParts   int     `json:"phiParts"`
	RParts     int     `json:"rParts"`
	RMax       float64 `json:"rMax"`
}

type Visualiser struct {
	objs          []*reader.Model
	ids           []int64
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
		v.cfg.CharWidth,
		v.cfg.CharHeight,
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
	v.ids = make([]int64, 0, 10)
	v.renderOptions = &renderer.RenderOptions{
		LightSources:    make([]reader.Vec3, 0, 10),
		LightSourcesIds: make([]int64, 0, 10),
	}
	v.renderOptions.Cam = &renderer.Camera{
		Fov:    fov,
		ZNear:  zNear,
		ZFar:   zFar,
		Aspect: float64(v.cfg.CharWidth) / float64(v.cfg.CharHeight),
	}

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

func (v *Visualiser) AddObj(obj *reader.Model, id int64) {
	v.objs = append(v.objs, obj)
	v.ids = append(v.ids, id)
}

func (v *Visualiser) DeleteObj(id int64) error {
	for i, vid := range v.ids {
		if vid == id {
			v.objs = append(v.objs[:i], v.objs[i+1:]...)
			v.ids = append(v.ids[:i], v.ids[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("object with id %d not found", id)
}

func (v *Visualiser) TranslateObj(id int64, tx, ty, tz float64) error {
	for i, vid := range v.ids {
		if vid == id {
			transformer.Translate(v.objs[i], tx, ty, tz)
			return nil
		}
	}
	return fmt.Errorf("object with id %d not found", id)
}

func (v *Visualiser) ScaleObj(id int64, sx, sy, sz float64) error {
	for i, vid := range v.ids {
		if vid == id {
			transformer.Scale(v.objs[i], sx, sy, sz)
			return nil
		}
	}
	return fmt.Errorf("object with id %d not found", id)
}

func (v *Visualiser) RotateObj(id int64, angle float64, axis int) error {
	for i, vid := range v.ids {
		if vid == id {
			transformer.Rotate(v.objs[i], angle, axis)
			return nil
		}
	}
	return fmt.Errorf("object with id %d not found", id)
}

func (v *Visualiser) Reconvert() (asciiser.ASCIImtx, error) {
	canvas := renderer.RenderModels(v.objs, v.renderOptions)

	cells, err := asciiser.SplitToCells(canvas, v.cfg.CharWidth, v.cfg.CharHeight)
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

func (v *Visualiser) AddLightSource(x, y, z float64, id int64) {
	v.renderOptions.LightSources = append(v.renderOptions.LightSources, reader.Vec3{X: x, Y: y, Z: z})
	v.renderOptions.LightSourcesIds = append(v.renderOptions.LightSourcesIds, id)
}

func (v *Visualiser) DeleteLightSource(id int64) error {
	for i, lid := range v.renderOptions.LightSourcesIds {
		if lid == id {
			v.renderOptions.LightSources = append(v.renderOptions.LightSources[:i], v.renderOptions.LightSources[i+1:]...)
			v.renderOptions.LightSourcesIds = append(v.renderOptions.LightSourcesIds[:i], v.renderOptions.LightSourcesIds[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("light source with id %d not found", id)
}

func (v *Visualiser) Resize(w, h int) {
	v.renderOptions.Width = v.cfg.CharWidth * w
	v.renderOptions.Height = v.cfg.CharHeight * h
}

func (v *Visualiser) OptimizeCamera() {
	if len(v.objs) == 1 {
		z := renderer.OptimalCameraDist(v.objs[0], v.renderOptions.Cam.Fov, v.renderOptions.Cam.Aspect)
		v.renderOptions.Cam.Z = z
	}
}

func (v *Visualiser) MoveCam(d float64) {
	v.renderOptions.Cam.Z += d
}
