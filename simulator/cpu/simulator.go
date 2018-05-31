package cpu

import (
	"log"
	"github.com/galaco/gRAD/bsp"
)

type Simulator struct {
}

// NewSimulator
// Create a new Radiosity simulator
// Also sends data to the gpu
func NewSimulator(tracer *RayTracer, vradBsp *bsp.Bsp) (*Simulator,error) {
	log.Printf("        Using Default CPU Simulator        \n")

	return &Simulator{
	}, nil
}

// ComputeDirectLighting
func (rad Simulator) ComputeDirectLighting() {

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