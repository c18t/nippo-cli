package service

import (
	"os"
	"path"
	"text/template"

	"github.com/c18t/nippo-cli/internal/core"
)

type ITemplateService interface {
	SaveTo(f *os.File, templateName string, data any) error
}

type templateService struct {
	t *template.Template
}

func NewTemplateService() ITemplateService {
	return &templateService{}
}

func (s *templateService) SaveTo(f *os.File, templateName string, data any) error {
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
	t, err := template.ParseGlob(path.Join(core.Cfg.GetDataDir(), "templates", "*.html"))
	if err != nil {
		return err
	}
	s.t = t
	return nil
}
