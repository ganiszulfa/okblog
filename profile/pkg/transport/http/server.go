package http

import (
	"context"
	"net/http"

	"github.com/ganis/okblog/profile/pkg/service"
	"github.com/gorilla/mux"
)

type Server struct {
	svc    service.Service
	router *mux.Router
}

func NewServer(svc service.Service) *Server {
	s := &Server{
		svc:    svc,
		router: mux.NewRouter(),
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) routes() {
	endpoints := MakeEndpoints(s.svc)

	s.router.HandleFunc("/profiles", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		req, err := DecodeCreateProfileRequest(context.Background(), r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := endpoints.CreateProfile.ServeHTTP(context.Background(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		EncodeResponse(context.Background(), w, response)
	}).Methods(http.MethodPost)

	s.router.HandleFunc("/profiles/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		switch r.Method {
		case http.MethodGet:
			response, err := endpoints.GetProfile.ServeHTTP(context.Background(), id)
			if err != nil {
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
			response, err := endpoints.UpdateProfile.ServeHTTP(context.Background(), req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			EncodeResponse(context.Background(), w, response)

		case http.MethodDelete:
			err := endpoints.DeleteProfile.ServeHTTP(context.Background(), id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)
}
