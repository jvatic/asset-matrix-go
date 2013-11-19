package matrix

type InputManifest struct {
	AssetRoots      []*AssetDir
	InputDirs       []string
	OutputDir       string
	DirPathMapping  map[string]*AssetDir
	DirNameMapping  map[string]*AssetDir
	FilePathMapping map[string]*AssetFile
	FileNameMapping map[string]*AssetFile
}

func NewInputManifest(inputDirs []string, outputDir string) *InputManifest {
	return &InputManifest{InputDirs: inputDirs, OutputDir: outputDir, DirPathMapping: make(map[string]*AssetDir), DirNameMapping: make(map[string]*AssetDir), FilePathMapping: make(map[string]*AssetFile), FileNameMapping: make(map[string]*AssetFile)}
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
