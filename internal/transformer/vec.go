package transformer

import (
	"math"

	"github.com/pai0id/CgCourseProject/internal/reader"
)

func Subtract(a, b reader.Vec3) reader.Vec3 {
	return reader.Vec3{X: a.X - b.X, Y: a.Y - b.Y, Z: a.Z - b.Z}
}

func Dot(a, b reader.Vec3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func Cross(a, b reader.Vec3) reader.Vec3 {
	return reader.Vec3{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func Normalize(v reader.Vec3) reader.Vec3 {
	length := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	return reader.Vec3{X: v.X / length, Y: v.Y / length, Z: v.Z / length}
}

func InterpolateVec3(v1, v2 reader.Vec3, t float64) reader.Vec3 {
	return reader.Vec3{
		X: v1.X + t*(v2.X-v1.X),
		Y: v1.Y + t*(v2.Y-v1.Y),
		Z: v1.Z + t*(v2.Z-v1.Z),
	}
}
