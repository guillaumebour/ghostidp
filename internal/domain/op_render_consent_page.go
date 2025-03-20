package domain

import (
	"context"
	"errors"
	"fmt"
	"io"
)

func RenderConsentPage(ctx context.Context, consentChallenge string, h HydraClient, ir IdentityRepository, tp TemplateProvider) (RenderablePageFn, string, error) {
	// Validate that the consent challenge is not empty
	if consentChallenge == "" {
		return nil, "", NewValidationError("consent_challenge", "Required")
	}

	// Fetch the consent request from Hydra
	consentReq, err := h.GetConsentRequest(consentChallenge)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get consent request: %w", err)
	}

	trustedClientSkipConsent := false
	if consentReq.Client.SkipConsent != nil {
		trustedClientSkipConsent = *consentReq.Client.SkipConsent
	}

	// Check if we can skip the consent
	// 1) Either the user has accepted the consent previously (consent.Skip)
	// 2) The client is trusted (consent.Client.SkipConsent)
	if consentReq.Skip || trustedClientSkipConsent {
		// Fetch the user
		identity, err := ir.FindIdentityByUsername(ctx, consentReq.Subject)
		if err != nil {
			if errors.Is(err, ErrIdentityRepositoryIdentityNotFound) {
				return nil, "", NewNotFoundError("user does not exist")
			}
			return nil, "", fmt.Errorf("failed to find identity: %w", err)
		}

		res, err := h.AcceptConsent(consentChallenge, consentReq.RequestedScope, identity.AllClaims())
		if err != nil {
			return nil, "", fmt.Errorf("failed to accept consent: %w", err)
		}

		return nil, res.RedirectTo, nil
	}

	consentTmpl, err := tp.GetTemplate(ConsentPage)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get template: %w", err)
	}

	clientName := "Unknown"
	if consentReq.Client.ClientName != nil {
		clientName = *consentReq.Client.ClientName
	}

	// Render the consent page
	return func(w io.Writer) error {
		return consentTmpl.Execute(w, tp.BuildTemplateEnv(map[string]any{
			"Challenge":  consentChallenge,
			"ClientName": clientName,
			"Scopes":     consentReq.RequestedScope,
			"User":       consentReq.Subject,
		}))
	}, "", nil

}
