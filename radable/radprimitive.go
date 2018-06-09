package radable

import (
	"github.com/galaco/bsp/primitives/plane"
	"github.com/galaco/bsp/primitives/texdata"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/bsp/primitives/visibility"
	"github.com/galaco/bsp/primitives/texinfo"
	"github.com/galaco/bsp/primitives/face"
	"github.com/galaco/bsp/primitives/leaf"
	"github.com/galaco/bsp/primitives/model"
	"github.com/galaco/bsp/primitives/brush"
	"github.com/galaco/bsp/primitives/brushside"
	"github.com/galaco/bsp/primitives/area"
	"github.com/galaco/bsp/primitives/areaportal"
	"github.com/galaco/bsp/primitives/mapflags"
	"github.com/galaco/bsp/primitives/vertnormal"
	"github.com/galaco/source-tools-common/entity"
	"github.com/galaco/source-tools-common/texdatastringtable"
	"github.com/galaco/gRAD/radable/light"
	"github.com/galaco/bsp/flags"
	"github.com/galaco/bsp/primitives/leafambientindex"
	"github.com/galaco/bsp/primitives/leafambientlighting"
	"github.com/galaco/bsp/primitives/common"
)

type RadPrimitive struct {
	// 1:1 Lump data extracts
	Planes []plane.Plane
	texData []texdata.TexData
	Vertexes []mgl32.Vec3
	visibility visibility.Vis
	texInfo []texinfo.TexInfo
	faces []face.Face
	facesHDR []face.Face
	leafs []leaf.Leaf
	leafFaces []uint16
	leafBrushes []uint16
	Edges [][2]uint16
	SurfEdges []int32
	Models []model.Model
	brushes []brush.Brush
	brushSides []brushside.BrushSide
	areas []area.Area
	areaPortals []areaportal.AreaPortal
	mapFlags mapflags.MapFlags
	vertNormals []vertnormal.VertNormal
	vertNormalIndices []uint16

	// Modified/derived lump datas
	Entities entity.List
	TexDataStringTable texdatastringtable.TexDataStringTable
	WorldLights []light.DirectLight
	AmbientLight light.DirectLight

	// Ambient lighting info
	ambientLightIndices []leafambientindex.LeafAmbientIndex
	ambientLightIndicesHDR []leafambientindex.LeafAmbientIndex
	ambientLightSamples []leafambientlighting.LeafAmbientLighting
	ambientLightSamplesHDR []leafambientlighting.LeafAmbientLighting

	LightSamples []mgl32.Vec3

	IsHDR bool
}

// PrepareAmbientSamples
func (f *RadPrimitive) PrepareAmbientSamples() {
	const SAMPLE_SPACING_X = 128.0
	const SAMPLE_SPACING_Y = 128.0
	const SAMPLE_SPACING_Z = 256.0

	ambientLightIndices := []leafambientindex.LeafAmbientIndex{}
	ambientLightSamples := []leafambientlighting.LeafAmbientLighting{}


	for _,leaf := range f.leafs {
		if 0 != (leaf.Contents & flags.CONTENTS_SOLID) {
			ambientLightIndices = append(ambientLightIndices,
				leafambientindex.LeafAmbientIndex{
					AmbientSampleCount: 0,
					FirstAmbientSample: 0,
				},
			)
			continue
		}

		leafSize := mgl32.Vec3{
			float32(leaf.Maxs[0] - leaf.Mins[0]),
			float32(leaf.Maxs[1] - leaf.Mins[1]),
			float32(leaf.Maxs[2] - leaf.Mins[2]),
		}

		numSamplesX := leafSize[0] / SAMPLE_SPACING_X
		numSamplesY := leafSize[1] / SAMPLE_SPACING_Y
		numSamplesZ := leafSize[2] / SAMPLE_SPACING_Z

		numSamples := numSamplesX * numSamplesY * numSamplesZ

		ambientLightIndices = append(ambientLightIndices,
			leafambientindex.LeafAmbientIndex{
				AmbientSampleCount: uint16(numSamples),
				FirstAmbientSample: uint16(len(ambientLightSamples)),
			})

		for i := 0; float32(i) < numSamplesZ; i++ {
			z := (float32(i) + 0.5) / numSamplesZ * 255

			for j := 0; float32(j) < numSamplesY; j++ {
				y := (float32(j) + 0.5) / numSamplesY * 255

				for k := 0; float32(k) < numSamplesX; k++ {
					x := (float32(k) + 0.5) / numSamplesX * 255

					ambientLightSamples = append(
						ambientLightSamples,
						leafambientlighting.LeafAmbientLighting{
							Cube: leafambientlighting.CompressedLightCube{
								Color: [6]common.ColorRGBExponent32{
									{0, 255, 0, 0},
									{0, 255, 0, 0},
									{0, 255, 0, 0},
									{0, 255, 0, 0},
									{0, 255, 0, 0},
									{0, 255, 0, 0},
								},
							},
							X: uint8(x),
							Y: uint8(y),
							Z: uint8(z),
							Pad: 0x0,
						},
					)
				}
			}
		}
	}

	if f.IsHDR == true {
		f.ambientLightIndicesHDR = ambientLightIndices
		f.ambientLightSamplesHDR = ambientLightSamples
	} else {
		f.ambientLightIndices = ambientLightIndices
		f.ambientLightSamples = ambientLightSamples
	}
}