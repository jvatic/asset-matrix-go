package matrix

import (
	"fmt"
	"io"
)

type SyncHandler struct {
	Reader io.Reader
	w      io.Writer
	ext    string
}

func NewSyncHandler(ext string) *SyncHandler {
	r, w := io.Pipe()
	return &SyncHandler{r, w, ext}
}

func (handler *SyncHandler) Handle(in io.Reader, out io.Writer, name *string, exts *[]string) (err error) {
	// TODO: copy in to handler.w, copy in to out
	return
}

func (handler *SyncHandler) OutputExt() string {
	return handler.ext
}

func (handler *SyncHandler) String() string {
	return fmt.Sprintf("SyncHandler{%p}(%s)", handler, handler.ext)
}
