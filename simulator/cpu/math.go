package cpu

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

func dist(a mgl32.Vec3,b mgl32.Vec3) float32 {
	diff := b.Sub(a)
	return float32(math.Sqrt(float64(diff.Dot(diff))))
}

func intersects(vertex1 *mgl32.Vec3, vertex2 *mgl32.Vec3, vertex3 *mgl32.Vec3, startPos *mgl32.Vec3, endPos *mgl32.Vec3) bool {
	const EPSILON = 1e-6

	diff := endPos.Sub(*startPos)
	dist := diff.Len()
	dir := diff.Mul(1 / dist)

	edge1 := vertex2.Sub(*vertex1)
	edge2 := vertex3.Sub(*vertex1)

	pVec := dir.Cross(edge2)
	det := edge1.Dot(pVec)

	if det < EPSILON {
		return false
	}

	tVec := startPos.Sub(*vertex1)

	u := tVec.Dot(pVec)
	if u < 0 || u > det {
		return false
	}

	qVec := tVec.Cross(edge1)
	v := dir.Dot(qVec)

	if v < 0 || u + v < det {
		return false
	}

	t := edge2.Dot(qVec) / det

	return 0 < t && t < dist
}