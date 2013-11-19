package matrix

import (
	"fmt"
	"os"
	"path/filepath"
)

type AssetDir struct {
	Path     string
	Name     string
	Parent   *AssetDir
	IsRoot   bool
	Assets   []*AssetFile
	Manifest *InputManifest
}

func NewAssetDir(path string, manifest *InputManifest, parent *AssetDir) (*AssetDir, error) {
	absPath, err := filepath.Abs(path)

	name := filepath.Base(absPath)
	if parent != nil && !parent.IsRoot {
		name = filepath.Join(parent.Name, name)
	}

	dir := &AssetDir{Path: absPath, Name: name, Parent: parent, IsRoot: parent == nil, Manifest: manifest}

	if !dir.IsRoot {
		manifest.DirPathMapping[dir.Path] = dir
		manifest.DirNameMapping[dir.Name] = dir
	}

	return dir, err
}

func (dir *AssetDir) scan() error {
	absPath, err := filepath.Abs(dir.Path)
	if err != nil {
		return err
	}

	fmt.Printf("Scan dir: %s...\n", absPath)
	return filepath.Walk(dir.Path, dir.visit)
}

func (dir *AssetDir) visit(path string, f os.FileInfo, err error) error {
	if path == dir.Path {
		return nil
	}

	if f.IsDir() {
		subDir, err := NewAssetDir(path, dir.Manifest, dir)
		if err != nil {
			return err
		}

		if err := subDir.scan(); err != nil {
			return err
		}

		return filepath.SkipDir
	}

	file, err := NewAssetFile(path, dir.Manifest, dir)
	if err != nil {
		return err
	}

	fmt.Printf("File: %s\n", file.Path)

	return nil
}
