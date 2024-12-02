package transformer

import (
	"math"

	"github.com/pai0id/CgCourseProject/internal/reader"
)

const (
	XAxis = 1
	YAxis = 2
	ZAxis = 3
)

func Translate(model *reader.Model, tx, ty, tz float64) {
	for i, v := range model.Vertices {
		model.Vertices[i] = reader.Vertex{
			X: v.X + tx,
			Y: v.Y + ty,
			Z: v.Z + tz,
		}
	}
}

func Scale(model *reader.Model, sx, sy, sz float64) {
	for i, v := range model.Vertices {
		model.Vertices[i] = reader.Vertex{
			X: v.X * sx,
			Y: v.Y * sy,
			Z: v.Z * sz,
		}
	}
}

func Rotate(model *reader.Model, angle float64, axis int) {
	rad := angle * math.Pi / 180
	sin, cos := math.Sin(rad), math.Cos(rad)

	for i, v := range model.Vertices {
		switch axis {
		case XAxis:
			model.Vertices[i] = reader.Vertex{
				X: v.X,
				Y: v.Y*cos - v.Z*sin,
				Z: v.Y*sin + v.Z*cos,
			}
		case YAxis:
			model.Vertices[i] = reader.Vertex{
				X: v.Z*sin + v.X*cos,
				Y: v.Y,
				Z: v.Z*cos - v.X*sin,
			}
		case ZAxis:
			model.Vertices[i] = reader.Vertex{
				X: v.X*cos - v.Y*sin,
				Y: v.X*sin + v.Y*cos,
				Z: v.Z,
			}
		default:
			panic("Invalid axis specified. Use XAxis, YAxis, or ZAxis")
		}
	}
}
