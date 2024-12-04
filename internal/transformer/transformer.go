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

	for _, face := range model.Faces {
		for i, v := range face.Vertices {
			face.Vertices[i] = reader.Vec3{
				X: v.X + tx,
				Y: v.Y + ty,
				Z: v.Z + tz,
			}
		}
	}

	model.Center.X += tx
	model.Center.Y += ty
	model.Center.Z += tz
}

func Scale(model *reader.Model, sx, sy, sz float64) {
	cx, cy, cz := model.Center.X, model.Center.Y, model.Center.Z

	for _, face := range model.Faces {
		for i, v := range face.Vertices {
			face.Vertices[i] = reader.Vec3{
				X: cx + (v.X-cx)*sx,
				Y: cy + (v.Y-cy)*sy,
				Z: cz + (v.Z-cz)*sz,
			}
		}
	}

}

func Rotate(model *reader.Model, angle float64, axis int) {
	cx, cy, cz := model.Center.X, model.Center.Y, model.Center.Z
	rad := angle * math.Pi / 180
	sin, cos := math.Sin(rad), math.Cos(rad)

	for _, face := range model.Faces {
		for i, v := range face.Vertices {
			x, y, z := v.X-cx, v.Y-cy, v.Z-cz

			switch axis {
			case XAxis:
				face.Vertices[i] = reader.Vec3{
					X: cx + x,
					Y: cy + y*cos - z*sin,
					Z: cz + y*sin + z*cos,
				}
			case YAxis:
				face.Vertices[i] = reader.Vec3{
					X: cx + z*sin + x*cos,
					Y: cy + y,
					Z: cz + z*cos - x*sin,
				}
			case ZAxis:
				face.Vertices[i] = reader.Vec3{
					X: cx + x*cos - y*sin,
					Y: cy + x*sin + y*cos,
					Z: cz + z,
				}
			default:
				panic("Invalid axis specified. Use XAxis, YAxis, or ZAxis")
			}
		}
	}

	for _, face := range model.Faces {
		for i := range face.Normals {
			switch axis {
			case XAxis:
				face.Normals[i] = reader.Vec3{
					X: cx,
					Y: cy*cos - cz*sin,
					Z: cy*sin + cz*cos,
				}
			case YAxis:
				face.Normals[i] = reader.Vec3{
					X: cz*sin + cx*cos,
					Y: cy,
					Z: cz*cos - cx*sin,
				}
			case ZAxis:
				face.Normals[i] = reader.Vec3{
					X: cx*cos - cy*sin,
					Y: cx*sin + cy*cos,
					Z: cz,
				}
			default:
				panic("Invalid axis specified. Use XAxis, YAxis, or ZAxis")
			}
		}
	}
}

func Project(model *reader.Model, scale, cameraDist float64) *reader.Model {
	projectedVertices := make([]reader.Face, len(model.Faces))
	for i, face := range model.Faces {
		projectedVertices[i] = reader.Face{
			Vertices: make([]reader.Vec3, len(face.Vertices)),
			Normals:  face.Normals,
		}
		for j, vertex := range face.Vertices {
			projectedVertices[i].Vertices[j] = vertex
			// projectedVertices[i].Vertices[j] = perspectiveProject(vertex, scale, cameraDist)
		}
	}
	return &reader.Model{Faces: projectedVertices, Center: model.Center}
}

func perspectiveProject(vertex reader.Vec3, scale, cameraDist float64) reader.Vec3 {
	z := vertex.Z + cameraDist
	if z == 0 {
		z = 0.0001
	}
	return reader.Vec3{X: vertex.X * scale / z, Y: vertex.Y * scale / z, Z: vertex.Z}
}
