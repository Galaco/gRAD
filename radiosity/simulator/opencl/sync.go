package opencl

import (
	"log"
	"time"
	"github.com/samuel/go-opencl/cl"
	"github.com/galaco/gRAD/radiosity/raytracer"
	"unsafe"
)

// send RayTracer to Target Device
func sendToGPU(device *cl.Device, context *cl.Context, tracer *raytracer.RayTracer) (*cl.CommandQueue,error) {
	log.Printf("Initialising environment...\n")
	setupStart := time.Now().UnixNano() / int64(time.Millisecond)

	var deviceRayTracer *raytracer.RayTracer

	// Create buffer
	buffer, err := context.CreateBufferUnsafe(
		cl.MemReadWrite,
		int(unsafe.Sizeof(tracer)),
		unsafe.Pointer(deviceRayTracer))
	if err != nil {
		return nil,err
	}

	// Create command queue for device
	queue,err := context.CreateCommandQueue(device, 0)
	if err != nil {
		return nil,err
	}

	// Pass raytracer into device
	log.Printf("Syncing data with device...\n")
	_,err = queue.EnqueueWriteBuffer(buffer, true, 0, int(unsafe.Sizeof(tracer)), unsafe.Pointer(tracer), nil)
	if err != nil {
		return nil,err
	}

	setupEnd := time.Now().UnixNano() / int64(time.Millisecond)
	log.Printf("Done (%f seconds)\n\n", float32(setupEnd-setupStart) / 1000)

	return queue,nil
}
