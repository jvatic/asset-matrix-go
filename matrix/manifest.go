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
}

func NewManifest(inputDirs []string, outputDir string) *Manifest {
	return &Manifest{InputDirs: inputDirs, OutputDir: outputDir, DirPathMapping: make(map[string]*Dir), FilePathMapping: make(map[string]*File), NameMapping: make(map[string]*AssetMap)}
}

func (manifest *Manifest) AddDir(dir *Dir) {
	manifest.DirPathMapping[dir.Path()] = dir

	if manifest.NameMapping[dir.Name()] == nil {
		manifest.NameMapping[dir.Name()] = &AssetMap{Dir: dir}
	} else {
		manifest.NameMapping[dir.Name()].Dir = dir
	}
}

func (manifest *Manifest) AddFile(asset *File) {
	manifest.FilePathMapping[asset.Path()] = asset

	if manifest.NameMapping[asset.Name()] == nil {
		manifest.NameMapping[asset.Name()] = &AssetMap{}
	}
	if manifest.NameMapping[asset.Name()].Files == nil {
		manifest.NameMapping[asset.Name()].Files = make(map[string]*File)
	}
	manifest.NameMapping[asset.Name()].Files[asset.Ext()] = asset
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
		for _, asset := range files {
			return asset
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

		for _, asset := range assetMap.Files {
			if err := asset.EvaluateDirectives(); err != nil {
				return err
			}
		}
	}

	return nil
}
