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
	return &AssetFile{Path: absPath}, err
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
