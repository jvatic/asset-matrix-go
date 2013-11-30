package matrix

import (
	"io"
)

type DefaultHandler struct {
	Handler
}

func (handler *DefaultHandler) Handle(in io.Reader, out io.Writer, inputName string, inputExts []string) (name string, exts []string, err error) {
	// TODO: copy(in, out)
	return inputName, inputExts, nil
}
