package service

import (
	"github.com/guillaumebour/ghostidp/internal/domain"
	"github.com/sirupsen/logrus"
)

type IdentityManagerParams struct {
	Log                *logrus.Logger
	IdentityRepository domain.IdentityRepository
	TemplateProvider   domain.TemplateProvider
	HydraClient        domain.HydraClient
}

func NewIdentityManager(p *IdentityManagerParams) (domain.IdentityManager, error) {
	return &identityManager{
		log: p.Log,
		ir:  p.IdentityRepository,
		tp:  p.TemplateProvider,
		hc:  p.HydraClient,
	}, nil
}
