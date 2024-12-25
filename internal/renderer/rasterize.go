package renderer

import (
	"github.com/pai0id/CgCourseProject/internal/asciiser"
)

const LineEPS = 0.005

func Barycentric(p Point, a, b, c Point) (float64, float64, float64) {
	v0 := Point{X: b.X - a.X, Y: b.Y - a.Y}
	v1 := Point{X: c.X - a.X, Y: c.Y - a.Y}
	v2 := Point{X: p.X - a.X, Y: p.Y - a.Y}

	d00 := float64(v0.X*v0.X + v0.Y*v0.Y)
	d01 := float64(v0.X*v1.X + v0.Y*v1.Y)
	d11 := float64(v1.X*v1.X + v1.Y*v1.Y)
	d20 := float64(v2.X*v0.X + v2.Y*v0.Y)
	d21 := float64(v2.X*v1.X + v2.Y*v1.Y)

	denom := d00*d11 - d01*d01
	if denom == 0 {
		return -1, -1, -1
	}

	v := (d11*d20 - d01*d21) / denom
	w := (d00*d21 - d01*d20) / denom
	u := 1 - v - w

	return u, v, w
}

func BoundingBox(vertices []Point) (int, int, int, int) {
	xMin, xMax := vertices[0].X, vertices[0].X
	yMin, yMax := vertices[0].Y, vertices[0].Y

	for _, v := range vertices {
		if v.X < xMin {
			xMin = v.X
		}
		if v.X > xMax {
			xMax = v.X
		}
		if v.Y < yMin {
			yMin = v.Y
		}
		if v.Y > yMax {
			yMax = v.Y
		}
	}

	return xMin, xMax, yMin, yMax
}

func rasterize(in <-chan *Polygon, img asciiser.Image, zb zBuffer) {
	for p := range in {
		rasterizePolygon(p, img, zb)
	}
}

func rasterizePolygon(p *Polygon, img asciiser.Image, zb zBuffer) {
	xMin, xMax, yMin, yMax := BoundingBox(p.Vertices)

	for y := yMin; y <= yMax; y++ {
		for x := xMin; x <= xMax; x++ {
			if x >= 0 && x < len(img[0]) && y >= 0 && y < len(img) {
				pt := Point{X: x, Y: y}
				u, v, w := Barycentric(pt, p.Vertices[0], p.Vertices[1], p.Vertices[2])

				if u >= 0 && v >= 0 && w >= 0 {
					z := u*p.Vertices[0].Z + v*p.Vertices[1].Z + w*p.Vertices[2].Z
					if zb[y][x] > z {
						lighting := u*p.Intensities[0] + v*p.Intensities[1] + w*p.Intensities[2]
						if p.Skeletonize && (u < LineEPS || v < LineEPS || w < LineEPS) {
							img[y][x] = asciiser.Pixel{IsLine: true}
						} else {
							img[y][x] = asciiser.Pixel{Brightness: lighting, IsPolygon: true}
						}
						zb[y][x] = z
					}
				}
			}
		}
	}
}
