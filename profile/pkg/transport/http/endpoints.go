package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/ganis/okblog/profile/pkg/model"
	"github.com/ganis/okblog/profile/pkg/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Endpoints struct {
	RegisterProfile endpoint.Endpoint
	Login           endpoint.Endpoint
	ValidateToken   endpoint.Endpoint
	GetProfile      endpoint.Endpoint
	UpdateProfile   endpoint.Endpoint
	DeleteProfile   endpoint.Endpoint
}

// EndpointLoggingMiddleware returns an endpoint middleware that logs endpoint performance
func EndpointLoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				logger.Log(
					"endpoint_time", time.Since(begin),
					"err", err,
				)
			}(time.Now())
			return next(ctx, request)
		}
	}
}

func MakeEndpoints(svc service.Service, logger log.Logger) Endpoints {
	// Create a middleware for all endpoints
	loggingMiddleware := EndpointLoggingMiddleware(logger)

	return Endpoints{
		RegisterProfile: loggingMiddleware(makeRegisterProfileEndpoint(svc)),
		Login:           loggingMiddleware(makeLoginEndpoint(svc)),
		ValidateToken:   loggingMiddleware(makeValidateTokenEndpoint(svc)),
		GetProfile:      loggingMiddleware(makeGetProfileEndpoint(svc)),
		UpdateProfile:   loggingMiddleware(makeUpdateProfileEndpoint(svc)),
		DeleteProfile:   loggingMiddleware(makeDeleteProfileEndpoint(svc)),
	}
}

func makeRegisterProfileEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.RegisterProfileRequest)
		profile, err := svc.RegisterProfile(ctx, req)
		if err != nil {
			return nil, err
		}
		return profile, nil
	}
}

func makeLoginEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.LoginRequest)
		loginResponse, err := svc.Login(ctx, req)
		if err != nil {
			return nil, err
		}
		return loginResponse, nil
	}
}

func makeValidateTokenEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.TokenValidationRequest)
		claims, err := svc.ValidateToken(ctx, req.Token)
		if err != nil {
			if err == service.ErrInvalidToken {
				return nil, errors.New("unauthorized: invalid token")
			}
			return nil, err
		}
		return model.TokenValidationResponse{Valid: true, Claims: claims}, nil
	}
}

func makeGetProfileEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(string)
		profile, err := svc.GetProfile(ctx, id)
		if err != nil {
			return nil, err
		}
		return profile, nil
	}
}

func makeUpdateProfileEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(struct {
			ID   string                     `json:"id"`
			Data model.UpdateProfileRequest `json:"data"`
		})
		profile, err := svc.UpdateProfile(ctx, req.ID, req.Data)
		if err != nil {
			return nil, err
		}
		return profile, nil
	}
}

func makeDeleteProfileEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(string)
		err := svc.DeleteProfile(ctx, id)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

func DecodeRegisterProfileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req model.RegisterProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeValidateTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("missing Authorization header")
	}

	// Check if the header has the Bearer prefix
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		return nil, errors.New("invalid Authorization header format, expected 'Bearer token'")
	}

	token := tokenParts[1]
	if token == "" {
		return nil, errors.New("empty token in Authorization header")
	}

	return model.TokenValidationRequest{Token: token}, nil
}

func DecodeGetProfileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return r.URL.Query().Get("id"), nil
}

func DecodeUpdateProfileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req struct {
		ID   string                     `json:"id"`
		Data model.UpdateProfileRequest `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeDeleteProfileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return r.URL.Query().Get("id"), nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
