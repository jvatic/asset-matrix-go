package matrix

import (
	"io"
)

type DefaultHandler struct {
	Handler

	inputExt     string
	inputReader  io.Reader
	inputCloser  io.Closer
	outputWriter io.Writer
}

func NewDefaultHandler(inputExt string, inputReader io.Reader, outputWriter io.Writer) (Handler, bool) {
	return &DefaultHandler{inputExt: inputExt, inputReader: inputReader, outputWriter: outputWriter}, true
}

func (handler *DefaultHandler) HandlerInputOutputs() []*HandlerInputOutput {
	exts := make([]*HandlerInputOutput, 0)
	return append(exts, &HandlerInputOutput{Input: handler.inputExt, Output: handler.inputExt, OutputMode: OutputModeReplace})
}

func (handler *DefaultHandler) SetInputCloser(inputCloser io.Closer) {
	handler.inputCloser = inputCloser
}

func (handler *DefaultHandler) IsDefaultHandler() bool {
	return true
}
