package matrix

import (
	"fmt"
)

type AssetDirective struct {
	String string
	Name   string
	Value  string
}

func NewAssetDirective(str string) (*AssetDirective, error) {
	// parse name and value
	match := directiveRegex.FindAllStringSubmatch(str, -1)
	if len(match) < 1 || len(match[0]) < 3 {
		return nil, fmt.Errorf("matrix: invalid directive string: %s", str)
	}

	return &AssetDirective{String: str, Name: match[0][1], Value: match[0][2]}, nil
}
