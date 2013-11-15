package matrix

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type AssetFile struct {
	Path           string
	Directives     []*AssetDirective
	dataByteOffset int
}

func NewAssetFile(path string) (*AssetFile, error) {
	absPath, _ := filepath.Abs(path)
	return &AssetFile{Path: absPath}, nil
}

var directiveRegexPattern = "\\A\\S+=\\s*(\\S+)\\s*(.*)\\z"

func (file *AssetFile) ParseDirectives() {
	fileDesc, _ := os.Open(file.Path)

	var directives []*AssetDirective

	scanner := bufio.NewScanner(fileDesc)

	bytesRead := 0
	for scanner.Scan() {
		line := scanner.Text()

		bytesRead = bytesRead + len(line)

		// Ignore empty lines
		lineEmpty, _ := regexp.MatchString("\\A\\s*\\z", line)
		if lineEmpty {
			continue
		}

		// Only read directives
		// which are always at the top of the file
		matched, _ := regexp.MatchString(directiveRegexPattern, line)
		if !matched {
			break
		}

		directive, _ := NewAssetDirective(line)

		directives = append(directives, directive)

		fmt.Println(directive.Name, directive.Value)
	}

	file.dataByteOffset = bytesRead
	file.Directives = directives

	fmt.Printf("dataByteOffset: %d\n", file.dataByteOffset)
}
