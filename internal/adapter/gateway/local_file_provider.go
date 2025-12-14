package gateway

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/domain/repository"
	"github.com/samber/do/v2"
)

type LocalFileProvider interface {
	List(param *repository.QueryListParam) ([]os.DirEntry, error)
	Read(baseDir, filePath string) ([]byte, error)
	Write(filePath string, content []byte) error
	Copy(baseDir, destPath, srcPath string) error
}

type localFileProvider struct {
}

func NewLocalFileProvider(_ do.Injector) (LocalFileProvider, error) {
	return &localFileProvider{}, nil
}

func (g *localFileProvider) List(param *repository.QueryListParam) ([]os.DirEntry, error) {
	if len(param.Folders) > 1 {
		panic("param.Folders > 1")
	}

	files, err := os.ReadDir(param.Folders[0])
	if err != nil {
		return nil, err
	}

	var fileList []os.DirEntry
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		hasSuffix := len(param.FileExtensions) == 0
		for _, ext := range param.FileExtensions {
			if strings.HasSuffix(file.Name(), "."+ext) {
				hasSuffix = true
				continue
			}
		}
		if !hasSuffix {
			continue
		}
		fileList = append(fileList, file)
	}

	return fileList, nil
}

func (g *localFileProvider) Read(baseDir, filePath string) (content []byte, err error) {
	file, err := core.SafeOpen(baseDir, filePath)
	if err != nil {
		return
	}
	defer func() { _ = file.Close() }()

	_, err = file.Read(content)
	return
}

func (g *localFileProvider) Write(filePath string, content []byte) (err error) {
	outDir := filepath.Dir(filePath)
	err = os.MkdirAll(outDir, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	dest, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := dest.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	_, err = dest.Write(content)
	return err
}

func (g *localFileProvider) Copy(baseDir, destPath, srcPath string) (err error) {
	outDir := filepath.Dir(destPath)
	err = os.MkdirAll(outDir, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := dest.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	src, err := core.SafeOpen(baseDir, srcPath)
	if err != nil {
		return err
	}
	defer func() { _ = src.Close() }()

	_, err = io.Copy(dest, src)
	return err
}
