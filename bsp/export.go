package bsp

import (
	"github.com/galaco/bsp"
	"os"
)

// ExportToFile
// Exports the Radiosity Bsp structure back out to a file
func ExportToFile(filename string, vradBsp *Bsp) error {
	// Load bsp back into memory
	file,err := os.Open(filename)
	if err != nil {
		return err
	}
	// Parse bsp into format
	reader := bsp.NewReader(file)
	baseFile,err := reader.Read()
	if err != nil {
		return err
	}
	file.Close()

	// Inject modified lumps into target bsp
	baseFile,err = updateLumps(baseFile, vradBsp)
	if err != nil {
		return err
	}

	// Write out new bsp
	writer := bsp.NewWriter()
	writer.SetBsp(*baseFile)
	writer.Write()

	return nil
}

// updateLumps
// Insert modified lumps into bsp
func updateLumps(base *bsp.Bsp, target *Bsp) (*bsp.Bsp,error) {


	return base,nil
}
