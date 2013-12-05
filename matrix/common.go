package matrix

import (
	"regexp"
)

var directiveRegex = regexp.MustCompile("\\A\\S+=\\s*(\\S+)\\s*(.*)\\z")
var emptyLineRegex = regexp.MustCompile("\\A\\s*\\z")

// Used by File and Dir types
type AssetPointer interface {
	// Abs path of file or dir
	Path() string
	// Relative path from the file root excluding the file extension(s)
	Name() string
	// The directory the file or dir lives in
	Dir() *Dir
	// The asset root the file or dir lives in
	RootDir() *Dir
	// The manifest used to scan for files
	Manifest() *Manifest
	// Is true if it's an asset root
	IsRoot() bool
}

// ByLenParentHandlerChain sort.Interface
type ByLenHandlerChain []*FileHandler

func (a ByLenHandlerChain) Len() int      { return len(a) }
func (a ByLenHandlerChain) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByLenHandlerChain) Less(i, j int) bool {
	return len(a[i].HandlerChain) < len(a[j].HandlerChain)
}

// ByLenParentHandlersReversed implements sort.Interface
type ByLenParentHandlersReversed []*FileHandler

func (a ByLenParentHandlersReversed) Len() int      { return len(a) }
func (a ByLenParentHandlersReversed) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByLenParentHandlersReversed) Less(i, j int) bool {
	return len(a[j].ParentHandlers) < len(a[i].ParentHandlers)
}
