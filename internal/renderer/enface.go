package renderer

import (
	"github.com/pai0id/CgCourseProject/internal/object"
)

func enface(models []*object.Object, out chan<- *object.Face) {
	viewDirection := object.Vec3{X: 0, Y: 0, Z: 1}
	for _, m := range models {
		for _, f := range m.Faces {
			n := object.SumNormals(f.Normals)
			if n.Dot(viewDirection) > 0 {
				continue
			}
			out <- f.DeepCopy()
		}
	}
}
