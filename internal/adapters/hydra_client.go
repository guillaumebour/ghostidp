package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/guillaumebour/ghostidp/internal/domain"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

const (
	getLoginRequestEndpoint   = "%s/oauth2/auth/requests/login?login_challenge=%s"
	acceptLoginEndpoint       = "%s/oauth2/auth/requests/login/accept?login_challenge=%s"
	getConsentRequestEndpoint = "%s/oauth2/auth/requests/consent?consent_challenge=%s"
	acceptConsentEndpoint     = "%s/oauth2/auth/requests/consent/accept?consent_challenge=%s"
	rejectConsentEndpoint     = "%s/oauth2/auth/requests/consent/reject?consent_challenge=%s"
)

type hydraClient struct {
	log      *logrus.Entry
	adminURL string
	client   *http.Client
}

type HydraClientParams struct {
	Log      *logrus.Logger
	AdminURL string
	client   *http.Client
}

func NewHydraClient(p *HydraClientParams) domain.HydraClient {
	httpClient := http.DefaultClient
	if p.client != nil {
		httpClient = p.client
	}
	return &hydraClient{
		log:      p.Log.WithField("category", "hydra-client"),
		adminURL: p.AdminURL,
		client:   httpClient,
	}
}

func (h *hydraClient) GetLoginRequest(loginChallenge string) (*domain.LoginRequest, error) {
	url := fmt.Sprintf(getLoginRequestEndpoint, h.adminURL, loginChallenge)
	resp, err := h.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		h.log.Errorf("failed to get login request: %s", responseBody)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var loginReq *domain.LoginRequest
	if err := json.NewDecoder(resp.Body).Decode(&loginReq); err != nil {
		return nil, err
	}

	return loginReq, nil
}

type acceptLoginRequest struct {
	Acr                       string   `json:"acr"`
	Amr                       []string `json:"amr"`
	Context                   string   `json:"context"`
	ExtendSessionLifespan     bool     `json:"extend_session_lifespan"`
	ForceSubjectIdentifier    string   `json:"force_subject_identifier"`
	IdentityProviderSessionId string   `json:"identity_provider_session_id"`
	Remember                  bool     `json:"remember"`
	RememberFor               int      `json:"remember_for"`
	Subject                   string   `json:"subject"`
}

func (h *hydraClient) AcceptLogin(loginChallenge string, subject string) (*domain.RedirectToResponse, error) {
	url := fmt.Sprintf(acceptLoginEndpoint, h.adminURL, loginChallenge)

	loginAccept := acceptLoginRequest{
		Subject:     subject,
		Remember:    false,
		RememberFor: 0,
	}

	body, err := json.Marshal(loginAccept)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response *domain.RedirectToResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

func (h *hydraClient) GetConsentRequest(consentChallenge string) (*domain.ConsentRequest, error) {
	url := fmt.Sprintf(getConsentRequestEndpoint, h.adminURL, consentChallenge)
	resp, err := h.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var consentReq *domain.ConsentRequest
	if err := json.NewDecoder(resp.Body).Decode(&consentReq); err != nil {
		return nil, err
	}

	return consentReq, nil
}

type SessionInfo struct {
	AccessToken map[string]any `json:"access_token"`
	IdToken     map[string]any `json:"id_token"`
}

type acceptConsentRequest struct {
	Context                  string      `json:"context"`
	GrantAccessTokenAudience []string    `json:"grant_access_token_audience"`
	GrantScope               []string    `json:"grant_scope"`
	HandledAt                time.Time   `json:"handled_at"`
	Remember                 bool        `json:"remember"`
	RememberFor              int         `json:"remember_for"`
	Session                  SessionInfo `json:"session"`
}

func (h *hydraClient) AcceptConsent(consentChallenge string, scopes []string, claims map[string]any) (*domain.RedirectToResponse, error) {
	url := fmt.Sprintf(acceptConsentEndpoint, h.adminURL, consentChallenge)

	consentAccept := acceptConsentRequest{
		GrantScope:  scopes,
		HandledAt:   time.Now(),
		Remember:    false,
		RememberFor: 0,
		Session: SessionInfo{
			AccessToken: claims,
			IdToken:     claims,
		},
	}

	body, err := json.Marshal(consentAccept)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response *domain.RedirectToResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

func (h *hydraClient) RejectConsent(consentChallenge string, reason string) (*domain.RedirectToResponse, error) {
	url := fmt.Sprintf(rejectConsentEndpoint, h.adminURL, consentChallenge)

	reject := map[string]string{
		"error":             "User rejected consent",
		"error_description": reason,
	}

	body, err := json.Marshal(reject)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response *domain.RedirectToResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}
