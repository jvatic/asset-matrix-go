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
	cCmdLimit := flag.Int("ccmd", 10, "max concurrent exec commands")
	flag.Parse()
	inputPaths := flag.Args()

	inputManifest := matrix.NewManifest(inputPaths, *outputDir, os.Stdout)
	inputManifest.SetCCmdLimit(*cCmdLimit)
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
