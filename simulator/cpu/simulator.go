package cpu

import (
	"log"
	"github.com/galaco/gRAD/filesystem"
	"time"
)

type Simulator struct {
	tracer *RayTracer
	bsp *filesystem.Bsp
}

// NewSimulator
// Create a new Radiosity simulator
// Also sends data to the gpu
func NewSimulator(tracer *RayTracer, vradBsp *filesystem.Bsp) (*Simulator,error) {
	log.Printf("        Using Default CPU Simulator        \n")

	return &Simulator{
		tracer: tracer,
		bsp: vradBsp,
	}, nil
}

// ComputeDirectLighting
func (rad Simulator) ComputeDirectLighting() {
	var facesCompleted = 0

	const BLOCK_WIDTH = 4
	const BLOCK_HEIGHT = 4
	//numFaces := len(*rad.bsp.GetFaces())

	//KERNEL_LAUNCH(
	//	DirectLighting::map_faces,
	//	numFaces, blockDim,
	//	pCudaBSP, const_cast<size_t*>(pDeviceFacesCompleted)
	//);
	setupStart := time.Now().UnixNano() / int64(time.Millisecond)
	for facesCompleted = 0; facesCompleted < len(*rad.bsp.GetFaces()); facesCompleted++ {
		rad.tracer.mapFaces(rad.bsp, &facesCompleted, [2]int{1, 1}, [2]int{1, 1})

		if facesCompleted % 100 == 0 {
			log.Printf("Processed %d faces\n", facesCompleted)
		}
	}

	//for bw := 1; bw < BLOCK_WIDTH; bw++ {
	//	for bh := 1; bh < BLOCK_HEIGHT; bh++ {
	//		for tw := 1; tw < BLOCK_DIMENSIONS[0]; tw++ {
	//			for th := 1; th < BLOCK_DIMENSIONS[1]; th++ {
	//				rad.tracer.mapFaces(rad.bsp, &facesCompleted, [2]int{bw, bh}, [2]int{tw, th})
	//				facesCompleted++
	//			}
	//		}
	//	}
	//}

	setupEnd := time.Now().UnixNano() / int64(time.Millisecond)
	log.Printf("Processed %d faces\n", facesCompleted)
	log.Printf("Done (%f seconds)\n\n", float32(setupEnd-setupStart) / 1000)
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