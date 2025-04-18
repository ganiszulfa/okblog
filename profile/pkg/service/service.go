package service

import (
	"context"
	"errors"
	"time"

	"github.com/ganis/okblog/profile/pkg/model"
	"github.com/ganis/okblog/profile/pkg/repository"
	"github.com/go-kit/log"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrProfileNotFound    = errors.New("profile not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrHashingFailed      = errors.New("password hashing failed")
)

// Service defines the interface for profile operations
type Service interface {
	RegisterProfile(ctx context.Context, req model.RegisterProfileRequest) (*model.Profile, error)
	Login(ctx context.Context, req model.LoginRequest) (*model.Profile, error)
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

func (s *profileService) RegisterProfile(ctx context.Context, req model.RegisterProfileRequest) (*model.Profile, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, ErrInvalidInput
	}

	// Hash the password with bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Log("err", err, "msg", "Failed to hash password")
		return nil, ErrHashingFailed
	}

	now := time.Now()
	profile := model.Profile{
		ID:        uuid.New().String(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword), // Store the hashed password
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Bio:       req.Bio,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save to the repository
	err = s.repo.CreateProfile(ctx, profile)
	if err != nil {
		return nil, err
	}

	// Don't return the password in the response
	profile.Password = ""

	return &profile, nil
}

func (s *profileService) Login(ctx context.Context, req model.LoginRequest) (*model.Profile, error) {
	if req.Username == "" || req.Password == "" {
		return nil, ErrInvalidInput
	}

	// Get profile by username
	profile, err := s.repo.GetProfileByUsername(ctx, req.Username)
	if err != nil {
		if err.Error() == "profile not found" {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(profile.Password), []byte(req.Password))
	if err != nil {
		s.logger.Log("err", err, "msg", "Password comparison failed")
		return nil, ErrInvalidCredentials
	}

	// Don't return the password in the response
	profile.Password = ""

	return profile, nil
}

func (s *profileService) GetProfile(ctx context.Context, id string) (*model.Profile, error) {
	profile, err := s.repo.GetProfile(ctx, id)
	if err != nil {
		if err.Error() == "profile not found" {
			return nil, ErrProfileNotFound
		}
		return nil, err
	}

	// Don't return the password in the response
	profile.Password = ""

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

	// Don't return the password in the response
	profile.Password = ""

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
