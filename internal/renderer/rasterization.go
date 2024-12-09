package renderer

import (
	"fmt"
	"math"
	"sync"

	"github.com/pai0id/CgCourseProject/internal/reader"
	"github.com/pai0id/CgCourseProject/internal/transformer"
)

const shininessC = 32.0
const intensityC = 0.5

func calculateLighting(point, normal reader.Vec3, lightSources []reader.Vec3) float64 {
	normal = transformer.Normalize(normal)
	viewDirection := reader.Vec3{X: 0, Y: 0, Z: 1}

	totalLight := 0.0

	for _, light := range lightSources {
		lightDir := transformer.Normalize(transformer.Subtract(light, point))

		fmt.Println(transformer.Dot(normal, lightDir))
		if transformer.Dot(normal, lightDir) <= 0 {
			continue
		}

		diffuseFactor := math.Max(0, transformer.Dot(normal, lightDir))
		diffuse := intensityC * diffuseFactor

		reflection := transformer.Normalize(transformer.Subtract(transformer.MultiplyScalar(normal, 2*diffuseFactor), lightDir))
		specularFactor := math.Pow(math.Max(0, transformer.Dot(viewDirection, reflection)), shininessC)
		specular := intensityC * specularFactor

		totalLight += diffuse + specular
	}

	return math.Min(1, math.Max(0, totalLight))
}

func rasterization(in <-chan *reader.Model, out chan<- *face, wg *sync.WaitGroup, projectionMatrix transformer.Mat4, lights []reader.Vec3, width, height int) {
	defer wg.Done()
	for m := range in {
		faceLightings := make([][]float64, len(m.Faces))
		for k, f := range m.Faces {
			faceLightings[k] = make([]float64, len(f.Vertices))
			for i := range f.Vertices {
				faceLightings[k][i] = calculateLighting(f.Vertices[i], f.Normals[i], lights)
			}
		}
		m := transformer.ProjectModel(m, projectionMatrix)
		for k, f := range m.Faces {
			screenFace := face{skeletonize: m.Skeletonize, vertexLightings: faceLightings[k]}
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
	xScreen := int(vertex.X) + width/2
	yScreen := int(vertex.Y) + height/2
	return point{x: xScreen, y: yScreen, z: vertex.Z}
}
