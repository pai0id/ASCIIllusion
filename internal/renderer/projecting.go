package renderer

import (
	"sync"

	"github.com/pai0id/CgCourseProject/internal/reader"
	"github.com/pai0id/CgCourseProject/internal/transformer"
)

func project(in <-chan *reader.Model, out chan<- *reader.Model, wg *sync.WaitGroup, viewMatrix, projectionMatrix transformer.Mat4) {
	defer wg.Done()
	for m := range in {
		out <- transformer.TransformModelToCamera(m, viewMatrix, projectionMatrix)
	}
	if out != nil {
		close(out)
	}
}
