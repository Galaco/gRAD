package main

import (
	"log"
	"github.com/urfave/cli"
	"github.com/galaco/gRAD/bsp"
	"github.com/galaco/gRAD/simulator"
	"github.com/galaco/gRAD/simulator/opencl"
	"github.com/galaco/gRAD/simulator/cpu"
)

var useHDR = false

// Rad
// Main rad function
func Rad(c *cli.Context) error {
	log.Printf("     Galaco's Radiosity Simulator     \n")
	log.Printf("See: https://github.com/galaco/gRAD\n\n")

	// Step 1: Load files
	filename := c.Args().Get(0)
	file,err := bsp.ImportFromFile(filename)
	file.IsHDR = useHDR

	if err != nil {
		return err
	}
	// Extract lights. Either hdr or ldr
	file.ExtractLights()

	// Step 2: Prepare environment
	file.PrepareAmbientSamples()

	var runner simulator.ISimulator

	// Determines which simulator to use
	switch c.Args().Get(1) {
	case "opencl":
		tracer := opencl.NewRayTracer()
		tracer.SetupAccelerationStructure(file)
		runner,err = opencl.NewSimulator(tracer, file)
	default:
		tracer := cpu.NewRayTracer()
		tracer.SetupAccelerationStructure(file)
		runner,err = cpu.NewSimulator(tracer, file)
	}

	if err != nil {
		log.Fatal(err)
	}

	// Step 3: Run
	runner.ComputeDirectLighting()
//	runner.AntialiasLightmap(5)
//	runner.AntialiasDirectLighting()
//	runner.BounceLighting()
//	runner.ComputeAmbientLighting()
//	runner.ConvertLightSamples()

	// Step 4: Export
	// if GPU return data to host
	// Write bsp back to file
//	bsp.ExportToFile(filename, file)


	return nil
}