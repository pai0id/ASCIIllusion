package renderer

import "math"

func calculateVertexNormals(obj *object) map[point]normal {
	vertexNormals := make(map[point]normal)
	vertexCounts := make(map[point]int) // Track how many faces share each vertex

	// Accumulate normals for each vertex
	for _, face := range obj.faces {
		for _, vertex := range face.vertices {
			n := vertexNormals[vertex]
			n.x += face.normal.x
			n.y += face.normal.y
			n.z += face.normal.z
			vertexNormals[vertex] = n

			vertexCounts[vertex]++
		}
	}

	// Normalize the accumulated normals
	for vertex, normal := range vertexNormals {
		count := float64(vertexCounts[vertex])
		normal.x /= count
		normal.y /= count
		normal.z /= count
		vertexNormals[vertex] = normalizeNormal(normal)
	}

	return vertexNormals
}

// Normalize a normal vector
func normalizeNormal(n normal) normal {
	length := math.Sqrt(n.x*n.x + n.y*n.y + n.z*n.z)
	if length == 0 {
		return normal{0, 0, 0} // Avoid division by zero
	}
	return normal{n.x / length, n.y / length, n.z / length}
}
