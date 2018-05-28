package main

import (
	"github.com/urfave/cli"
	"github.com/galaco/goRAD/bsp"
	"log"
)

var useHDR = false

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

	return nil
}