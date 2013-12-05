package matrix

import (
	"io"
)

type InputMode int

const (
	InputModeFlow InputMode = iota // output replaces input in the current file stream (e.g. coffee -> js)
	InputModeFork InputMode = iota // forks input into a new file stream (e.g. js -> [js, js.gz])
)

type OutputMode int

const (
	OutputModeFlow  OutputMode = iota // single output for each input (e.g. coffee -> js)
	OutputModeUnite OutputMode = iota // single output for the collection of all inputs (e.g. * -> manifest.json)
)

type Handler interface {
	Handle(in io.Reader, out io.Writer, inputName string, inputExts []string) (name string, exts []string, err error)
}

type HandlerOptions struct {
	InputMode  InputMode
	OutputMode OutputMode
}

type RegisteredHandler struct {
	Handler Handler
	Options *HandlerOptions
}

// map[input ext]map[output ext]{handler, options}
var registeredHandlers map[string]map[string]*RegisteredHandler = make(map[string]map[string]*RegisteredHandler)

func Register(inExt string, outExt string, handler Handler, options *HandlerOptions) {
	if registeredHandlers[inExt] == nil {
		registeredHandlers[inExt] = make(map[string]*RegisteredHandler)
	}
	registeredHandlers[inExt][outExt] = &RegisteredHandler{handler, options}
}

func FindHandlers(inExt string) (handlers map[string]*RegisteredHandler) {
	handlers = registeredHandlers[inExt]
	if handlers != nil {
		return handlers
	}

	return registeredHandlers["*"]
}
