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

func ApplyTransformation(model *reader.Model, transformation Mat4, withNormals bool) {
	for _, face := range model.Faces {
		for i, v := range face.Vertices {
			face.Vertices[i] = transformation.MultiplyVec3(v)
		}
		if withNormals {
			for i, normal := range face.Normals {
				face.Normals[i] = transformation.TransformNormal(normal)
			}
		}
	}
}

func TranslateMatrix(tx, ty, tz float64) Mat4 {
	m := IdentityMatrix()
	m[0][3] = tx
	m[1][3] = ty
	m[2][3] = tz
	return m
}

func ScaleMatrix(sx, sy, sz float64) Mat4 {
	m := IdentityMatrix()
	m[0][0] = sx
	m[1][1] = sy
	m[2][2] = sz
	return m
}

func RotateMatrix(angle float64, axis int) Mat4 {
	rad := angle * math.Pi / 180
	sin, cos := math.Sin(rad), math.Cos(rad)

	m := IdentityMatrix()
	switch axis {
	case XAxis:
		m[1][1], m[1][2] = cos, -sin
		m[2][1], m[2][2] = sin, cos
	case YAxis:
		m[0][0], m[0][2] = cos, sin
		m[2][0], m[2][2] = -sin, cos
	case ZAxis:
		m[0][0], m[0][1] = cos, -sin
		m[1][0], m[1][1] = sin, cos
	default:
		panic("Invalid axis specified. Use XAxis, YAxis, or ZAxis")
	}
	return m
}

func Translate(model *reader.Model, tx, ty, tz float64) {
	transformation := TranslateMatrix(tx, ty, tz)
	ApplyTransformation(model, transformation, false)
	model.Center = transformation.MultiplyVec3(model.Center)
}

func Scale(model *reader.Model, sx, sy, sz float64) {
	translateToOrigin := TranslateMatrix(-model.Center.X, -model.Center.Y, -model.Center.Z)
	scaleMatrix := ScaleMatrix(sx, sy, sz)
	translateBack := TranslateMatrix(model.Center.X, model.Center.Y, model.Center.Z)

	transformation := multiplyMatrices(translateBack, multiplyMatrices(scaleMatrix, translateToOrigin))
	ApplyTransformation(model, transformation, true)
}

func Rotate(model *reader.Model, angle float64, axis int) {
	translateToOrigin := TranslateMatrix(-model.Center.X, -model.Center.Y, -model.Center.Z)
	rotateMatrix := RotateMatrix(angle, axis)
	translateBack := TranslateMatrix(model.Center.X, model.Center.Y, model.Center.Z)

	transformation := multiplyMatrices(translateBack, multiplyMatrices(rotateMatrix, translateToOrigin))
	ApplyTransformation(model, transformation, true)
}

func ViewMatrix(cameraZ float64) Mat4 {
	return Mat4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, -cameraZ},
		{0, 0, 0, 1},
	}
}

func PerspectiveMatrix(fov, aspect, near, far float64) Mat4 {
	return Mat4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
	// fovRad := fov * math.Pi / 180.0
	// tanHalfFov := math.Tan(fovRad / 2.0)

	// m := Mat4{}
	// m[0][0] = 1.0 / (tanHalfFov * aspect)
	// m[1][1] = 1.0 / tanHalfFov
	// m[2][2] = -(far + near) / (far - near)
	// m[2][3] = -(2.0 * far * near) / (far - near)
	// m[3][2] = -1.0
	// m[3][3] = 0.0
	// return m
}

func TransformModelToCamera(model *reader.Model, viewMatrix, projectionMatrix Mat4) *reader.Model {
	transformedModel := reader.Model{Skeletonize: model.Skeletonize}

	for _, face := range model.Faces {
		transformedFace := reader.Face{}
		transformedFace.Vertices = append(transformedFace.Vertices, face.Vertices...)
		transformedFace.Normals = append(transformedFace.Normals, face.Normals...)
		transformedModel.Faces = append(transformedModel.Faces, transformedFace)
	}

	ApplyTransformation(&transformedModel, viewMatrix, true)
	ApplyTransformation(&transformedModel, projectionMatrix, false)

	return &transformedModel
}
