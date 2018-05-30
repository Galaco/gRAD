package opencl

import (
	"github.com/samuel/go-opencl/cl"
	"log"
	"strings"
	"github.com/galaco/gRAD/radiosity/simulator/raytracer"
	"github.com/galaco/gRAD/bsp"
)

type Simulator struct {
	device *cl.Device
	context *cl.Context
	queue *cl.CommandQueue
}

// NewSimulator
// Create a new Radiosity simulator
// Also sends data to the gpu
func NewSimulator(tracer *raytracer.RayTracer, vradBsp *bsp.Bsp) (*Simulator,error) {
	platforms,err := cl.GetPlatforms()
	//device,err := blackcl.GetDefaultDevice()
	if err != nil {
		return nil, err
	}
	devices,err := cl.GetDevices(platforms[0], cl.DeviceTypeAll)
	if err != nil {
		return nil, err
	}
	context,err := cl.CreateContext(devices)
	if err != nil {
		return nil, err
	}

	// Create command queue for device
	queue,err := context.CreateCommandQueue(devices[0], 0)
	if err != nil {
		return nil,err
	}

	log.Printf("        Using OpenCL Simulator        \n")
	log.Printf("Using device: %s\n", strings.TrimRight(devices[0].Name(), "\x00"))

	err = sendRayTracerToGPU(context, queue, tracer)
	if err != nil {
		return nil,err
	}
	err = sendBspToGPU(context, queue, vradBsp)
	if err != nil {
		return nil,err
	}

	return &Simulator{
		device: devices[0],
		context: context,
		queue: queue,
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