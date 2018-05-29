package main

import (
	"log"
	"github.com/urfave/cli"
	"github.com/galaco/gRAD/bsp"
	"github.com/galaco/gRAD/radiosity"
	"github.com/galaco/gRAD/radiosity/simulator"
	"github.com/galaco/gRAD/radiosity/simulator/opencl"
)

var useHDR = false
var useGPU = true

// Rad
// Main rad function
func Rad(c *cli.Context) error {
	log.Printf("     Galaco's Radiosity Simulator     \n")
	log.Printf("See: https://github.com/galaco/RADiant\n\n")

	// Step 1: Load files
	filename := c.Args().Get(0)
	file,err := bsp.ImportFromFile(filename)

	if err != nil {
		return err
	}
	// Extract lights. Either hdr or ldr
	file.ExtractLights(useHDR)

	// Step 2: Prepare environment
	file.PrepareAmbientSamples(useHDR)
	tracer := radiosity.SetupAccelerationStructure(file, useHDR)

	var runner simulator.ISimulator
	if useGPU == true {
		runner,err = opencl.NewSimulator(tracer)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("GPU must be enabled...")
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