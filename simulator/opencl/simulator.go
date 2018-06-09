package opencl

import (
	"github.com/samuel/go-opencl/cl"
	"log"
	"strings"
	"github.com/galaco/gRAD/filesystem"
	"time"
)

type Simulator struct {
	device *cl.Device
	context *cl.Context
	queue *cl.CommandQueue
}

// NewSimulator
// Create a new Radiosity simulator
// Also sends data to the gpu
func NewSimulator(tracer *RayTracer, vradBsp *filesystem.Bsp) (*Simulator,error) {
	platforms,err := cl.GetPlatforms()
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
	log.Printf("Using device: %s\n\n", strings.TrimRight(devices[0].Name(), "\x00"))

	log.Printf("Start device preparation\n")
	setupStart := time.Now().UnixNano() / int64(time.Millisecond)

	_,err = sendBspToDevice(context, queue, vradBsp)
	if err != nil {
		return nil,err
	}
	_,err = sendRayTracerToDevice(context, queue, tracer)
	if err != nil {
		return nil,err
	}

	setupEnd := time.Now().UnixNano() / int64(time.Millisecond)
	log.Printf("Done (%f seconds)\n\n", float32(setupEnd-setupStart) / 1000)

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