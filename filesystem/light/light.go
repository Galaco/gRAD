package light

import (
	"github.com/galaco/bsp/primitives/worldlight"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/source-tools-common/entity"
)

var lightScale = float32(1.0)
//IncrementalLightID = uint16

type DirectLight struct {
	Light worldlight.WorldLight

	PVS []byte		// accumulated domain of the light
	FaceNum int		// domain of attached lights
	TexData int		// texture source of traced lights

	SNormal mgl32.Vec3
	TNormal mgl32.Vec3
	SScale float32
	TScale float32
	SOffset float32
	TOffset float32

	DoRecalc int // position, vector, spot angle, etc.
	IncrementalID uint16

	// hard-falloff lights (lights that fade to an actual zero). between m_flStartFadeDistance and
	// m_flEndFadeDistance, a smoothstep to zero will be done, so that the light goes to zero at
	// the end.
	StartFadeDistance float32
	EndFadeDistance float32
	CapDistance float32										// max distance to feed in
}

// NewDirectLight
func NewDirectLight(entity *entity.Entity) *DirectLight {
	dl := DirectLight{
		EndFadeDistance: -1.0,
		StartFadeDistance: 0.0,
		CapDistance: 1.0e22,
	}

	dl.Light.Origin = entity.VectorForKey("origin")

	//dl.Light.Cluster = int32(clustertable.ClusterFromPoint(&dl.Light.Origin))
	//SetDLightVis( dl, int(dl.Light.Cluster) )

	dl.FaceNum = -1

	return &dl
}
