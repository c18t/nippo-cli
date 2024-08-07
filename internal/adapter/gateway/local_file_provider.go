package gateway

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/c18t/nippo-cli/internal/domain/repository"
	"github.com/samber/do/v2"
)

type LocalFileProvider interface {
	List(param *repository.QueryListParam) ([]os.DirEntry, error)
	Read(filePath string) ([]byte, error)
	Write(filePath string, content []byte) error
	Copy(destPath string, srcPath string) error
}

type localFileProvider struct {
}

func NewLocalFileProvider(i do.Injector) (LocalFileProvider, error) {
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

func (g *localFileProvider) Read(filePath string) (content []byte, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.Read(content)
	return
}

func (g *localFileProvider) Write(filePath string, content []byte) error {
	outDir := filepath.Dir(filePath)
	err := os.MkdirAll(outDir, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	dest, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = dest.Write(content)
	if err != nil {
		return err
	}
	return nil
}

func (g *localFileProvider) Copy(destPath string, srcPath string) error {
	outDir := filepath.Dir(destPath)
	err := os.MkdirAll(outDir, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	_, err = io.Copy(dest, src)
	if err != nil {
		return err
	}
	return nil
}
