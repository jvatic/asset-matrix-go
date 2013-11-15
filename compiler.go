package main

import (
	"flag"
	"fmt"
	"github.com/tent/asset-matrix-go/matrix"
)

func main() {
	fmt.Println("Compiling assets...")

	var outputDir string
	flag.StringVar(&outputDir, "output", "./output", "path to output directory")

	flag.Parse()

	inputPaths := flag.Args()

	inputManifest, _ := matrix.NewInputManifest(inputPaths, outputDir)

	inputManifest.ScanInputDirs()
}
