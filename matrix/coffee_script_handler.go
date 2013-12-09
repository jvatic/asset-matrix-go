package matrix

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

type CoffeeHandler struct {
	Handler
}

func init() {
	Register("coffee", "js", new(CoffeeHandler), &HandlerOptions{InputMode: InputModeFlow, OutputMode: OutputModeFlow})
}

func (handler *CoffeeHandler) Handle(in io.Reader, out io.Writer, inputName string, inputExts []string) (name string, exts []string, err error) {
	cmd := exec.Command("coffee", "--compile", "--stdio")

	cmdIn, err := cmd.StdinPipe()
	if err != nil {
		return inputName, inputExts, err
	}

	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		return inputName, inputExts, err
	}

	cmdErr, err := cmd.StderrPipe()
	if err != nil {
		return inputName, inputExts, err
	}

	go func() {
		io.Copy(cmdIn, in)
		cmdIn.Close()
	}()
	go io.Copy(out, cmdOut)

	var errBytes []byte
	errBuf := bytes.NewBuffer(errBytes)
	go func() {
		io.Copy(errBuf, cmdErr)
		cmdErr.Close()
	}()

	if err := cmd.Run(); err != nil {
		return inputName, inputExts, err
	}

	if errBuf.Len() > 0 {
		return inputName, inputExts, fmt.Errorf("%v", errBuf)
	}

	for _, inExt := range inputExts {
		if inExt != "coffee" && inExt != "js" {
			exts = append(exts, inExt)
		}
	}
	exts = append(exts, "js")

	return inputName, exts, nil
}

func (handler *CoffeeHandler) OutputExt() string {
	return "js"
}

func (handler *CoffeeHandler) String() string {
	return "CoffeeHandler"
}
