package kd

import "github.com/go-gl/mathgl/mgl32"

const NODETYPE_NODE = 0
const NODETYPE_LEAF = 1

const AXIS_X = 0
const AXIS_Y = 1
const AXIS_Z = 2

type Node struct {
	Type int
	Axis int
	Pos float32

	TMin mgl32.Vec3
	TMax mgl32.Vec3

	TriangleIds []int
	NumTris int

	Children []Node
}
