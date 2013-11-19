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
	absPath, err := filepath.Abs(path)
	return &AssetDir{Path: absPath}, err
}

func (dir *AssetDir) scan() error {
	absPath, err := filepath.Abs(dir.Path)

	if err != nil {
		return err
	}

	fmt.Printf("Scan dir: %s...\n", absPath)
	err = filepath.Walk(dir.Path, dir.visit)

	return err
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

	file, err := NewAssetFile(path)

	if err != nil {
		return err
	}

	fmt.Printf("File: %s\n", file.Path)

	err = file.ParseDirectives()

	return err
}
