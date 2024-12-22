package renderer

import (
	"sync"

	"github.com/pai0id/CgCourseProject/internal/object"
)

func toScreen(vertex object.Vec3, width, height int) point {
	xScreen := -int(vertex.X) + width/2
	yScreen := -int(vertex.Y) + height/2
	return point{x: xScreen, y: yScreen, z: vertex.Z}
}

func screen(in <-chan *object.Face, out chan<- *polygon, wg *sync.WaitGroup, width, height int) {
	defer wg.Done()
	for f := range in {
		p := polygon{vertices: make([]point, 3), normals: f.Normals, intensities: f.Intensities, skeletonize: f.Skeletonize}
		for i, v := range f.Vertices {
			p.vertices[i] = toScreen(v, width, height)
		}
		out <- &p
	}
	if out != nil {
		close(out)
	}
}
