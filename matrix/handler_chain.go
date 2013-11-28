package matrix

import (
	"io"
	"os"
)

type HandlerChain struct {
	Handlers []Handler
	Outputs  map[*HandlerInputOutput]io.Reader
}

func NewHandlerChain(file *File, handlerConstructors []HandlerConstructor) (*HandlerChain, error) {
	handlerChain := &HandlerChain{Handlers: make([]Handler, 0), Outputs: make(map[*HandlerInputOutput]io.Reader)}
	err := handlerChain.evaluate(file, handlerConstructors)
	return handlerChain, err
}

func (chain *HandlerChain) evaluate(file *File, handlerConstructors []HandlerConstructor) error {
	f, err := os.Open(file.Path())
	if err != nil {
		return err
	}

	r, w := io.Pipe()
	handler := chain.addHandler(file.Ext(), f, w, handlerConstructors)
	handler.SetInputCloser(f)

	if handler.IsDefaultHandler() {
		chain.Outputs[handler.HandlerInputOutputs()[0]] = r // DefaultHandler only has a single HandlerInputOutput
	} else {
		chain.addHandlerFromHandler(handler, r, handlerConstructors)
	}

	return nil
}

func (chain *HandlerChain) addHandlerFromHandler(ref Handler, input io.Reader, handlerConstructors []HandlerConstructor) {
	for _, inOut := range ref.HandlerInputOutputs() {
		r, w := io.Pipe()
		h := chain.addHandler(inOut.Output, input, w, handlerConstructors)
		if h.IsDefaultHandler() {
			chain.Outputs[inOut] = r
		} else {
			chain.addHandlerFromHandler(h, r, handlerConstructors)
		}
	}
}

func (chain *HandlerChain) addHandler(ext string, input io.Reader, output io.Writer, handlerConstructors []HandlerConstructor) Handler {
	for _, c := range handlerConstructors {
		if handler, canHandle := c(ext, input, output); canHandle {
			chain.Handlers = append(chain.Handlers, handler)
			return handler
		}
	}

	handler, _ := NewDefaultHandler(ext, input, output)
	chain.Handlers = append(chain.Handlers, handler)
	return handler
}
