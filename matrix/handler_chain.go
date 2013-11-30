package matrix

import ()

type HandlerChain struct {
	Handlers []Handler
}

func NewHandlerChain(file *File) (*HandlerChain, error) {
	handlerChain := &HandlerChain{Handlers: make([]Handler, 0)}
	err := handlerChain.evaluate(file)
	return handlerChain, err
}

func (chain *HandlerChain) evaluate(file *File) error {
	return nil
}
