package matrix

import (
	"sync"
)

type FileParents struct {
	Files []*File
	lock  sync.Mutex
}

func (p *FileParents) Lock() {
	p.lock.Lock()
}

func (p *FileParents) Unlock() {
	p.lock.Unlock()
}

func (p *FileParents) AddFile(file *File) {
	p.Lock()
	defer p.Unlock()

	p.Files = append(p.Files, file)
}

func (p *FileParents) RemoveFile(file *File) {
	p.Lock()
	defer p.Unlock()

	files := make([]*File, 0, len(p.Files))
	for _, f := range p.Files {
		if f != file {
			files = append(files, f)
		}
	}
	p.Files = files
}

func (p *FileParents) Dedup() {
	p.Lock()
	files := make([]*File, len(p.Files))
	copy(files, p.Files)
	p.Unlock()
	for _, f := range files {
		f.dedupFileParents(p)
	}
}
