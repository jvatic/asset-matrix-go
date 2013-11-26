package matrix

type DefaultHandler struct {
	Handler

	file *File
}

func NewDefaultHandler(file *File) (Handler, bool) {
	return &DefaultHandler{file: file}, true
}

func (handler *DefaultHandler) FileExts() []*FileExt {
	return make([]*FileExt, 0)
}

func (handler *DefaultHandler) InputFile() *File {
	return handler.file
}
