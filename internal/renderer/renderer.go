package renderer

import (
	"math"
	"sync"

	"github.com/pai0id/CgCourseProject/internal/asciiser"
	"github.com/pai0id/CgCourseProject/internal/reader"
	"github.com/pai0id/CgCourseProject/internal/transformer"
)

type point struct {
	x, y int
	z    float64
}

type normal struct {
	x, y, z float64
}

type face struct {
	vertices    []point
	normals     []normal
	skeletonize bool
}

type polygon map[point]asciiser.Pixel

type Camera struct {
	Fov   float64
	Z     float64
	ZFar  float64
	ZNear float64
}

type RenderOptions struct {
	Width           int
	Height          int
	Cam             *Camera
	LightSources    []reader.Vec3
	LightSourcesIds []int64
}

type zBuffer [][]float64

func newZBuffer(width, height int) zBuffer {
	zBuffer := make([][]float64, height)
	for i := range zBuffer {
		zBuffer[i] = make([]float64, width)
		for j := range zBuffer[i] {
			zBuffer[i][j] = math.MaxFloat64
		}
	}
	return zBuffer
}

func RenderModels(models []*reader.Model, options *RenderOptions) asciiser.Image {
	if options.Width == 0 || options.Height == 0 {
		return nil
	}
	image := asciiser.NewImage(options.Width, options.Height)
	zb := newZBuffer(options.Width, options.Height)

	viewMatrix := transformer.ViewMatrix(options.Cam.Z)
	aspect := float64(options.Width) / float64(options.Height)
	projectionMatrix := transformer.PerspectiveMatrix(options.Cam.Fov, aspect, options.Cam.ZNear, options.Cam.ZFar)

	projectQueue := make(chan *reader.Model, 10)
	clippingQueue := make(chan *reader.Model, 10)
	screenQueue := make(chan *reader.Model, 10)
	gouraudQueue := make(chan *face, 100)
	resQueue := make(chan polygon, 100)

	var polygons []polygon

	var wg sync.WaitGroup

	wg.Add(1)
	go start(models, projectQueue, &wg)

	wg.Add(1)
	go project(projectQueue, clippingQueue, &wg, viewMatrix, projectionMatrix)

	wg.Add(1)
	go clip(clippingQueue, screenQueue, &wg, options.Cam.ZNear, options.Cam.ZFar)

	wg.Add(1)
	go screen(screenQueue, gouraudQueue, &wg, options.Width, options.Height)

	wg.Add(1)
	go shading(gouraudQueue, resQueue, &wg, options.LightSources)

	wg.Add(1)
	go end(resQueue, &polygons, &wg)

	wg.Wait()

	render(polygons, image, zb)
	return image
}

func start(models []*reader.Model, out chan<- *reader.Model, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, m := range models {
		out <- m
	}

	if out != nil {
		close(out)
	}
}

func end(in <-chan polygon, res *[]polygon, wg *sync.WaitGroup) {
	defer wg.Done()

	for p := range in {
		*res = append(*res, p)
	}
}
