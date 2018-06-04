package cpu

const NODETYPE_NODE = 0
const NODETYPE_LEAF = 1

const AXIS_X = 0
const AXIS_Y = 1
const AXIS_Z = 2

type KDNode struct {
	Type int
	Axis int
	Pos float32

	TMin float32
	TMax float32

	TriangleIds []int
	NumTris int

	Children *[]KDNode
}
