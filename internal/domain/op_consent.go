package domain

import (
	"context"
	"errors"
	"fmt"
)

func Consent(ctx context.Context, consentChallenge string, consentGranted bool, grantedScopes []string, ir IdentityRepository, h HydraClient) (string, error) {
	if consentChallenge == "" {
		return "", NewValidationError("consent_challenge", "Required")
	}

	// Fetch the consent request from Hydra
	consentReq, err := h.GetConsentRequest(consentChallenge)
	if err != nil {
		return "", fmt.Errorf("failed to get consent request: %w", err)
	}

	// If the user did not grant consent, we reject the consent
	if !consentGranted {
		redirect, err := h.RejectConsent(consentChallenge, "User denied access")
		if err != nil {
			return "", err
		}
		return redirect.RedirectTo, nil
	}

	// Fetch the identity
	identity, err := ir.FindIdentityByUsername(ctx, consentReq.Subject)
	if err != nil {
		if errors.Is(err, ErrIdentityRepositoryIdentityNotFound) {
			return "", NewNotFoundError(fmt.Sprintf("user %s not found", consentReq.Subject))
		}
		return "", err
	}

	// Accept the consent
	res, err := h.AcceptConsent(consentChallenge, grantedScopes, identity.AllClaims())
	if err != nil {
		return "", err
	}

	return res.RedirectTo, nil
}
