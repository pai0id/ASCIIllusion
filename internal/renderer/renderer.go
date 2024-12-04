package renderer

import (
	"math"
	"sort"
	"sync"

	"github.com/pai0id/CgCourseProject/internal/asciiser"
	"github.com/pai0id/CgCourseProject/internal/reader"
	"github.com/pai0id/CgCourseProject/internal/transformer"
)

const (
	ambientC           = 0.1
	diffuseIntensityC  = 0.0
	specularIntensityC = 0.0
)

type point struct {
	x, y int
	z    float64
}

type normal struct {
	x, y, z float64
}

type face struct {
	vertices []point
	normal   normal
}

type object struct {
	faces []*face
}

type polygon struct {
	vertices    []point
	intensities []float64
}

type filledPolygon map[point]asciiser.Pixel

type RenderOptions struct {
	Width           int
	Height          int
	Fov             float64
	CameraDist      float64
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

	centerX := float64(options.Width) / 2
	centerY := float64(options.Height) / 2
	scale := centerX / math.Tan(options.Fov*math.Pi/360)

	projectQueue := make(chan *reader.Model, 10)
	vertIntensitiesQueue := make(chan *object, 10)
	fillQueue := make(chan *polygon, 100)
	resQueue := make(chan *filledPolygon, 100)

	var polygons []*filledPolygon

	var wg sync.WaitGroup

	wg.Add(1)
	go start(models, projectQueue, &wg)

	wg.Add(1)
	go project(projectQueue, vertIntensitiesQueue, &wg, scale, options.CameraDist, point{x: int(centerX), y: int(centerY), z: options.CameraDist})

	wg.Add(1)
	go vertIntensities(vertIntensitiesQueue, fillQueue, &wg, options.LightSources, options.CameraDist)

	wg.Add(1)
	go fill(fillQueue, resQueue, &wg)

	wg.Add(1)
	go end(resQueue, &polygons, &wg)

	wg.Wait()

	renderPolygons(polygons, image, zb)

	return image
}

