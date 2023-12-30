package service

type TemplateService interface {
	SaveTo(filePath string, templateName string, data any) error
}
