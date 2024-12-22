package renderer

import (
	"sync"

	"github.com/pai0id/CgCourseProject/internal/object"
	"github.com/pai0id/CgCourseProject/internal/transformer"
)

func project(in <-chan *object.Face, out chan<- *object.Face, wg *sync.WaitGroup, projectionMatrix transformer.Mat4) {
	defer wg.Done()
	for f := range in {
		for i, v := range f.Vertices {
			f.Vertices[i] = projectionMatrix.MultiplyVec3(v)
		}
		out <- f
	}
	if out != nil {
		close(out)
	}
}
