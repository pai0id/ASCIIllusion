package mapping

import "github.com/pai0id/CgCourseProject/internal/fontparser"

type Cell interface {
	GetData() [][]bool
}

// Временно здесь, потом перенести
func FontMapToCellSlice(cms []fontparser.CharMatrix) []Cell {
	cells := make([]Cell, len(cms))
	for i, cm := range cms {
		cells[i] = cm
	}
	return cells
}
