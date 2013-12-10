package matrix

import (
	"fmt"
	"io"
)

type ForkHandler struct {
	Handler

	fileHandler *FileHandler
	ext         string
}

func NewForkHandler(fileHandler *FileHandler, ext string) (handler *ForkHandler) {
	return &ForkHandler{fileHandler: fileHandler, ext: ext}
}

func (handler *ForkHandler) Handle(in io.Reader, out io.Writer, name string, exts []string) (string, []string, error) {
	// TODO: feed copy of input stream into fileHandler's handler chain

	return name, exts, nil
}

func (handler *ForkHandler) OutputExt() string {
	return handler.ext
}

func (handler *ForkHandler) String() string {
	return fmt.Sprintf("ForkHandler(%s)", handler.OutputExt())
}
