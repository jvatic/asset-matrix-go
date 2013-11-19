package matrix

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

type AssetFile struct {
	Path           string
	Directives     []*AssetDirective
	dataByteOffset int
}

func NewAssetFile(path string) (*AssetFile, error) {
	absPath, err := filepath.Abs(path)

	if err != nil {
		return nil, err
	}

	return &AssetFile{Path: absPath}, nil
}

func (asset *AssetFile) ParseDirectives() error {
	file, err := os.Open(asset.Path)

	defer file.Close()

	if err != nil {
		return err
	}

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

		directive := NewAssetDirective(string(line))

		directives = append(directives, directive)

		fmt.Println(directive.Name, directive.Value)
	}

	asset.dataByteOffset = bytesRead
	asset.Directives = directives

	fmt.Printf("dataByteOffset: %d\n", asset.dataByteOffset)

	return nil
}
