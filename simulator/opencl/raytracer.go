package opencl

import (
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"time"
	"github.com/galaco/bsp/flags"
	"github.com/galaco/gRAD/bsp"
)

const MAX_TRIANGLES = 256000

type Edge struct {
	Vertex1 mgl32.Vec3
	Vertex2 mgl32.Vec3
}

type Triangle struct {
	Vertices [3]mgl32.Vec3
}

// RayTracer
// Important. No pointers here!
// This is a generic raytracer useable by cpu and gpu
type RayTracer struct {
	Triangles [MAX_TRIANGLES]Triangle
	NumTriangles int
}

func NewRayTracer() *RayTracer {
	return &RayTracer{}
}

func (tracer *RayTracer) SetupAccelerationStructure(vradBsp *bsp.Bsp) {
	log.Printf("Setting up ray-trace acceleration structure...\n")
	// Time preparation
	setupStart := time.Now().UnixNano() / int64(time.Millisecond)

	// Create raytracer

	// Create triangles from bsp faces
	triangles := [MAX_TRIANGLES]Triangle{}
	tIndex := 0

	// Add triangles to raytracer
	for _,face := range *vradBsp.GetFaces() {
		texFlags := (*vradBsp.GetTexInfo())[face.TexInfo].Flags

		// Skip translucent faces, but keep nodraw faces.
		if (texFlags & flags.SURF_TRANS) != 0 && 0 == (texFlags & flags.SURF_NODRAW) {
			continue
		}

		edges := []Edge{}
		firstEdge := int(face.FirstEdge)
		lastEdge := firstEdge + int(face.NumEdges)
		for i := firstEdge; i < lastEdge; i++ {
			surfEdge := (*vradBsp.GetSurfEdges())[i]
			firstToSecond := surfEdge >= 0

			if !firstToSecond {
				surfEdge *= -1
			}

			bspEdge := (*vradBsp.GetEdges())[surfEdge]
			edge := Edge{}
			vertices := *vradBsp.GetVertexes()

			if firstToSecond {
				edge.Vertex1 = vertices[bspEdge[0]]
				edge.Vertex2 = vertices[bspEdge[1]]
			} else {
				edge.Vertex1 = vertices[bspEdge[1]]
				edge.Vertex2 = vertices[bspEdge[0]]
			}

			edges = append(edges, edge)
		}

		iterator := 0
		var vertex1 = edges[iterator].Vertex1
		iterator++
		var vertex2 mgl32.Vec3
		var vertex3 = edges[iterator].Vertex1
		iterator++
		hasRun := false

		for ; iterator < len(edges) || hasRun == false; iterator++ {
			vertex2 = vertex3
			vertex3 = edges[iterator].Vertex1

			tri := Triangle{
				Vertices: [3]mgl32.Vec3{
					vertex1,
					vertex2,
					vertex3,
				},
			}

			triangles[tIndex] = tri
			tIndex++
			hasRun = true
		}
	}

	tracer.Triangles = triangles
	tracer.NumTriangles = tIndex

	// report elapsed time
	setupEnd := time.Now().UnixNano() / int64(time.Millisecond)
	log.Printf("Found %d triangles\n", tIndex)
	log.Printf("Done (%f seconds)\n\n", float32(setupEnd-setupStart) / 1000)
}
