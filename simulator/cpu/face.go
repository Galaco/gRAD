package cpu

import (
	"github.com/galaco/bsp/primitives/face"
	"github.com/galaco/bsp/primitives/plane"
	"github.com/galaco/bsp/primitives/texinfo"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/gRAD/filesystem"
	"unsafe"
)

type FaceInfo struct {
	Face face.Face
	Plane plane.Plane
	TexInfo texinfo.TexInfo
	Ainv mgl32.Mat3 // actually a 3x3
	//CUDAMatrix::CUDAMatrix<double, 3, 3> Ainv;
	FaceNorm mgl32.Vec3
	TotalLight mgl32.Vec3
	AvgLight mgl32.Vec3
	FaceIndex int
	LightmapWidth int
	LightmapHeight int
	LightmapSize int
	LightmapStartIndex int
	patchStartIndex int
}

func NewFaceInfo(vradBsp*filesystem.Bsp, faceIndex int) *FaceInfo{
	face := (*vradBsp.GetFaces())[faceIndex]
	plane := vradBsp.Planes[face.Planenum]
	lightmapWidth := int(face.LightmapTextureSizeInLuxels[0] + 1)
	lightmapHeight := int(face.LightmapTextureSizeInLuxels[1] + 1)
	faceInfo := FaceInfo{
		Face: face,
		Plane: plane,
		TexInfo: (*vradBsp.GetTexInfo())[face.TexInfo],
		//Ainv: cudaBSP.xyzMatrices[faceIndex]
		FaceNorm: plane.Normal,
		LightmapWidth: lightmapWidth,
		LightmapHeight: lightmapHeight,
		LightmapSize: lightmapWidth * lightmapHeight,
		LightmapStartIndex: int(face.Lightofs) / int(unsafe.Sizeof(face.Lightofs)), // / sizeof(BSP::RGBExp32)),
		TotalLight: mgl32.Vec3{0, 0, 0},
	}
	return &faceInfo
}

func (faceInfo *FaceInfo) XYXFromST(s float32, t float32) mgl32.Vec3 {
	sOffset := faceInfo.TexInfo.LightmapVecsLuxelsPerWorldUnits[0][3]
	tOffset := faceInfo.TexInfo.LightmapVecsLuxelsPerWorldUnits[1][3]

	sMin := faceInfo.Face.LightmapTextureMinsInLuxels[0]
	tMin := faceInfo.Face.LightmapTextureMinsInLuxels[1]

	B := mgl32.Mat3{}

	B[0] = float32(s - float32(sOffset + float32(sMin)))
	B[4] = float32(t - float32(tOffset + float32(tMin)))
	B[8] = faceInfo.Plane.Distance

	result := faceInfo.Ainv.Mul3(B)

	return mgl32.Vec3{
		result[0],
		result[4],
		result[8],
	}
}

