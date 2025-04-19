package domain

import "context"

type IdentityManager interface {
	GetLoginInformation(ctx context.Context, loginChallenge string) (*LoginInformationResponse, error)
	Login(ctx context.Context, loginChallenge string, username string) (string, error)

	GetConsentInformation(ctx context.Context, consentChallenge string) (*ConsentInformationResponse, error)
	Consent(ctx context.Context, consentChallenge string, consentGranted bool, grantedScopes []string) (string, error)
}
