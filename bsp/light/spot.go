package light

import (
	"github.com/galaco/source-tools-common/entity"
	"log"
)

// ParseLightSpot
func ParseLightSpot(e *entity.Entity, dl *DirectLight, list *entity.List) {
	target := e.ValueForKey("target")
	if target != "" {    // point towards target
		e2 := list.FindByKeyValue("targetname", target)
		if e2 == nil {
			log.Printf("WARNING: light at (%d %d %d) has missing target\n",
				int(dl.Light.Origin[0]), int(dl.Light.Origin[1]), int(dl.Light.Origin[2]))
		} else {
			dest := e2.VectorForKey("origin")
			dl.Light.Normal = dest.Sub(dl.Light.Origin)
			dl.Light.Normal = dl.Light.Normal.Normalize()
		}
	}
}

