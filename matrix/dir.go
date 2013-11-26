package matrix

import (
	"os"
	"path/filepath"
)

type Dir struct {
	Files []*File
	Dirs  []*Dir

	AssetPointer

	path     string
	name     string
	parent   *Dir
	rootDir  *Dir
	isRoot   bool
	manifest *Manifest
}

func NewDir(path string, manifest *Manifest, parent *Dir) (*Dir, error) {
	absPath, err := filepath.Abs(path)

	name := filepath.Base(absPath)
	if parent != nil && !parent.IsRoot() {
		name = filepath.Join(parent.Name(), name)
	}

	dir := &Dir{path: absPath, name: name, parent: parent, isRoot: parent == nil, manifest: manifest, Files: make([]*File, 0), Dirs: make([]*Dir, 0)}

	if dir.IsRoot() {
		dir.rootDir = dir
	} else {
		dir.rootDir = parent.RootDir()

		manifest.AddDir(dir)
	}

	return dir, err
}

func (dir *Dir) Path() string {
	return dir.path
}

func (dir *Dir) Name() string {
	return dir.name
}

func (dir *Dir) Dir() *Dir {
	return dir.parent
}

func (dir *Dir) RootDir() *Dir {
	return dir.rootDir
}

func (dir *Dir) Manifest() *Manifest {
	return dir.manifest
}

func (dir *Dir) IsRoot() bool {
	return dir.isRoot
}

func (dir *Dir) AddFile(file *File) {
	dir.Files = append(dir.Files, file)
}

func (dir *Dir) AddSubDir(subDir *Dir) {
	dir.Dirs = append(dir.Dirs, subDir)
}

func (dir *Dir) scan() error {
	return filepath.Walk(dir.Path(), dir.visit)
}

func (dir *Dir) visit(path string, f os.FileInfo, err error) error {
	if path == dir.Path() {
		return nil
	}

	if f.IsDir() {
		subDir, err := NewDir(path, dir.Manifest(), dir)
		if err != nil {
			return err
		}

		dir.AddSubDir(subDir)

		if err := subDir.scan(); err != nil {
			return err
		}

		return filepath.SkipDir
	}

	file, err := NewFile(path, dir)
	if err != nil {
		dir.AddFile(file)
	}
	return err
}
