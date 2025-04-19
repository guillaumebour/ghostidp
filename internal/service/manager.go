package service

import (
	"context"
	"github.com/guillaumebour/ghostidp/internal/domain"
	"github.com/sirupsen/logrus"
)

type identityManager struct {
	log *logrus.Logger
	ir  domain.IdentityRepository
	hc  domain.HydraClient
}

func (i *identityManager) GetLoginInformation(ctx context.Context, loginChallenge string) (*domain.LoginInformationResponse, error) {
	return domain.GetLoginInformation(ctx, loginChallenge, i.hc, i.ir)
}

func (i *identityManager) Login(ctx context.Context, loginChallenge string, username string) (string, error) {
	return domain.Login(ctx, loginChallenge, username, i.ir, i.hc)
}

func (i *identityManager) GetConsentInformation(ctx context.Context, consentChallenge string) (*domain.ConsentInformationResponse, error) {
	return domain.GetConsentInformation(ctx, consentChallenge, i.hc, i.ir)
}

func (i *identityManager) Consent(ctx context.Context, consentChallenge string, consentGranted bool, grantedScopes []string) (string, error) {
	return domain.Consent(ctx, consentChallenge, consentGranted, grantedScopes, i.ir, i.hc)
}
