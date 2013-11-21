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

	file := &File{path: absPath, name: name, dir: dir}
	file.Manifest().AddFile(file)

	return file, err
}

func (file *File) Path() string {
	return file.path
}

func (file *File) Name() string {
	return file.name
}

func (file *File) Dir() *Dir {
	return file.dir
}

func (file *File) RootDir() *Dir {
	return file.dir.RootDir()
}

func (file *File) Manifest() *Manifest {
	return file.dir.Manifest()
}

func (file *File) IsRoot() bool {
	return false
}

func (file *File) Ext() string {
	return filepath.Ext(file.Path())
}

func (file *File) EvaluateDirectives() error {
	if !file.directivesParsed {
		if err := file.parseDirectives(); err != nil {
			return err
		}
	}

	for _, directive := range file.Directives {
		if err := directive.Evaluate(); err != nil {
			return err
		}
	}

	return nil
}

func (file *File) parseDirectives() error {
	fileRef, err := os.Open(file.Path())
	if err != nil {
		return err
	}

	defer fileRef.Close()

	var directives []*Directive

	scanner := bufio.NewScanner(fileRef)

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

		directive, err := NewDirective(file, string(line))
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

	file.Directives = directives
	file.dataByteOffset = bytesRead
	file.directivesParsed = true

	return nil
}
