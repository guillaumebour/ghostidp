package domain

import "context"

type IdentityManager interface {
	RenderLoginPage(ctx context.Context, loginChallenge string) (RenderablePageFn, string, error)
	Login(ctx context.Context, loginChallenge string, username string) (RenderablePageFn, string, error)

	RenderConsentPage(ctx context.Context, consentChallenge string) (RenderablePageFn, string, error)
	Consent(ctx context.Context, consentChallenge string, consentGranted bool, grantedScopes []string) (RenderablePageFn, string, error)
}
