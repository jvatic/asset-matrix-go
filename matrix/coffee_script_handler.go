package matrix

import (
	"io"
)

type CoffeeScriptHandler struct {
	Handler

	inputExt     string
	inputReader  io.Reader
	inputCloser  io.Closer
	outputWriter io.Writer
}

func NewCoffeeScriptHandler(inputExt string, inputReader io.Reader, outputWriter io.Writer) (Handler, bool) {
	handler := &CoffeeScriptHandler{inputExt: inputExt, inputReader: inputReader, outputWriter: outputWriter}
	return handler, handler.canHandleInput()
}

func (handler *CoffeeScriptHandler) HandlerInputOutputs() []*HandlerInputOutput {
	exts := make([]*HandlerInputOutput, 0)
	return append(exts, &HandlerInputOutput{Input: "coffee", Output: "js", OutputMode: OutputModeReplace})
}

func (handler *CoffeeScriptHandler) canHandleInput() bool {
	for _, inOut := range handler.HandlerInputOutputs() {
		if inOut.Input == handler.inputExt {
			return true
		}
	}
	return false
}

func (handler *CoffeeScriptHandler) SetInputCloser(inputCloser io.Closer) {
	handler.inputCloser = inputCloser
}

func (handler *CoffeeScriptHandler) IsDefaultHandler() bool {
	return false
}
