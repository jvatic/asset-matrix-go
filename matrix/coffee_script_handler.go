package matrix

import (
	"io"
)

type CoffeeHandler struct {
	Handler
}

func init() {
	Register("coffee", "js", new(CoffeeHandler), &HandlerOptions{InputMode: InputModeFlow, OutputMode: OutputModeFlow})
}

func (handler *CoffeeHandler) Handle(in io.Reader, out io.Writer, inputName string, inputExts []string) (name string, exts []string, err error) {
	// TODO: in -> exec coffee -> out
	return inputName, append(exts, "coffee"), nil
}
