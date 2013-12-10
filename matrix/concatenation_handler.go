package matrix

import (
	"fmt"
	"io"
)

type ConcatenationMode int

const (
	ConcatenationModePrepend ConcatenationMode = iota
	ConcatenationModeAppend
)

type ConcatenationHandler struct {
	Handler

	parent *FileHandler
	child  *FileHandler
	mode   ConcatenationMode
	ext    string
}

func NewConcatenationHandler(parent *FileHandler, child *FileHandler, mode ConcatenationMode, ext string) (handler *ConcatenationHandler) {
	return &ConcatenationHandler{parent: parent, child: child, mode: mode, ext: ext}
}

func (handler *ConcatenationHandler) Handle(in io.Reader, out io.Writer, inName string, inExts []string) (name string, exts []string, err error) {
	name, exts = inName, inExts

	switch handler.mode {
	case ConcatenationModePrepend:
		_, _, err = handler.child.Handle(in, out, inName, inExts)
		if err != nil {
			return
		}
		_, err = io.Copy(out, in)
		if err != nil {
			return
		}
	case ConcatenationModeAppend:
		_, err = io.Copy(out, in)
		if err != nil {
			return
		}
		_, _, err = handler.child.Handle(in, out, inName, inExts)
		if err != nil {
			return
		}
	}

	return
}

func (handler *ConcatenationHandler) OutputExt() string {
	return handler.ext
}

func (handler *ConcatenationHandler) String() string {
	var name string
	if handler.child.File != nil {
		name = handler.child.File.Name()
	}
	return fmt.Sprintf("ConcatenationHandler(%s — %v)", name, handler.child.HandlerChain)
}
