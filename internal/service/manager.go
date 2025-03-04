package service

import (
	"context"
	"github.com/guillaumebour/ghostidp/internal/domain"
	"github.com/sirupsen/logrus"
)

type identityManager struct {
	log *logrus.Logger
	ir  domain.IdentityRepository
	tp  domain.TemplateProvider
	hc  domain.HydraClient
}

func (i *identityManager) RenderLoginPage(ctx context.Context, loginChallenge string) (domain.RenderablePageFn, string, error) {
	return domain.RenderLoginPage(ctx, loginChallenge, i.hc, i.ir, i.tp)
}

func (i *identityManager) Login(ctx context.Context, loginChallenge string, username string) (domain.RenderablePageFn, string, error) {
	return domain.Login(ctx, loginChallenge, username, i.ir, i.hc, i.tp)
}

func (i *identityManager) RenderConsentPage(ctx context.Context, consentChallenge string) (domain.RenderablePageFn, string, error) {
	return domain.RenderConsentPage(ctx, consentChallenge, i.hc, i.ir, i.tp)
}

func (i *identityManager) Consent(ctx context.Context, consentChallenge string, consentGranted bool, grantedScopes []string) (domain.RenderablePageFn, string, error) {
	return domain.Consent(ctx, consentChallenge, consentGranted, grantedScopes, i.ir, i.hc)
}
