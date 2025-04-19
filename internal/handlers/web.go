package handlers

import (
	"embed"
	"fmt"
	chilogrus "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/guillaumebour/ghostidp/internal/application"
	"github.com/guillaumebour/ghostidp/internal/domain"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

const (
	loginRoute   = "/login"
	consentRoute = "/consent"
)

type server struct {
	log *logrus.Entry
	im  domain.IdentityManager
	tp  TemplateProvider
}

//go:embed assets
var assets embed.FS

// CreateWebServer creates the Chi Router of the web application
func CreateWebServer(app *application.Application, tp TemplateProvider) (*chi.Mux, error) {
	s := &server{
		log: app.Log.WithField("component", "webserver"),
		im:  app.IdentityManager,
		tp:  tp,
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(chilogrus.Logger("web", app.Log))
	r.Use(middleware.Recoverer)

	// Healthcheck
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Static assets
	r.Handle("/assets/*", http.FileServer(http.FS(assets)))

	// Application
	r.Get(loginRoute, s.handleGetLogin())
	r.Post(loginRoute, s.handlePostLogin())
	r.Get(consentRoute, s.handleGetConsent())
	r.Post(consentRoute, s.handlePostConsent())

	return r, nil
}

func (a *server) renderPage(w http.ResponseWriter, page TemplateName, data map[string]any) {
	loginTpl, err := a.tp.GetTemplate(page)
	if err != nil {
		a.log.WithError(err).Errorf("failed to get template")
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if err = loginTpl.Execute(w, a.tp.BuildTemplateEnv(data)); err != nil {
		a.log.WithError(err).Errorf("failed to execute template")
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}

func (a *server) handleError(w http.ResponseWriter, err error) {
	errCode := http.StatusInternalServerError
	msg := "Something went wrong"

	switch {
	case domain.IsNotFoundError(err):
		errCode = http.StatusNotFound
		msg = "Not found"
	case domain.IsValidationError(err):
		errCode = http.StatusBadRequest
		msg = "Validation error: " + err.Error()
	}

	http.Error(w, msg, errCode)
}

func (a *server) handleGetLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get login_challenge from query
		loginChallenge := r.URL.Query().Get("login_challenge")

		// Get login information from domain
		loginInformation, err := a.im.GetLoginInformation(r.Context(), loginChallenge)
		if err != nil {
			a.log.WithError(err).Error("failed to get login information")
			a.handleError(w, err)
			return
		}

		// If the RedirectURL is provided, we don't need to render the login page
		if loginInformation.RedirectURL != "" {
			http.Redirect(w, r, loginInformation.RedirectURL, http.StatusFound)
			return
		}

		a.renderPage(w, LoginPage, map[string]any{
			"Challenge": loginChallenge,
			"Users":     loginInformation.Identities,
		})
	}
}

func (a *server) handleGetConsent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get consent_challenge from query
		consentChallenge := r.URL.Query().Get("consent_challenge")

		// Get consent information from domain
		consentInformation, err := a.im.GetConsentInformation(r.Context(), consentChallenge)
		if err != nil {
			a.log.WithError(err).Error("error rendering consent page")
			a.handleError(w, err)
			return
		}

		// If the RedirectURL is provided, we don't need to render the consent page
		if consentInformation.RedirectURL != "" {
			http.Redirect(w, r, consentInformation.RedirectURL, http.StatusFound)
			return
		}

		a.renderPage(w, ConsentPage, map[string]any{
			"Challenge":  consentChallenge,
			"ClientName": consentInformation.ClientName,
			"Scopes":     consentInformation.Scopes,
			"User":       consentInformation.User,
		})
	}
}

func (a *server) handlePostLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get consent_challenge from query
		loginChallenge := r.URL.Query().Get("login_challenge")

		// Get username from body
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		redirect, err := a.im.Login(r.Context(), loginChallenge, r.Form.Get("username"))
		if err != nil {
			a.log.WithError(err).Error("error performing login")
			a.handleError(w, err)
			return
		}

		if redirect != "" {
			http.Redirect(w, r, redirect, http.StatusFound)
			return
		}

		// safeguard
		a.handleError(w, fmt.Errorf("missing redirect URL"))
	}
}

func (a *server) handlePostConsent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get consent_challenge from query
		consentChallenge := r.URL.Query().Get("consent_challenge")

		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form data", http.StatusBadRequest)
			return
		}

		grantedAccess := r.Form.Get("consent") == "allow"

		// Fetch all "scope_" values that are set to "on" from the body
		var grantedScopes []string
		for key := range r.Form {
			if strings.HasPrefix(key, "scope_") {
				if r.Form.Get(key) == "on" {
					scope := strings.Replace(key, "scope_", "", -1)
					grantedScopes = append(grantedScopes, scope)
				}
			}
		}

		redirect, err := a.im.Consent(r.Context(), consentChallenge, grantedAccess, grantedScopes)
		if err != nil {
			a.log.WithError(err).Error("error accepting consent")
			a.handleError(w, err)
			return
		}

		if redirect != "" {
			http.Redirect(w, r, redirect, http.StatusFound)
			return
		}

		// safeguard
		a.handleError(w, fmt.Errorf("missing redirect URL"))
	}
}
