package renderer

import (
	"math"

	"github.com/pai0id/CgCourseProject/internal/asciiser"
	"github.com/pai0id/CgCourseProject/internal/reader"
)

type RenderOptions struct {
	Width      int
	Height     int
	Fov        float64
	CameraDist float64
}

func RenderModels(models []*reader.Model, options *RenderOptions) asciiser.Image {
	if options.Width == 0 || options.Height == 0 {
		return nil
	}
	image := asciiser.NewImage(options.Width, options.Height)

	for _, model := range models {
		renderModel(model, options, image)
	}

	return image
}

func renderModel(model *reader.Model, options *RenderOptions, image asciiser.Image) {
	centerX := float64(options.Width) / 2
	centerY := float64(options.Height) / 2
	scale := centerX / math.Tan(options.Fov*math.Pi/360)

	projectedVertices := make([][2]int, len(model.Vertices))
	for i, v := range model.Vertices {
		px, py := perspectiveProject(v, scale, options.CameraDist)
		projectedVertices[i] = [2]int{
			int(math.Round(centerX + px)),
			int(math.Round(centerY - py)),
		}
	}

	for _, face := range model.Faces {
		for i := 0; i < len(face.VertexIndices); i++ {
			v1 := projectedVertices[face.VertexIndices[i]]
			v2 := projectedVertices[face.VertexIndices[(i+1)%len(face.VertexIndices)]]
			rasterizeLine(image, v1[0], v1[1], v2[0], v2[1])
		}
	}
}

func perspectiveProject(vertex reader.Vertex, scale, cameraDist float64) (float64, float64) {
	z := vertex.Z + cameraDist
	if z == 0 {
		z = 0.0001
	}
	return vertex.X * scale / z, vertex.Y * scale / z
}

func rasterizeLine(image asciiser.Image, x0, y0, x1, y1 int) {
	dx := int(math.Abs(float64(x1 - x0)))
	dy := int(math.Abs(float64(y1 - y0)))
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx - dy

	for {
		if x0 >= 0 && x0 < len(image[0]) && y0 >= 0 && y0 < len(image) {
			image[y0][x0].IsLine = true
		}
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func OptimalCameraDist(models []*reader.Model, options *RenderOptions) float64 {
	maxDist := 0.0
	for _, model := range models {
		var minX, minY, minZ = math.MaxFloat64, math.MaxFloat64, math.MaxFloat64
		var maxX, maxY, maxZ = -math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64

		for _, v := range model.Vertices {
			if v.X < minX {
				minX = v.X
			}
			if v.X > maxX {
				maxX = v.X
			}
			if v.Y < minY {
				minY = v.Y
			}
			if v.Y > maxY {
				maxY = v.Y
			}
			if v.Z < minZ {
				minZ = v.Z
			}
			if v.Z > maxZ {
				maxZ = v.Z
			}
		}

		width := maxX - minX
		height := maxY - minY

		centerX := float64(options.Width) / 2

		scale := centerX / math.Tan(options.Fov*math.Pi/360)

		distX := width * scale / float64(options.Width)
		distY := height * scale / float64(options.Height)

		res := 1.2 * (math.Max(distX, distY) + maxZ)
		if res > maxDist {
			maxDist = res
		}
	}
	return maxDist
}
