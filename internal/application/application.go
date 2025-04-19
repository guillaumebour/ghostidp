package application

import (
	"context"
	"github.com/guillaumebour/ghostidp/internal/adapters"
	"github.com/guillaumebour/ghostidp/internal/domain"
	"github.com/guillaumebour/ghostidp/internal/service"
	"github.com/sirupsen/logrus"
)

type Application struct {
	IdentityManager domain.IdentityManager
	Log             *logrus.Logger
}

type Params struct {
	Log *logrus.Logger

	HydraAdminURL string
	UsersFile     string
}

func NewApplication(ctx context.Context, p *Params) (*Application, func(), error) {
	appCtx, cancelAppCtx := context.WithCancel(ctx)

	// Adapters
	hydraClient := adapters.NewHydraClient(&adapters.HydraClientParams{
		Log:      p.Log,
		AdminURL: p.HydraAdminURL,
	})
	identitiesRepo, err := adapters.NewIdentityFileRepository(appCtx, &adapters.IdentityFileRepositoryParams{
		Logger:         p.Log,
		ConfigFilepath: p.UsersFile,
	})
	if err != nil {
		return nil, cancelAppCtx, err
	}

	// Identity Manager
	identityMgr, err := service.NewIdentityManager(&service.IdentityManagerParams{
		Log:                p.Log,
		IdentityRepository: identitiesRepo,
		HydraClient:        hydraClient,
	})
	if err != nil {
		return nil, cancelAppCtx, err
	}

	return &Application{
		Log:             p.Log,
		IdentityManager: identityMgr,
	}, cancelAppCtx, nil
}
