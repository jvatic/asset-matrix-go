package matrix

import (
	"bytes"
	"fmt"
	"io"
)

type ConcatenationMode int

const (
	ConcatenationModePrepend ConcatenationMode = iota
	ConcatenationModeAppend  ConcatenationMode = iota
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

func (handler *ConcatenationHandler) Handle(in io.Reader, out io.Writer, inputName string, inputExts []string) (name string, exts []string, err error) {
	var childOutBytes []byte
	childBuf := bytes.NewBuffer(childOutBytes)
	handler.child.Handle(in, childBuf, inputName, inputExts)

	switch handler.mode {
	case ConcatenationModePrepend:
		io.Copy(out, childBuf)
		io.Copy(out, in)
	case ConcatenationModeAppend:
		io.Copy(out, in)
		io.Copy(out, childBuf)
	}

	return inputName, inputExts, nil
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
