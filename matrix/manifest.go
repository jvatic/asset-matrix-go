package matrix

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type AssetMap struct {
	Dir   *Dir
	Files map[string]*File
}

type Manifest struct {
	AssetRoots      []*Dir
	InputDirs       []string
	OutputDir       string
	DirPathMapping  map[string]*Dir
	FilePathMapping map[string]*File
	NameMapping     map[string]*AssetMap
	log             *log.Logger
}

func NewManifest(inputDirs []string, outputDir string, logOut io.Writer) *Manifest {
	manifest := &Manifest{InputDirs: inputDirs, OutputDir: outputDir, DirPathMapping: make(map[string]*Dir), FilePathMapping: make(map[string]*File), NameMapping: make(map[string]*AssetMap)}
	manifest.log = log.New(logOut, "matrix: ", 0)
	return manifest
}

func (manifest *Manifest) AddDir(dir *Dir) {
	manifest.DirPathMapping[dir.Path()] = dir

	if manifest.NameMapping[dir.Name()] == nil {
		manifest.NameMapping[dir.Name()] = &AssetMap{Dir: dir}
	} else {
		manifest.NameMapping[dir.Name()].Dir = dir
	}
}

func (manifest *Manifest) AddFile(file *File) {
	manifest.FilePathMapping[file.Path()] = file

	if manifest.NameMapping[file.Name()] == nil {
		manifest.NameMapping[file.Name()] = &AssetMap{}
	}
	if manifest.NameMapping[file.Name()].Files == nil {
		manifest.NameMapping[file.Name()].Files = make(map[string]*File)
	}
	manifest.NameMapping[file.Name()].Files[file.Ext()] = file
}

func (manifest *Manifest) FindDirName(name string) *Dir {
	assetMap := manifest.NameMapping[name]
	if assetMap == nil {
		return nil
	}

	return assetMap.Dir
}

func (manifest *Manifest) FindFileName(name string, ext string) *File {
	assetMap := manifest.NameMapping[name]
	if assetMap == nil {
		return nil
	}

	files := assetMap.Files
	if files == nil {
		return nil
	}

	if ext != "" {
		return files[ext]
	} else {
		for _, file := range files {
			return file
		}
		return nil
	}
}

func (manifest *Manifest) ScanInputDirs() error {
	manifest.AssetRoots = make([]*Dir, len(manifest.InputDirs))
	for i, path := range manifest.InputDirs {
		dir, err := NewDir(path, manifest, nil)
		if err != nil {
			return err
		}

		if err := dir.Scan(); err != nil {
			return err
		}

		manifest.AssetRoots[i] = dir
	}

	return nil
}

func (manifest *Manifest) EvaluateDirectives() error {
	for _, assetMap := range manifest.NameMapping {
		if assetMap == nil {
			continue
		}
		if assetMap.Files == nil {
			continue
		}

		for _, file := range assetMap.Files {
			if err := file.EvaluateDirectives(); err != nil {
				return err
			}

			file.LinkChildren()
		}
	}

	return nil
}

func (manifest *Manifest) ConfigureHandlers() error {
	for _, f := range manifest.FilePathMapping {
		f.DedupParents()
	}
	for _, f := range manifest.FilePathMapping {
		f.BuildHandlerChain(f.Ext())
	}
	for _, f := range manifest.FilePathMapping {
		if err := f.InsertConcatenationHandlers(); err != nil {
			return err
		}
	}

	return nil
}

func (manifest *Manifest) outFilePath(name string, exts []string) (string, error) {
	path, err := filepath.Abs(filepath.Join(manifest.OutputDir, name))
	if err != nil {
		return "", err
	}
	parts := append([]string{path}, exts...)
	return strings.Join(parts, "."), nil
}

func (manifest *Manifest) WriteOutput() (err error) {
	// Loop through fileHandlers in reverse order (least to most ParentHandlers)
	done := make(chan struct{}, len(manifest.FilePathMapping))

	writeFileOutput := func(file *File) {
		go func() {
			defer func() { done <- struct{}{} }()
			err = manifest.writeFileOutput(file)
		}()
	}

	for _, f := range manifest.FilePathMapping {
		writeFileOutput(f)
	}

	for i := 0; i < len(manifest.FilePathMapping); i++ {
		<-done
	}

	return
}

func (manifest *Manifest) writeFileOutput(file *File) (err error) {
	out := new(bytes.Buffer)
	var (
		name    string
		exts    []string
		outPath string
		outFile *os.File
	)

	manifest.log.Printf("chain: %v\n", file.HandlerChain)

	name = file.Name()
	exts = file.Exts()
	if err = file.Handle(out, &name, &exts); err != nil {
		return
	}

	if !shouldOpenFD(1) {
		waitFD(1)
	}
	defer fdClosed(1)

	manifest.log.Printf("Writing %s\n", file.Name())

	outPath, err = manifest.outFilePath(name, exts)
	if err != nil {
		return
	}
	if err = os.MkdirAll(filepath.Dir(outPath), os.ModePerm); err != nil {
		return
	}
	outFile, err = os.Create(outPath)
	if err != nil {
		return
	}
	_, err = io.Copy(outFile, out)
	if err != nil {
		return
	}

	if err = outFile.Close(); err != nil {
		return
	}

	return
}
