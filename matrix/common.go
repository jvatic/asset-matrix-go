package matrix

import (
	"regexp"
)

var directiveRegex = regexp.MustCompile("\\A\\S+=\\s*(\\S+)\\s*(.*)\\z")
var emptyLineRegex = regexp.MustCompile("\\A\\s*\\z")

// Used by File and Dir types
type AssetPointer interface {
	// Abs path of asset or dir
	Path() string
	// Relative path from the asset root excluding the file extension(s)
	Name() string
	// The directory the asset or dir lives in
	Dir() *Dir
	// The asset root the asset or dir lives in
	RootDir() *Dir
	// The manifest used to scan for assets
	Manifest() *Manifest
	// Is true if it's an asset root
	IsRoot() bool
}
