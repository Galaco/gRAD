package radable

import (
	"github.com/galaco/bsp"
	"strings"
	"github.com/galaco/vmf"
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
	"github.com/galaco/source-tools-common/entity"
	"github.com/galaco/source-tools-common/texdatastringtable"
	"log"
	"github.com/galaco/gRAD/radable/light"
	"github.com/galaco/gRAD/radable/kd"
)

func GenerateRadiosityPrimitives(file *bsp.Bsp) (*RadPrimitive, error) {
	radBSP := RadPrimitive{}

	// Fetch unmodified primitives
	radBSP.Planes = (*file.GetLump(bsp.LUMP_PLANES).GetContents()).GetData().([]plane.Plane)
	radBSP.texData = *(*file.GetLump(bsp.LUMP_TEXDATA).GetContents()).GetData().(*[]texdata.TexData)
	radBSP.Vertexes = *(*file.GetLump(bsp.LUMP_VERTEXES).GetContents()).GetData().(*[]mgl32.Vec3)
	radBSP.visibility = *(*file.GetLump(bsp.LUMP_VISIBILITY).GetContents()).GetData().(*visibility.Vis)
	radBSP.texInfo = *(*file.GetLump(bsp.LUMP_TEXINFO).GetContents()).GetData().(*[]texinfo.TexInfo)
	radBSP.faces = *(*file.GetLump(bsp.LUMP_FACES).GetContents()).GetData().(*[]face.Face)
	radBSP.leafs = *(*file.GetLump(bsp.LUMP_LEAFS).GetContents()).GetData().(*[]leaf.Leaf)
	radBSP.leafFaces = *(*file.GetLump(bsp.LUMP_LEAFFACES).GetContents()).GetData().(*[]uint16)
	radBSP.leafBrushes = *(*file.GetLump(bsp.LUMP_LEAFBRUSHES).GetContents()).GetData().(*[]uint16)
	radBSP.Edges = *(*file.GetLump(bsp.LUMP_EDGES).GetContents()).GetData().(*[][2]uint16)
	radBSP.SurfEdges = *(*file.GetLump(bsp.LUMP_SURFEDGES).GetContents()).GetData().(*[]int32)
	radBSP.Models = *(*file.GetLump(bsp.LUMP_MODELS).GetContents()).GetData().(*[]model.Model)
	radBSP.brushes = *(*file.GetLump(bsp.LUMP_BRUSHES).GetContents()).GetData().(*[]brush.Brush)
	radBSP.brushSides = *(*file.GetLump(bsp.LUMP_BRUSHSIDES).GetContents()).GetData().(*[]brushside.BrushSide)
	radBSP.areas = *(*file.GetLump(bsp.LUMP_AREAS).GetContents()).GetData().(*[]area.Area)
	radBSP.areaPortals = *(*file.GetLump(bsp.LUMP_AREAPORTALS).GetContents()).GetData().(*[]areaportal.AreaPortal)
	radBSP.mapFlags = *(*file.GetLump(bsp.LUMP_MAP_FLAGS).GetContents()).GetData().(*mapflags.MapFlags)
	radBSP.facesHDR = *(*file.GetLump(bsp.LUMP_FACES_HDR).GetContents()).GetData().(*[]face.Face)

	//Entities
	var err error
	entData := *(*file.GetLump(bsp.LUMP_ENTITIES).GetContents()).GetData().(*string)
	radBSP.Entities,err = createEntityList(entData)
	if err != nil {
		return nil, err
	}

	// TexDataStringTable
	radBSP.TexDataStringTable = createStringTable(
		*(*file.GetLump(bsp.LUMP_TEXDATA_STRING_DATA).GetContents()).GetData().(*string),
		*(*file.GetLump(bsp.LUMP_TEXDATA_STRING_TABLE).GetContents()).GetData().(*[]int32))

	return &radBSP, nil
}

// createEntityList
// Creates a set of entities from entdata string lump
func createEntityList(entData string) (entity.List, error) {
	entDataReader := strings.NewReader(entData)
	vmfReader := vmf.NewReader(entDataReader)
	vmfEntities, err := vmfReader.Read()
	if err != nil {
		return entity.List{}, err
	}
	return entity.FromVmfNodeTree(vmfEntities.Unclassified), nil
}

// createStringTable
// create TexDataStringTable from TexData and TexInfo
func createStringTable(stringData string, tableData []int32) texdatastringtable.TexDataStringTable {
	return *texdatastringtable.NewTable(stringData, tableData)
}

// ExtractLights
func ExtractLights(primitives *RadPrimitive, useHDR bool) {
	log.Printf("Extracting lights from entdata...")
	var worldLights []light.DirectLight
	var ambientLight light.DirectLight

	var numLights = 0
	for i := 0; i < primitives.Entities.Length(); i++ {
		e := primitives.Entities.Get(i)
		classname := e.ValueForKey("classname")
		if strings.Contains(classname, "light") == false {
			continue
		}
		numLights++
		l := light.NewDirectLight(e)

		if classname == "light" {
			light.ParseLightGeneric(e, l, useHDR)
			worldLights = append(worldLights, *l)
			continue
		}
		if classname == "light_environment" {
			light.ParseLightGeneric(e, l, useHDR)
			light.ParseLightEnvironment(e, l, useHDR)
			ambientLight = *l
			continue
		}
		if classname == "light_spot" {
			light.ParseLightGeneric(e, l, useHDR)
			light.ParseLightSpot(e, l, &primitives.Entities)
			worldLights = append(worldLights, *l)
			continue
		}

	}

	log.Printf("Found %d Facelights\n\n", numLights)

	primitives.WorldLights = worldLights
	primitives.AmbientLight = ambientLight
}


// GenerateKDTree
// Build KD Tree to handling faces per leaf
func GenerateKDTree(primitives *RadPrimitive) {
	log.Printf("Building visibility tree...")
	kd.BuildTree()
}