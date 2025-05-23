package adapters

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/guillaumebour/ghostidp/internal/domain"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

type DisplayOptions struct {
	Description string `yaml:"description"`
	AvatarText  string `yaml:"avatar_text"`
	AvatarColor string `yaml:"avatar_color"`
}

type fileIdentity struct {
	Username     string                 `yaml:"username"`
	Password     string                 `yaml:"password"`
	Email        string                 `yaml:"email"`
	GivenName    string                 `yaml:"given_name"`
	FamilyName   string                 `yaml:"family_name"`
	Display      *DisplayOptions        `yaml:"display"`
	CustomClaims map[string]interface{} `yaml:"custom_claims"`
}

type identitiesConfig struct {
	Users []*fileIdentity `yaml:"users"`
}

type identityFileRepository struct {
	log *logrus.Entry

	// Identities
	iLock           sync.RWMutex
	configFilepath  string
	identitiesAsMap map[string]*fileIdentity
	identities      []*fileIdentity
}

type IdentityFileRepositoryParams struct {
	Logger         *logrus.Logger
	ConfigFilepath string
}

func NewIdentityFileRepository(ctx context.Context, p *IdentityFileRepositoryParams) (domain.IdentityRepository, error) {
	repo := &identityFileRepository{
		log: p.Logger.WithFields(logrus.Fields{
			"category": "identities-file-repository",
		}),

		iLock:           sync.RWMutex{},
		identitiesAsMap: make(map[string]*fileIdentity),
		identities:      make([]*fileIdentity, 0),
		configFilepath:  p.ConfigFilepath,
	}

	// load users for the first time
	if err := repo.loadIdentities(); err != nil {
		return nil, err
	}

	// Start a watcher to reload identities when the file is updated
	if err := repo.watchIdentitiesFile(ctx); err != nil {
		return nil, err
	}

	return repo, nil
}

func (i *identityFileRepository) loadIdentities() error {
	i.iLock.Lock()
	defer i.iLock.Unlock()

	data, err := os.ReadFile(i.configFilepath)
	if err != nil {
		return fmt.Errorf("failed to read identities config: %w", err)
	}

	var config identitiesConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse identities config: %w", err)
	}

	// Reset and reload users
	i.identities = config.Users
	i.identitiesAsMap = make(map[string]*fileIdentity)
	for _, user := range config.Users {
		i.identitiesAsMap[user.Username] = user
	}

	i.log.Infof("loaded %d identities", len(i.identitiesAsMap))

	return nil
}

func (i *identityFileRepository) reloadIdentities() error {
	if err := i.loadIdentities(); err != nil {
		// We log the error but do not return it
		i.log.Errorf("failed to reload identities: %v", err)
	}
	return nil
}

func (i *identityFileRepository) watchIdentitiesFile(ctx context.Context) error {
	// Creating a file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	// Start a dedicated go routine to watch changes to the users file
	// This goroutine will stop when the parent context is cancelled
	go func() {
		for {
			select {
			case <-ctx.Done():
				i.log.Infof("stop watching identities file")
				watcher.Close()
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					i.log.Debugf("changes detected to identities file %s", event.Name)
					_ = i.reloadIdentities()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				i.log.Errorf("error watching file: %v", err)
			}
		}
	}()

	if watcher.Add(i.configFilepath) != nil {
		return fmt.Errorf("failed to watch identities file %s", i.configFilepath)
	}

	return nil
}

func (i *identityFileRepository) FindIdentityByUsername(_ context.Context, username string) (*domain.Identity, error) {
	i.iLock.RLock()
	defer i.iLock.RUnlock()

	identity, exists := i.identitiesAsMap[username]
	if !exists {
		return nil, domain.ErrIdentityRepositoryIdentityNotFound
	}

	return fileIdentityToDomain(identity), nil
}

func (i *identityFileRepository) ListIdentities(_ context.Context) ([]*domain.Identity, error) {
	i.iLock.RLock()
	defer i.iLock.RUnlock()
	dIdentities := make([]*domain.Identity, 0, len(i.identitiesAsMap))
	for _, identity := range i.identities {
		dIdentities = append(dIdentities, fileIdentityToDomain(identity))
	}
	return dIdentities, nil
}

func fileIdentityToDomain(i *fileIdentity) *domain.Identity {
	var do *domain.DisplayOptions
	if i.Display != nil {
		do = &domain.DisplayOptions{
			Description: i.Display.Description,
			AvatarText:  i.Display.AvatarText,
			AvatarColor: i.Display.AvatarColor,
		}
	}

	return &domain.Identity{
		Username:     i.Username,
		Email:        i.Email,
		GivenName:    i.GivenName,
		FamilyName:   i.FamilyName,
		CustomClaims: i.CustomClaims,
		DisplayOpts:  do,
	}
}
