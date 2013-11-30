package matrix

type FileHandler struct {
	HandlerChain   *HandlerChain
	FileSet        []*FileHandler
	ParentHandlers []*FileHandler
}

func NewFileHandler(file *File) (*FileHandler, error) {
	fileHandler := new(FileHandler)
	chain, err := NewHandlerChain(file)
	fileHandler.HandlerChain = chain
	return fileHandler, err
}
