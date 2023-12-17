package model

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Nippo struct {
	Date     NippoDate
	FilePath string
}

type NippoDate string

func NewNippo(filePath string) (Nippo, error) {
	nippo := Nippo{}

	if err := checkNippoIsExist(filePath); err != nil {
		return nippo, err
	}

	nippo.Date = NewNippoDate(filePath)
	nippo.FilePath = filePath
	return nippo, nil
}

func NewNippoDate(filePath string) NippoDate {
	return NippoDate(strings.TrimSuffix(path.Base(filePath), ".md"))
}

func (n *Nippo) GetMarkdown() ([]byte, error) {
	f, err := os.Open(n.FilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}

func (n *Nippo) GetHtml() ([]byte, error) {
	data, err := n.GetMarkdown()
	if err != nil {
		return nil, err
	}

	extensions := parser.CommonExtensions | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(data)

	htmlFlags := html.CommonFlags
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer), nil
}

func checkNippoIsExist(filePath string) error {
	if f, err := os.Stat(filePath); os.IsNotExist(err) || f.IsDir() {
		return fmt.Errorf("nippo not found: filePath=%v", filePath)
	}
	return nil
}
