package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ganis/okblog/profile/pkg/model"
	"github.com/ganis/okblog/profile/pkg/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Endpoints struct {
	CreateProfile endpoint.Endpoint
	GetProfile    endpoint.Endpoint
	UpdateProfile endpoint.Endpoint
	DeleteProfile endpoint.Endpoint
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
		CreateProfile: loggingMiddleware(makeCreateProfileEndpoint(svc)),
		GetProfile:    loggingMiddleware(makeGetProfileEndpoint(svc)),
		UpdateProfile: loggingMiddleware(makeUpdateProfileEndpoint(svc)),
		DeleteProfile: loggingMiddleware(makeDeleteProfileEndpoint(svc)),
	}
}

func makeCreateProfileEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.CreateProfileRequest)
		profile, err := svc.CreateProfile(ctx, req)
		if err != nil {
			return nil, err
		}
		return profile, nil
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

func DecodeCreateProfileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req model.CreateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
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
