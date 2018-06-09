package light

import (
	"github.com/galaco/source-tools-common/entity"
)

// ParseLightEnvironment
func ParseLightEnvironment(e *entity.Entity, dl *DirectLight, useHDR bool) {
/*
	angleStr := e.ValueForKeyWithDefault("SunSpreadAngle", "")
	if angleStr != "" {
		sunAngularExtent,_ = strconv.ParseFloat(angleStr, 32)
		sunAngularExtent = math.Sin((vmath.PI/180.0) * sunAngularExtent)
		log.Printf("sun extent from map=%f\n", sunAngularExtent)
	}
	if nil == globalSkyLight {
		// Sky light.
		globalSkyLight = dl
		dl.Light.Type = worldlight.EMIT_SKYLIGHT

		// Sky ambient light.
		ambient := AllocDLight(&dl.Light.Origin, false)
		ambient.Light.Type = worldlight.EMIT_SKYAMBIENT;
		var err error
		ambient.Light.Intensity,err = e.LightForKey("_ambientHDR", useHDR, lightScale)
		if useHDR && err == nil {
			// we have a valid HDR ambient light value
		} else {
			ambient.Light.Intensity,err = e.LightForKey("_ambient", useHDR, lightScale)
			if err == nil {
				vector.Scale(&dl.Light.Intensity, 0.5, &ambient.Light.Intensity)
			}
		}
		if useHDR == true {
			vector.Scale(&ambient.Light.Intensity,
				e.FloatForKeyWithDefault("_AmbientScaleHDR", 1.0),
				&ambient.Light.Intensity)
		}

		BuildVisForLightEnvironment()

		// Add sky and sky ambient lights to the list.
		AddDLightToActiveList(globalSkyLight)
		AddDLightToActiveList(ambient)
	}
*/
}