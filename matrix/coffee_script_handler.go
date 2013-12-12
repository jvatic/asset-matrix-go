package matrix

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

type CoffeeHandler struct{}

func init() {
	Register("coffee", "js", new(CoffeeHandler), &HandlerOptions{InputMode: InputModeFlow, OutputMode: OutputModeFlow})
}

func (handler *CoffeeHandler) RequiredFds() int {
	return 3 // stdin, stdout, stderr
}

func (handler *CoffeeHandler) Handle(in io.Reader, out io.Writer, name *string, exts *[]string) (err error) {
	cmd := exec.Command("coffee", "--compile", "--stdio")

	cmdIn, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	cmdErr, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	errChan := make(chan error, 4)

	go func() {
		_, err := io.Copy(cmdIn, in)
		if err != nil {
			errChan <- err
		}
		if err := cmdIn.Close(); err != nil {
			errChan <- err
		}
	}()
	go func() {
		_, err := io.Copy(out, cmdOut)
		if err != nil {
			errChan <- err
		}
	}()

	errBuf := new(bytes.Buffer)
	go func() {
		_, err := io.Copy(errBuf, cmdErr)
		if err != nil {
			errChan <- err
		}
	}()

	if err = cmd.Run(); err != nil {
		err = fmt.Errorf("%v:\n%v", err, errBuf)
		return
	}

	var outExts []string
	for _, inExt := range *exts {
		if inExt != "coffee" && inExt != "js" {
			outExts = append(outExts, inExt)
		}
	}
	outExts = append(outExts, "js")
	*exts = outExts

	return
}

func (handler *CoffeeHandler) OutputExt() string {
	return "js"
}

func (handler *CoffeeHandler) String() string {
	return "CoffeeHandler"
}
