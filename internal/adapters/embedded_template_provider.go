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
	templates    map[domain.TemplateName]*template.Template
	templatesEnv map[string]any
}

type EmbeddedTemplateProviderParams struct {
	Badge       string
	Version     string
	AccentColor string
}

func NewEmbeddedTemplateProvider(p *EmbeddedTemplateProviderParams) domain.TemplateProvider {
	templates := make(map[domain.TemplateName]*template.Template)
	templates[domain.LoginPage] = template.Must(template.New(string(domain.LoginPage)).Parse(loginPageTemplate))
	templates[domain.ConsentPage] = template.Must(template.New(string(domain.ConsentPage)).Parse(consentPageTemplate))

	return &embeddedTemplateProvider{
		templates: templates,
		templatesEnv: map[string]any{
			"BadgeContent": p.Badge,
			"Version":      p.Version,
			"AccentColor":  p.AccentColor,
		},
	}
}

func (e *embeddedTemplateProvider) GetTemplate(template domain.TemplateName) (*template.Template, error) {
	t, ok := e.templates[template]
	if !ok {
		return nil, domain.TemplateProviderErrTemplateNotFound
	}
	return t, nil
}

func (e *embeddedTemplateProvider) BuildTemplateEnv(values map[string]any) map[string]any {
	v := make(map[string]any)

	// Add the pre-defined values
	for key, value := range e.templatesEnv {
		v[key] = value
	}

	// Add the remaining values
	for key, value := range values {
		v[key] = value
	}

	return v
}
