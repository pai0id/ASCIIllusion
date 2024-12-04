package renderer

import "math"

type vec3 struct {
	x, y, z float64
}

func dot(v1, v2 vec3) float64 {
	return v1.x*v2.x + v1.y*v2.y + v1.z*v2.z
}

func magnitude(v vec3) float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

func normalize(v vec3) vec3 {
	mag := magnitude(v)
	return vec3{v.x / mag, v.y / mag, v.z / mag}
}

func subtract(v1, v2 vec3) vec3 {
	return vec3{v1.x - v2.x, v1.y - v2.y, v1.z - v2.z}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Linear interpolation between two values
func lerp(start, end, t float64) float64 {
	return start + t*(end-start)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
