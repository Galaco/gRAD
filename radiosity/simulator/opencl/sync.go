package opencl

import (
	"log"
	"time"
	"github.com/samuel/go-opencl/cl"
	"github.com/galaco/gRAD/radiosity/raytracer"
)

func sendToGPU(context *cl.Context, tracer *raytracer.RayTracer) {
	log.Printf("Syncing data with device...\n")
	setupStart := time.Now().UnixNano() / int64(time.Millisecond)

	/*
		RayTracer::CUDARayTracer* pDeviceRayTracer;

        CUDA_CHECK_ERROR(
            cudaMalloc(&pDeviceRayTracer, sizeof(RayTracer::CUDARayTracer))
        );
        CUDA_CHECK_ERROR(
            cudaMemcpy(
                pDeviceRayTracer, g_pRayTracer.get(),
                sizeof(RayTracer::CUDARayTracer),
                cudaMemcpyHostToDevice
            )
        );
        CUDA_CHECK_ERROR(
            cudaMemcpyToSymbol(
                g_pDeviceRayTracer, &pDeviceRayTracer,
                sizeof(RayTracer::CUDARayTracer*), 0,
                cudaMemcpyHostToDevice
            )
        );
	 */


	setupEnd := time.Now().UnixNano() / int64(time.Millisecond)
	log.Printf("Done (%f seconds)\n\n", float32(setupEnd-setupStart) / 1000)
}
