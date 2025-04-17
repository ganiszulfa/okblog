package service

import (
	"context"
	"errors"
	"time"

	"github.com/ganis/okblog/profile/pkg/model"
	"github.com/google/uuid"
)

var (
	ErrProfileNotFound = errors.New("profile not found")
	ErrInvalidInput    = errors.New("invalid input")
)

// Service defines the interface for profile operations
type Service interface {
	CreateProfile(ctx context.Context, req model.CreateProfileRequest) (*model.Profile, error)
	GetProfile(ctx context.Context, id string) (*model.Profile, error)
	UpdateProfile(ctx context.Context, id string, req model.UpdateProfileRequest) (*model.Profile, error)
	DeleteProfile(ctx context.Context, id string) error
}

// profileService implements the Service interface
type profileService struct {
	// In a real implementation, this would have a repository/database client
	profiles map[string]*model.Profile
}

// NewService creates a new instance of the profile service
func NewService() Service {
	return &profileService{
		profiles: make(map[string]*model.Profile),
	}
}

func (s *profileService) CreateProfile(ctx context.Context, req model.CreateProfileRequest) (*model.Profile, error) {
	if req.Username == "" || req.Email == "" {
		return nil, ErrInvalidInput
	}

	now := time.Now()
	profile := &model.Profile{
		ID:        uuid.New().String(),
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Bio:       req.Bio,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.profiles[profile.ID] = profile
	return profile, nil
}

func (s *profileService) GetProfile(ctx context.Context, id string) (*model.Profile, error) {
	profile, exists := s.profiles[id]
	if !exists {
		return nil, ErrProfileNotFound
	}
	return profile, nil
}

func (s *profileService) UpdateProfile(ctx context.Context, id string, req model.UpdateProfileRequest) (*model.Profile, error) {
	profile, exists := s.profiles[id]
	if !exists {
		return nil, ErrProfileNotFound
	}

	if req.FirstName != "" {
		profile.FirstName = req.FirstName
	}
	if req.LastName != "" {
		profile.LastName = req.LastName
	}
	if req.Bio != "" {
		profile.Bio = req.Bio
	}

	profile.UpdatedAt = time.Now()
	s.profiles[id] = profile
	return profile, nil
}

func (s *profileService) DeleteProfile(ctx context.Context, id string) error {
	if _, exists := s.profiles[id]; !exists {
		return ErrProfileNotFound
	}

	delete(s.profiles, id)
	return nil
}
