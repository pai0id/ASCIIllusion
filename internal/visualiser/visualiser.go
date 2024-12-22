package visualiser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pai0id/CgCourseProject/internal/asciiser"
	"github.com/pai0id/CgCourseProject/internal/asciiser/mapping"
	"github.com/pai0id/CgCourseProject/internal/fontparser"
	"github.com/pai0id/CgCourseProject/internal/object"
	"github.com/pai0id/CgCourseProject/internal/reader"
	"github.com/pai0id/CgCourseProject/internal/renderer"
	"github.com/pai0id/CgCourseProject/internal/transformer"
)

const fov = 2
const zNear = 0.01
const zFar = 1000.0

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
	objs          map[int64]*object.Object
	dctx          *asciiser.DrawContext
	renderOptions *renderer.RenderOptions
	cfg           *visualiserConfig
}

func getObjs(objMap map[int64]*object.Object) []*object.Object {
	result := make([]*object.Object, len(objMap))
	i := 0
	for _, obj := range objMap {
		result[i] = obj
		i++
	}
	return result
}

func NewVisualiser(cfgFileName, sliceFileName, fontFileName string) (*Visualiser, error) {
	v := &Visualiser{}
	err := v.readConfig(cfgFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var chars []fontparser.Char
	if strings.HasSuffix(sliceFileName, ".json") {
		chars, err = reader.ReadCharsJson(sliceFileName)
		if err != nil {
			return nil, fmt.Errorf("failed to read slice: %w", err)
		}
	} else if strings.HasSuffix(sliceFileName, ".txt") {
		chars, err = reader.ReadCharsTxt(sliceFileName)
		if err != nil {
			return nil, fmt.Errorf("failed to read slice: %w", err)
		}
	} else {
		return nil, fmt.Errorf("unsupported slice file format: %s", sliceFileName)
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
	v.objs = make(map[int64]*object.Object, 10)
	v.renderOptions = &renderer.RenderOptions{
		LightSources: make(map[int64]renderer.Light, 10),
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

func (v *Visualiser) AddObj(obj *object.Object) int64 {
	var maxId int64 = 0
	for k := range v.objs {
		if k > maxId {
			maxId = k
		}
	}
	maxId++
	v.objs[maxId] = obj
	return maxId
}

func (v *Visualiser) DeleteObj(id int64) error {
	if _, ok := v.objs[id]; ok {
		delete(v.objs, id)
		return nil
	}
	return fmt.Errorf("object with id %d not found", id)
}

func (v *Visualiser) TranslateObj(id int64, tx, ty, tz float64) error {
	if obj, ok := v.objs[id]; ok {
		transformer.Translate(obj, tx, ty, tz)
		return nil
	}
	return fmt.Errorf("object with id %d not found", id)
}

func (v *Visualiser) ScaleObj(id int64, sx, sy, sz float64) error {
	if obj, ok := v.objs[id]; ok {
		transformer.Scale(obj, sx, sy, sz)
		return nil
	}
	return fmt.Errorf("object with id %d not found", id)
}

func (v *Visualiser) RotateObj(id int64, angle float64, axis int) error {
	if obj, ok := v.objs[id]; ok {
		transformer.Rotate(obj, angle, axis)
		return nil
	}
	return fmt.Errorf("object with id %d not found", id)
}

func (v *Visualiser) Reconvert() (asciiser.ASCIImtx, error) {
	canvas := renderer.RenderModels(getObjs(v.objs), v.renderOptions)

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

func (v *Visualiser) AddLightSource(x, y, z, intensity float64) int64 {
	var maxId int64 = 0
	for k := range v.renderOptions.LightSources {
		if k > maxId {
			maxId = k
		}
	}
	maxId++
	v.renderOptions.LightSources[maxId] = renderer.Light{Position: object.Vec3{X: x, Y: y, Z: z}, Intensity: intensity}
	return maxId
}

func (v *Visualiser) DeleteLightSource(id int64) error {
	if _, ok := v.renderOptions.LightSources[id]; ok {
		delete(v.renderOptions.LightSources, id)
		return nil
	}
	return fmt.Errorf("light source with id %d not found", id)
}

func (v *Visualiser) Resize(w, h int) {
	v.renderOptions.Width = v.cfg.CharWidth * w
	v.renderOptions.Height = v.cfg.CharHeight * h
}

func (v *Visualiser) OptimizeCamera() {
	if len(v.objs) == 1 {
		for _, obj := range v.objs {
			z := obj.OptimalCameraDist()
			v.renderOptions.Cam.Z = z
		}
	}
}

func (v *Visualiser) MoveCam(d float64) {
	v.renderOptions.Cam.Z += d
}

func (v *Visualiser) GetCam() float64 {
	return v.renderOptions.Cam.Z
}
