package handlers

import (
	"embed"
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
}

//go:embed assets
var assets embed.FS

// CreateWebServer creates the Chi Router of the web application
func CreateWebServer(app *application.Application) (*chi.Mux, error) {
	s := &server{
		log: app.Log.WithField("component", "webserver"),
		im:  app.IdentityManager,
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

func (a *server) handleGetLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get login_challenge from query
		loginChallenge := r.URL.Query().Get("login_challenge")

		// Get login page from domain
		renderingFn, redirect, err := a.im.RenderLoginPage(r.Context(), loginChallenge)
		if err != nil {
			a.log.WithError(err).Error("error rendering login page")
			http.Error(w, "something went wrong 1", http.StatusInternalServerError) // ToDo: handle other errors
			return
		}

		// We either have a renderingFn, or a redirect
		if redirect != "" {
			http.Redirect(w, r, redirect, http.StatusFound)
			return
		}

		if renderingFn != nil {
			if err := renderingFn(w); err != nil {
				a.log.WithError(err).Error("error rendering login page")
				http.Error(w, "something went wrong 2", http.StatusInternalServerError)
			}
			return
		}

		// We should not end up here
		http.Error(w, "something went wrong 3", http.StatusInternalServerError)
	}
}

func (a *server) handleGetConsent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get consent_challenge from query
		consentChallenge := r.URL.Query().Get("consent_challenge")

		// Get consent page from domain
		renderingFn, redirect, err := a.im.RenderConsentPage(r.Context(), consentChallenge)
		if err != nil {
			a.log.WithError(err).Error("error rendering consent page")
			http.Error(w, "something went wrong 4", http.StatusInternalServerError) // ToDo: handle other errors
			return
		}

		// We either have a renderingFn, or a redirect
		if redirect != "" {
			http.Redirect(w, r, redirect, http.StatusFound)
			return
		}

		if renderingFn != nil {
			if err := renderingFn(w); err != nil {
				a.log.WithError(err).Error("error rendering consent page")
				http.Error(w, "something went wrong 5", http.StatusInternalServerError)
			}
			return
		}

		// We should not end up here
		http.Error(w, "something went wrong 6", http.StatusInternalServerError)
	}
}

func (a *server) handlePostLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get consent_challenge from query
		loginChallenge := r.URL.Query().Get("login_challenge")

		// Get username from body
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form data", http.StatusBadRequest)
			return
		}

		username := r.Form.Get("username")

		renderingFn, redirect, err := a.im.Login(r.Context(), loginChallenge, username)
		if err != nil {
			a.log.WithError(err).Error("error login in")
			http.Error(w, "something went wrong 7", http.StatusInternalServerError) // ToDo: handle other errors
			return
		}

		// Happy path
		if redirect != "" {
			http.Redirect(w, r, redirect, http.StatusFound)
			return
		}

		if renderingFn != nil {
			if err := renderingFn(w); err != nil {
				a.log.WithError(err).Error("error rendering login page")
				http.Error(w, "something went wrong 8", http.StatusInternalServerError)
			}
			return
		}

		// We should not end up here
		http.Error(w, "something went wrong 9", http.StatusInternalServerError)
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

		renderingFn, redirect, err := a.im.Consent(r.Context(), consentChallenge, grantedAccess, grantedScopes)
		if err != nil {
			a.log.WithError(err).Error("error accepting consent")
			http.Error(w, "something went wrong 10", http.StatusInternalServerError) // ToDo: handle other errors
			return
		}

		// Happy path
		if redirect != "" {
			http.Redirect(w, r, redirect, http.StatusFound)
			return
		}

		if renderingFn != nil {
			if err := renderingFn(w); err != nil {
				a.log.WithError(err).Error("error rendering login page")
				http.Error(w, "something went wrong 11", http.StatusInternalServerError)
			}
			return
		}

		// We should not end up here
		http.Error(w, "something went wrong 12", http.StatusInternalServerError)
	}
}
