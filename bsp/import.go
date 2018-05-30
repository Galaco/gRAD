package bsp

import (
	"github.com/galaco/bsp"
	"os"
	"github.com/galaco/bsp/primitives/brush"
	"github.com/galaco/bsp/primitives/brushside"
	"github.com/galaco/bsp/primitives/plane"
	"github.com/galaco/bsp/primitives/texdata"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/bsp/primitives/visibility"
	"github.com/galaco/bsp/primitives/texinfo"
	"github.com/galaco/bsp/primitives/face"
	"github.com/galaco/bsp/primitives/leaf"
	"github.com/galaco/bsp/primitives/model"
	"github.com/galaco/bsp/primitives/area"
	"github.com/galaco/bsp/primitives/areaportal"
	"github.com/galaco/bsp/primitives/mapflags"
	"github.com/galaco/source-tools-common/entity"
	"github.com/galaco/vmf"
	"strings"
	"github.com/galaco/source-tools-common/texdatastringtable"
	"log"
)

// Load
// Import a BSP into a format containing everything needed for rad
func ImportFromFile(filename string) (*Bsp,error) {
	file,err := getRawBSP(filename)


	if err != nil {
		return nil, err
	}

	log.Printf("     Compile target info:        \n")
	log.Printf("Target: %s\n", filename)
	log.Printf("BSP version: %d\n", file.GetHeader().Version)
	log.Printf("File revision: %d\n\n", file.GetHeader().Revision)

	return transformRawBspToRadBsp(file)
}

// getRawBSP
// Read raw file to bsp package format.
func getRawBSP(filename string) (*bsp.Bsp, error) {
	// Read file
	handle,err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// Import file as bsp lib bsp
	reader := bsp.NewReader(handle)

	rawBSP,err := reader.Read()
	handle.Close()
	if err != nil {
		return nil, err
	}

	return rawBSP, nil
}

// transformRawBspToRadBsp
// Transform generic bsp to rad specific format.
// Essentially adds methods to structs, and only contains relevant lump data
func transformRawBspToRadBsp(file *bsp.Bsp) (*Bsp, error) {
	radBSP := Bsp{}
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
	entData := *(*file.GetLump(bsp.LUMP_ENTITIES).GetContents()).GetData().(*string)
	entDataReader := strings.NewReader(entData)
	vmfReader := vmf.NewReader(entDataReader)
	vmfEntities,err := vmfReader.Read()
	if err != nil {
		return nil, err
	}
	radBSP.entities = entity.FromVmfNodeTree(vmfEntities.Unclassified)

	// TexDataStringTable
	stringData := *(*file.GetLump(bsp.LUMP_TEXDATA_STRING_DATA).GetContents()).GetData().(*string)
	stringTable := *(*file.GetLump(bsp.LUMP_TEXDATA_STRING_TABLE).GetContents()).GetData().(*[]int32)
	radBSP.texDataStringTable = *texdatastringtable.NewTable(stringData, stringTable)


	return &radBSP, nil
}