package matrix

import (
	"fmt"
	"io"
)

type ForkHandler struct {
	ext string
}

func NewForkHandler(ext string) (handler *ForkHandler) {
	return &ForkHandler{ext: ext}
}

func (handler *ForkHandler) Handle(in io.Reader, out io.Writer, name *string, exts *[]string) (err error) {
	// TODO: feed copy of input stream into file's handler chain
	return
}

func (handler *ForkHandler) OutputExt() string {
	return handler.ext
}

func (handler *ForkHandler) String() string {
	return fmt.Sprintf("ForkHandler(%s)", handler.OutputExt())
}
