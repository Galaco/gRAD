package opencl

import (
	"log"
	"github.com/samuel/go-opencl/cl"
	"unsafe"
	"github.com/galaco/gRAD/filesystem"
)

// sendRayTracerToGPU
// send RayTracer to Target Device
func sendRayTracerToDevice(context *cl.Context, queue *cl.CommandQueue, tracer *RayTracer) (*cl.MemObject,error) {
	log.Printf("Sending Raytracer to device target...\n")

	// Pass raytracer into device
	deviceRayTracer, err := context.CreateEmptyBuffer(cl.MemReadWrite, int(unsafe.Sizeof(tracer)))
	if err != nil {
		return nil,err
	}
	_,err = queue.EnqueueWriteBuffer(deviceRayTracer, true, 0, int(unsafe.Sizeof(tracer)), unsafe.Pointer(tracer), nil)
	if err != nil {
		return nil,err
	}

	return deviceRayTracer,nil
}

//sendBspToGPU
func sendBspToDevice(context *cl.Context, queue *cl.CommandQueue, vradBsp *filesystem.Bsp) (*DeviceBsp,error) {
	log.Printf("Sending BSP to device target...\n")

	// Whats gonna happen here?
	// Probably need to sync data lump by lump onto the device
	deviceBsp := DeviceBsp{
		NumModels: len(vradBsp.Models),
		NumPlanes: len(vradBsp.Planes),
		NumVertices: len(vradBsp.Vertexes),
		NumEdges: len(vradBsp.Edges),
		NumSurfEdges: len(vradBsp.SurfEdges),
		NumFaces: len(*vradBsp.GetFaces()),
	}

	modelsSize := int(unsafe.Sizeof(vradBsp.Models[0])) * deviceBsp.NumModels
	planesSize := int(unsafe.Sizeof(vradBsp.Planes[0])) * deviceBsp.NumPlanes
	verticesSize := int(unsafe.Sizeof(vradBsp.Vertexes[0])) * deviceBsp.NumVertices
	edgesSize := int(unsafe.Sizeof(vradBsp.Edges[0])) * deviceBsp.NumEdges
	surfEdgesSize := int(unsafe.Sizeof(vradBsp.SurfEdges[0])) * deviceBsp.NumSurfEdges
	facesSize := int(unsafe.Sizeof((*vradBsp.GetFaces())[0])) * deviceBsp.NumFaces

	// Copy data item-by-item to device
	var err error
	deviceBsp.Models, err = context.CreateEmptyBuffer(cl.MemReadWrite, modelsSize)
	if err != nil {
		return nil,err
	}
	_,err = queue.EnqueueWriteBuffer(deviceBsp.Models, true, 0, modelsSize, unsafe.Pointer(&vradBsp.Models[0]), nil)
	if err != nil {
		return nil,err
	}

	deviceBsp.Planes, err = context.CreateEmptyBuffer(cl.MemReadWrite, planesSize)
	if err != nil {
		return nil,err
	}
	_,err = queue.EnqueueWriteBuffer(deviceBsp.Planes, true, 0, planesSize, unsafe.Pointer(&vradBsp.Planes[0]), nil)
	if err != nil {
		return nil,err
	}

	deviceBsp.Vertices, err = context.CreateEmptyBuffer(cl.MemReadWrite, verticesSize)
	if err != nil {
		return nil,err
	}
	_,err = queue.EnqueueWriteBuffer(deviceBsp.Vertices, true, 0, verticesSize, unsafe.Pointer(&vradBsp.Vertexes[0]), nil)
	if err != nil {
		return nil,err
	}

	deviceBsp.Edges, err = context.CreateEmptyBuffer(cl.MemReadWrite, edgesSize)
	if err != nil {
		return nil,err
	}
	_,err = queue.EnqueueWriteBuffer(deviceBsp.Edges, true, 0, edgesSize, unsafe.Pointer(&vradBsp.Edges[0]), nil)
	if err != nil {
		return nil,err
	}

	deviceBsp.SurfEdges, err = context.CreateEmptyBuffer(cl.MemReadWrite, surfEdgesSize)
	if err != nil {
		return nil,err
	}
	_,err = queue.EnqueueWriteBuffer(deviceBsp.SurfEdges, true, 0, surfEdgesSize, unsafe.Pointer(&vradBsp.SurfEdges[0]), nil)
	if err != nil {
		return nil,err
	}

	deviceBsp.Faces, err = context.CreateEmptyBuffer(cl.MemReadWrite, facesSize)
	if err != nil {
		return nil,err
	}
	_,err = queue.EnqueueWriteBuffer(deviceBsp.Faces, true, 0, facesSize, unsafe.Pointer(&(*vradBsp.GetFaces())[0]), nil)
	if err != nil {
		return nil,err
	}

	return &deviceBsp,nil
}