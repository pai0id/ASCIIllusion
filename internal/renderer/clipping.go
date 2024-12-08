package renderer

import (
	"sync"

	"github.com/pai0id/CgCourseProject/internal/reader"
)

// Функция обрезки по zNear, zFar, cameraWidth, cameraHeight
func clipping(in <-chan *reader.Model, out chan<- *reader.Model, wg *sync.WaitGroup, zNear, zFar float64) {
	defer wg.Done()
	for m := range in {
		// clippedModel := reader.Model{Skeletonize: m.Skeletonize}
		// for _, face := range m.Faces {
		// 	clippedFaces := ClipFace(face, zNear, zFar)
		// 	clippedModel.Faces = append(clippedModel.Faces, clippedFaces...)
		// }
		// out <- &clippedModel
		out <- m
	}
	if out != nil {
		close(out)
	}
}
