package domain

import "context"

type IdentityRepositoryError string

func (e IdentityRepositoryError) Error() string {
	return string(e)
}

const (
	ErrIdentityRepositoryIdentityNotFound = IdentityRepositoryError("identity not found")
)

type IdentityRepository interface {
	FindIdentityByUsername(ctx context.Context, username string) (*Identity, error)
	ListIdentities(ctx context.Context) ([]*Identity, error)
}
