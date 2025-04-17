package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ganis/okblog/profile/pkg/model"
	"github.com/ganis/okblog/profile/pkg/service"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateProfile endpoint.Endpoint
	GetProfile    endpoint.Endpoint
	UpdateProfile endpoint.Endpoint
	DeleteProfile endpoint.Endpoint
}

func MakeEndpoints(svc service.Service) Endpoints {
	return Endpoints{
		CreateProfile: makeCreateProfileEndpoint(svc),
		GetProfile:    makeGetProfileEndpoint(svc),
		UpdateProfile: makeUpdateProfileEndpoint(svc),
		DeleteProfile: makeDeleteProfileEndpoint(svc),
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
