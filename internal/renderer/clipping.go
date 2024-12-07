package renderer

import (
	"sync"

	"github.com/pai0id/CgCourseProject/internal/reader"
)

func ClipFace(face reader.Face, zNear, zFar float64) []reader.Face {
	const ndcMin, ndcMax = -100.0, 100.0

	clippedFaces := []reader.Face{}
	currentVertices := face.Vertices
	currentNormals := face.Normals

	for axis := 0; axis < 3; axis++ {
		nextVertices := []reader.Vec3{}
		nextNormals := []reader.Vec3{}
		for i := 0; i < len(currentVertices); i++ {
			current := currentVertices[i]
			prev := currentVertices[(i+len(currentVertices)-1)%len(currentVertices)]
			currentNormal := currentNormals[i]
			prevNormal := currentNormals[(i+len(currentNormals)-1)%len(currentNormals)]

			insideCurrent := inside(current, axis, ndcMin, ndcMax, zNear, zFar)
			insidePrev := inside(prev, axis, ndcMin, ndcMax, zNear, zFar)

			if insideCurrent {
				if !insidePrev {
					intersection, normal := intersectWithNormal(prev, current, prevNormal, currentNormal, axis, ndcMin, ndcMax, zNear, zFar)
					nextVertices = append(nextVertices, intersection)
					nextNormals = append(nextNormals, normal)
				}
				nextVertices = append(nextVertices, current)
				nextNormals = append(nextNormals, currentNormal)
			} else if insidePrev {
				intersection, normal := intersectWithNormal(prev, current, prevNormal, currentNormal, axis, ndcMin, ndcMax, zNear, zFar)
				nextVertices = append(nextVertices, intersection)
				nextNormals = append(nextNormals, normal)
			}
		}
		currentVertices = nextVertices
		currentNormals = nextNormals
	}

	if len(currentVertices) >= 3 {
		clippedFaces = append(clippedFaces, reader.Face{Vertices: currentVertices, Normals: currentNormals})
	}
	return clippedFaces
}

func inside(vertex reader.Vec3, axis int, ndcMin, ndcMax, zNear, zFar float64) bool {
	switch axis {
	case 0:
		return vertex.X >= ndcMin && vertex.X <= ndcMax
	case 1:
		return vertex.Y >= ndcMin && vertex.Y <= ndcMax
	case 2:
		return vertex.Z >= zNear && vertex.Z <= zFar
	}
	return false
}

func intersectWithNormal(v1, v2, n1, n2 reader.Vec3, axis int, ndcMin, ndcMax, zNear, zFar float64) (reader.Vec3, reader.Vec3) {
	t := 0.0
	switch axis {
	case 0:
		if v2.X != v1.X {
			if v1.X < ndcMin {
				t = (ndcMin - v1.X) / (v2.X - v1.X)
			} else {
				t = (ndcMax - v1.X) / (v2.X - v1.X)
			}
		}
	case 1:
		if v2.Y != v1.Y {
			if v1.Y < ndcMin {
				t = (ndcMin - v1.Y) / (v2.Y - v1.Y)
			} else {
				t = (ndcMax - v1.Y) / (v2.Y - v1.Y)
			}
		}
	case 2:
		if v2.Z != v1.Z {
			if v1.Z < zNear {
				t = (zNear - v1.Z) / (v2.Z - v1.Z)
			} else {
				t = (zFar - v1.Z) / (v2.Z - v1.Z)
			}
		}
	}
	interpolatedVertex := reader.Vec3{
		X: v1.X + t*(v2.X-v1.X),
		Y: v1.Y + t*(v2.Y-v1.Y),
		Z: v1.Z + t*(v2.Z-v1.Z),
	}
	interpolatedNormal := reader.Vec3{
		X: n1.X + t*(n2.X-n1.X),
		Y: n1.Y + t*(n2.Y-n1.Y),
		Z: n1.Z + t*(n2.Z-n1.Z),
	}
	return interpolatedVertex, interpolatedNormal
}

func clip(in <-chan *reader.Model, out chan<- *reader.Model, wg *sync.WaitGroup, zNear, zFar float64) {
	defer wg.Done()
	for m := range in {
		// clippedModel := reader.Model{Skeletonize: m.Skeletonize}
		// for _, face := range m.Faces {
		// 	clippedFaces := ClipFace(face, zNear, zFar)
		// 	clippedModel.Faces = append(clippedModel.Faces, clippedFaces...)
		// }
		// out <- &clippedModel
		out <- m
	}
	if out != nil {
		close(out)
	}
}
