package renderer

import (
	"math"
	"runtime"
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
	screenQueue := make(chan *object.Face, 100)
	rasterizeQueue := make(chan *Polygon, 100)

	var wg sync.WaitGroup

	runtime.GOMAXPROCS(runtime.NumCPU())

	wg.Add(1)
	go func() {
		defer wg.Done()
		enface(models, calcIntensityQueue)
		close(calcIntensityQueue)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		calcIntensity(calcIntensityQueue, camerizeQueue, getLights(options.LightSources))
		close(camerizeQueue)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		camerize(camerizeQueue, projectQueue, viewMatrix)
		close(projectQueue)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		project(projectQueue, screenQueue, projectionMatrix)
		close(screenQueue)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		screen(screenQueue, rasterizeQueue, options.Width, options.Height)
		close(rasterizeQueue)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		rasterize(rasterizeQueue, image, zb)
	}()

	wg.Wait()

	return image
}
