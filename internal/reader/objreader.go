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

type Normal struct {
	X, Y, Z float64
}

type Face struct {
	VertexIndices []int
}

type Model struct {
	Vertices []Vertex
	Normals  []Normal
	Faces    []Face
	Center   Vertex
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
		switch {
		case strings.HasPrefix(line, "v "):
			vertex, err := parseVertex(line)
			if err != nil {
				return nil, fmt.Errorf("failed to parse vertex: %w", err)
			}
			model.Vertices = append(model.Vertices, vertex)
		case strings.HasPrefix(line, "vn "):
			normal, err := parseNormal(line)
			if err != nil {
				return nil, fmt.Errorf("failed to parse normal: %w", err)
			}
			model.Normals = append(model.Normals, normal)
		case strings.HasPrefix(line, "f "):
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

func parseNormal(line string) (Normal, error) {
	parts := strings.Fields(line)
	if len(parts) < 4 {
		return Normal{}, fmt.Errorf("invalid normal line: %s", line)
	}
	x, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return Normal{}, fmt.Errorf("invalid x coordinate: %w", err)
	}
	y, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return Normal{}, fmt.Errorf("invalid y coordinate: %w", err)
	}
	z, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return Normal{}, fmt.Errorf("invalid z coordinate: %w", err)
	}
	return Normal{X: x, Y: y, Z: z}, nil
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

func (model *Model) CalculateCenter() {
	if len(model.Vertices) == 0 {
		model.Center = Vertex{X: 0, Y: 0, Z: 0}
		return
	}

	var sumX, sumY, sumZ float64
	for _, vertex := range model.Vertices {
		sumX += vertex.X
		sumY += vertex.Y
		sumZ += vertex.Z
	}

	count := float64(len(model.Vertices))
	model.Center = Vertex{
		X: sumX / count,
		Y: sumY / count,
		Z: sumZ / count,
	}
}
