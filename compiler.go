package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jvatic/asset-matrix-go/matrix"
)

func main() {
	fmt.Println("Compiling assets...")

	outputDir := flag.String("output", "./output", "path to output directory")
	fdLimit := flag.Int("fdlimit", 10, "max open file descriptors allowed")
	flag.Parse()
	inputPaths := flag.Args()

	matrix.SetFDLimit(*fdLimit)

	inputManifest := matrix.NewManifest(inputPaths, *outputDir, os.Stdout)
	if err := inputManifest.ScanInputDirs(); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	if err := inputManifest.EvaluateDirectives(); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	if err := inputManifest.ConfigureHandlers(); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	if err := inputManifest.WriteOutput(); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}
