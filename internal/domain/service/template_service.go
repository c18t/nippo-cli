package service

import (
	"os"
)

type TemplateService interface {
	SaveTo(f *os.File, templateName string, data any) error
}
