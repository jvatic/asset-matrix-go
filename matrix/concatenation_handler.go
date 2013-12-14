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
	childIn io.Reader
	mode    ConcatenationMode
	ext     string
}

func NewConcatenationHandler(childIn io.Reader, mode ConcatenationMode, ext string) (handler *ConcatenationHandler) {
	return &ConcatenationHandler{childIn, mode, ext}
}

func (handler *ConcatenationHandler) Handle(in io.Reader, out io.Writer, name *string, exts *[]string) (err error) {
	switch handler.mode {
	case ConcatenationModePrepend:
		// TODO
	case ConcatenationModeAppend:
		// TODO
	}

	return
}

func (handler *ConcatenationHandler) OutputExt() string {
	return handler.ext
}

func (handler *ConcatenationHandler) String() string {
	return fmt.Sprintf("ConcatenationHandler{%p}(%s)", handler, handler.ext)
}
