package simulator

// Generic Simulator Interface
// Maybe we want various different calculators in the future
type ISimulator interface {
	ComputeDirectLighting()
	AntialiasLightmap(int)
	AntialiasDirectLighting()
	BounceLighting()
	ComputeAmbientLighting()
	ConvertLightSamples()
}
