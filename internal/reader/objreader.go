package reader

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pai0id/CgCourseProject/internal/object"
)

func LoadOBJ(filepath string) (*object.Object, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var vertices []object.Vec3
	var normals []object.Vec3
	model := &object.Object{}
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
				return nil, fmt.Errorf("failed to parse object.Face: %w", err)
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

func parseVertex(line string) (object.Vec3, error) {
	parts := strings.Fields(line)
	if len(parts) < 4 {
		return object.Vec3{}, fmt.Errorf("invalid vertex line: %s", line)
	}
	x, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return object.Vec3{}, fmt.Errorf("invalid x coordinate: %w", err)
	}
	y, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return object.Vec3{}, fmt.Errorf("invalid y coordinate: %w", err)
	}
	z, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return object.Vec3{}, fmt.Errorf("invalid z coordinate: %w", err)
	}
	return object.Vec3{X: x, Y: y, Z: z}, nil
}

func parseNormal(line string) (object.Vec3, error) {
	parts := strings.Fields(line)
	if len(parts) < 4 {
		return object.Vec3{}, fmt.Errorf("invalid normal line: %s", line)
	}
	x, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return object.Vec3{}, fmt.Errorf("invalid x coordinate: %w", err)
	}
	y, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return object.Vec3{}, fmt.Errorf("invalid y coordinate: %w", err)
	}
	z, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return object.Vec3{}, fmt.Errorf("invalid z coordinate: %w", err)
	}
	return object.Vec3{X: x, Y: y, Z: z}, nil
}

func parseFace(line string, vertices []object.Vec3, normals []object.Vec3) (object.Face, error) {
	parts := strings.Fields(line)
	if len(parts) != 4 {
		return object.Face{}, fmt.Errorf("invalid object.Face line: %s, must have exactly 3 vertices", line)
	}

	var faceVertices = make([]object.Vec3, 3)
	var faceNormals = make([]object.Vec3, 3)

	for i, part := range parts[1:] {
		indices := strings.Split(part, "/")
		if len(indices) < 1 {
			return object.Face{}, fmt.Errorf("invalid object.Face element: %s", part)
		}

		vertexIndex, err := strconv.Atoi(indices[0])
		if err != nil || vertexIndex < 1 || vertexIndex > len(vertices) {
			return object.Face{}, fmt.Errorf("invalid vertex index: %w", err)
		}
		faceVertices[i] = vertices[vertexIndex-1]

		normalIndex, err := strconv.Atoi(indices[2])
		if err != nil || normalIndex < 1 || normalIndex > len(normals) {
			return object.Face{}, fmt.Errorf("invalid normal index: %w", err)
		}
		faceNormals[i] = normals[normalIndex-1]
	}

	return object.Face{
		Vertices: faceVertices,
		Normals:  faceNormals,
	}, nil
}
