package raytracer

import "github.com/go-gl/mathgl/mgl32"

type Edge struct {
	Vertex1 mgl32.Vec3
	Vertex2 mgl32.Vec3
}

type Triangle struct {
	Vertices [3]mgl32.Vec3
}
