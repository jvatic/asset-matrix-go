package matrix

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
	Handlers        []HandlerConstructor
}

func NewManifest(inputDirs []string, outputDir string) *Manifest {
	return &Manifest{InputDirs: inputDirs, OutputDir: outputDir, DirPathMapping: make(map[string]*Dir), FilePathMapping: make(map[string]*File), NameMapping: make(map[string]*AssetMap), Handlers: make([]HandlerConstructor, 0)}
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

func (manifest *Manifest) AddHandler(fn HandlerConstructor) {
	manifest.Handlers = append(manifest.Handlers, fn)
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

		if err := dir.scan(); err != nil {
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
		}
	}

	return nil
}

func (manifest *Manifest) ConfigureHandlers() error {
	for _, file := range manifest.FilePathMapping {
		fileHandler, err := NewFileHandler(file, manifest.Handlers)
		if err != nil {
			return err
		}
		file.FileHandler = fileHandler
	}
	return nil
}
