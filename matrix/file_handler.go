package matrix

import (
	"bytes"
	"fmt"
	"io"
	"sort"
)

type FileHandler struct {
	HandlerChain   []Handler
	FileSet        []*FileHandler
	ParentHandlers []*FileHandler
	File           *File
}

func NewFileHandler(inExt string) *FileHandler {
	fileHandler := new(FileHandler)
	fileHandler.buildHandlerChain(inExt)
	return fileHandler
}

func (fileHandler *FileHandler) buildHandlerChain(inExt string) {
	handlers := FindHandlers(inExt)
	if handlers == nil && len(fileHandler.HandlerChain) == 0 {
		fileHandler.HandlerChain = append(fileHandler.HandlerChain, NewDefaultHandler(inExt))
		return
	}

	canAppendFlow := true
	for outExt, rh := range handlers {
		if rh.Options.InputMode == InputModeFlow && canAppendFlow {
			canAppendFlow = false
			fileHandler.HandlerChain = append(fileHandler.HandlerChain, rh.Handler)
			fileHandler.buildHandlerChain(outExt)
		} else {
			fh := NewForkHandler(NewFileHandler(inExt), inExt)
			fileHandler.HandlerChain = append(fileHandler.HandlerChain, fh)
		}
	}
}

func (fileHandler *FileHandler) addHandlerAfterIndex(handler Handler, index int) {
	chain := make([]Handler, 0)
	for i, h := range fileHandler.HandlerChain {
		chain = append(chain, h)
		if i == index {
			chain = append(chain, handler)
		}
	}
	fileHandler.HandlerChain = chain
}

func (fileHandler *FileHandler) AddFileHandler(fh *FileHandler) {
	if fh != fileHandler {
		fh.AddParentFileHandler(fileHandler)
	}
	fileHandler.FileSet = append(fileHandler.FileSet, fh)
}

func (fileHandler *FileHandler) AddParentFileHandler(fh *FileHandler) {
	fileHandler.ParentHandlers = append(fileHandler.ParentHandlers, fh)
}

func (parent *FileHandler) concatinateAtIndex(child *FileHandler, handlerIndex int) {
	mode := ConcatenationModePrepend
	for _, fh := range parent.FileSet {
		if fh == child {
			break
		}

		if fh == parent {
			mode = ConcatenationModeAppend
			break
		}
	}

	ext := parent.HandlerChain[handlerIndex].OutputExt()
	parent.addHandlerAfterIndex(NewConcatenationHandler(parent, child, mode, ext), handlerIndex)
}

func (fileHandler *FileHandler) MergeWithParents() error {
	// sort parent handlers by lowest len(fh.HandlerChain)
	sort.Sort(ByLenHandlerChain(fileHandler.ParentHandlers))
	for _, parent := range fileHandler.ParentHandlers {
		// ensure the last handler in each chain have the same out ext
		index, err := removeIncompatibleHandlers(fileHandler.HandlerChain, parent.HandlerChain)
		if err != nil {
			return err
		}

		// add concatination handler to parent
		parent.concatinateAtIndex(fileHandler, index)
	}

	return nil
}

func removeIncompatibleHandlers(a []Handler, b []Handler) (int, error) {
	for i := len(a) - 1; i >= 0; i-- {
		for j := len(b) - 1; j >= 0; j-- {
			if a[i].OutputExt() == b[j].OutputExt() {
				a = a[0:i]
				return j, nil
			}
		}
	}

	return 0, fmt.Errorf("matrix: FileHandler: incompatible handler chains: %v, %v", a, b)
}

func (fileHandler *FileHandler) Handle(in io.Reader, out io.Writer, name string, exts []string) (string, []string, error) {
	var (
		err     error
		inData  []byte
		outData []byte
	)
	inDataBuf := bytes.NewBuffer(inData)
	outDataBuf := bytes.NewBuffer(outData)
	io.Copy(inDataBuf, in)
	for _, handler := range fileHandler.HandlerChain {
		name, exts, err = handler.Handle(inDataBuf, outDataBuf, name, exts)
		if err != nil {
			return name, exts, err
		}
		inDataBuf.Reset()
		io.Copy(inDataBuf, outDataBuf)
		outDataBuf.Reset()
	}
	io.Copy(out, inDataBuf)
	return name, exts, err
}
