package matrix

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type AssetFile struct {
	Path             string
	Name             string
	Directives       []*AssetDirective
	Dir              *AssetDir
	Manifest         *InputManifest
	dataByteOffset   int
	directivesParsed bool
}

var fileNameRegex = regexp.MustCompile("([^.]+)\\.?.*?\\z")

func BuildAssetName(path string) (string, error) {
	match := fileNameRegex.FindAllStringSubmatch(filepath.Base(path), -1)
	if len(match) < 1 || len(match[0]) < 2 {
		return "", fmt.Errorf("matrix: invalid path string: %s", path)
	}
	return match[0][1], nil
}

func NewAssetFile(path string, manifest *InputManifest, dir *AssetDir) (*AssetFile, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	name, err := BuildAssetName(absPath)
	if err != nil {
		return nil, err
	}
	if !dir.IsRoot {
		name = filepath.Join(dir.Name, name)
	}

	asset := &AssetFile{Path: absPath, Name: name, Manifest: manifest, Dir: dir}

	manifest.FilePathMapping[asset.Path] = asset
	manifest.FileNameMapping[asset.Name] = asset

	return asset, err
}

func (asset *AssetFile) EvaluateDirectives() error {
	if !asset.directivesParsed {
		if err := asset.parseDirectives(); err != nil {
			return err
		}
	}

	for _, directive := range asset.Directives {
		if err := directive.Evaluate(); err != nil {
			return err
		}
	}

	return nil
}

func (asset *AssetFile) parseDirectives() error {
	file, err := os.Open(asset.Path)
	if err != nil {
		return err
	}

	defer file.Close()

	var directives []*AssetDirective

	scanner := bufio.NewScanner(file)

	bytesRead := 0
	for scanner.Scan() {
		line := scanner.Bytes()

		bytesRead = bytesRead + len(line)

		// Ignore empty lines
		if emptyLineRegex.Match(line) {
			continue
		}

		// Only read directives
		// which are always at the top of the file
		if !directiveRegex.Match(line) {
			break
		}

		directive, err := NewAssetDirective(asset, string(line))
		if err != nil {
			return err
		}

		if err := directive.Evaluate(); err != nil {
			return err
		}

		directives = append(directives, directive)
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	asset.Directives = directives
	asset.dataByteOffset = bytesRead
	asset.directivesParsed = true

	return nil
}
