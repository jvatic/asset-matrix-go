package matrix

type DefaultHandler struct {
	Handler

	file *File
}

func NewDefaultHandler(file *File) (Handler, bool) {
	return &DefaultHandler{file: file}, true
}

func (handler *DefaultHandler) FileExts() []*FileExt {
	ext := handler.file.Ext()
	exts := make([]*FileExt, 0)
	return append(exts, &FileExt{Input: ext, Output: ext, OutputMode: OM_Replace})
}

func (handler *DefaultHandler) InputFile() *File {
	return handler.file
}
