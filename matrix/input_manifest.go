package matrix

type InputManifest struct {
	AssetRoots      []*AssetDir
	InputDirs       []string
	OutputDir       string
	DirPathMapping  map[string]string
	DirNameMapping  map[string]string
	FilePathMapping map[string]string
	FileNameMapping map[string]string
}

func NewInputManifest(inputDirs []string, outputDir string) *InputManifest {
	return &InputManifest{InputDirs: inputDirs, OutputDir: outputDir}
}

func (manifest *InputManifest) ScanInputDirs() {
	manifest.AssetRoots = make([]*AssetDir, len(manifest.InputDirs))
	for i, path := range manifest.InputDirs {
		dir, _ := NewAssetDir(path)
		dir.scan()

		manifest.AssetRoots[i] = dir
	}
}
