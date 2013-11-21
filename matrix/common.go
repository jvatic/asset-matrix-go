package matrix

import (
	"regexp"
)

var directiveRegex = regexp.MustCompile("\\A\\S+=\\s*(\\S+)\\s*(.*)\\z")
var emptyLineRegex = regexp.MustCompile("\\A\\s*\\z")

/*
  Used by AssetFile and AssetDir types
*/
type AssetPointer interface {
	Path() string
	Name() string
	Dir() *AssetDir
	RootDir() *AssetDir
	Manifest() *InputManifest
	IsRoot() bool
}
