package reader

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Vertex struct {
	X, Y, Z float64
}

type Face struct {
	Vertices []Vertex
	Normal   Vertex
}

type Model struct {
	Faces  []Face
	Center Vertex
}

func LoadOBJ(filepath string) (*Model, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var vertices []Vertex
	var normals []Vertex
	model := &Model{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case strings.HasPrefix(line, "v "):
			vertex, err := parseVertex(line)
			if err != nil {
				return nil, fmt.Errorf("failed to parse vertex: %w", err)
			}
			vertices = append(vertices, vertex)
		case strings.HasPrefix(line, "vn "):
			normal, err := parseNormal(line)
			if err != nil {
				return nil, fmt.Errorf("failed to parse normal: %w", err)
			}
			normals = append(normals, normal)
		case strings.HasPrefix(line, "f "):
			face, err := parseFace(line, vertices, normals)
			if err != nil {
				return nil, fmt.Errorf("failed to parse face: %w", err)
			}
			model.Faces = append(model.Faces, face)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	model.CalculateCenter()

	return model, nil
}

func parseVertex(line string) (Vertex, error) {
	parts := strings.Fields(line)
	if len(parts) < 4 {
		return Vertex{}, fmt.Errorf("invalid vertex line: %s", line)
	}
	x, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return Vertex{}, fmt.Errorf("invalid x coordinate: %w", err)
	}
	y, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return Vertex{}, fmt.Errorf("invalid y coordinate: %w", err)
	}
	z, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return Vertex{}, fmt.Errorf("invalid z coordinate: %w", err)
	}
	return Vertex{X: x, Y: y, Z: z}, nil
}

func parseNormal(line string) (Vertex, error) {
	parts := strings.Fields(line)
	if len(parts) < 4 {
		return Vertex{}, fmt.Errorf("invalid normal line: %s", line)
	}
	x, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return Vertex{}, fmt.Errorf("invalid x coordinate: %w", err)
	}
	y, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return Vertex{}, fmt.Errorf("invalid y coordinate: %w", err)
	}
	z, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return Vertex{}, fmt.Errorf("invalid z coordinate: %w", err)
	}
	return Vertex{X: x, Y: y, Z: z}, nil
}

func parseFace(line string, vertices []Vertex, normals []Vertex) (Face, error) {
	parts := strings.Fields(line)
	if len(parts) != 4 {
		return Face{}, fmt.Errorf("invalid face line: %s, must have exactly 3 vertices", line)
	}

	var faceVertices = make([]Vertex, 3)
	var faceNormal Vertex
	normalSet := false

	for i, part := range parts[1:] {
		indices := strings.Split(part, "/")
		if len(indices) < 1 {
			return Face{}, fmt.Errorf("invalid face element: %s", part)
		}

		vertexIndex, err := strconv.Atoi(indices[0])
		if err != nil || vertexIndex < 1 || vertexIndex > len(vertices) {
			return Face{}, fmt.Errorf("invalid vertex index: %w", err)
		}
		faceVertices[i] = vertices[vertexIndex-1]

		if len(indices) > 2 && len(indices[2]) > 0 {
			if !normalSet {
				normalIndex, err := strconv.Atoi(indices[2])
				if err != nil || normalIndex < 1 || normalIndex > len(normals) {
					return Face{}, fmt.Errorf("invalid normal index: %w", err)
				}
				faceNormal = normals[normalIndex-1]
				normalSet = true
			}
		}
	}

	if !normalSet {
		return Face{}, fmt.Errorf("face lacks normal index")
	}

	return Face{
		Vertices: faceVertices,
		Normal:   faceNormal,
	}, nil
}

func (model *Model) CalculateCenter() {
	if len(model.Faces) == 0 {
		model.Center = Vertex{X: 0, Y: 0, Z: 0}
		return
	}

	var sumX, sumY, sumZ float64
	var count int

	for _, face := range model.Faces {
		for _, vertex := range face.Vertices {
			sumX += vertex.X
			sumY += vertex.Y
			sumZ += vertex.Z
			count++
		}
	}

	model.Center = Vertex{
		X: sumX / float64(count),
		Y: sumY / float64(count),
		Z: sumZ / float64(count),
	}
}
