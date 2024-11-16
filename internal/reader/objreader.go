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
	VertexIndices []int
}

type Model struct {
	Vertices []Vertex
	Faces    []Face
}

func LoadOBJ(filepath string) (*Model, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	model := &Model{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "v ") {
			vertex, err := parseVertex(line)
			if err != nil {
				return nil, fmt.Errorf("failed to parse vertex: %w", err)
			}
			model.Vertices = append(model.Vertices, vertex)
		} else if strings.HasPrefix(line, "f ") {
			face, err := parseFace(line)
			if err != nil {
				return nil, fmt.Errorf("failed to parse face: %w", err)
			}
			model.Faces = append(model.Faces, face)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

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

func parseFace(line string) (Face, error) {
	parts := strings.Fields(line)
	if len(parts) < 4 {
		return Face{}, fmt.Errorf("invalid face line: %s", line)
	}
	indices := make([]int, 0, len(parts)-1)
	for _, part := range parts[1:] {

		vertexIndexStr := strings.Split(part, "/")[0]
		vertexIndex, err := strconv.Atoi(vertexIndexStr)
		if err != nil {
			return Face{}, fmt.Errorf("invalid vertex index: %w", err)
		}
		indices = append(indices, vertexIndex-1)
	}
	return Face{VertexIndices: indices}, nil
}
