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

	inputManifest := matrix.NewManifest(inputPaths, *outputDir)
	if err := inputManifest.ScanInputDirs(); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	if err := inputManifest.EvaluateDirectives(); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}
