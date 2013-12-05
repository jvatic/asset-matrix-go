package matrix

import (
	"io"
)

type ConcatinationMode int

const (
	ConcatinationModePrepend ConcatinationMode = iota
	ConcatinationModeAppend  ConcatinationMode = iota
)

type ConcatinationHandler struct {
	Handler

	parent *FileHandler
	child  *FileHandler
	mode   ConcatinationMode
	ext    string
}

func NewConcatinationHandler(parent *FileHandler, child *FileHandler, mode ConcatinationMode, ext string) (handler *ConcatinationHandler) {
	return &ConcatinationHandler{parent: parent, child: child, mode: mode, ext: ext}
}

func (handler *ConcatinationHandler) Handle(in io.Reader, out io.Writer, inputName string, inputExts []string) (name string, exts []string, err error) {
	switch handler.mode {
	case ConcatinationModePrepend:
		// TODO: write(out, handle(child) + in)
	case ConcatinationModeAppend:
		// TODO: write(out, in + handle(child))
	}

	return inputName, inputExts, nil
}

func (handler *ConcatinationHandler) OutputExt() string {
	return handler.ext
}
