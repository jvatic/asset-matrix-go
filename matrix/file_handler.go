package matrix

type FileHandler struct {
	HandlerChain   *HandlerChain
	FileSet        []*FileHandler
	ParentHandlers []*FileHandler
}

func NewFileHandler(file *File, handlerConstructors []HandlerConstructor) (*FileHandler, error) {
	fileHandler := new(FileHandler)
	chain, err := NewHandlerChain(file, handlerConstructors)
	fileHandler.HandlerChain = chain
	return fileHandler, err
}
