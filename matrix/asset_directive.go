package matrix

import (
	"fmt"
	"path/filepath"
)

type AssetDirective struct {
	Asset    *AssetFile
	String   string
	Name     string
	Value    string
	AssetRef *AssetFile
	DirRef   *AssetDir
}

func NewAssetDirective(asset *AssetFile, str string) (*AssetDirective, error) {
	// parse name and value
	match := directiveRegex.FindAllStringSubmatch(str, -1)
	if len(match) < 1 || len(match[0]) < 3 {
		return nil, fmt.Errorf("matrix: invalid directive string: %s", str)
	}

	return &AssetDirective{Asset: asset, String: str, Name: match[0][1], Value: match[0][2]}, nil
}

func (directive *AssetDirective) Evaluate() error {
	switch directive.Name {
	case "require":
		// TODO Handle URL to file
		name := directive.evaluateName(directive.Value)
		asset := directive.Asset.Manifest.FileNameMapping[name]

		if asset == nil {
			return fmt.Errorf("matrix: require: file not found: %s", name)
		}

		directive.AssetRef = asset
	case "require_self":
		directive.AssetRef = directive.Asset
	case "require_tree":
		name := directive.evaluateName(directive.Value)

		dir := directive.Asset.Manifest.DirNameMapping[name]
		if dir == nil {
			return fmt.Errorf("matrix: require_tree: dir not found: %s", name)
		}

		directive.DirRef = dir
	default:
		return fmt.Errorf("matrix: unknown directive \"%s\"", directive.Name)
	}

	return nil
}

func (directive *AssetDirective) evaluateName(path string) string {
	if string(path[0]) == "." {
		if directive.Asset.Dir.IsRoot {
			return string(filepath.Join("/", path)[1:])
		} else {
			return filepath.Join(directive.Asset.Dir.Name, path)
		}
	} else {
		return path
	}
}
