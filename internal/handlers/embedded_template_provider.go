package handlers

import (
	_ "embed"
	"html/template"
	"strings"
)

//go:embed templates/login.tmpl
var loginPageTemplate string

//go:embed templates/consent.tmpl
var consentPageTemplate string

type embeddedTemplateProvider struct {
	templates    map[TemplateName]*template.Template
	templatesEnv map[string]any
}

type EmbeddedTemplateProviderParams struct {
	HeaderText    string
	HeaderLogoURL string
	Badge         string
	Version       string
	AccentColor   string
}

func NewEmbeddedTemplateProvider(p *EmbeddedTemplateProviderParams) TemplateProvider {
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
	}

	templates := make(map[TemplateName]*template.Template)
	templates[LoginPage] = template.Must(template.New(string(LoginPage)).Funcs(funcMap).Parse(loginPageTemplate))
	templates[ConsentPage] = template.Must(template.New(string(ConsentPage)).Funcs(funcMap).Parse(consentPageTemplate))

	return &embeddedTemplateProvider{
		templates: templates,
		templatesEnv: map[string]any{
			"HeaderText":    p.HeaderText,
			"HeaderLogoURL": p.HeaderLogoURL,
			"BadgeContent":  p.Badge,
			"Version":       p.Version,
			"AccentColor":   p.AccentColor,
		},
	}
}

func (e *embeddedTemplateProvider) GetTemplate(template TemplateName) (*template.Template, error) {
	t, ok := e.templates[template]
	if !ok {
		return nil, TemplateProviderErrTemplateNotFound
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
