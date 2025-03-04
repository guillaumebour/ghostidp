package adapters

import (
	_ "embed"
	"github.com/guillaumebour/ghostidp/internal/domain"
	"html/template"
)

//go:embed templates/login.tmpl
var loginPageTemplate string

//go:embed templates/consent.tmpl
var consentPageTemplate string

type embeddedTemplateProvider struct {
	templates map[domain.TemplateName]*template.Template
}

func NewEmbeddedTemplateProvider() domain.TemplateProvider {
	templates := make(map[domain.TemplateName]*template.Template)
	templates[domain.LoginPage] = template.Must(template.New(string(domain.LoginPage)).Parse(loginPageTemplate))
	templates[domain.ConsentPage] = template.Must(template.New(string(domain.ConsentPage)).Parse(consentPageTemplate))

	return &embeddedTemplateProvider{
		templates: templates,
	}
}

func (e *embeddedTemplateProvider) GetTemplate(template domain.TemplateName) (*template.Template, error) {
	t, ok := e.templates[template]
	if !ok {
		return nil, domain.TemplateProviderErrTemplateNotFound
	}
	return t, nil
}
