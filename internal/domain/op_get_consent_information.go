package domain

import (
	"context"
	"errors"
	"fmt"
)

type ConsentInformationResponse struct {
	RedirectURL string
	ClientName  string
	Scopes      []string
	User        string
}

func GetConsentInformation(ctx context.Context, consentChallenge string, h HydraClient, ir IdentityRepository) (*ConsentInformationResponse, error) {
	// Validate that the consent challenge is not empty
	if consentChallenge == "" {
		return nil, NewValidationError("consent_challenge", "Required")
	}

	// Fetch the consent request from Hydra
	consentReq, err := h.GetConsentRequest(consentChallenge)
	if err != nil {
		return nil, fmt.Errorf("failed to get consent request: %w", err)
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
				return nil, NewNotFoundError("user does not exist")
			}
			return nil, fmt.Errorf("failed to find identity: %w", err)
		}

		res, err := h.AcceptConsent(consentChallenge, consentReq.RequestedScope, identity.AllClaims())
		if err != nil {
			return nil, fmt.Errorf("failed to accept consent: %w", err)
		}

		return &ConsentInformationResponse{
			RedirectURL: res.RedirectTo,
		}, nil
	}
	clientName := "Unknown"
	if consentReq.Client.ClientName != nil {
		clientName = *consentReq.Client.ClientName
	}

	return &ConsentInformationResponse{
		ClientName: clientName,
		Scopes:     consentReq.RequestedScope,
		User:       consentReq.Subject,
	}, nil
}
