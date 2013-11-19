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
}

var fileNameRegex = regexp.MustCompile("([^.]+)\\.?.*?\\z")

func NewAssetFile(path string, manifest *InputManifest, dir *AssetDir) (*AssetFile, error) {
	absPath, err := filepath.Abs(path)

	nameMatch := fileNameRegex.FindAllStringSubmatch(filepath.Base(absPath), -1)
	if len(nameMatch) < 1 || len(nameMatch[0]) < 2 {
		return nil, fmt.Errorf("matrix: invalid path string: %s", path)
	}
	name := nameMatch[0][1]
	if !dir.IsRoot {
		name = filepath.Join(dir.Name, name)
	}

	asset := &AssetFile{Path: absPath, Name: name, Manifest: manifest, Dir: dir}

	manifest.FilePathMapping[asset.Path] = asset
	manifest.FileNameMapping[asset.Name] = asset

	return asset, err
}

func (asset *AssetFile) ParseDirectives() error {
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

		directive, err := NewAssetDirective(string(line))
		if err != nil {
			return err
		}

		directives = append(directives, directive)

		fmt.Println(directive.Name, directive.Value)
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	asset.dataByteOffset = bytesRead
	asset.Directives = directives

	fmt.Printf("dataByteOffset: %d\n", asset.dataByteOffset)

	return nil
}
