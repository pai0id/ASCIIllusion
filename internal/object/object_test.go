package object_test

import (
	"testing"

	"github.com/pai0id/CgCourseProject/internal/object"
)

func TestDeepCopy(t *testing.T) {
	face := object.Face{
		Vertices:    []object.Vec3{{X: 1, Y: 2, Z: 3}},
		Normals:     []object.Vec3{{X: 0, Y: 1, Z: 0}},
		Intensities: []float64{0.5},
		Skeletonize: false,
	}

	copy := face.DeepCopy()

	// Ensure deep copy is successful and independent
	if &face.Vertices[0] == &copy.Vertices[0] {
		t.Fatalf("Vertices were not deep copied")
	}
	if &face.Normals[0] == &copy.Normals[0] {
		t.Fatalf("Normals were not deep copied")
	}
	if &face.Intensities[0] == &copy.Intensities[0] {
		t.Fatalf("Intensities were not deep copied")
	}

	// Ensure data is preserved
	if face.Skeletonize != copy.Skeletonize {
		t.Fatalf("Skeletonize value was not copied")
	}
}

func TestSkeletonize(t *testing.T) {
	obj := object.Object{
		Faces: []object.Face{
			{Skeletonize: false},
			{Skeletonize: false},
		},
	}

	obj.Skeletonize()

	for _, face := range obj.Faces {
		if !face.Skeletonize {
			t.Fatalf("Face was not skeletonized")
		}
	}
}

func TestCalculateBoundingBox(t *testing.T) {
	obj := object.Object{
		Faces: []object.Face{
			{Vertices: []object.Vec3{
				{X: 1, Y: 2, Z: 3},
				{X: -1, Y: -2, Z: -3},
			}},
		},
	}

	min, max := obj.CalculateBoundingBox()

	expectedMin := object.Vec3{X: -1, Y: -2, Z: -3}
	expectedMax := object.Vec3{X: 1, Y: 2, Z: 3}

	if min != expectedMin || max != expectedMax {
		t.Fatalf("Bounding box incorrect. Got min: %v, max: %v; Expected min: %v, max: %v", min, max, expectedMin, expectedMax)
	}
}

func TestCalculateCenter(t *testing.T) {
	obj := object.Object{
		Faces: []object.Face{
			{Vertices: []object.Vec3{
				{X: 1, Y: 1, Z: 1},
				{X: 3, Y: 3, Z: 3},
			}},
		},
	}

	obj.CalculateCenter()

	expectedCenter := object.Vec3{X: 2, Y: 2, Z: 2}
	if obj.Center != expectedCenter {
		t.Fatalf("Center calculation incorrect. Got: %v, Expected: %v", obj.Center, expectedCenter)
	}
}

func TestCalculateCenterEmpty(t *testing.T) {
	obj := object.Object{}

	obj.CalculateCenter()

	expectedCenter := object.Vec3{X: 0, Y: 0, Z: 0}
	if obj.Center != expectedCenter {
		t.Fatalf("Center calculation incorrect for empty object. Got: %v, Expected: %v", obj.Center, expectedCenter)
	}
}
