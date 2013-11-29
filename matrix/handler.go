package matrix

import (
	"io"
)

type OutputMode int

const (
	OutputModeReplace OutputMode = iota // replace input extension
	OutputModePrepend OutputMode = iota // prepend to input extension
	OutputModeAppend  OutputMode = iota // append to input extension
	OutputModeDiscard OutputMode = iota // there is no output
)

type HandlerInputOutput struct {
	// May be set to "*" to catch all or a string such as "js" to match all files with the "js" file extention
	// Is set to an empty string if the output does not directly correlate to an input
	Input string
	// Is true if handler takes all matching inputs at once
	InputMulti bool
	// May be set to a string such as "js"
	// Is set to an empty string if the input does not directly correlate to an output
	Output string
	// Specifies how the output suffix should be placed relative to the input suffix and if there is an output
	OutputMode OutputMode
}

type Handler interface {
	// File extensions supported for input/output
	HandlerInputOutputs() []*HandlerInputOutput
	InputReader() io.Reader
	SetInputCloser(io.Closer)
	IsDefaultHandler() bool
}

type HandlerConstructor func(inputExt string, inputReader io.Reader, outputWriter io.Writer) (handler Handler, canHandle bool)
