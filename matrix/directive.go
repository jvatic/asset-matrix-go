package matrix

import (
	"fmt"
	"path/filepath"
)

type Directive struct {
	Asset    *File
	String   string
	Name     string
	Value    string
	AssetRef *File
	DirRef   *Dir
}

func NewDirective(asset *File, str string) (*Directive, error) {
	// parse name and value
	match := directiveRegex.FindAllStringSubmatch(str, -1)
	if len(match) < 1 || len(match[0]) < 3 {
		return nil, fmt.Errorf("matrix: invalid directive string: %s", str)
	}

	return &Directive{Asset: asset, String: str, Name: match[0][1], Value: match[0][2]}, nil
}

func (directive *Directive) Evaluate() error {
	switch directive.Name {
	case "require":
		// TODO Handle URL to file
		name := directive.evaluateName(directive.Value)
		ext := directive.evaluateExt(directive.Value)
		asset := directive.Asset.Manifest().FindFileName(name, ext)

		if asset == nil {
			return fmt.Errorf("matrix: require: file not found: %s — %s", name, directive.Asset.Path())
		}

		directive.AssetRef = asset
	case "require_self":
		directive.AssetRef = directive.Asset
	case "require_tree":
		name := directive.evaluateName(directive.Value)

		dir := directive.Asset.Manifest().FindDirName(name)
		if dir == nil {
			return fmt.Errorf("matrix: require_tree: dir not found: %s", name)
		}

		directive.DirRef = dir
	default:
		return fmt.Errorf("matrix: unknown directive \"%s\" — %s", directive.Name, directive.Asset.Path())
	}

	return nil
}

func (directive *Directive) evaluateName(path string) string {
	if string(path[0]) == "." {
		if directive.Asset.Dir().IsRoot() {
			return string(filepath.Join("/", path)[1:])
		} else {
			return filepath.Join(directive.Asset.Dir().Name(), path)
		}
	} else {
		return path
	}
}

func (directive *Directive) evaluateExt(path string) string {
	return filepath.Ext(path)
}
