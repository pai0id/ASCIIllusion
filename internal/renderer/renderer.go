package renderer

import (
	"math"
	"sync"

	"github.com/pai0id/CgCourseProject/internal/asciiser"
	"github.com/pai0id/CgCourseProject/internal/object"
	"github.com/pai0id/CgCourseProject/internal/transformer"
)

type Point struct {
	X, Y int
	Z    float64
}

type Polygon struct {
	Vertices    []Point
	Normals     []object.Vec3
	Intensities []float64
	Skeletonize bool
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

func RenderModels(models []*object.Object, options *RenderOptions) asciiser.Image {
	if options.Width == 0 || options.Height == 0 {
		return nil
	}
	image := asciiser.NewImage(options.Width, options.Height)
	zb := newZBuffer(options.Width, options.Height)

	viewMatrix := transformer.ViewMatrix(options.Cam.Z)
	projectionMatrix := transformer.PerspectiveMatrix(options.Cam.Fov, options.Cam.Aspect, options.Cam.ZNear, options.Cam.ZFar)

	calcIntensityQueue := make(chan *object.Face, 100)
	camerizeQueue := make(chan *object.Face, 100)
	projectQueue := make(chan *object.Face, 100)
	// clipQueue := make(chan *object.Face, 100)
	screenQueue := make(chan *object.Face, 100)
	rasterizeQueue := make(chan *Polygon, 100)

	var wg sync.WaitGroup

	wg.Add(1)
	go enface(models, calcIntensityQueue, &wg)

	wg.Add(1)
	go calcIntensity(calcIntensityQueue, camerizeQueue, &wg, getLights(options.LightSources))

	wg.Add(1)
	go camerize(camerizeQueue, projectQueue, &wg, viewMatrix)

	wg.Add(1)
	go project(projectQueue, screenQueue, &wg, projectionMatrix)

	// wg.Add(1)
	// go clip(clipQueue, screenQueue, &wg, options.Cam.ZNear, options.Cam.ZFar, float64(options.Width), float64(options.Height))

	wg.Add(1)
	go screen(screenQueue, rasterizeQueue, &wg, options.Width, options.Height)

	wg.Add(1)
	go rasterize(rasterizeQueue, &wg, image, zb)

	wg.Wait()
	return image
}
