package matrix

import (
	"fmt"
	"io"
)

type DefaultHandler struct {
	Handler

	ext string
}

func NewDefaultHandler(ext string) *DefaultHandler {
	return &DefaultHandler{ext: ext}
}

func (handler *DefaultHandler) Handle(in io.Reader, out io.Writer, name string, exts []string) (string, []string, error) {
	_, err := io.Copy(out, in)
	return name, exts, err
}

func (handler *DefaultHandler) OutputExt() string {
	return handler.ext
}

func (handler *DefaultHandler) String() string {
	return fmt.Sprintf("DefaultHandler(%s)", handler.OutputExt())
}
