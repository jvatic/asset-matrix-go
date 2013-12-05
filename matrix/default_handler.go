package matrix

import (
	"io"
)

type DefaultHandler struct {
	Handler

	ext string
}

func NewDefaultHandler(ext string) *DefaultHandler {
	return &DefaultHandler{ext: ext}
}

func (handler *DefaultHandler) Handle(in io.Reader, out io.Writer, inputName string, inputExts []string) (name string, exts []string, err error) {
	// TODO: copy(in, out)
	return inputName, inputExts, nil
}

func (handler *DefaultHandler) OutputExt() string {
	return handler.ext
}
