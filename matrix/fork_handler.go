package matrix

import (
	"io"
)

type ForkHandler struct {
	Handler

	fileHandler *FileHandler
}

func NewForkHandler(fileHandler *FileHandler) (handler *ForkHandler) {
	return &ForkHandler{fileHandler: fileHandler}
}

func (handler *ForkHandler) Handle(in io.Reader, out io.Writer, inputName string, inputExts []string) (name string, exts []string, err error) {
	// TODO: feed copy of input stream into fileHandler's handler chain

	return inputName, inputExts, nil
}
