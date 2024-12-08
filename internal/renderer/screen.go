package renderer

import (
	"sync"

	"github.com/pai0id/CgCourseProject/internal/reader"
)

func rasterization(in <-chan *reader.Model, out chan<- *face, wg *sync.WaitGroup, width, height int) {
	defer wg.Done()
	for m := range in {
		for _, f := range m.Faces {
			screenFace := face{skeletonize: m.Skeletonize}
			for i, vertex := range f.Vertices {
				screenPoint := NDCToScreen(vertex, width, height)
				screenFace.vertices = append(screenFace.vertices, screenPoint)
				screenFace.normals = append(screenFace.normals, normal{x: f.Normals[i].X, y: f.Normals[i].Y, z: f.Normals[i].Z})
			}
			out <- &screenFace
		}
	}

	if out != nil {
		close(out)
	}
}

func NDCToScreen(vertex reader.Vec3, width, height int) point {
	xScreen := -int(vertex.X) + width/2
	yScreen := -int(vertex.Y) + height/2
	return point{x: xScreen, y: yScreen, z: vertex.Z}
}
