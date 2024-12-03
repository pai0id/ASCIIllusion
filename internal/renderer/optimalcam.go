package renderer

import (
	"math"

	"github.com/pai0id/CgCourseProject/internal/reader"
)

func OptimalCameraDist(models []*reader.Model, options *RenderOptions) float64 {
	maxDist := 0.0
	for _, model := range models {
		var minX, minY, minZ = math.MaxFloat64, math.MaxFloat64, math.MaxFloat64
		var maxX, maxY, maxZ = -math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64

		for _, face := range model.Faces {
			for _, vertex := range face.Vertices {
				if vertex.X < minX {
					minX = vertex.X
				}
				if vertex.X > maxX {
					maxX = vertex.X
				}
				if vertex.Y < minY {
					minY = vertex.Y
				}
				if vertex.Y > maxY {
					maxY = vertex.Y
				}
				if vertex.Z < minZ {
					minZ = vertex.Z
				}
				if vertex.Z > maxZ {
					maxZ = vertex.Z
				}
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