func renderPolygons(fps []*filledPolygon, image asciiser.Image, zb zBuffer) {
	for _, fp := range fps {
		for pt, px := range *fp {
			if pt.x >= 0 && pt.x < len(image[0]) && pt.y >= 0 && pt.y < len(image) {

				if pt.z < zb[pt.y][pt.x] {
					zb[pt.y][pt.x] = pt.z
					image[pt.y][pt.x] = px
				}
			}
		}
	}
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

func project(in <-chan *reader.Model, out chan<- *object, wg *sync.WaitGroup, scale, cameraDist float64, screenCenter point) {
	defer wg.Done()
	for m := range in {

		if out != nil {
			newM := transformer.Project(m, scale, cameraDist)
			out <- mapToScreenCoords(newM, screenCenter)
		}
	}

	if out != nil {
		close(out)
	}
}

func mapToScreenCoords(m *reader.Model, screenCenter point) *object {
	res := &object{
		faces: make([]*face, len(m.Faces)),
	}
	for i, f := range m.Faces {
		res.faces[i] = &face{
			vertices: make([]point, len(f.Vertices)),
			normal:   normal{f.Normals[0].X, f.Normals[0].Y, f.Normals[0].Z},
		}
		for j, v := range f.Vertices {
			res.faces[i].vertices[j].x = int(math.Round(float64(screenCenter.x) + v.X))
			res.faces[i].vertices[j].y = int(math.Round(float64(screenCenter.y) - v.Y))
			res.faces[i].vertices[j].z = v.Z
		}
	}
	return res
}

func vertIntensities(in <-chan *object, out chan<- *polygon, wg *sync.WaitGroup, lightSources []reader.Vec3, cameraDist float64) {
	defer wg.Done()

	for obj := range in {
		its := calculateVertexIntensities(obj, lightSources, cameraDist)
		for _, f := range obj.faces {
			out <- mapToPolygons(f, its)
		}
	}

	if out != nil {
		close(out)
	}
}

func mapToPolygons(f *face, its map[point]float64) *polygon {
	res := &polygon{
		vertices:    make([]point, len(f.vertices)),
		intensities: make([]float64, len(f.vertices)),
	}

	for i, v := range f.vertices {
		res.vertices[i] = point{x: v.x, y: v.y, z: v.z}
		res.intensities[i] = its[v]
	}

	return res
}

func calculateVertexIntensities(obj *object, lightSources []reader.Vec3, cameraDist float64) map[point]float64 {
	vertexColors := make(map[point]float64)
	normals := calculateVertexNormals(obj)

	for _, face := range obj.faces {
		for _, vertex := range face.vertices {
			if _, exists := vertexColors[vertex]; !exists {
				vertexColors[vertex] = calculateLighting(vertex, normals[vertex], lightSources, cameraDist)
			}
		}
	}

	return vertexColors
}

func calculateLighting(vertex point, normal normal, lightSources []reader.Vec3, cameraDist float64) float64 {
	ambient := ambientC
	diffuseIntensity := diffuseIntensityC
	specularIntensity := specularIntensityC

	n := normalize(vec3(normal))
	v := vec3{float64(vertex.x), float64(vertex.y), vertex.z}

	for _, light := range lightSources {
		lightPos := vec3{light.X, light.Y, light.Z}
		lightDir := normalize(subtract(lightPos, v))

		diffuse := math.Max(0, dot(n, lightDir))
		diffuseIntensity += diffuse

		reflection := subtract(vec3{2 * dot(n, lightDir) * n.x, 2 * dot(n, lightDir) * n.y, 2 * dot(n, lightDir) * n.z}, lightDir)
		viewDir := normalize(subtract(vec3{0, 0, cameraDist}, v))
		specular := math.Pow(math.Max(0, dot(reflection, viewDir)), 32)
		specularIntensity += specular
	}

	return ambient + 0.8*diffuseIntensity + 0.5*specularIntensity
}

func fill(in <-chan *polygon, out chan<- *filledPolygon, wg *sync.WaitGroup) {
	defer wg.Done()

	for p := range in {
		fp := rasterizePolygon(p)
		rasterizeLines(p, fp)
		out <- fp

	}

	if out != nil {
		close(out)
	}
}

func rasterizePolygon(p *polygon) *filledPolygon {
	filled := make(filledPolygon)

	minY, maxY := math.MaxInt, math.MinInt
	for _, v := range p.vertices {
		minY = min(minY, v.y)
		maxY = max(maxY, v.y)
	}

	for y := minY; y <= maxY; y++ {
		var intersections []struct {
			x         int
			z         float64
			intensity float64
		}

		for i := 0; i < len(p.vertices); i++ {
			v1 := p.vertices[i]
			v2 := p.vertices[(i+1)%len(p.vertices)]

			if (v1.y <= y && v2.y > y) || (v2.y <= y && v1.y > y) {

				t := float64(y-v1.y) / float64(v2.y-v1.y)
				x := int(lerp(float64(v1.x), float64(v2.x), t))
				z := lerp(v1.z, v2.z, t)
				intensity := lerp(p.intensities[i], p.intensities[(i+1)%len(p.vertices)], t)

				intersections = append(intersections, struct {
					x         int
					z         float64
					intensity float64
				}{x: x, z: z, intensity: intensity})
			}
		}

		sort.Slice(intersections, func(i, j int) bool {
			return intersections[i].x < intersections[j].x
		})

		for i := 0; i < len(intersections)-1; i += 2 {
			xStart, xEnd := intersections[i].x, intersections[i+1].x
			zStart, zEnd := intersections[i].z, intersections[i+1].z
			intensityStart, intensityEnd := intersections[i].intensity, intersections[i+1].intensity

			for x := xStart; x <= xEnd; x++ {
				t := float64(x-xStart) / float64(xEnd-xStart)
				z := lerp(zStart, zEnd, t)
				intensity := lerp(intensityStart, intensityEnd, t)
				filled[point{x: x, y: y, z: z}] = asciiser.Pixel{Brightness: intensity, IsLine: false}
			}
		}
	}

	return &filled
}

func rasterizeLines(p *polygon, fp *filledPolygon) {

	for i := 0; i < len(p.vertices); i++ {
		v1 := p.vertices[i]
		v2 := p.vertices[(i+1)%len(p.vertices)]

		dx := abs(v2.x - v1.x)
		dy := abs(v2.y - v1.y)
		sx := 1
		if v1.x > v2.x {
			sx = -1
		}
		sy := 1
		if v1.y > v2.y {
			sy = -1
		}

		err := dx - dy

		x, y, z := v1.x, v1.y, v1.z
		for {

			(*fp)[point{x: x, y: y, z: z}] = asciiser.Pixel{Brightness: 1.0, IsLine: true}

			if x == v2.x && y == v2.y {
				break
			}

			e2 := 2 * err
			if e2 > -dy {
				err -= dy
				x += sx
			}
			if e2 < dx {
				err += dx
				y += sy
			}

			t := float64(abs(x-v1.x)) / float64(dx+1)
			z = lerp(v1.z, v2.z, t)
		}
	}
}

func end(in <-chan *filledPolygon, res *([]*filledPolygon), wg *sync.WaitGroup) {
	defer wg.Done()

	for p := range in {
		*res = append(*res, p)
	}
}
