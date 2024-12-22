package renderer

import "github.com/pai0id/CgCourseProject/internal/object"

type Camera struct {
	Fov    float64
	Z      float64
	ZFar   float64
	ZNear  float64
	Aspect float64
}

type Light struct {
	Position  object.Vec3
	Intensity float64
}

type RenderOptions struct {
	Width        int
	Height       int
	Cam          *Camera
	LightSources map[int64]Light
}

func getLights(lMap map[int64]Light) []Light {
	lights := make([]Light, len(lMap))
	i := 0
	for _, light := range lMap {
		lights[i] = light
	}
	return lights
}
