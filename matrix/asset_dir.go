package matrix

import (
	"os"
	"path/filepath"
)

type AssetDir struct {
	Assets []*AssetFile

	AssetPointer

	path     string
	name     string
	parent   *AssetDir
	rootDir  *AssetDir
	isRoot   bool
	manifest *InputManifest
}

func NewAssetDir(path string, manifest *InputManifest, parent *AssetDir) (*AssetDir, error) {
	absPath, err := filepath.Abs(path)

	name := filepath.Base(absPath)
	if parent != nil && !parent.IsRoot() {
		name = filepath.Join(parent.Name(), name)
	}

	dir := &AssetDir{path: absPath, name: name, parent: parent, isRoot: parent == nil, manifest: manifest}

	if dir.IsRoot() {
		dir.rootDir = dir
	} else {
		dir.rootDir = parent.RootDir()

		manifest.AddDir(dir)
	}

	return dir, err
}

func (dir *AssetDir) Path() string {
	return dir.path
}

func (dir *AssetDir) Name() string {
	return dir.name
}

func (dir *AssetDir) Dir() *AssetDir {
	return dir.parent
}

func (dir *AssetDir) RootDir() *AssetDir {
	return dir.rootDir
}

func (dir *AssetDir) Manifest() *InputManifest {
	return dir.manifest
}

func (dir *AssetDir) IsRoot() bool {
	return dir.isRoot
}

func (dir *AssetDir) scan() error {
	return filepath.Walk(dir.Path(), dir.visit)
}

func (dir *AssetDir) visit(path string, f os.FileInfo, err error) error {
	if path == dir.Path() {
		return nil
	}

	if f.IsDir() {
		subDir, err := NewAssetDir(path, dir.Manifest(), dir)
		if err != nil {
			return err
		}

		if err := subDir.scan(); err != nil {
			return err
		}

		return filepath.SkipDir
	}

	_, err = NewAssetFile(path, dir)
	return err
}
