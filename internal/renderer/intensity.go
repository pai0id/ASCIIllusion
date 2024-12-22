package renderer

import (
	"math"
	"sync"

	"github.com/pai0id/CgCourseProject/internal/object"
)

const (
	shininessC     = 32.0
	ambientC       = 0.0
	attenConst     = 0.01
	attenLinear    = 0.009
	attenQuadratic = 0.0032
)

func calculateVertexIntensity(point, normal object.Vec3, lightSources []Light) float64 {
	normal = normal.Normalize()
	viewDirection := object.Vec3{X: 0, Y: 0, Z: 1}

	totalLight := ambientC

	for _, light := range lightSources {
		lightDir := light.Position.Subtract(point).Normalize()
		distance := light.Position.Subtract(point).Length()

		attenuation := 1.0 / (attenConst + attenLinear*distance + attenQuadratic*distance*distance)

		if normal.Dot(lightDir) <= 0 {
			continue
		}

		diffuseFactor := math.Max(0, normal.Dot(lightDir))
		diffuse := light.Intensity * diffuseFactor * attenuation

		reflection := object.MultiplyScalar(normal, 2*diffuseFactor).Subtract(lightDir).Normalize()
		specularFactor := math.Pow(math.Max(0, viewDirection.Dot(reflection)), shininessC)
		specular := light.Intensity * specularFactor * attenuation

		totalLight += diffuse + specular
	}

	return math.Min(1, math.Max(0, totalLight))
}

func calcIntensity(in <-chan *object.Face, out chan<- *object.Face, wg *sync.WaitGroup, lightSrc []Light) {
	defer wg.Done()
	for f := range in {
		f.Intensities = make([]float64, 3)
		for i, v := range f.Vertices {
			f.Intensities[i] = calculateVertexIntensity(v, f.Normals[i], lightSrc)
		}
		out <- f
	}
	if out != nil {
		close(out)
	}
}
