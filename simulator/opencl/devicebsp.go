package opencl

import (
	"github.com/samuel/go-opencl/cl"
)

// Device independent bsp
type DeviceBsp struct {
	Models *cl.MemObject
	Planes *cl.MemObject
	Vertices *cl.MemObject
	Edges *cl.MemObject
	SurfEdges *cl.MemObject
	Faces *cl.MemObject

	// Track size of data structures
	NumModels int
	NumPlanes int
	NumVertices int
	NumEdges int
	NumSurfEdges int
	NumFaces int
}