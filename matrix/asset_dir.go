package matrix

import (
	"fmt"
	"os"
	"path/filepath"
)

type AssetDir struct {
	Path     string
	Assets   []*AssetFile
	manifest InputManifest
}

func NewAssetDir(path string) (*AssetDir, error) {
	absPath, _ := filepath.Abs(path)
	return &AssetDir{Path: absPath}, nil
}

func (dir *AssetDir) scan() error {
	absPath, _ := filepath.Abs(dir.Path)
	fmt.Printf("Scan dir: %s...\n", absPath)
	filepath.Walk(dir.Path, dir.visit)

	return nil
}

func (dir *AssetDir) visit(path string, f os.FileInfo, err error) error {
	if path == dir.Path {
		return nil
	}

	if f.IsDir() {
		subDir, _ := NewAssetDir(path)
		subDir.scan()

		return filepath.SkipDir
	}

	file, _ := NewAssetFile(path)

	fmt.Printf("File: %s\n", file.Path)

	file.ParseDirectives()

	return nil
}
