package model

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"google.golang.org/api/drive/v3"
)

type Nippo struct {
	Date       NippoDate
	FilePath   string
	Content    []byte
	RemoteFile *drive.File
}

type NippoDate interface {
	PathString() string
	FileString() string
	TitleString() string

	Year() int
	Month() time.Month
	Day() int
	Weekday() time.Weekday
}
type nippoDate struct {
	time time.Time
}

func NewNippo(filePath string) (*Nippo, error) {
	nippo := &Nippo{}

	if err := checkNippoIsExist(filePath); err != nil {
		return nippo, err
	}

	nippo.Date = NewNippoDate(filePath)
	nippo.FilePath = filePath
	return nippo, nil
}

func NewNippoDate(filePath string) NippoDate {
	date, err := time.Parse("2006-01-02", filepath.Base(filePath)[:10])
	if err != nil {
		panic(err)
	}
	return &nippoDate{date}
}

func (n *Nippo) GetMarkdown() ([]byte, error) {
	if len(n.Content) > 0 {
		return n.Content, nil
	}
	f, err := os.Open(n.FilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	n.Content, err = io.ReadAll(f)
	return n.Content, err
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

func (date *nippoDate) String() string {
	return date.PathString()
}

func (date *nippoDate) FileString() string {
	return fmt.Sprintf("%04d-%02d-%02d", date.time.Year(), date.time.Month(), date.time.Day())
}

func (date *nippoDate) PathString() string {
	return fmt.Sprintf("%04d%02d%02d", date.time.Year(), date.time.Month(), date.time.Day())
}

func (date *nippoDate) TitleString() string {
	return fmt.Sprintf("%02d/%02d %s", date.time.Month(), date.time.Day(),
		strings.ToLower(date.time.Weekday().String()[:3]))
}

func (date *nippoDate) Year() int {
	return date.time.Year()
}

func (date *nippoDate) Month() time.Month {
	return date.time.Month()
}

func (date *nippoDate) Day() int {
	return date.time.Day()
}

func (date *nippoDate) Weekday() time.Weekday {
	return date.time.Weekday()
}
