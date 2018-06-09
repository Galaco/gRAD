package cpu

import (
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"time"
	"github.com/galaco/bsp/flags"
	"github.com/galaco/gRAD/filesystem"
	"github.com/galaco/gRAD/simulator/kd"
)

type Edge struct {
	Vertex1 mgl32.Vec3
	Vertex2 mgl32.Vec3
}

type Triangle struct {
	Vertices [3]mgl32.Vec3
}

// RayTracer
type RayTracer struct {
	Triangles []Triangle
	NumTriangles int
	TreeRoot *kd.Node
}

func NewRayTracer() *RayTracer {
	return &RayTracer{}
}

func (tracer *RayTracer) SetupAccelerationStructure(vradBsp *filesystem.Bsp) {
	log.Printf("Setting up ray-trace acceleration structure...\n")
	// Time preparation
	setupStart := time.Now().UnixNano() / int64(time.Millisecond)

	// Create raytracer

	// Create triangles from bsp faces
	triangles := []Triangle{}

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

			triangles = append(triangles, tri)
			hasRun = true
		}
	}

	tracer.Triangles = triangles
	tracer.NumTriangles = len(triangles)

	// report elapsed time
	setupEnd := time.Now().UnixNano() / int64(time.Millisecond)
	log.Printf("Found %d triangles\n", len(triangles))
	log.Printf("Done (%f seconds)\n\n", float32(setupEnd-setupStart) / 1000)
}


func (tracer *RayTracer) LOS_blocked(startPos mgl32.Vec3, endPos mgl32.Vec3) bool {
	EPSILON := float32(1e-6)

	dir := endPos.Sub(startPos).Normalize()
	invDir := mgl32.Vec3{0, 0, 0}
	if dir.X() < 0 {
		invDir[0] = 1.0 / (dir.X() + -EPSILON)
	} else {
		invDir[0] = 1.0 / (dir.X() + EPSILON)
	}
	if dir.Y() < 0 {
		invDir[1] = 1.0 / (dir.Y() + -EPSILON)
	} else {
		invDir[1] = 1.0 / (dir.Y() + EPSILON)
	}
	if dir.Z() < 0 {
		invDir[2] = 1.0 / (dir.Z() + -EPSILON)
	} else {
		invDir[2] = 1.0 / (dir.Z() + EPSILON)
	}

	type StackEntry struct {
		pNode *kd.Node
		start mgl32.Vec3
		end mgl32.Vec3
	}

	var stack [1024]StackEntry   // empty ascending stack
	stackSize := 0

	stack[stackSize] = StackEntry{
		tracer.TreeRoot,
		startPos,
		endPos,
	}
	//stackSize++

	for stackSize > 0 {
		if stackSize >= 1024 {
			log.Printf("ALERT: Stack size too big!!!\n")
			return false
		}

		//stackSize--
		entry := &stack[stackSize]

		pNode := entry.pNode
		start := entry.start
		end := entry.end

		length := dist(start, end)

		children := &pNode.Children

		var t float32

		switch pNode.Type {
		case kd.NODETYPE_LEAF:
			for ti := 0; ti < pNode.NumTris; ti++ {
				tri := &tracer.Triangles[pNode.TriangleIds[ti]]

				// The M-T intersection algorithm uses CCW vertex
				// winding, but Source uses CW winding. So, we need to
				// pass the vertices in reverse order to get backface
				// culling to work correctly.
				isLOSBlocked := intersects(
					&tri.Vertices[2], &tri.Vertices[1], &tri.Vertices[0],
					&startPos, &endPos)

				if isLOSBlocked {
					return true
				}
			}

			break

		case kd.NODETYPE_NODE:
			var dirPositive bool

			switch pNode.Axis {
			case kd.AXIS_X:
				t = (pNode.Pos - start.X()) * invDir.X()
				dirPositive = dir.X() >= 0.0
				break

			case kd.AXIS_Y:
				t = (pNode.Pos - start.Y()) * invDir.Y()
				dirPositive = dir.Y() >= 0.0
				break

			case kd.AXIS_Z:
				t = (pNode.Pos - start.Z()) * invDir.Z()
				dirPositive = dir.Z() >= 0.0
				break
			}

			if t < 0.0 {
				// Plane is "behind" the line start.
				// Recurse on the right side if dir is positive.
				// Recurse on the left side if dir is negative.

				cIndex := 0
				if dirPositive {
					cIndex = 1
				}
				stack[stackSize] = StackEntry{
					&(*children)[cIndex],
					start,
					end,
				}
				stackSize++
			} else if t >= length {
				// Plane is "ahead" of the line end.
				// Recurse on the left side if dir is positive.
				// Recurse on the right side if dir is negative.

				cIndex := 0
				if dirPositive {
					cIndex = 1
				}
				stack[stackSize] = StackEntry{
					&(*children)[cIndex],
					start,
					end,
				}
				stackSize++
			} else {
				// The line segment straddles the plane.
				// Clip the line and recurse on both sides.

				clipPoint := start.Add(dir.Mul(t))

				cIndex := 0
				if dirPositive {
					cIndex = 1
				}
				stack[stackSize] = StackEntry{
					&(*children)[cIndex],
					start,
					clipPoint,
				}
				stackSize++

				if stackSize >= 1024 {
					log.Printf("ALERT: Stack size too big!!!\n");
					return false
				}

				stack[stackSize] = StackEntry{
					&(*children)[cIndex],
					clipPoint,
					end,
				}
				stackSize++
			}
			break
		}
	}
	return false
}