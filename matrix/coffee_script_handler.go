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

	errChan := make(chan error, 2)

	go func() {
		defer cmdIn.Close()
		_, err := io.Copy(cmdIn, in)
		if err != nil {
			errChan <- err
		}
	}()
	go func() {
		_, err := io.Copy(out, cmdOut)
		if err != nil {
			errChan <- err
		}
	}()

	var errBytes []byte
	errBuf := bytes.NewBuffer(errBytes)
	go func() {
		defer cmdErr.Close()
		_, err := io.Copy(errBuf, cmdErr)
		if err != nil {
			errChan <- err
		}
	}()

	if err := cmd.Run(); err != nil {
		return inputName, inputExts, err
	}

	if errBuf.Len() > 0 {
		return inputName, inputExts, fmt.Errorf("%v", errBuf)
	}

	select {
	case err := <-errChan:
		return inputName, inputExts, err
	default:
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
