package reader

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Vec3 struct {
	X, Y, Z float64
}

type Face struct {
	Vertices []Vec3
	Normals  []Vec3
}

type Model struct {
	Faces       []Face
	Center      Vec3
	Skeletonize bool
}

func LoadOBJ(filepath string) (*Model, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var vertices []Vec3
	var normals []Vec3
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

func parseVertex(line string) (Vec3, error) {
	parts := strings.Fields(line)
	if len(parts) < 4 {
		return Vec3{}, fmt.Errorf("invalid vertex line: %s", line)
	}
	x, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return Vec3{}, fmt.Errorf("invalid x coordinate: %w", err)
	}
	y, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return Vec3{}, fmt.Errorf("invalid y coordinate: %w", err)
	}
	z, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return Vec3{}, fmt.Errorf("invalid z coordinate: %w", err)
	}
	return Vec3{X: x, Y: y, Z: z}, nil
}

func parseNormal(line string) (Vec3, error) {
	parts := strings.Fields(line)
	if len(parts) < 4 {
		return Vec3{}, fmt.Errorf("invalid normal line: %s", line)
	}
	x, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return Vec3{}, fmt.Errorf("invalid x coordinate: %w", err)
	}
	y, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return Vec3{}, fmt.Errorf("invalid y coordinate: %w", err)
	}
	z, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return Vec3{}, fmt.Errorf("invalid z coordinate: %w", err)
	}
	return Vec3{X: x, Y: y, Z: z}, nil
}

func parseFace(line string, vertices []Vec3, normals []Vec3) (Face, error) {
	parts := strings.Fields(line)
	if len(parts) != 4 {
		return Face{}, fmt.Errorf("invalid face line: %s, must have exactly 3 vertices", line)
	}

	var faceVertices = make([]Vec3, 3)
	var faceNormals = make([]Vec3, 3)

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

		normalIndex, err := strconv.Atoi(indices[2])
		if err != nil || normalIndex < 1 || normalIndex > len(normals) {
			return Face{}, fmt.Errorf("invalid normal index: %w", err)
		}
		faceNormals[i] = normals[normalIndex-1]
	}

	return Face{
		Vertices: faceVertices,
		Normals:  faceNormals,
	}, nil
}

func (model *Model) CalculateCenter() {
	if len(model.Faces) == 0 {
		model.Center = Vec3{X: 0, Y: 0, Z: 0}
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

	model.Center = Vec3{
		X: math.Round(sumX / float64(count)),
		Y: math.Round(sumY / float64(count)),
		Z: math.Round(sumZ / float64(count)),
	}
}
