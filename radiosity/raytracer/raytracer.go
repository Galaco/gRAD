package raytracer

const MAX_TRIANGLES = 256000

// RayTracer
// Important. No pointers here!
// This is a generic raytracer useable by cpu and gpu
type RayTracer struct {
	Triangles [MAX_TRIANGLES]Triangle
	NumTriangles int
}

