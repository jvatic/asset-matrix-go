package matrix

import (
	"fmt"
	"path/filepath"
)

type Directive struct {
	File    *File
	Name    string
	Value   string
	FileRef *File
	DirRef  *Dir
}

var DirectiveExts = []string{"js", "coffee", "css", "scss", "sass"}

func NewDirective(file *File, str string) (*Directive, error) {
	// parse name and value
	match := directiveRegex.FindAllStringSubmatch(str, -1)
	if len(match) < 1 || len(match[0]) < 3 {
		return nil, fmt.Errorf("matrix: invalid directive string: %s", str)
	}

	return &Directive{File: file, Name: match[0][1], Value: match[0][2]}, nil
}

func (directive *Directive) Files() []*File {
	var files []*File

	if directive.FileRef != nil {
		files = append(files, directive.FileRef)
	} else if directive.DirRef != nil {
		for _, file := range dirFiles(directive.DirRef) {
			files = append(files, file)
		}
	}

	return files
}

func (directive *Directive) Evaluate() error {
	switch directive.Name {
	case "require":
		// TODO Handle URL to file
		name := directive.evaluateName(directive.Value)
		ext := directive.evaluateExt(directive.Value)
		file := directive.File.Manifest().FindFileName(name, ext)

		if file == nil {
			return fmt.Errorf("matrix: require: file not found: %s — %s", name, directive.File.Path())
		}

		directive.FileRef = file
	case "require_self":
		directive.FileRef = directive.File
	case "require_tree":
		name := directive.evaluateName(directive.Value)

		dir := directive.File.Manifest().FindDirName(name)
		if dir == nil {
			return fmt.Errorf("matrix: require_tree: dir not found: %s", name)
		}

		directive.DirRef = dir
	default:
		return fmt.Errorf("matrix: unknown directive \"%s\" — %s", directive.Name, directive.File.Path())
	}

	return nil
}

func (directive *Directive) evaluateName(path string) string {
	if string(path[0]) == "." {
		if directive.File.Dir().IsRoot() {
			return string(filepath.Join("/", path)[1:])
		} else {
			return filepath.Join(directive.File.Dir().Name(), path)
		}
	} else {
		return path
	}
}

func (directive *Directive) evaluateExt(path string) string {
	return filepath.Ext(path)
}

func dirFiles(dir *Dir) (files []*File) {
	for _, file := range dir.Files {
		files = append(files, file)
	}

	for _, dir := range dir.Dirs {
		for _, file := range dirFiles(dir) {
			files = append(files, file)
		}
	}

	return
}

func (directive *Directive) String() string {
	return fmt.Sprintf("Directive(%s %s)", directive.Name, directive.Value)
}
