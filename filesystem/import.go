package filesystem

import (
	"github.com/galaco/bsp"
	"os"
)

// Load
// Import a BSP into a format containing everything needed for rad
func ImportFromFile(filename string) (*bsp.Bsp,error) {
	return getRawBSP(filename)
}

// getRawBSP
// Read raw file to bsp package format.
func getRawBSP(filename string) (*bsp.Bsp, error) {
	// Read file
	handle,err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// Import file as bsp lib bsp
	reader := bsp.NewReader(handle)

	rawBSP,err := reader.Read()
	handle.Close()
	if err != nil {
		return nil, err
	}

	return rawBSP, nil
}