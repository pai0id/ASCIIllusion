package renderer

import (
	"github.com/pai0id/CgCourseProject/internal/asciiser"
)

func render(polygons []polygon, img asciiser.Image, zb zBuffer) {
	h := len(img)
	w := len(img[0])
	for _, p := range polygons {
		for i, px := range p {
			x, y := i.x, i.y
			if x >= 0 && x < w && y >= 0 && y < h {
				if zb[y][x] > i.z {
					img[y][x] = px
					zb[y][x] = i.z
				}
			}
		}
	}
}
