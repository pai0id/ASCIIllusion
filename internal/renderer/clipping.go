package renderer

import (
	"sync"

	"github.com/pai0id/CgCourseProject/internal/reader"
	"github.com/pai0id/CgCourseProject/internal/transformer"
)

type plane struct {
	A, B, C, D float64
}

func pointPlaneDistance(point reader.Vec3, plane plane) float64 {
	return plane.A*point.X + plane.B*point.Y + plane.C*point.Z + plane.D
}

func intersectEdgePlane(
	start, end reader.Vec3,
	startNormal, endNormal reader.Vec3,
	plane plane,
) (reader.Vec3, reader.Vec3) {
	startDist := pointPlaneDistance(start, plane)
	endDist := pointPlaneDistance(end, plane)
	t := startDist / (startDist - endDist)

	intersection := transformer.InterpolateVec3(start, end, t)
	interpolatedNormal := transformer.InterpolateVec3(startNormal, endNormal, t)
	return intersection, interpolatedNormal
}

func clipFace(face reader.Face, plane plane) reader.Face {
	clippedVertices := make([]reader.Vec3, 0, 4)
	clippedNormals := make([]reader.Vec3, 0, 4)

	for i := 0; i < len(face.Vertices); i++ {
		currentVertex := face.Vertices[i]
		previousVertex := face.Vertices[(i+len(face.Vertices)-1)%len(face.Vertices)]
		currentNormal := face.Normals[i]
		previousNormal := face.Normals[(i+len(face.Normals)-1)%len(face.Normals)]

		currentDist := pointPlaneDistance(currentVertex, plane)
		previousDist := pointPlaneDistance(previousVertex, plane)

		if currentDist >= 0 {

			if previousDist < 0 {

				intersection, interpolatedNormal := intersectEdgePlane(
					previousVertex, currentVertex, previousNormal, currentNormal, plane,
				)
				clippedVertices = append(clippedVertices, intersection)
				clippedNormals = append(clippedNormals, interpolatedNormal)
			}

			clippedVertices = append(clippedVertices, currentVertex)
			clippedNormals = append(clippedNormals, currentNormal)
		} else if previousDist >= 0 {

			intersection, interpolatedNormal := intersectEdgePlane(
				previousVertex, currentVertex, previousNormal, currentNormal, plane,
			)
			clippedVertices = append(clippedVertices, intersection)
			clippedNormals = append(clippedNormals, interpolatedNormal)
		}
	}

	return reader.Face{
		Vertices: clippedVertices,
		Normals:  clippedNormals,
	}
}

func triangulate(f reader.Face) []reader.Face {
	if len(f.Vertices) < 3 {
		return nil
	}

	triangles := make([]reader.Face, 0, 4)
	anchor := f.Vertices[0]
	anchorNor := f.Normals[0]

	for i := 1; i < len(f.Vertices)-1; i++ {
		currVer := []reader.Vec3{anchor, f.Vertices[i], f.Vertices[i+1]}
		currNor := []reader.Vec3{anchorNor, f.Normals[i], f.Normals[i+1]}
		triangles = append(triangles, reader.Face{Vertices: currVer, Normals: currNor})
	}

	return triangles
}

func ClipAndTriangulate(f reader.Face, zNear, zFar, cameraWidth, cameraHeight float64) []reader.Face {
	halfWidth := cameraWidth / 2
	halfHeight := cameraHeight / 2

	planes := []plane{
		{-1, 0, 0, halfWidth},
		{1, 0, 0, halfWidth},
		{0, -1, 0, halfHeight},
		{0, 1, 0, halfHeight},
		{0, 0, 1, -zNear},
		{0, 0, -1, zFar},
	}

	clipped := f
	for _, p := range planes {
		clipped = clipFace(clipped, p)
		if len(clipped.Vertices) == 0 {
			break
		}
	}

	return triangulate(clipped)
}

func clipping(in <-chan *reader.Model, out chan<- *reader.Model, wg *sync.WaitGroup, zNear, zFar, cameraWidth, cameraHeight float64) {
	defer wg.Done()
	for m := range in {
		clippedModel := reader.Model{Skeletonize: m.Skeletonize}
		for _, face := range m.Faces {
			clippedFaces := ClipAndTriangulate(face, zNear, zFar, cameraWidth, cameraHeight)
			clippedModel.Faces = append(clippedModel.Faces, clippedFaces...)
		}
		out <- &clippedModel
	}
	if out != nil {
		close(out)
	}
}
