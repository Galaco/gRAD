package cpu

import (
	radLight "github.com/galaco/gRAD/filesystem/light"
	"github.com/galaco/gRAD/filesystem"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/bsp/primitives/worldlight"
)

func attenuate(light *radLight.DirectLight, dist float32) float32 {
	c := light.Light.ConstantAttenuation
	l := light.Light.LinearAttenuation
	q := light.Light.QuadraticAttenuation

	return c + l * dist + q * dist * dist
}

func (tracer *RayTracer) mapFaces(vradBsp *filesystem.Bsp, facesCompleted *int, blockId [2]int, threadId [2]int) {
	primaryThread := true//(threadId[0] == 0 && threadId[1] == 0)

	//if (pCudaBSP->tag != CUDABSP::TAG) {
	//	if (primaryThread) {
	//		printf("Invalid CUDABSP Tag: %x\n", pCudaBSP->tag);
	//	}
	//	return;
	//}

	faceInfo := &FaceInfo{}

	if primaryThread {
		// Map block numbers to faces.
		faceInfo = NewFaceInfo(vradBsp, 0)

		//printf(
		//    "Processing Face %u...\n",
		//    static_cast<unsigned int>(faceInfo.faceIndex)
		//);
	}

	//__syncthreads()

	/* Take a sample at each lightmap luxel. */
	for i := 0; i < faceInfo.LightmapHeight; i += blockId[1] {
		t := i + threadId[1]

		if t >= faceInfo.LightmapHeight {
			continue
		}

		for j := 0; j < faceInfo.LightmapWidth; j += blockId[0] {
			s := j + threadId[0]

			if s >= faceInfo.LightmapWidth {
				continue
			}

			 colour := tracer.sampleAtFaceInfo(vradBsp, faceInfo, float32(s), float32(t))

			lightmapStart := &faceInfo.LightmapStartIndex
			sampleIndex := t * faceInfo.LightmapWidth + s

			vradBsp.LightSamples[*lightmapStart + sampleIndex] = colour

			faceInfo.TotalLight[0] = colour.X()
			faceInfo.TotalLight[1] = colour.Y()
			faceInfo.TotalLight[2] = colour.Z()
			//atomicAdd(&faceInfo.TotalLight.X(), color.X())
			//atomicAdd(&faceInfo.TotalLight.Y(), color.Y())
			//atomicAdd(&faceInfo.TotalLight.Z(), color.Z())
		}
	}

	//__syncthreads()

	if primaryThread {
		faceInfo.AvgLight = faceInfo.TotalLight
		faceInfo.AvgLight = faceInfo.AvgLight.Mul(float32(1 / faceInfo.LightmapSize))

		vradBsp.LightSamples[faceInfo.LightmapStartIndex - 1] = faceInfo.AvgLight

		// Still have no idea how this works. But if we don't do this,
		// EVERYTHING becomes a disaster...
		faceInfo.Face.Styles[0] = 0x00
		faceInfo.Face.Styles[1] = 0xFF
		faceInfo.Face.Styles[2] = 0xFF
		faceInfo.Face.Styles[3] = 0xFF

		/* Copy our changes back to the CUDABSP. */
		(*vradBsp.GetFaces())[faceInfo.FaceIndex] = faceInfo.Face

		//*facesCompleted += 1
		//atomicAdd(reinterpret_cast<unsigned int*>(facesCompleted), 1)
		//__threadfence_system()
	}

	//printf(
	//    "Lightmap offset for face %u: %u\n",
	//    static_cast<unsigned int>(faceIndex),
	//    static_cast<unsigned int>(lightmapStartIndex)
	//);

	//printf("%u\n", static_cast<unsigned int>(*pFacesCompleted));
}

