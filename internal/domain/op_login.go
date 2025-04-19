package domain

import (
	"context"
	"errors"
	"fmt"
)

func Login(ctx context.Context, loginChallenge string, username string, ir IdentityRepository, h HydraClient) (string, error) {
	if loginChallenge == "" {
		return "", NewValidationError("login_challenge", "Required")
	}

	if username == "" {
		return "", NewValidationError("username", "Required")
	}

	// Authenticate the user
	identity, err := ir.FindIdentityByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, ErrIdentityRepositoryIdentityNotFound) {
			return "", NewNotFoundError(fmt.Sprintf("user %s does not exist", username))
		}
		return "", err
	}

	// Accept the Login Request
	redirect, err := h.AcceptLogin(loginChallenge, identity.Username)
	if err != nil {
		return "", fmt.Errorf("failed to accept login: %w", err)
	}

	return redirect.RedirectTo, nil
}
