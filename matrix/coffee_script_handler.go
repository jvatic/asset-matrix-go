package matrix

type CoffeeScriptHandler struct {
	Handler

	file *File
}

func NewCoffeeScriptHandler(file *File) (Handler, bool) {
	handler := new(CoffeeScriptHandler)

	if !handler.canHandleFile(file) {
		return nil, false
	}

	handler.file = file

	return handler, true
}

func (handler *CoffeeScriptHandler) HandlerInputOutputs() []*HandlerInputOutput {
	exts := make([]*HandlerInputOutput, 0)
	return append(exts, &HandlerInputOutput{Input: "coffee", Output: "js", OutputMode: OM_Replace})
}

func (handler *CoffeeScriptHandler) InputFile() *File {
	return handler.file
}

func (handler *CoffeeScriptHandler) canHandleFile(file *File) bool {
	exts := handler.HandlerInputOutputs()
	fileExt := file.Ext()
	for _, ext := range exts {
		if fileExt == ext.Input {
			return true
		}
	}
	return false
}
