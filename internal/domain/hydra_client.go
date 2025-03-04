package domain

import "time"

type OIDCContext struct {
	AcrValues         []string       `json:"acr_values,omitempty"`
	Display           *string        `json:"display,omitempty"`
	IDTokenHintClaims map[string]any `json:"id_token_hint_claims,omitempty"`
	LoginHint         *string        `json:"login_hint,omitempty"`
	UiLocales         []string       `json:"ui_locales,omitempty"`
}

type OAuth2Client struct {
	AccessTokenStrategy                        *string    `json:"access_token_strategy,omitempty"`
	AllowedCorsOrigins                         []string   `json:"allowed_cors_origins,omitempty"`
	Audience                                   []string   `json:"audience,omitempty"`
	AuthorizationCodeGrantAccessTokenLifespan  *string    `json:"authorization_code_grant_access_token_lifespan,omitempty"`
	AuthorizationCodeGrantIdTokenLifespan      *string    `json:"authorization_code_grant_id_token_lifespan,omitempty"`
	AuthorizationCodeGrantRefreshTokenLifespan *string    `json:"authorization_code_grant_refresh_token_lifespan,omitempty"`
	BackchannelLogoutSessionRequired           *bool      `json:"backchannel_logout_session_required,omitempty"`
	BackchannelLogoutUri                       *string    `json:"backchannel_logout_uri,omitempty"`
	ClientCredentialsGrantAccessTokenLifespan  *string    `json:"client_credentials_grant_access_token_lifespan,omitempty"`
	ClientId                                   *string    `json:"client_id,omitempty"`
	ClientName                                 *string    `json:"client_name,omitempty"`
	ClientSecret                               *string    `json:"client_secret,omitempty"`
	ClientSecretExpiresAt                      *int64     `json:"client_secret_expires_at,omitempty"`
	ClientUri                                  *string    `json:"client_uri,omitempty"`
	Contacts                                   []string   `json:"contacts,omitempty"`
	CreatedAt                                  *time.Time `json:"created_at,omitempty"`
	FrontchannelLogoutSessionRequired          *bool      `json:"frontchannel_logout_session_required,omitempty"`
	FrontchannelLogoutUri                      *string    `json:"frontchannel_logout_uri,omitempty"`
	GrantTypes                                 []string   `json:"grant_types,omitempty"`
	ImplicitGrantAccessTokenLifespan           *string    `json:"implicit_grant_access_token_lifespan,omitempty"`
	ImplicitGrantIdTokenLifespan               *string    `json:"implicit_grant_id_token_lifespan,omitempty"`
	Jwks                                       any        `json:"jwks,omitempty"`
	JwksUri                                    *string    `json:"jwks_uri,omitempty"`
	JwtBearerGrantAccessTokenLifespan          *string    `json:"jwt_bearer_grant_access_token_lifespan,omitempty"`
	LogoUri                                    *string    `json:"logo_uri,omitempty"`
	Metadata                                   any        `json:"metadata,omitempty"`
	Owner                                      *string    `json:"owner,omitempty"`
	PolicyUri                                  *string    `json:"policy_uri,omitempty"`
	PostLogoutRedirectUris                     []string   `json:"post_logout_redirect_uris,omitempty"`
	RedirectUris                               []string   `json:"redirect_uris,omitempty"`
	RefreshTokenGrantAccessTokenLifespan       *string    `json:"refresh_token_grant_access_token_lifespan,omitempty"`
	RefreshTokenGrantIdTokenLifespan           *string    `json:"refresh_token_grant_id_token_lifespan,omitempty"`
	RefreshTokenGrantRefreshTokenLifespan      *string    `json:"refresh_token_grant_refresh_token_lifespan,omitempty"`
	RegistrationAccessToken                    *string    `json:"registration_access_token,omitempty"`
	RegistrationClientUri                      *string    `json:"registration_client_uri,omitempty"`
	RequestObjectSigningAlg                    *string    `json:"request_object_signing_alg,omitempty"`
	RequestUris                                []string   `json:"request_uris,omitempty"`
	ResponseTypes                              []string   `json:"response_types,omitempty"`
	Scope                                      *string    `json:"scope,omitempty"`
	SectorIdentifierUri                        *string    `json:"sector_identifier_uri,omitempty"`
	SkipConsent                                *bool      `json:"skip_consent,omitempty"`
	SkipLogoutConsent                          *bool      `json:"skip_logout_consent,omitempty"`
	SubjectType                                *string    `json:"subject_type,omitempty"`
	TokenEndpointAuthMethod                    *string    `json:"token_endpoint_auth_method,omitempty"`
	TokenEndpointAuthSigningAlg                *string    `json:"token_endpoint_auth_signing_alg,omitempty"`
	TosUri                                     *string    `json:"tos_uri,omitempty"`
	UpdatedAt                                  *time.Time `json:"updated_at,omitempty"`
	UserinfoSignedResponseAlg                  *string    `json:"userinfo_signed_response_alg,omitempty"`
}

type LoginRequest struct {
	Challenge                    string       `json:"challenge"`
	Client                       OAuth2Client `json:"client"`
	OIDCContext                  *OIDCContext `json:"oidc_context,omitempty"`
	RequestUrl                   string       `json:"request_url"`
	RequestedAccessTokenAudience []string     `json:"requested_access_token_audience,omitempty"`
	RequestedScope               []string     `json:"requested_scope,omitempty"`
	SessionId                    *string      `json:"session_id,omitempty"`
	Skip                         bool         `json:"skip"`
	Subject                      string       `json:"subject"`
}

type ConsentRequest struct {
	RequestedAccessTokenAudience []string     `json:"requested_access_token_audience"`
	LoginChallenge               string       `json:"login_challenge"`
	Subject                      string       `json:"subject"`
	Amr                          []string     `json:"amr"`
	OidcContext                  OIDCContext  `json:"oidc_context"`
	DeviceChallengeId            string       `json:"device_challenge_id"`
	Skip                         bool         `json:"skip"`
	RequestUrl                   string       `json:"request_url"`
	Acr                          string       `json:"acr"`
	Context                      string       `json:"context"`
	Challenge                    string       `json:"challenge"`
	Client                       OAuth2Client `json:"client"`
	LoginSessionId               string       `json:"login_session_id"`
	RequestedScope               []string     `json:"requested_scope"`
}

type Session struct {
	AccessToken string `json:"access_token"`
	IdToken     string `json:"id_token"`
}
type AcceptConsentResponse struct {
	Context                  string    `json:"context"`
	GrantAccessTokenAudience []string  `json:"grant_access_token_audience"`
	GrantScope               []string  `json:"grant_scope"`
	HandledAt                time.Time `json:"handled_at"`
	Remember                 bool      `json:"remember"`
	RememberFor              int       `json:"remember_for"`
	Session                  Session   `json:"session"`
}

type RedirectToResponse struct {
	RedirectTo string `json:"redirect_to"`
}

type HydraClient interface {
	GetLoginRequest(loginChallenge string) (*LoginRequest, error)
	AcceptLogin(loginChallenge string, subject string) (*RedirectToResponse, error)

	GetConsentRequest(consentChallenge string) (*ConsentRequest, error)
	AcceptConsent(consentChallenge string, scopes []string, claims map[string]any) (*RedirectToResponse, error)
	RejectConsent(consentChallenge string, reason string) (*RedirectToResponse, error)
}
