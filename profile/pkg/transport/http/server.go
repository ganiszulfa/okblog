package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/ganis/okblog/profile/pkg/service"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Server struct {
	svc      service.Service
	router   *mux.Router
	logger   log.Logger
	newRelic *newrelic.Application
}

func NewServer(svc service.Service, logger log.Logger, newRelicApp *newrelic.Application) *Server {
	s := &Server{
		svc:      svc,
		router:   mux.NewRouter(),
		logger:   logger,
		newRelic: newRelicApp,
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler = s.router

	handler = LoggingMiddleware(s.logger)(handler)

	if s.newRelic != nil {
		handler = NewRelicMiddleware(s.newRelic, s.logger)(handler)
	}

	handler.ServeHTTP(w, r)
}

func (s *Server) routes() {
	endpoints := MakeEndpoints(s.svc, s.logger)

	s.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		EncodeResponse(context.Background(), w, "ok")
	}).Methods(http.MethodGet, http.MethodOptions)

	s.router.HandleFunc("/api/profiles/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost && r.Method != http.MethodOptions {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Skip processing for OPTIONS requests
		if r.Method == http.MethodOptions {
			return
		}

		req, err := DecodeRegisterProfileRequest(context.Background(), r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := endpoints.RegisterProfile(context.Background(), req)
		if err != nil {
			if err == service.ErrInvalidInput {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		EncodeResponse(context.Background(), w, response)
	}).Methods(http.MethodPost, http.MethodOptions)

	s.router.HandleFunc("/api/profiles/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost && r.Method != http.MethodOptions {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Skip processing for OPTIONS requests
		if r.Method == http.MethodOptions {
			return
		}

		req, err := DecodeLoginRequest(context.Background(), r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := endpoints.Login(context.Background(), req)
		if err != nil {
			if err == service.ErrInvalidInput || err == service.ErrInvalidCredentials {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		EncodeResponse(context.Background(), w, response)
	}).Methods(http.MethodPost, http.MethodOptions)

	s.router.HandleFunc("/api/profiles/validate-token", func(w http.ResponseWriter, r *http.Request) {
		// Skip processing for OPTIONS requests
		if r.Method == http.MethodOptions {
			return
		}

		req, err := DecodeValidateTokenRequest(context.Background(), r)
		if err != nil {
			// Handle Authorization header errors with 401 Unauthorized
			if strings.Contains(err.Error(), "Authorization header") || strings.Contains(err.Error(), "token") {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := endpoints.ValidateToken(context.Background(), req)
		if err != nil {
			if err == service.ErrInvalidInput {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			// Return 401 for invalid token errors
			if strings.Contains(err.Error(), "unauthorized: invalid token") {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		EncodeResponse(context.Background(), w, response)
	}).Methods(http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodOptions)

	s.router.HandleFunc("/api/profiles/{id}", func(w http.ResponseWriter, r *http.Request) {
		// Skip processing for OPTIONS requests
		if r.Method == http.MethodOptions {
			return
		}

		vars := mux.Vars(r)
		id := vars["id"]

		switch r.Method {
		case http.MethodGet:
			response, err := endpoints.GetProfile(context.Background(), id)
			if err != nil {
				if err == service.ErrProfileNotFound {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			EncodeResponse(context.Background(), w, response)

		case http.MethodPut:
			req, err := DecodeUpdateProfileRequest(context.Background(), r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			response, err := endpoints.UpdateProfile(context.Background(), req)
			if err != nil {
				if err == service.ErrProfileNotFound {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			EncodeResponse(context.Background(), w, response)

		case http.MethodDelete:
			_, err := endpoints.DeleteProfile(context.Background(), id)
			if err != nil {
				if err == service.ErrProfileNotFound {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}).Methods(http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodOptions)
}
