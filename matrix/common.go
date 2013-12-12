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

// byLenParentHandlerChain sort.Interface
type byLenHandlerChain []*FileHandler

func (a byLenHandlerChain) Len() int      { return len(a) }
func (a byLenHandlerChain) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byLenHandlerChain) Less(i, j int) bool {
	return len(a[i].HandlerChain) < len(a[j].HandlerChain)
}

// byLenParentHandlersReversed implements sort.Interface
type byLenParentHandlersReversed []*FileHandler

func (a byLenParentHandlersReversed) Len() int      { return len(a) }
func (a byLenParentHandlersReversed) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byLenParentHandlersReversed) Less(i, j int) bool {
	return len(a[j].ParentHandlers) < len(a[i].ParentHandlers)
}

var commandBucket = make(chan struct{}, 10)

func setCommandBucketLimit(limit int) {
	commandBucket = make(chan struct{}, limit)
}

func shouldExecCommand() bool {
	select {
	case commandBucket <- struct{}{}:
		return true
	default:
		return false
	}
}

func waitCommand() {
	commandBucket <- struct{}{}
}

func commandDone() {
	<-commandBucket
}
