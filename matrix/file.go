package matrix

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type File struct {
	Directives   []*Directive
	Children     []*File
	Parents      *FileParents
	HandlerChain []Handler

	AssetPointer

	path             string
	name             string
	dir              *Dir
	dataByteOffset   int
	directivesParsed bool
}

var fileNameRegex = regexp.MustCompile("([^.]+)\\.?.*?\\z")

func BuildAssetName(path string) (string, error) {
	match := fileNameRegex.FindAllStringSubmatch(filepath.Base(path), -1)
	if len(match) < 1 || len(match[0]) < 2 {
		return "", fmt.Errorf("matrix: invalid path string: %s", path)
	}
	return match[0][1], nil
}

func NewFile(path string, dir *Dir) (*File, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	name, err := BuildAssetName(absPath)
	if err != nil {
		return nil, err
	}
	if !dir.IsRoot() {
		name = filepath.Join(dir.Name(), name)
	}

	file := &File{path: absPath, name: name, dir: dir, Parents: new(FileParents)}
	file.Manifest().AddFile(file)

	return file, err
}

func (file *File) Path() string {
	return file.path
}

func (file *File) Name() string {
	return file.name
}

func (file *File) Dir() *Dir {
	return file.dir
}

func (file *File) RootDir() *Dir {
	return file.dir.RootDir()
}

func (file *File) Manifest() *Manifest {
	return file.dir.Manifest()
}

func (file *File) IsRoot() bool {
	return false
}

func (file *File) Ext() string {
	return filepath.Ext(file.Path())[1:]
}

func (file *File) Exts() []string {
	return strings.Split(filepath.Base(file.Path()), ".")[1:]
}

func (file *File) EvaluateDirectives() error {
	// only parse directives on compatible files
	fileExt := file.Ext()
	directiveCompatible := false
	for _, ext := range DirectiveExts {
		if ext == fileExt {
			directiveCompatible = true
			break
		}
	}
	if !directiveCompatible {
		return nil
	}

	if !file.directivesParsed {
		if err := file.parseDirectives(); err != nil {
			return err
		}
	}

	for i := len(file.Directives); i > 0; i-- {
		directive := file.Directives[i-1]
		if err := directive.Evaluate(); err != nil {
			return err
		}

		for _, f := range directive.FileRefs {
			file.Children = append(file.Children, f)
		}
	}

	return nil
}

func (file *File) AddParent(f *File) {
	file.Parents.AddFile(f)
}

func (file *File) DedupParents() {
	file.Parents.Dedup()
}

func (file *File) LinkChildren() {
	for _, c := range file.Children {
		c.AddParent(file)
	}
}

func (file *File) BuildHandlerChain(inExt string) {
	handlers := FindHandlers(inExt)
	if handlers == nil && len(file.HandlerChain) == 0 {
		file.HandlerChain = append(file.HandlerChain, NewDefaultHandler(inExt))
		return
	}

	canAppendFlow := true
	for outExt, rh := range handlers {
		if rh.Options.InputMode == InputModeFlow && canAppendFlow {
			canAppendFlow = false
			file.HandlerChain = append(file.HandlerChain, rh.Handler)
			file.BuildHandlerChain(outExt)
		} else {
			fh := NewForkHandler(inExt) // TODO: build new handler chain for forked stream
			file.HandlerChain = append(file.HandlerChain, fh)
		}
	}
}

func (file *File) InsertConcatenationHandlers() error {
	mode := ConcatenationModePrepend
	for _, c := range file.Children {
		if c == file {
			mode = ConcatenationModeAppend
			continue
		}

		fi, ci, err := file.ConcatenationIndices(file.HandlerChain, c.HandlerChain)
		if err != nil {
			fmt.Printf("Warning: %v\n", err) // TODO use logger
			continue
		}

		ext := file.HandlerChain[ci].OutputExt()

		sh := NewSyncHandler(ext)
		ch := NewConcatenationHandler(sh.Reader, mode, ext)

		file.AddHandlerAfterIndex(ch, fi)
		c.AddHandlerAfterIndex(sh, ci)
	}

	return nil
}

func (file *File) ConcatenationIndices(a []Handler, b []Handler) (ai int, bi int, err error) {
	findIndex := func() (int, int) {
		for i := len(a) - 1; i >= 0; i-- {
			for j := len(b) - 1; j >= 0; j-- {
				if _, ok := b[j].(*ConcatenationHandler); ok {
					if j > 0 {
						continue
					}
				}
				if a[i].OutputExt() == b[j].OutputExt() {
					return j, i
				}
			}
		}
		return -1, -1
	}
	ai, bi = findIndex()
	if ai == -1 || bi == -1 {
		err = fmt.Errorf("matrix: File: incompatible handler chains: %v, %v", a, b)
	}
	return
}

func (file *File) AddHandlerAfterIndex(handler Handler, index int) {
	chain := make([]Handler, 0, len(file.HandlerChain)+1)
	for i, h := range file.HandlerChain {
		chain = append(chain, h)

		if i == index {
			chain = append(chain, handler)
		}
	}
	file.HandlerChain = chain
}

func (file *File) parseDirectives() error {
	fileRef, err := os.Open(file.Path())
	if err != nil {
		return err
	}

	defer fileRef.Close()

	var directives []*Directive

	scanner := bufio.NewScanner(fileRef)

	bytesRead := 0
	for scanner.Scan() {
		line := scanner.Bytes()

		bytesRead = bytesRead + len(line)

		// Ignore empty lines
		if emptyLineRegex.Match(line) {
			continue
		}

		// Only read directives
		// which are always at the top of the file
		if !directiveRegex.Match(line) {
			break
		}

		directive, err := NewDirective(file, string(line))
		if err != nil {
			return err
		}

		if err := directive.Evaluate(); err != nil {
			return err
		}

		directives = append(directives, directive)
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	file.Directives = directives
	file.dataByteOffset = bytesRead
	file.directivesParsed = true

	return nil
}

func (file *File) dedupFileParents(p *FileParents) {
	for _, c := range file.Children {
		p.RemoveFile(c)
	}
}

func (file *File) Handle(out io.Writer, name *string, exts *[]string) (err error) {
	if !shouldOpenFD(1) {
		waitFD(1)
	}

	f, err := os.Open(file.Path())
	if err != nil {
		return
	}

	r, w := io.Pipe()

	go func() {
		_, err = io.Copy(w, f)

		w.CloseWithError(err)
		f.Close()
		fdClosed(1)
	}()

	handlerFn := func(handler Handler, in io.Reader) *io.PipeReader {
		r, w := io.Pipe()
		go func() {
			if fdHandler, ok := handler.(FDHandler); ok {
				// handler requires file descriptors
				nFds := fdHandler.RequiredFds()

				if nFds > 0 && !shouldOpenFD(nFds) {
					data := new(bytes.Buffer)

					_, err := io.Copy(data, in)
					if err != nil {
						w.CloseWithError(err)
						return
					}

					in = data

					waitFD(nFds)
				}
				defer fdClosed(nFds)
			}

			err = handler.Handle(in, w, name, exts)
			w.CloseWithError(err)
		}()
		return r
	}

	for _, handler := range file.HandlerChain {
		r = handlerFn(handler, r)
	}

	_, err = io.Copy(out, r)
	return
}
