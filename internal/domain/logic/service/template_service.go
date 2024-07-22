package service

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/c18t/nippo-cli/internal/core"
	i "github.com/c18t/nippo-cli/internal/domain/service"
	"github.com/samber/do/v2"
)

type templateService struct {
	t *template.Template
}

func NewTemplateService(i do.Injector) (i.TemplateService, error) {
	return &templateService{}, nil
}

func (s *templateService) SaveTo(filePath string, templateName string, data any) error {
	outputDir := filepath.Dir(filePath)
	err := os.MkdirAll(outputDir, 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println(err)
		return nil
	}

	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer f.Close()

	tmpl, err := s.template().Lookup("layout").Clone()
	if err != nil {
		return err
	}
	tmpl, err = tmpl.AddParseTree("content", s.template().Lookup(templateName).Tree)
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(f, "layout", data)
}

func (s *templateService) template() *template.Template {
	if s.t == nil {
		s.lazyLoadTemplate()
	}
	return s.t
}

func (s *templateService) lazyLoadTemplate() error {
	t, err := template.ParseGlob(filepath.Join(core.Cfg.GetDataDir(), "templates", "*.html"))
	if err != nil {
		return err
	}
	s.t = t
	return nil
}
