package renderer_test

import (
	"testing"

	"github.com/pai0id/CgCourseProject/internal/object"
	"github.com/pai0id/CgCourseProject/internal/renderer"
)

func TestBarycentric(t *testing.T) {
	a := renderer.Point{X: 0, Y: 0}
	b := renderer.Point{X: 5, Y: 0}
	c := renderer.Point{X: 0, Y: 5}
	p := renderer.Point{X: 2, Y: 2}

	u, v, w := renderer.Barycentric(p, a, b, c)

	if u < 0 || v < 0 || w < 0 {
		t.Errorf("Expected positive barycentric coordinates, got u=%f, v=%f, w=%f", u, v, w)
	}
}

func TestBoundingBox(t *testing.T) {
	vertices := []renderer.Point{
		{X: 1, Y: 1},
		{X: 3, Y: 4},
		{X: 2, Y: 2},
	}

	xMin, xMax, yMin, yMax := renderer.BoundingBox(vertices)

	if xMin != 1 || xMax != 3 || yMin != 1 || yMax != 4 {
		t.Errorf("Bounding box incorrect, got xMin=%d, xMax=%d, yMin=%d, yMax=%d", xMin, xMax, yMin, yMax)
	}
}

func TestCalculateVertexIntensity(t *testing.T) {
	point := object.Vec3{X: 0, Y: 0, Z: 0}
	normal := object.Vec3{X: 0, Y: 0, Z: 1}
	lightSources := []renderer.Light{
		{
			Position:  object.Vec3{X: 1, Y: 1, Z: 1},
			Intensity: 1.0,
		},
	}

	intensity := renderer.CalculateVertexIntensity(point, normal, lightSources)

	if intensity <= 0 || intensity > 1 {
		t.Errorf("Expected intensity between 0 and 1, got %f", intensity)
	}
}

func TestToScreen(t *testing.T) {
	vertex := object.Vec3{X: 1.0, Y: 1.0, Z: 1.0}
	width, height := 800, 600
	result := renderer.ToScreen(vertex, width, height)

	expectedX := -1 + width/2
	expectedY := -1 + height/2

	if result.X != expectedX || result.Y != expectedY {
		t.Errorf("Expected (%d, %d), got (%d, %d)", expectedX, expectedY, result.X, result.Y)
	}
}
