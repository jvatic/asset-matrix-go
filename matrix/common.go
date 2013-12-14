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

var fdBucket = make(chan struct{}, 150)

func SetFDLimit(limit int) {
	fdBucket = make(chan struct{}, limit)
}

func fdLimit() int {
	return cap(fdBucket)
}

func shouldOpenFD(n int) bool {
	j := 0
	for {
		select {
		case fdBucket <- struct{}{}:
			j++
			if j == n {
				return true
			}
		default:
			fdClosed(j)
			return false
		}
	}
}

func waitFD(n int) {
	for i := 0; i < n; i++ {
		fdBucket <- struct{}{}
	}
}

func fdClosed(n int) {
	for i := 0; i < n; i++ {
		<-fdBucket
	}
}
