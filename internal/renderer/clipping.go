package renderer

// import (
// 	"sync"

// 	"github.com/pai0id/CgCourseProject/internal/object"
// )

// type plane struct {
// 	A, B, C, D float64
// }

// func pointPlaneDistance(point object.Vec3, plane plane) float64 {
// 	return plane.A*point.X + plane.B*point.Y + plane.C*point.Z + plane.D
// }

// func intersectEdgePlane(
// 	start, end object.Vec3,
// 	startNormal, endNormal object.Vec3,
// 	plane plane,
// ) (object.Vec3, object.Vec3) {
// 	startDist := pointPlaneDistance(start, plane)
// 	endDist := pointPlaneDistance(end, plane)
// 	t := startDist / (startDist - endDist)

// 	intersection := object.InterpolateVec3(start, end, t)
// 	interpolatedNormal := object.InterpolateVec3(startNormal, endNormal, t)
// 	return intersection, interpolatedNormal
// }

// func clipFace(face object.Face, plane plane) object.Face {
// 	clippedVertices := make([]object.Vec3, 0, 4)
// 	clippedNormals := make([]object.Vec3, 0, 4)

// 	for i := 0; i < len(face.Vertices); i++ {
// 		currentVertex := face.Vertices[i]
// 		previousVertex := face.Vertices[(i+len(face.Vertices)-1)%len(face.Vertices)]
// 		currentNormal := face.Normals[i]
// 		previousNormal := face.Normals[(i+len(face.Normals)-1)%len(face.Normals)]

// 		currentDist := pointPlaneDistance(currentVertex, plane)
// 		previousDist := pointPlaneDistance(previousVertex, plane)

// 		if currentDist >= 0 {

// 			if previousDist < 0 {

// 				intersection, interpolatedNormal := intersectEdgePlane(
// 					previousVertex, currentVertex, previousNormal, currentNormal, plane,
// 				)
// 				clippedVertices = append(clippedVertices, intersection)
// 				clippedNormals = append(clippedNormals, interpolatedNormal)
// 			}

// 			clippedVertices = append(clippedVertices, currentVertex)
// 			clippedNormals = append(clippedNormals, currentNormal)
// 		} else if previousDist >= 0 {

// 			intersection, interpolatedNormal := intersectEdgePlane(
// 				previousVertex, currentVertex, previousNormal, currentNormal, plane,
// 			)
// 			clippedVertices = append(clippedVertices, intersection)
// 			clippedNormals = append(clippedNormals, interpolatedNormal)
// 		}
// 	}

// 	return object.Face{
// 		Vertices: clippedVertices,
// 		Normals:  clippedNormals,
// 	}
// }

// func triangulate(f object.Face) []object.Face {
// 	if len(f.Vertices) < 3 {
// 		return nil
// 	}

// 	triangles := make([]object.Face, 0, 4)
// 	anchor := f.Vertices[0]
// 	anchorNor := f.Normals[0]

// 	for i := 1; i < len(f.Vertices)-1; i++ {
// 		currVer := []object.Vec3{anchor, f.Vertices[i], f.Vertices[i+1]}
// 		currNor := []object.Vec3{anchorNor, f.Normals[i], f.Normals[i+1]}
// 		triangles = append(triangles, object.Face{Vertices: currVer, Normals: currNor})
// 	}

// 	return triangles
// }

// func ClipAndTriangulate(f object.Face, zNear, zFar, cameraWidth, cameraHeight float64) []object.Face {
// 	halfWidth := cameraWidth / 2
// 	halfHeight := cameraHeight / 2

// 	planes := []plane{
// 		{-1, 0, 0, halfWidth},
// 		{1, 0, 0, halfWidth},
// 		{0, -1, 0, halfHeight},
// 		{0, 1, 0, halfHeight},
// 		{0, 0, 1, -zNear},
// 		{0, 0, -1, zFar},
// 	}

// 	clipped := f
// 	for _, p := range planes {
// 		clipped = clipFace(clipped, p)
// 		if len(clipped.Vertices) == 0 {
// 			break
// 		}
// 	}

// 	return triangulate(clipped)
// }

// func clipping(in <-chan *object.Object, out chan<- *object.Object, wg *sync.WaitGroup, zNear, zFar, cameraWidth, cameraHeight float64) {
// 	defer wg.Done()
// 	for m := range in {
// 		clippedModel := object.Object{Skeletonize: m.Skeletonize}
// 		for _, face := range m.Faces {
// 			clippedFaces := ClipAndTriangulate(face, zNear, zFar, cameraWidth, cameraHeight)
// 			clippedModel.Faces = append(clippedModel.Faces, clippedFaces...)
// 		}
// 		out <- &clippedModel
// 	}
// 	if out != nil {
// 		close(out)
// 	}
// }
