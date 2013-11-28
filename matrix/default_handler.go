package matrix

type DefaultHandler struct {
	Handler

	file *File
}

func NewDefaultHandler(file *File) (Handler, bool) {
	return &DefaultHandler{file: file}, true
}

func (handler *DefaultHandler) HandlerInputOutputs() []*HandlerInputOutput {
	ext := handler.file.Ext()
	exts := make([]*HandlerInputOutput, 0)
	return append(exts, &HandlerInputOutput{Input: ext, Output: ext, OutputMode: OM_Replace})
}

func (handler *DefaultHandler) InputFile() *File {
	return handler.file
}
