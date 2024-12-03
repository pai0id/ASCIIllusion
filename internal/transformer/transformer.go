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
	model.Center.X += tx
	model.Center.Y += ty
	model.Center.Z += tz
}

func Scale(model *reader.Model, sx, sy, sz float64) {
	cx, cy, cz := model.Center.X, model.Center.Y, model.Center.Z
	for i, v := range model.Vertices {
		model.Vertices[i] = reader.Vertex{
			X: cx + (v.X-cx)*sx,
			Y: cy + (v.Y-cy)*sy,
			Z: cz + (v.Z-cz)*sz,
		}
	}
}

func Rotate(model *reader.Model, angle float64, axis int) {
	cx, cy, cz := model.Center.X, model.Center.Y, model.Center.Z
	rad := angle * math.Pi / 180
	sin, cos := math.Sin(rad), math.Cos(rad)

	for i, v := range model.Vertices {
		x, y, z := v.X-cx, v.Y-cy, v.Z-cz

		switch axis {
		case XAxis:
			model.Vertices[i] = reader.Vertex{
				X: cx + x,
				Y: cy + y*cos - z*sin,
				Z: cz + y*sin + z*cos,
			}
		case YAxis:
			model.Vertices[i] = reader.Vertex{
				X: cx + z*sin + x*cos,
				Y: cy + y,
				Z: cz + z*cos - x*sin,
			}
		case ZAxis:
			model.Vertices[i] = reader.Vertex{
				X: cx + x*cos - y*sin,
				Y: cy + x*sin + y*cos,
				Z: cz + z,
			}
		default:
			panic("Invalid axis specified. Use XAxis, YAxis, or ZAxis")
		}
	}
}
