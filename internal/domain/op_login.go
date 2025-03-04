package domain

import (
	"context"
	"errors"
	"fmt"
	"io"
)

func Login(ctx context.Context, loginChallenge string, username string, ir IdentityRepository, h HydraClient, tp TemplateProvider) (RenderablePageFn, string, error) {
	if loginChallenge == "" {
		return nil, "", NewValidationError("login_challenge", "Required")
	}

	if username == "" {
		return nil, "", NewValidationError("username", "Required")
	}

	loginPageTemplate, err := tp.GetTemplate(LoginPage)
	if err != nil {
		return nil, "", err
	}

	// Authenticate the user
	identity, err := ir.FindIdentityByUsername(ctx, username)
	if err != nil {
		// If the user is not found we show the login page again with an error
		// Should not happen when the users are hard-coded
		if errors.Is(err, ErrIdentityRepositoryIdentityNotFound) {
			return func(w io.Writer) error {
				return loginPageTemplate.Execute(w, map[string]any{
					"Challenge": loginChallenge,
					"Error":     "User not found",
				})
			}, "", NewNotFoundError(fmt.Sprintf("user %s does not exist", username))
		}
		return nil, "", err
	}

	// Accept the Login Request
	redirect, err := h.AcceptLogin(loginChallenge, identity.Username)
	if err != nil {
		return nil, "", fmt.Errorf("failed to accept login: %w", err)
	}

	return nil, redirect.RedirectTo, nil
}
