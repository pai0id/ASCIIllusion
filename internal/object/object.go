package object

import "math"

type Face struct {
	Vertices    []Vec3
	Normals     []Vec3
	Intensities []float64
	Skeletonize bool
}

func (f *Face) DeepCopy() *Face {
	verticesCopy := make([]Vec3, len(f.Vertices))
	copy(verticesCopy, f.Vertices)

	normalsCopy := make([]Vec3, len(f.Normals))
	copy(normalsCopy, f.Normals)

	intensitiesCopy := make([]float64, len(f.Intensities))
	copy(intensitiesCopy, f.Intensities)

	return &Face{
		Vertices:    verticesCopy,
		Normals:     normalsCopy,
		Intensities: intensitiesCopy,
		Skeletonize: f.Skeletonize,
	}
}

type Object struct {
	Faces  []Face
	Center Vec3
}

func (o *Object) Skeletonize() {
	for i := range o.Faces {
		o.Faces[i].Skeletonize = true
	}
}

func (o *Object) CalculateBoundingBox() (Vec3, Vec3) {
	min := Vec3{X: math.MaxFloat64, Y: math.MaxFloat64, Z: math.MaxFloat64}
	max := Vec3{X: -math.MaxFloat64, Y: -math.MaxFloat64, Z: -math.MaxFloat64}

	for _, face := range o.Faces {
		for _, vertex := range face.Vertices {
			if vertex.X < min.X {
				min.X = vertex.X
			}
			if vertex.Y < min.Y {
				min.Y = vertex.Y
			}
			if vertex.Z < min.Z {
				min.Z = vertex.Z
			}
			if vertex.X > max.X {
				max.X = vertex.X
			}
			if vertex.Y > max.Y {
				max.Y = vertex.Y
			}
			if vertex.Z > max.Z {
				max.Z = vertex.Z
			}
		}
	}

	return min, max
}

func (o *Object) CalculateCenter() {
	if len(o.Faces) == 0 {
		o.Center = Vec3{X: 0, Y: 0, Z: 0}
		return
	}

	var sumX, sumY, sumZ float64
	var count int

	for _, face := range o.Faces {
		for _, vertex := range face.Vertices {
			sumX += vertex.X
			sumY += vertex.Y
			sumZ += vertex.Z
			count++
		}
	}

	o.Center = Vec3{
		X: math.Round(sumX / float64(count)),
		Y: math.Round(sumY / float64(count)),
		Z: math.Round(sumZ / float64(count)),
	}
}
