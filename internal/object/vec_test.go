package object_test

import (
	"math"
	"testing"

	"github.com/pai0id/CgCourseProject/internal/object"
)

func TestSubtract(t *testing.T) {
	a := object.Vec3{X: 3, Y: 4, Z: 5}
	b := object.Vec3{X: 1, Y: 2, Z: 3}
	expected := object.Vec3{X: 2, Y: 2, Z: 2}

	result := a.Subtract(b)
	if result != expected {
		t.Fatalf("Subtract incorrect. Got: %v, Expected: %v", result, expected)
	}
}

func TestDot(t *testing.T) {
	a := object.Vec3{X: 1, Y: 2, Z: 3}
	b := object.Vec3{X: 4, Y: 5, Z: 6}
	expected := 32.0

	result := a.Dot(b)
	if result != expected {
		t.Fatalf("Dot product incorrect. Got: %v, Expected: %v", result, expected)
	}
}

func TestCross(t *testing.T) {
	a := object.Vec3{X: 1, Y: 2, Z: 3}
	b := object.Vec3{X: 4, Y: 5, Z: 6}
	expected := object.Vec3{X: -3, Y: 6, Z: -3}

	result := a.Cross(b)
	if result != expected {
		t.Fatalf("Cross product incorrect. Got: %v, Expected: %v", result, expected)
	}
}

func TestNormalize(t *testing.T) {
	v := object.Vec3{X: 3, Y: 4, Z: 0}
	expected := object.Vec3{X: 0.6, Y: 0.8, Z: 0}

	result := v.Normalize()
	if math.Abs(result.X-expected.X) > 1e-6 || math.Abs(result.Y-expected.Y) > 1e-6 || math.Abs(result.Z-expected.Z) > 1e-6 {
		t.Fatalf("Normalize incorrect. Got: %v, Expected: %v", result, expected)
	}
}

func TestInterpolateVec3(t *testing.T) {
	v1 := object.Vec3{X: 1, Y: 2, Z: 3}
	v2 := object.Vec3{X: 4, Y: 5, Z: 6}
	tVal := 0.5
	expected := object.Vec3{X: 2.5, Y: 3.5, Z: 4.5}

	result := object.InterpolateVec3(v1, v2, tVal)
	if result != expected {
		t.Fatalf("InterpolateVec3 incorrect. Got: %v, Expected: %v", result, expected)
	}
}

func TestLength(t *testing.T) {
	v := object.Vec3{X: 3, Y: 4, Z: 12}
	expected := 13.0

	result := v.Length()
	if math.Abs(result-expected) > 1e-6 {
		t.Fatalf("Length incorrect. Got: %v, Expected: %v", result, expected)
	}
}

func TestMultiplyScalar(t *testing.T) {
	v := object.Vec3{X: 1, Y: 2, Z: 3}
	scalar := 2.0
	expected := object.Vec3{X: 2, Y: 4, Z: 6}

	result := object.MultiplyScalar(v, scalar)
	if result != expected {
		t.Fatalf("MultiplyScalar incorrect. Got: %v, Expected: %v", result, expected)
	}
}

func TestSumNormals(t *testing.T) {
	normals := []object.Vec3{
		{X: 1, Y: 0, Z: 0},
		{X: 0, Y: 1, Z: 0},
	}
	expected := object.Vec3{X: math.Sqrt(2) / 2, Y: math.Sqrt(2) / 2, Z: 0}

	result := object.SumNormals(normals)
	if math.Abs(result.X-expected.X) > 1e-6 || math.Abs(result.Y-expected.Y) > 1e-6 || math.Abs(result.Z-expected.Z) > 1e-6 {
		t.Fatalf("SumNormals incorrect. Got: %v, Expected: %v", result, expected)
	}
}
