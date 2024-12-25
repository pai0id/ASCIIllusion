package renderer

import (
	"github.com/pai0id/CgCourseProject/internal/object"
	"github.com/pai0id/CgCourseProject/internal/transformer"
)

func camerize(in <-chan *object.Face, out chan<- *object.Face, viewMatrix transformer.Mat4) {
	for f := range in {
		for i, v := range f.Vertices {
			f.Vertices[i] = viewMatrix.MultiplyVec3(v)
		}
		out <- f
	}
}
