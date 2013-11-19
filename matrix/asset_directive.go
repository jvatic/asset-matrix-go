package matrix

type AssetDirective struct {
	String string
	Name   string
	Value  string
}

func NewAssetDirective(str string) *AssetDirective {
	// parse name and value
	parts := directiveRegex.FindAllStringSubmatch(str, -1)[0]

	return &AssetDirective{String: str, Name: parts[1], Value: parts[2]}
}
