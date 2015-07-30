package assetmatrix

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Config struct {
	Paths     []*AssetRoot
	Outputs   []string
	OutputDir string
}

type Matrix struct {
	config            *Config
	cacheBreaker      string
	transformerJSPath string
	scssJSPath        string
	erbRBPath         string
	Manifest          *Manifest
}

type Manifest struct {
	Assets map[string]string `json:"assets"`
}

func New(config *Config) *Matrix {
	m := &Matrix{
		config: config,
	}
	for _, r := range config.Paths {
		r.findAsset = m.findAsset
	}
	return m
}

func (m *Matrix) Build() error {
	defer m.cleanupTempfiles()

	startedAt := time.Now()
	hashData, err := time.Now().MarshalBinary()
	if err != nil {
		return err
	}
	m.cacheBreaker = fmt.Sprintf("%x", md5.Sum(hashData))

	m.Manifest = &Manifest{
		Assets: make(map[string]string, 0),
	}

	log.Println("Installing dependencies...")
	if err := m.installDeps(); err != nil {
		return err
	}
	if err := m.createTempfiles(); err != nil {
		return err
	}
	for _, r := range m.config.Paths {
		r.transformerJSPath = m.transformerJSPath
		r.scssJSPath = m.scssJSPath
		r.erbRBPath = m.erbRBPath
	}

	log.Println("Cloning external repos...")
	for _, r := range m.config.Paths {
		if r.GitRepo != "" {
			if err := r.CloneRepo(); err != nil {
				return err
			}
		}
	}

	log.Println("Enumerating assets...")
	if err := m.enumerateAssets(); err != nil {
		return err
	}
	log.Println("Building output trees...")
	trees, err := m.buildOutputTrees()
	if err != nil {
		return err
	}
	if len(m.config.Outputs) != 0 {
		log.Println("Filtering output trees...")
		matchers := make([]*regexp.Regexp, 0, len(m.config.Outputs))
		for _, pattern := range m.config.Outputs {
			pattern = "^" + strings.Replace(strings.Replace(pattern, ".", "\\.", -1), "*", ".*", -1) + "$"
			matchers = append(matchers, regexp.MustCompile(pattern))
		}
		filteredTrees := make([][]Asset, 0, len(matchers))
		for _, t := range trees {
			a := t[len(t)-1]
			name, err := a.RelPath()
			if err != nil {
				return err
			}
			for _, r := range matchers {
				if r.MatchString(name) {
					log.Println(name)
					filteredTrees = append(filteredTrees, t)
					break
				}
			}
		}
		trees = filteredTrees
	}
	log.Println("Compiling output trees...")
	if err := m.compileTrees(trees); err != nil {
		return err
	}

	log.Println("Writing manifest.json...")
	manifestJSONPath := filepath.Join(m.config.OutputDir, "manifest.json")
	f, err := os.Create(manifestJSONPath)
	defer f.Close()
	if err != nil {
		return err
	}
	e := json.NewEncoder(f)
	if err := e.Encode(m.Manifest); err != nil {
		return err
	}

	duration := time.Since(startedAt)
	log.Printf("Completed in %s", duration)
	return nil
}

func installNpmPackages(names []string) error {
	for _, n := range names {
		if _, err := os.Stat("node_modules/" + n); err == nil {
			continue
		}
		cmd := exec.Command("npm", "install", n)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Matrix) installDeps() error {
	return installNpmPackages([]string{"recast", "es6-promise", "node-sass", "react-tools"})
}

func (m *Matrix) createTempfiles() error {
	f, err := os.Create("transformer.js")
	if err != nil {
		return err
	}
	m.transformerJSPath = f.Name()
	if _, err := f.WriteString(transformerJS); err != nil {
		return err
	}

	f, err = os.Create("scss.js")
	if err != nil {
		return err
	}
	m.scssJSPath = f.Name()
	if _, err := f.WriteString(scssJS); err != nil {
		return err
	}

	f, err = os.Create("erb.rb")
	if err != nil {
		return err
	}
	m.erbRBPath = f.Name()
	if _, err := f.WriteString(erbRB); err != nil {
		return err
	}

	return nil
}

func (m *Matrix) cleanupTempfiles() {
	os.Remove(m.transformerJSPath)
	os.Remove(m.scssJSPath)
	os.Remove(m.erbRBPath)
}

func (m *Matrix) enumerateAssets() error {
	errChan := make(chan error)
	enumerateAssets := func(r *AssetRoot) {
		errChan <- r.enumerateAssets()
	}
	for _, r := range m.config.Paths {
		r.SetCacheBreaker(m.cacheBreaker)
		go enumerateAssets(r)
	}
	for _ = range m.config.Paths {
		if err := <-errChan; err != nil {
			return err
		}
	}
	return nil
}

func (m *Matrix) buildOutputTrees() ([][]Asset, error) {
	errChan := make(chan error)
	treesChan := make(chan [][]Asset)
	trees := [][]Asset{}
	buildOutputTrees := func(r *AssetRoot) {
		trees, err := r.buildOutputTrees()
		go func() {
			errChan <- err
		}()
		go func() {
			treesChan <- trees
		}()
	}
	for _, r := range m.config.Paths {
		go buildOutputTrees(r)
	}
	for _ = range m.config.Paths {
		if err := <-errChan; err != nil {
			return nil, err
		}
		for _, t := range <-treesChan {
			trees = append(trees, t)
		}
	}
	return trees, nil
}

func (m *Matrix) compileTrees(trees [][]Asset) error {
	errChan := make(chan error)
	compileTree := func(t []Asset) {
		errChan <- m.compileTree(t)
	}
	for _, t := range trees {
		go compileTree(t)
	}
	for _ = range trees {
		if err := <-errChan; err != nil {
			return err
		}
	}
	return nil
}

var AssetNotFoundError = errors.New("Asset not found")

func (m *Matrix) findAsset(key string) (Asset, error) {
	for _, r := range m.config.Paths {
		if a, ok := r.assetIndex[key]; ok {
			return a, nil
		}
	}
	log.Printf("Asset not found: %#v", key)
	return nil, AssetNotFoundError
}

func (m *Matrix) compileTree(tree []Asset) error {
	readers := make([]io.Reader, len(tree))
	errChan := make(chan error)
	compileAsset := func(i int, a Asset) {
		r, err := a.Compile()
		readers[i] = r
		errChan <- err
	}
	for i, a := range tree {
		go compileAsset(i, a)
	}
	for _ = range tree {
		if err := <-errChan; err != nil {
			return err
		}
	}
	outputPath := tree[len(tree)-1].OutputPath()
	manifestKey := outputPath
	ext := filepath.Ext(outputPath)
	outputPath = strings.TrimSuffix(outputPath, ext) + "-" + m.cacheBreaker + ext
	manifestVal := outputPath
	outputPath = filepath.Join(m.config.OutputDir, outputPath)
	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		return err
	}
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	m.Manifest.Assets[manifestKey] = manifestVal
	log.Printf("Writing %s", outputPath)
	defer file.Close()
	var offset int64
	for i, r := range readers {
		if i > 0 {
			n, err := file.WriteAt([]byte("\n"), offset)
			if err != nil {
				return err
			}
			offset += int64(n)
		}
		n, err := io.Copy(file, r)
		if err != nil {
			return err
		}
		offset += n
	}
	return nil
}
