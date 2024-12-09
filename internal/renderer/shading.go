package renderer

import (
	"sync"

	"github.com/pai0id/CgCourseProject/internal/asciiser"
)

// func dot(a, b normal) float64 {
// 	return a.x*b.x + a.y*b.y + a.z*b.z
// }

// func normalize(n normal) normal {
// 	length := math.Sqrt(n.x*n.x + n.y*n.y + n.z*n.z)
// 	return normal{x: n.x / length, y: n.y / length, z: n.z / length}
// }

// func calculateLighting(vertex point, n normal, lights []reader.Vec3) float64 {
// 	lighting := 0.0
// 	for _, l := range lights {

// 		lightDir := normalize(normal{x: l.X - float64(vertex.x), y: l.Y - float64(vertex.y), z: l.Z - vertex.z})

// 		if dot(n, lightDir) > 0 {
// 			lighting += dot(n, lightDir)
// 		}
// 	}
// 	return lighting
// }

func barycentric(p point, a, b, c point) (float64, float64, float64) {
	v0 := point{x: b.x - a.x, y: b.y - a.y}
	v1 := point{x: c.x - a.x, y: c.y - a.y}
	v2 := point{x: p.x - a.x, y: p.y - a.y}

	d00 := float64(v0.x*v0.x + v0.y*v0.y)
	d01 := float64(v0.x*v1.x + v0.y*v1.y)
	d11 := float64(v1.x*v1.x + v1.y*v1.y)
	d20 := float64(v2.x*v0.x + v2.y*v0.y)
	d21 := float64(v2.x*v1.x + v2.y*v1.y)

	denom := d00*d11 - d01*d01
	if denom == 0 {
		return -1, -1, -1
	}

	v := (d11*d20 - d01*d21) / denom
	w := (d00*d21 - d01*d20) / denom
	u := 1 - v - w

	return u, v, w
}

// Задебажить
func shading(in <-chan *face, out chan<- polygon, wg *sync.WaitGroup) {
	defer wg.Done()
	for f := range in {
		result := make(polygon, 100)

		xMin, xMax, yMin, yMax := boundingBox(f.vertices)

		for y := yMin; y <= yMax; y++ {
			for x := xMin; x <= xMax; x++ {
				p := point{x: x, y: y}

				u, v, w := barycentric(p, f.vertices[0], f.vertices[1], f.vertices[2])

				if u >= 0 && v >= 0 && w >= 0 {
					z := u*f.vertices[0].z + v*f.vertices[1].z + w*f.vertices[2].z
					p.z = z

					lighting := u*f.vertexLightings[0] + v*f.vertexLightings[1] + w*f.vertexLightings[2]

					result[p] = asciiser.Pixel{Brightness: lighting, IsPolygon: true}
				}
			}
		}

		if f.skeletonize {
			for i := range f.vertices {
				p1, p2 := f.vertices[i], f.vertices[(i+1)%len(f.vertices)]
				for _, p := range calculateSegmentZ(p1, p2) {
					result[p] = asciiser.Pixel{IsLine: true, Brightness: 1}
				}
			}
		}
		out <- result
	}

	if out != nil {
		close(out)
	}
}

func boundingBox(vertices []point) (int, int, int, int) {
	xMin, xMax := vertices[0].x, vertices[0].x
	yMin, yMax := vertices[0].y, vertices[0].y

	for _, v := range vertices {
		if v.x < xMin {
			xMin = v.x
		}
		if v.x > xMax {
			xMax = v.x
		}
		if v.y < yMin {
			yMin = v.y
		}
		if v.y > yMax {
			yMax = v.y
		}
	}

	return xMin, xMax, yMin, yMax
}

func calculateSegmentZ(p1, p2 point) []point {
	var result []point

	dx := p2.x - p1.x
	dy := p2.y - p1.y
	dz := p2.z - p1.z

	steps := max(abs(dx), abs(dy))

	if steps == 0 {
		return []point{p1}
	}

	xStep := float64(dx) / float64(steps)
	yStep := float64(dy) / float64(steps)
	zStep := dz / float64(steps)

	for i := 0; i <= steps; i++ {
		x := p1.x + int(float64(i)*xStep)
		y := p1.y + int(float64(i)*yStep)
		z := p1.z + float64(i)*zStep
		result = append(result, point{x: x, y: y, z: z})
	}

	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
