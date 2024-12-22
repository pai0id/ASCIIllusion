package transformer_test

import (
	"testing"

	"github.com/pai0id/CgCourseProject/internal/object"
	"github.com/pai0id/CgCourseProject/internal/transformer"
)

func TestIdentityMatrix(t *testing.T) {
	identity := transformer.IdentityMatrix()
	expected := transformer.Mat4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}

	if identity != expected {
		t.Errorf("Expected %v, got %v", expected, identity)
	}
}

func TestMultiplyVec3(t *testing.T) {
	matrix := transformer.IdentityMatrix()
	vec := object.Vec3{X: 1, Y: 2, Z: 3}

	result := matrix.MultiplyVec3(vec)

	if result != vec {
		t.Errorf("Expected %v, got %v", vec, result)
	}
}

func TestMultiplyNormal(t *testing.T) {
	matrix := transformer.IdentityMatrix()
	normal := object.Vec3{X: 1, Y: 0, Z: 0}
	result := matrix.MultiplyNormal(normal)

	if result != normal {
		t.Errorf("Expected %v, got %v", normal, result)
	}
}

func TestMultiplyMatrices(t *testing.T) {
	a := transformer.IdentityMatrix()
	b := transformer.IdentityMatrix()

	result := transformer.MultiplyMatrices(a, b)
	expected := transformer.IdentityMatrix()

	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestTransformNormal(t *testing.T) {
	matrix := transformer.IdentityMatrix()
	normal := object.Vec3{X: 1, Y: 0, Z: 0}
	result := matrix.TransformNormal(normal)

	expected := normal.Normalize()

	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
