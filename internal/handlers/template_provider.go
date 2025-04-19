package handlers

import "html/template"

type TemplateName string

const (
	LoginPage   TemplateName = "login-page"
	ConsentPage TemplateName = "consent-page"
)

type TemplateProviderError string

func (e TemplateProviderError) Error() string {
	return string(e)
}

const (
	TemplateProviderErrTemplateNotFound = TemplateProviderError("template not found")
)

type TemplateProvider interface {
	GetTemplate(template TemplateName) (*template.Template, error)
	BuildTemplateEnv(values map[string]any) map[string]any
}
