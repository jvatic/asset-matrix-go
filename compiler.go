package main

import (
	"flag"
	"fmt"

	"github.com/tent/asset-matrix-go/matrix"
)

func main() {
	fmt.Println("Compiling assets...")

	outputDir := flag.String("output", "./output", "path to output directory")
	flag.Parse()
	inputPaths := flag.Args()

	inputManifest := matrix.NewInputManifest(inputPaths, *outputDir)
	err := inputManifest.ScanInputDirs()
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
}
