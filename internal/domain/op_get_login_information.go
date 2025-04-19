package domain

import (
	"context"
	"errors"
	"fmt"
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

func IsValidationError(err error) bool {
	var validationError *ValidationError
	ok := errors.As(err, &validationError)
	return ok
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

func IsNotFoundError(err error) bool {
	var notFoundErr *NotFoundError
	ok := errors.As(err, &notFoundErr)
	return ok
}

type LoginInformationResponse struct {
	RedirectURL string
	Identities  []*Identity
}

// GetLoginInformation builds a
func GetLoginInformation(ctx context.Context, loginChallenge string, h HydraClient, ir IdentityRepository) (*LoginInformationResponse, error) {
	// Validate that the login challenge is not empty
	if loginChallenge == "" {
		return nil, NewValidationError("login_challenge", "Required")
	}

	// Fetch the login request from Hydra
	loginReq, err := h.GetLoginRequest(loginChallenge)
	if err != nil {
		return nil, fmt.Errorf("failed to get login request: %w", err)
	}

	// Check if we can skip the login and accept the response directly
	if loginReq.Skip {
		res, err := h.AcceptLogin(loginChallenge, loginReq.Subject)
		if err != nil {
			return nil, fmt.Errorf("failed to accept login: %w", err)
		}
		return &LoginInformationResponse{
			RedirectURL: res.RedirectTo,
		}, nil
	}

	// If there is a login hint, we try to login the user directly
	if loginReq.OIDCContext != nil && loginReq.OIDCContext.LoginHint != nil {
		identity, err := ir.FindIdentityByUsername(ctx, *loginReq.OIDCContext.LoginHint)
		if err != nil {
			return nil, NewNotFoundError(fmt.Sprintf("failed to find identity for login hint %s", *loginReq.OIDCContext.LoginHint))
		}

		res, err := h.AcceptLogin(loginChallenge, identity.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to accept login: %w", err)
		}
		return &LoginInformationResponse{
			RedirectURL: res.RedirectTo,
		}, nil
	}

	// If not, we return the list of identities
	identities, err := ir.ListIdentities(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list identities: %w", err)
	}

	return &LoginInformationResponse{
		Identities: identities,
	}, nil
}