func (tracer *RayTracer) sampleAt(vradBsp *filesystem.Bsp, samplePos mgl32.Vec3, sampleNormal mgl32.Vec3) mgl32.Vec3 {
	result := mgl32.Vec3{0, 0, 0}

	for lightIndex := 0; lightIndex < len(*vradBsp.GetDirectLights()); lightIndex++ {
		light := (*vradBsp.GetDirectLights())[lightIndex]
		lightPos := light.Light.Origin

		diff := samplePos.Sub(lightPos)

		/*
		 * This light is on the wrong side of the current sample.
		 * There's no way it could possibly light it.
		 */
		if sampleNormal.Len() > 0.0 && diff.Dot(sampleNormal) >= 0.0 {
			continue
		}

		dist := diff.Len()
		dir := diff.Mul(1 / dist)

		penumbraScale := float32(1.0)

		if light.Light.Type == worldlight.EMIT_SPOTLIGHT {
			lightNorm := light.Light.Normal

			lightDot := dir.Dot(lightNorm)

			if lightDot < light.Light.Stopdot2 {
				/* This sample is outside the spotlight cone. */
				continue
			} else if lightDot < light.Light.Stopdot {
				/* This sample is within the spotlight's penumbra. */
				penumbraScale = (lightDot - light.Light.Stopdot2) / (light.Light.Stopdot - light.Light.Stopdot2);
				//penumbraScale = 100.0;
			}

			//if (lightIndex == cudaBSP.numWorldLights - 1) {
			//    printf(
			//        "(%f, %f, %f) is within spotlight!\n"
			//        "Pos: (%f, %f, %f)\n"
			//        "Norm: <%f, %f, %f> (<%f, %f, %f>)\n"
			//        "stopdot: %f; stopdot2: %f\n"
			//        "Dot between light and sample: %f\n",
			//        samplePos.x, samplePos.y, samplePos.z,
			//        lightPos.x, lightPos.y, lightPos.z,
			//        lightNorm.x, lightNorm.y, lightNorm.z,
			//        light.normal.x, light.normal.y, light.normal.z,
			//        light.stopdot, light.stopdot2,
			//        lightDot
			//    );
			//}
		}

		EPSILON := float32(1e-3)

		// Nudge the sample position towards the light slightly, to avoid
		// colliding with triangles that directly contain the sample
		// position.
		samplePos = samplePos.Sub(dir.Mul(EPSILON))

		lightBlocked := tracer.LOS_blocked(lightPos, samplePos)

		if lightBlocked {
			// This light can't be seen from the position of the sample.
			// Ignore it.
			continue
		}

		/* I CAN SEE THE LIGHT */
		attenuation := attenuate(&light, dist)

		lightContribution := light.Light.Intensity
		lightContribution = lightContribution.Mul(penumbraScale * 255.0 / attenuation)

		result = result.Add(lightContribution)
	}

	//printf(
	//    "Sample at (%u, %u) for Face %u: (%f, %f, %f)\n",
	//    static_cast<unsigned int>(s),
	//    static_cast<unsigned int>(t),
	//    static_cast<unsigned int>(faceIndex),
	//    result.x, result.y, result.z
	//);

	return result
}

func (tracer *RayTracer) sampleAtFaceInfo(vradBsp *filesystem.Bsp, faceInfo *FaceInfo, s float32, t float32) mgl32.Vec3 {
	samplePos := faceInfo.XYXFromST(s, t)
	return tracer.sampleAt(vradBsp, samplePos, faceInfo.FaceNorm)
}
//
//__global__ void map_faces(
//CUDABSP::CUDABSP* pCudaBSP,
//size_t* pFacesCompleted
//) {
//
//bool primaryThread = (threadIdx.x == 0 && threadIdx.y == 0);
//
//if (pCudaBSP->tag != CUDABSP::TAG) {
//if (primaryThread) {
//printf("Invalid CUDABSP Tag: %x\n", pCudaBSP->tag);
//}
//return;
//}
//
//__shared__ CUDARAD::FaceInfo faceInfo;
//
//if (primaryThread) {
//// Map block numbers to faces.
//faceInfo = CUDARAD::FaceInfo(*pCudaBSP, blockIdx.x);
//
////printf(
////    "Processing Face %u...\n",
////    static_cast<unsigned int>(faceInfo.faceIndex)
////);
//}
//
//__syncthreads();
//
///* Take a sample at each lightmap luxel. */
//for (size_t i=0; i<faceInfo.lightmapHeight; i+=blockDim.y) {
//size_t t = i + threadIdx.y;
//
//if (t >= faceInfo.lightmapHeight) {
//continue;
//}
//
//for (size_t j=0; j<faceInfo.lightmapWidth; j+=blockDim.x) {
//size_t s = j + threadIdx.x;
//
//if (s >= faceInfo.lightmapWidth) {
//continue;
//}
//
//float3 color = sample_at(
//*pCudaBSP, faceInfo,
//static_cast<float>(s),
//static_cast<float>(t)
//);
//
//size_t& lightmapStart = faceInfo.lightmapStartIndex;
//size_t sampleIndex = t * faceInfo.lightmapWidth + s;
//
//pCudaBSP->lightSamples[lightmapStart + sampleIndex] = color;
//
//atomicAdd(&faceInfo.totalLight.x, color.x);
//atomicAdd(&faceInfo.totalLight.y, color.y);
//atomicAdd(&faceInfo.totalLight.z, color.z);
//}
//}
//
//__syncthreads();
//
//if (primaryThread) {
//faceInfo.avgLight = faceInfo.totalLight;
//faceInfo.avgLight /= static_cast<float>(faceInfo.lightmapSize);
//
//pCudaBSP->lightSamples[faceInfo.lightmapStartIndex - 1]
//= faceInfo.avgLight;
//
//// Still have no idea how this works. But if we don't do this,
//// EVERYTHING becomes a disaster...
//faceInfo.face.styles[0] = 0x00;
//faceInfo.face.styles[1] = 0xFF;
//faceInfo.face.styles[2] = 0xFF;
//faceInfo.face.styles[3] = 0xFF;
//
///* Copy our changes back to the CUDABSP. */
//pCudaBSP->faces[faceInfo.faceIndex] = faceInfo.face;
//
//atomicAdd(reinterpret_cast<unsigned int*>(pFacesCompleted), 1);
//__threadfence_system();
//}
//
////printf(
////    "Lightmap offset for face %u: %u\n",
////    static_cast<unsigned int>(faceIndex),
////    static_cast<unsigned int>(lightmapStartIndex)
////);
//
////printf("%u\n", static_cast<unsigned int>(*pFacesCompleted));
//}
//}
