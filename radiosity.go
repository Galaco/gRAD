package main

import (
	"log"
	"github.com/urfave/cli"
	"github.com/galaco/gRAD/filesystem"
	"github.com/galaco/gRAD/radable"
)

var useHDR = false

// Rad
// Main rad function
func Rad(c *cli.Context) error {
	log.Printf("     Galaco's Radiosity Simulator     \n")
	log.Printf("See: https://github.com/galaco/gRAD\n\n")

	// Step 1: Load files
	filename := c.Args().Get(0)
	file,err := filesystem.ImportFromFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("     Compile target info:        \n")
	log.Printf("Target: %s\n", filename)
	log.Printf("BSP version: %d\n", file.GetHeader().Version)
	log.Printf("File revision: %d\n\n", file.GetHeader().Revision)

	primitives,err := radable.GenerateRadiosityPrimitives(file)
	if err != nil {
		log.Fatal(err)
	}

	radable.ExtractLights(primitives, useHDR)
	radable.GenerateKDTree(primitives)


	/*
	// Step 2: Prepare environment
	radPrimitive.PrepareAmbientSamples()

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
*/

	return nil
}