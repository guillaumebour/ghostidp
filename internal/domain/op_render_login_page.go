package domain

import (
	"context"
	"fmt"
	"io"
)

type ValidationError struct {
	field string
	msg   string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.field, e.msg)
}

func NewValidationError(field string, msg string) error {
	return &ValidationError{
		field: field,
		msg:   msg,
	}
}

type NotFoundError struct {
	msg string
}

func (e NotFoundError) Error() string {
	return e.msg
}

func NewNotFoundError(msg string) error {
	return &NotFoundError{
		msg: msg,
	}
}

type RenderablePageFn func(w io.Writer) error

func RenderLoginPage(ctx context.Context, loginChallenge string, h HydraClient, ir IdentityRepository, tp TemplateProvider) (RenderablePageFn, string, error) {
	// Validate that the login challenge is not empty
	if loginChallenge == "" {
		return nil, "", NewValidationError("login_challenge", "Required")
	}

	// Fetch the login request from Hydra
	loginReq, err := h.GetLoginRequest(loginChallenge)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get login request: %w", err)
	}

	// Check if we can skip the login and accept the response directly
	if loginReq.Skip {
		res, err := h.AcceptLogin(loginChallenge, loginReq.Subject)
		if err != nil {
			return nil, "", fmt.Errorf("failed to accept login: %w", err)
		}
		return nil, res.RedirectTo, nil
	}

	// If not, we need to build and return the login screen
	identities, err := ir.ListIdentities(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to list identities: %w", err)
	}

	// Get the template for the login page
	loginTpl, err := tp.GetTemplate(LoginPage)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get template: %w", err)
	}

	return func(w io.Writer) error {
		return loginTpl.Execute(w, tp.BuildTemplateEnv(map[string]any{
			"Challenge": loginChallenge,
			"Users":     identities,
		}))
	}, "", nil
}
