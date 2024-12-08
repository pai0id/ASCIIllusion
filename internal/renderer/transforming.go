package renderer

import (
	"sync"

	"github.com/pai0id/CgCourseProject/internal/reader"
	"github.com/pai0id/CgCourseProject/internal/transformer"
)

func transforming(in <-chan *reader.Model, out chan<- *reader.Model, wg *sync.WaitGroup, viewMatrix transformer.Mat4) {
	defer wg.Done()
	for m := range in {
		out <- transformer.TransformModelToCamera(m, viewMatrix)
	}
	if out != nil {
		close(out)
	}
}

func projecting(in <-chan *reader.Model, out chan<- *reader.Model, wg *sync.WaitGroup, projectionMatrix transformer.Mat4) {
	defer wg.Done()
	for m := range in {
		out <- transformer.ProjectModel(m, projectionMatrix)
	}
	if out != nil {
		close(out)
	}
}
