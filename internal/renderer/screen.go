package renderer

import (
	"github.com/pai0id/CgCourseProject/internal/object"
)

func ToScreen(vertex object.Vec3, width, height int) Point {
	xScreen := -int(vertex.X) + width/2
	yScreen := -int(vertex.Y) + height/2
	return Point{X: xScreen, Y: yScreen, Z: vertex.Z}
}

func screen(in <-chan *object.Face, out chan<- *Polygon, width, height int) {
	for f := range in {
		p := Polygon{Vertices: make([]Point, 3), Normals: f.Normals, Intensities: f.Intensities, Skeletonize: f.Skeletonize}
		for i, v := range f.Vertices {
			p.Vertices[i] = ToScreen(v, width, height)
		}
		out <- &p
	}
}
