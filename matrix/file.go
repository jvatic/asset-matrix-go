package matrix

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type File struct {
	Directives []*Directive

	AssetPointer

	path             string
	name             string
	dir              *Dir
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

func NewFile(path string, dir *Dir) (*File, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	name, err := BuildAssetName(absPath)
	if err != nil {
		return nil, err
	}
	if !dir.IsRoot() {
		name = filepath.Join(dir.Name(), name)
	}

	asset := &File{path: absPath, name: name, dir: dir}
	asset.Manifest().AddFile(asset)

	return asset, err
}

func (asset *File) Path() string {
	return asset.path
}

func (asset *File) Name() string {
	return asset.name
}

func (asset *File) Dir() *Dir {
	return asset.dir
}

func (asset *File) RootDir() *Dir {
	return asset.dir.RootDir()
}

func (asset *File) Manifest() *Manifest {
	return asset.dir.Manifest()
}

func (asset *File) IsRoot() bool {
	return false
}

func (asset *File) Ext() string {
	return filepath.Ext(asset.Path())
}

func (asset *File) EvaluateDirectives() error {
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

func (asset *File) parseDirectives() error {
	file, err := os.Open(asset.Path())
	if err != nil {
		return err
	}

	defer file.Close()

	var directives []*Directive

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

		directive, err := NewDirective(asset, string(line))
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
