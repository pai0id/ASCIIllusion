package object

import (
	"math"
)

type Vec3 struct {
	X, Y, Z float64
}

func (a Vec3) Subtract(b Vec3) Vec3 {
	return Vec3{X: a.X - b.X, Y: a.Y - b.Y, Z: a.Z - b.Z}
}

func (a Vec3) Dot(b Vec3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func (a Vec3) Cross(b Vec3) Vec3 {
	return Vec3{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func (v Vec3) Normalize() Vec3 {
	length := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	return Vec3{X: v.X / length, Y: v.Y / length, Z: v.Z / length}
}

func InterpolateVec3(v1, v2 Vec3, t float64) Vec3 {
	return Vec3{
		X: v1.X + t*(v2.X-v1.X),
		Y: v1.Y + t*(v2.Y-v1.Y),
		Z: v1.Z + t*(v2.Z-v1.Z),
	}
}

func (v Vec3) Add(other Vec3) Vec3 {
	return Vec3{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

func (v Vec3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func MultiplyScalar(v Vec3, scalar float64) Vec3 {
	return Vec3{
		X: v.X * scalar,
		Y: v.Y * scalar,
		Z: v.Z * scalar,
	}
}

func SumNormals(normals []Vec3) Vec3 {
	sum := Vec3{0, 0, 0}
	for _, n := range normals {
		sum = sum.Add(n)
	}
	return sum.Normalize()
}
