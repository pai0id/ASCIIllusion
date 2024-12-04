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

// Translate shifts the model in 3D space
func Translate(model *reader.Model, tx, ty, tz float64) {
	// Translate vertices
	for _, face := range model.Faces {
		for i, v := range face.Vertices {
			face.Vertices[i] = reader.Vertex{
				X: v.X + tx,
				Y: v.Y + ty,
				Z: v.Z + tz,
			}
		}
	}
	// Translate center
	model.Center.X += tx
	model.Center.Y += ty
	model.Center.Z += tz
}

// Scale applies scaling to the model relative to its center
func Scale(model *reader.Model, sx, sy, sz float64) {
	cx, cy, cz := model.Center.X, model.Center.Y, model.Center.Z
	// Scale vertices
	for _, face := range model.Faces {
		for i, v := range face.Vertices {
			face.Vertices[i] = reader.Vertex{
				X: cx + (v.X-cx)*sx,
				Y: cy + (v.Y-cy)*sy,
				Z: cz + (v.Z-cz)*sz,
			}
		}
	}

	// Do not scale normals (they should remain normalized)
}

// Rotate rotates the model around a specified axis
func Rotate(model *reader.Model, angle float64, axis int) {
	cx, cy, cz := model.Center.X, model.Center.Y, model.Center.Z
	rad := angle * math.Pi / 180
	sin, cos := math.Sin(rad), math.Cos(rad)

	// Rotate vertices
	for _, face := range model.Faces {
		for i, v := range face.Vertices {
			x, y, z := v.X-cx, v.Y-cy, v.Z-cz

			switch axis {
			case XAxis:
				face.Vertices[i] = reader.Vertex{
					X: cx + x,
					Y: cy + y*cos - z*sin,
					Z: cz + y*sin + z*cos,
				}
			case YAxis:
				face.Vertices[i] = reader.Vertex{
					X: cx + z*sin + x*cos,
					Y: cy + y,
					Z: cz + z*cos - x*sin,
				}
			case ZAxis:
				face.Vertices[i] = reader.Vertex{
					X: cx + x*cos - y*sin,
					Y: cy + x*sin + y*cos,
					Z: cz + z,
				}
			default:
				panic("Invalid axis specified. Use XAxis, YAxis, or ZAxis")
			}
		}
	}

	// Rotate normals (optional)
	for _, face := range model.Faces {
		switch axis {
		case XAxis:
			face.Normal = reader.Vertex{
				X: cx,
				Y: cy*cos - cz*sin,
				Z: cy*sin + cz*cos,
			}
		case YAxis:
			face.Normal = reader.Vertex{
				X: cz*sin + cx*cos,
				Y: cy,
				Z: cz*cos - cx*sin,
			}
		case ZAxis:
			face.Normal = reader.Vertex{
				X: cx*cos - cy*sin,
				Y: cx*sin + cy*cos,
				Z: cz,
			}
		default:
			panic("Invalid axis specified. Use XAxis, YAxis, or ZAxis")
		}
	}
}

func Project(model *reader.Model, scale, cameraDist float64) *reader.Model {
	projectedVertices := make([]reader.Face, len(model.Faces))
	for i, face := range model.Faces {
		projectedVertices[i] = reader.Face{
			Vertices: make([]reader.Vertex, len(face.Vertices)),
			Normal:   face.Normal,
		}
		for j, vertex := range face.Vertices {
			projectedVertices[i].Vertices[j] = perspectiveProject(vertex, scale, cameraDist)
		}
	}
	return &reader.Model{Faces: projectedVertices, Center: model.Center}
}

func perspectiveProject(vertex reader.Vertex, scale, cameraDist float64) reader.Vertex {
	z := vertex.Z + cameraDist
	if z == 0 {
		z = 0.0001
	}
	return reader.Vertex{X: vertex.X * scale / z, Y: vertex.Y * scale / z, Z: vertex.Z}
}
