package opencl

import (
	"log"
	"time"
	"github.com/samuel/go-opencl/cl"
	"github.com/galaco/gRAD/radiosity/simulator/raytracer"
	"unsafe"
	"github.com/galaco/gRAD/bsp"
)

// sendRayTracerToGPU
// send RayTracer to Target Device
func sendRayTracerToGPU(context *cl.Context, queue *cl.CommandQueue, tracer *raytracer.RayTracer) error {
	log.Printf("Initialising environment...\n")
	setupStart := time.Now().UnixNano() / int64(time.Millisecond)

	var deviceRayTracer *raytracer.RayTracer

	// Create buffer
	buffer, err := context.CreateBufferUnsafe(
		cl.MemReadWrite,
		int(unsafe.Sizeof(tracer)),
		unsafe.Pointer(deviceRayTracer))
	if err != nil {
		return err
	}

	// Pass raytracer into device
	log.Printf("Syncing data with device...\n")
	_,err = queue.EnqueueWriteBuffer(buffer, true, 0, int(unsafe.Sizeof(tracer)), unsafe.Pointer(tracer), nil)
	if err != nil {
		return err
	}

	setupEnd := time.Now().UnixNano() / int64(time.Millisecond)
	log.Printf("Done (%f seconds)\n\n", float32(setupEnd-setupStart) / 1000)

	return nil
}

//sendBspToGPU
func sendBspToGPU(context *cl.Context, queue *cl.CommandQueue, vradBsp *bsp.Bsp) error {
	log.Printf("Sending BSP to device target...\n")
	setupStart := time.Now().UnixNano() / int64(time.Millisecond)

	// Whats gonna happen here?
	// Probably need to sync data lump by lump onto the device


	setupEnd := time.Now().UnixNano() / int64(time.Millisecond)
	log.Printf("Done (%f seconds)\n\n", float32(setupEnd-setupStart) / 1000)
	return nil
}