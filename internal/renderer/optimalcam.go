package renderer

import (
	"math"

	"github.com/pai0id/CgCourseProject/internal/reader"
)

func CalculateBoundingBox(m *reader.Model) (reader.Vec3, reader.Vec3) {
	min := reader.Vec3{X: math.MaxFloat64, Y: math.MaxFloat64, Z: math.MaxFloat64}
	max := reader.Vec3{X: -math.MaxFloat64, Y: -math.MaxFloat64, Z: -math.MaxFloat64}

	for _, face := range m.Faces {
		for _, vertex := range face.Vertices {
			if vertex.X < min.X {
				min.X = vertex.X
			}
			if vertex.Y < min.Y {
				min.Y = vertex.Y
			}
			if vertex.Z < min.Z {
				min.Z = vertex.Z
			}
			if vertex.X > max.X {
				max.X = vertex.X
			}
			if vertex.Y > max.Y {
				max.Y = vertex.Y
			}
			if vertex.Z > max.Z {
				max.Z = vertex.Z
			}
		}
	}

	return min, max
}

func OptimalCameraDist(m *reader.Model, fov, aspect float64) float64 {
	min, max := CalculateBoundingBox(m)
	radFov := toRad(fov)
	tanHalf := math.Tan(radFov / 2)

	width := max.X - min.X
	height := max.Y - min.Y

	distanceVertical := height / (2 * tanHalf)
	horizontalFOV := 2 * math.Atan(tanHalf*aspect)
	tanHalfHor := math.Tan(horizontalFOV / 2)
	distanceHorizontal := width / (2 * tanHalfHor)

	return math.Max(0, math.Max(distanceVertical, distanceHorizontal))
}

func toRad(angle float64) float64 {
	return angle * math.Pi / 180.0
}
