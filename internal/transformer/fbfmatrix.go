package transformer

import "github.com/pai0id/CgCourseProject/internal/reader"

type Mat4 [4][4]float64

func IdentityMatrix() Mat4 {
	return Mat4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func (m Mat4) MultiplyVec3(v reader.Vec3) reader.Vec3 {
	x := m[0][0]*v.X + m[0][1]*v.Y + m[0][2]*v.Z + m[0][3]
	y := m[1][0]*v.X + m[1][1]*v.Y + m[1][2]*v.Z + m[1][3]
	z := m[2][0]*v.X + m[2][1]*v.Y + m[2][2]*v.Z + m[2][3]
	w := m[3][0]*v.X + m[3][1]*v.Y + m[3][2]*v.Z + m[3][3]

	if w != 0 {
		return reader.Vec3{X: x / w, Y: y / w, Z: z / w}
	}
	return reader.Vec3{X: x, Y: y, Z: z}
}

func (m Mat4) MultiplyNormal(v reader.Vec3) reader.Vec3 {
	x := m[0][0]*v.X + m[0][1]*v.Y + m[0][2]*v.Z
	y := m[1][0]*v.X + m[1][1]*v.Y + m[1][2]*v.Z
	z := m[2][0]*v.X + m[2][1]*v.Y + m[2][2]*v.Z
	return reader.Vec3{X: x, Y: y, Z: z}
}

func MultiplyMatrices(a, b Mat4) Mat4 {
	var result Mat4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return result
}

func (m Mat4) TransformNormal(normal reader.Vec3) reader.Vec3 {
	var invTranspose Mat4
	invTranspose[0][0] = m[0][0]
	invTranspose[0][1] = m[1][0]
	invTranspose[0][2] = m[2][0]

	invTranspose[1][0] = m[0][1]
	invTranspose[1][1] = m[1][1]
	invTranspose[1][2] = m[2][1]

	invTranspose[2][0] = m[0][2]
	invTranspose[2][1] = m[1][2]
	invTranspose[2][2] = m[2][2]

	x := invTranspose[0][0]*normal.X + invTranspose[0][1]*normal.Y + invTranspose[0][2]*normal.Z
	y := invTranspose[1][0]*normal.X + invTranspose[1][1]*normal.Y + invTranspose[1][2]*normal.Z
	z := invTranspose[2][0]*normal.X + invTranspose[2][1]*normal.Y + invTranspose[2][2]*normal.Z

	return Normalize(reader.Vec3{X: x, Y: y, Z: z})
}
