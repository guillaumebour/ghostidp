package domain

import "fmt"

type Identity struct {
	Username     string
	Password     string
	Email        string
	GivenName    string
	FamilyName   string
	Description  string
	CustomClaims map[string]any
}

func (i *Identity) StandardClaims() map[string]any {
	claims := map[string]any{
		"sub":         i.Username,
		"email":       i.Email,
		"given_name":  i.GivenName,
		"family_name": i.FamilyName,
		"name":        fmt.Sprintf("%s %s", i.GivenName, i.FamilyName),
	}

	return claims
}

func (i *Identity) AllClaims() map[string]any {
	claims := i.StandardClaims()

	for k, v := range i.CustomClaims {
		claims[k] = v
	}

	return claims
}
