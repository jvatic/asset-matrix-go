package matrix

type FileHandler struct {
	HandlerChain   []Handler
	FileSet        []*FileHandler
	ParentHandlers []*FileHandler
}

func NewFileHandler(inExt string) *FileHandler {
	fileHandler := &FileHandler{HandlerChain: make([]Handler, 0)}
	fileHandler.buildHandlerChain(inExt)
	return fileHandler
}

func (fileHandler *FileHandler) buildHandlerChain(inExt string) {
	handlers := FindHandlers(inExt)
	if handlers == nil && len(fileHandler.HandlerChain) == 0 {
		fileHandler.HandlerChain = append(fileHandler.HandlerChain, new(DefaultHandler))
		return
	}

	canAppendFlow := true
	for outExt, rh := range handlers {
		if rh.Options.InputMode == InputModeFlow && canAppendFlow {
			canAppendFlow = false
			fileHandler.HandlerChain = append(fileHandler.HandlerChain, rh.Handler)
			fileHandler.buildHandlerChain(outExt)
		} else {
			fh := NewForkHandler(NewFileHandler(inExt))
			fileHandler.HandlerChain = append(fileHandler.HandlerChain, fh)
		}
	}
}
