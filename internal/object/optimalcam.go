package object

func (m *Object) OptimalCameraDist() float64 {
	min, max := m.CalculateBoundingBox()

	depth := max.Z - min.Z

	return -depth + min.Z - 10
}
