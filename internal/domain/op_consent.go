package domain

import (
	"context"
	"errors"
	"fmt"
)

func Consent(ctx context.Context, consentChallenge string, consentGranted bool, grantedScopes []string, ir IdentityRepository, h HydraClient) (RenderablePageFn, string, error) {
	if consentChallenge == "" {
		return nil, "", NewValidationError("consent_challenge", "Required")
	}

	// Fetch the consent request from Hydra
	consentReq, err := h.GetConsentRequest(consentChallenge)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get consent request: %w", err)
	}

	// If the user did not grant consent, we reject the consent
	if !consentGranted {
		redirect, err := h.RejectConsent(consentChallenge, "User denied access")
		if err != nil {
			return nil, "", err
		}

		return nil, redirect.RedirectTo, nil
	}

	// Fetch the identity
	identity, err := ir.FindIdentityByUsername(ctx, consentReq.Subject)
	if err != nil {
		// If the user is not found we show the login page again with an error
		// Should not happen when the users are hard-coded
		if errors.Is(err, ErrIdentityRepositoryIdentityNotFound) {
			return nil, "", NewNotFoundError(fmt.Sprintf("user %s does not exist", consentReq.Subject))
		}
		return nil, "", err
	}

	// Accept the consent
	res, err := h.AcceptConsent(consentChallenge, grantedScopes, identity.AllClaims())
	if err != nil {
		return nil, "", err
	}

	return nil, res.RedirectTo, nil
}
