package cpu

import (
	"log"
	"github.com/galaco/gRAD/bsp"
)

type Simulator struct {
	tracer *RayTracer
	bsp *bsp.Bsp
}

// NewSimulator
// Create a new Radiosity simulator
// Also sends data to the gpu
func NewSimulator(tracer *RayTracer, vradBsp *bsp.Bsp) (*Simulator,error) {
	log.Printf("        Using Default CPU Simulator        \n")

	return &Simulator{
		tracer: tracer,
		bsp: vradBsp,
	}, nil
}

// ComputeDirectLighting
func (rad Simulator) ComputeDirectLighting() {
	var facesCompleted = 0

	const BLOCK_WIDTH = 16
	const BLOCK_HEIGHT = 16
	BLOCK_DIMENSIONS := [2]int{3, 3}
	//numFaces := len(*rad.bsp.GetFaces())

	//KERNEL_LAUNCH(
	//	DirectLighting::map_faces,
	//	numFaces, blockDim,
	//	pCudaBSP, const_cast<size_t*>(pDeviceFacesCompleted)
	//);

	for bw := 0; bw < BLOCK_WIDTH; bw++ {
		for bh := 0; bh < BLOCK_HEIGHT; bh++ {
			for tw := 0; tw < BLOCK_DIMENSIONS[0]; tw++ {
				for th := 0; th < BLOCK_DIMENSIONS[1]; th++ {
					rad.tracer.mapFaces(rad.bsp, &facesCompleted, [2]int{bw, bh}, [2]int{tw, th})
				}
			}
		}
	}

	//rad.tracer.mapFaces(rad.bsp, &facesCompleted)
}

// AntialiasLightmap
func (rad Simulator) AntialiasLightmap(numPasses int) {

}

// AntialiasDirectLighting
func (rad Simulator) AntialiasDirectLighting() {

}

// BounceLighting
func (rad Simulator) BounceLighting() {

}

// ComputeAmbientLighting
func (rad Simulator) ComputeAmbientLighting() {

}

// ConvertLightSamples
func (rad Simulator) ConvertLightSamples() {

}