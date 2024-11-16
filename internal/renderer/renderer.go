package renderer

// import (
// 	"math"
// )

// // Vertex represents a 3D point
// type Vertex struct {
// 	X, Y, Z float64
// }

// // Face represents a face defined by vertex indices
// type Face struct {
// 	VertexIndices []int
// }

// // Model contains the vertices and faces of the 3D model
// type Model struct {
// 	Vertices []Vertex
// 	Faces    []Face
// }

// // RenderConfig contains settings for rendering
// type RenderConfig struct {
// 	ImageWidth, ImageHeight int
// 	FieldOfView             float64 // Perspective field of view in degrees
// 	NearPlane, FarPlane     float64 // Clipping planes
// }

// // ProjectVertex applies a perspective transformation and projects a vertex to 2D
// func ProjectVertex(v Vertex, config RenderConfig) (int, int) {
// 	aspectRatio := float64(config.ImageWidth) / float64(config.ImageHeight)
// 	fovRad := 1 / math.Tan(config.FieldOfView*math.Pi/180/2)

// 	// Perspective transformation
// 	zScale := 1 / (v.Z + config.NearPlane + 1) // Prevent division by zero
// 	screenX := int((v.X * fovRad * aspectRatio) * zScale * float64(config.ImageWidth/2))
// 	screenY := int((v.Y * fovRad) * zScale * float64(config.ImageHeight/2))

// 	// Center the model in the image
// 	screenX += config.ImageWidth / 2
// 	screenY = config.ImageHeight/2 - screenY

// 	return screenX, screenY
// }

// // RenderModel generates a boolean grid representation of the model
// func RenderModel(model Model, config RenderConfig) [][]bool {
// 	// Create an empty grid
// 	grid := make([][]bool, config.ImageHeight)
// 	for i := range grid {
// 		grid[i] = make([]bool, config.ImageWidth)
// 	}

// 	// Project vertices to 2D
// 	projectedVertices := make([][2]int, len(model.Vertices))
// 	for i, vertex := range model.Vertices {
// 		projectedVertices[i] = [2]int{
// 			ProjectVertex(vertex, config),
// 		}
// 	}

// 	// Render edges of faces onto the grid
// 	for _, face := range model.Faces {
// 		numVertices := len(face.VertexIndices)
// 		for i := 0; i < numVertices; i++ {
// 			start := projectedVertices[face.VertexIndices[i]]
// 			end := projectedVertices[face.VertexIndices[(i+1)%numVertices]] // Connect back to the first vertex
// 			drawLine(grid, start[0], start[1], end[0], end[1])
// 		}
// 	}

// 	return grid
// }

// // drawLine draws a line on the grid between two points using Bresenham's algorithm
// func drawLine(grid [][]bool, x0, y0, x1, y1 int) {
// 	dx := math.Abs(float64(x1 - x0))
// 	dy := math.Abs(float64(y1 - y0))
// 	sx := -1
// 	if x0 < x1 {
// 		sx = 1
// 	}
// 	sy := -1
// 	if y0 < y1 {
// 		sy = 1
// 	}
// 	err := dx - dy

// 	for {
// 		if x0 >= 0 && x0 < len(grid[0]) && y0 >= 0 && y0 < len(grid) {
// 			grid[y0][x0] = true
// 		}
// 		if x0 == x1 && y0 == y1 {
// 			break
// 		}
// 		e2 := 2 * err
// 		if e2 > -dy {
// 			err -= dy
// 			x0 += sx
// 		}
// 		if e2 < dx {
// 			err += dx
// 			y0 += sy
// 		}
// 	}
// }
