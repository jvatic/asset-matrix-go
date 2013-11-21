package matrix

type AssetMap struct {
	Dir   *AssetDir
	Files map[string]*AssetFile
}

type InputManifest struct {
	AssetRoots      []*AssetDir
	InputDirs       []string
	OutputDir       string
	DirPathMapping  map[string]*AssetDir
	FilePathMapping map[string]*AssetFile
	NameMapping     map[string]*AssetMap
}

func NewInputManifest(inputDirs []string, outputDir string) *InputManifest {
	return &InputManifest{InputDirs: inputDirs, OutputDir: outputDir, DirPathMapping: make(map[string]*AssetDir), FilePathMapping: make(map[string]*AssetFile), NameMapping: make(map[string]*AssetMap)}
}

func (manifest *InputManifest) AddDir(dir *AssetDir) {
	manifest.DirPathMapping[dir.Path()] = dir

	if manifest.NameMapping[dir.Name()] == nil {
		manifest.NameMapping[dir.Name()] = &AssetMap{Dir: dir}
	} else {
		manifest.NameMapping[dir.Name()].Dir = dir
	}
}

func (manifest *InputManifest) AddFile(asset *AssetFile) {
	manifest.FilePathMapping[asset.Path()] = asset

	if manifest.NameMapping[asset.Name()] == nil {
		manifest.NameMapping[asset.Name()] = &AssetMap{}
	}
	if manifest.NameMapping[asset.Name()].Files == nil {
		manifest.NameMapping[asset.Name()].Files = make(map[string]*AssetFile)
	}
	manifest.NameMapping[asset.Name()].Files[asset.Ext()] = asset
}

func (manifest *InputManifest) FindDirName(name string) *AssetDir {
	assetMap := manifest.NameMapping[name]
	if assetMap == nil {
		return nil
	}

	return assetMap.Dir
}

func (manifest *InputManifest) FindFileName(name string, ext string) *AssetFile {
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

func (manifest *InputManifest) ScanInputDirs() error {
	manifest.AssetRoots = make([]*AssetDir, len(manifest.InputDirs))
	for i, path := range manifest.InputDirs {
		dir, err := NewAssetDir(path, manifest, nil)
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

func (manifest *InputManifest) EvaluateDirectives() error {
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
