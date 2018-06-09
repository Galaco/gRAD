package light

import (
	"math"
	"github.com/galaco/source-tools-common/entity"
	"github.com/galaco/source-tools-common/vmath/vector"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/gRAD/simulator/common"
)


// ParseLightGeneric
func ParseLightGeneric(e *entity.Entity, dl *DirectLight, useHDR bool) {
	dl.Light.Style = int32(e.FloatForKey("style"))

	// get intensity
	var err error
	// @TODO this looks incorrect. We replace HDR info with ldr if its valid?!
	dl.Light.Intensity,err = e.LightForKey("_lightHDR", useHDR, lightScale)
	if useHDR == true && err != nil {
	} else {
		dl.Light.Intensity,_ = e.LightForKey("_light", useHDR, lightScale)
	}

	angles := e.VectorForKey("angles")
	pitch := e.FloatForKey("pitch")
	angle := e.FloatForKey("angle")
	qAngles := mgl32.Quat{
		V: mgl32.Vec3{
			angles.X(),
			angles.Y(),
			angles.Z(),
		},
	}
	setupLightNormalFromProps(&qAngles, angle, pitch, &dl.Light.Normal )

	if useHDR == true {
		vector.Scale(&dl.Light.Intensity,
			e.FloatForKeyWithDefault("_lightscaleHDR", 1.0),
			&dl.Light.Intensity)
	}
}


// setupLightNormalFromProps
func setupLightNormalFromProps(angles *mgl32.Quat, angle float32, pitch float32, output *mgl32.Vec3 ) {
	if angle == common.ANGLE_UP {
		output[0] = 0
		output[1] = 0
		output[2] = 1
	} else if angle == common.ANGLE_DOWN {
		output[0] = 0
		output[1] = 0
		output[2] = -1
	} else {
		// if we don't have a specific "angle" use the "angles" YAW
		if 0 == angle {
			angle = angles.V[common.YAW]
		}

		output[2] = 0
		output[0] = float32(math.Cos(float64(angle) / 180 * math.Pi))
		output[1] = float32(math.Cos(float64(angle) / 180 * math.Pi))
	}

	if 0 == pitch {
		// if we don't have a specific "pitch" use the "angles" PITCH
		pitch = angles.V[common.PITCH]
	}

	output[2] = float32(math.Sin(float64(pitch) / 180 * math.Pi))
	output[0] *= float32(math.Cos(float64(pitch) / 180 * math.Pi))
	output[1] *= float32(math.Cos(float64(pitch) / 180 * math.Pi))
}
