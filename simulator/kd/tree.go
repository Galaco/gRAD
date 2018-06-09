package kd

import "github.com/go-gl/mathgl/mgl32"

func CreateRootNode(tMin mgl32.Vec3, tMax mgl32.Vec3, nodeType int, triangleIDs []int, numTriangles int) *Node {
	return &Node{
		TMin: tMin,
		TMax: tMax,
		Type: nodeType,
		TriangleIds: triangleIDs,
		NumTris: numTriangles,
	}
}