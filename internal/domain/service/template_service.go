package service

import (
	"os"
	"path"
	"text/template"

	"github.com/c18t/nippo-cli/internal/core"
)

type templateService struct {
	config   *core.Config
	template *template.Template
}

type ITemplateService interface {
	SaveTo(f *os.File, templateName string, data any) error
}

func NewTemplateService() (ITemplateService, error) {
	var err error
	s := &templateService{}
	t, err := template.ParseGlob(path.Join(core.Cfg.GetDataDir(), "templates", "*.html"))
	if err != nil {
		return s, err
	}
	s.config = core.Cfg
	s.template = t
	return s, nil
}

func (s *templateService) SaveTo(f *os.File, templateName string, data any) error {
	tmpl, err := s.template.Lookup("layout").Clone()
	if err != nil {
		return err
	}
	tmpl, err = tmpl.AddParseTree("content", s.template.Lookup(templateName).Tree)
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(f, "layout", data)
}
