package service

import (
	"context"
	"errors"
	"time"

	"github.com/ganis/okblog/profile/pkg/model"
	"github.com/ganis/okblog/profile/pkg/repository"
	"github.com/go-kit/log"
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
	repo   repository.Repository
	logger log.Logger
}

// NewService creates a new instance of the profile service
func NewService(repo repository.Repository, logger log.Logger) Service {
	return &profileService{
		repo:   repo,
		logger: logger,
	}
}

func (s *profileService) CreateProfile(ctx context.Context, req model.CreateProfileRequest) (*model.Profile, error) {
	if req.Username == "" || req.Email == "" {
		return nil, ErrInvalidInput
	}

	now := time.Now()
	profile := model.Profile{
		ID:        uuid.New().String(),
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Bio:       req.Bio,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save to the repository
	err := s.repo.CreateProfile(ctx, profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (s *profileService) GetProfile(ctx context.Context, id string) (*model.Profile, error) {
	profile, err := s.repo.GetProfile(ctx, id)
	if err != nil {
		if err.Error() == "profile not found" {
			return nil, ErrProfileNotFound
		}
		return nil, err
	}
	return profile, nil
}

func (s *profileService) UpdateProfile(ctx context.Context, id string, req model.UpdateProfileRequest) (*model.Profile, error) {
	// First fetch the profile
	profile, err := s.repo.GetProfile(ctx, id)
	if err != nil {
		if err.Error() == "profile not found" {
			return nil, ErrProfileNotFound
		}
		return nil, err
	}

	// Update the profile with new values
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

	// Save the updated profile
	err = s.repo.UpdateProfile(ctx, *profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *profileService) DeleteProfile(ctx context.Context, id string) error {
	err := s.repo.DeleteProfile(ctx, id)
	if err != nil {
		if err.Error() == "profile not found" {
			return ErrProfileNotFound
		}
		return err
	}
	return nil
}
