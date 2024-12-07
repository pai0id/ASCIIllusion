package renderer

import (
	"math"

	"github.com/pai0id/CgCourseProject/internal/reader"
)

func OptimalCameraDist(m *reader.Model, fov, aspect float64) float64 {
	bboxMin, bboxMax := reader.Vec3{}, reader.Vec3{}

	for _, f := range m.Faces {
		for _, v := range f.Vertices {
			bboxMin.X = math.Min(bboxMin.X, v.X)
			bboxMin.Y = math.Min(bboxMin.Y, v.Y)
			bboxMin.Z = math.Min(bboxMin.Z, v.Z)

			bboxMax.X = math.Max(bboxMax.X, v.X)
			bboxMax.Y = math.Max(bboxMax.Y, v.Y)
			bboxMax.Z = math.Max(bboxMax.Z, v.Z)
		}
	}

	return calculateCameraZ(bboxMin, bboxMax, fov, aspect)
}

func calculateCameraZ(bboxMin, bboxMax reader.Vec3, fov, aspect float64) float64 {
	cx := (bboxMin.X + bboxMax.X) / 2
	cy := (bboxMin.Y + bboxMax.Y) / 2
	cz := (bboxMin.Z + bboxMax.Z) / 2

	radius := 0.0
	vertices := [][3]float64{
		{bboxMin.X, bboxMin.Y, bboxMin.Z},
		{bboxMin.X, bboxMin.Y, bboxMax.Z},
		{bboxMin.X, bboxMax.Y, bboxMin.Z},
		{bboxMin.X, bboxMax.Y, bboxMax.Z},
		{bboxMax.X, bboxMin.Y, bboxMin.Z},
		{bboxMax.X, bboxMin.Y, bboxMax.Z},
		{bboxMax.X, bboxMax.Y, bboxMin.Z},
		{bboxMax.X, bboxMax.Y, bboxMax.Z},
	}

	for _, v := range vertices {
		dx, dy, dz := v[0]-cx, v[1]-cy, v[2]-cz
		dist := math.Sqrt(dx*dx + dy*dy + dz*dz)
		if dist > radius {
			radius = dist
		}
	}

	fovHRad := fov * math.Pi / 180

	zCameraH := radius / math.Tan(fovHRad/2)
	fovVRad := 2 * math.Atan(math.Tan(fovHRad/2)/aspect)
	zCameraV := radius / math.Tan(fovVRad/2)

	return math.Max(zCameraH, zCameraV)
}
