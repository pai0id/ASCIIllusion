package renderer

import (
	"sync"

	"github.com/pai0id/CgCourseProject/internal/object"
)

func ToScreen(vertex object.Vec3, width, height int) Point {
	xScreen := -int(vertex.X) + width/2
	yScreen := -int(vertex.Y) + height/2
	return Point{X: xScreen, Y: yScreen, Z: vertex.Z}
}

func screen(in <-chan *object.Face, out chan<- *Polygon, wg *sync.WaitGroup, width, height int) {
	defer wg.Done()
	for f := range in {
		p := Polygon{Vertices: make([]Point, 3), Normals: f.Normals, Intensities: f.Intensities, Skeletonize: f.Skeletonize}
		for i, v := range f.Vertices {
			p.Vertices[i] = ToScreen(v, width, height)
		}
		out <- &p
	}
	if out != nil {
		close(out)
	}
}
