package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ganis/okblog/profile/pkg/model"
	"github.com/go-kit/log"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of repository.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateProfile(ctx context.Context, profile model.Profile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockRepository) GetProfile(ctx context.Context, id string) (*model.Profile, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Profile), args.Error(1)
}

func (m *MockRepository) UpdateProfile(ctx context.Context, profile model.Profile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockRepository) DeleteProfile(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateProfile(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockRepository)
	// Create a noop logger
	logger := log.NewNopLogger()
	// Create a service with the mock repository
	svc := NewService(mockRepo, logger)

	// Setup test data
	req := model.CreateProfileRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Bio:       "This is a test user",
	}

	// Setup expectations
	mockRepo.On("CreateProfile", mock.Anything, mock.MatchedBy(func(p model.Profile) bool {
		return p.Username == req.Username &&
			p.Email == req.Email &&
			p.FirstName == req.FirstName &&
			p.LastName == req.LastName &&
			p.Bio == req.Bio
	})).Return(nil)

	// Call the method
	profile, err := svc.CreateProfile(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, req.Username, profile.Username)
	assert.Equal(t, req.Email, profile.Email)
	assert.Equal(t, req.FirstName, profile.FirstName)
	assert.Equal(t, req.LastName, profile.LastName)
	assert.Equal(t, req.Bio, profile.Bio)
	mockRepo.AssertExpectations(t)
}

func TestCreateProfile_InvalidInput(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockRepository)
	// Create a noop logger
	logger := log.NewNopLogger()
	// Create a service with the mock repository
	svc := NewService(mockRepo, logger)

	// Test cases for invalid input
	testCases := []struct {
		name string
		req  model.CreateProfileRequest
	}{
		{
			name: "Empty Username",
			req: model.CreateProfileRequest{
				Username:  "",
				Email:     "test@example.com",
				FirstName: "Test",
				LastName:  "User",
			},
		},
		{
			name: "Empty Email",
			req: model.CreateProfileRequest{
				Username:  "testuser",
				Email:     "",
				FirstName: "Test",
				LastName:  "User",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the method
			profile, err := svc.CreateProfile(context.Background(), tc.req)

			// Assertions
			assert.Error(t, err)
			assert.Equal(t, ErrInvalidInput, err)
			assert.Nil(t, profile)
		})
	}
}

func TestGetProfile(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockRepository)
	// Create a noop logger
	logger := log.NewNopLogger()
	// Create a service with the mock repository
	svc := NewService(mockRepo, logger)

	// Create a test profile
	id := uuid.New().String()
	profileData := &model.Profile{
		ID:        id,
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Bio:       "This is a test user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Setup expectations
	mockRepo.On("GetProfile", mock.Anything, id).Return(profileData, nil)

	// Call the method
	profile, err := svc.GetProfile(context.Background(), id)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, profileData, profile)
	mockRepo.AssertExpectations(t)
}

func TestGetProfile_NotFound(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockRepository)
	// Create a noop logger
	logger := log.NewNopLogger()
	// Create a service with the mock repository
	svc := NewService(mockRepo, logger)

	// Setup expectations
	id := "non-existent-id"
	mockRepo.On("GetProfile", mock.Anything, id).Return(nil, errors.New("profile not found"))

	// Call the method
	profile, err := svc.GetProfile(context.Background(), id)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, ErrProfileNotFound, err)
	assert.Nil(t, profile)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProfile(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockRepository)
	// Create a noop logger
	logger := log.NewNopLogger()
	// Create a service with the mock repository
	svc := NewService(mockRepo, logger)

	// Create a test profile
	id := uuid.New().String()
	profileData := &model.Profile{
		ID:        id,
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Bio:       "Original bio",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create update request
	updateReq := model.UpdateProfileRequest{
		FirstName: "Updated",
		LastName:  "Name",
		Bio:       "Updated bio",
	}

	// Setup expectations
	mockRepo.On("GetProfile", mock.Anything, id).Return(profileData, nil)
	mockRepo.On("UpdateProfile", mock.Anything, mock.MatchedBy(func(p model.Profile) bool {
		return p.ID == id &&
			p.FirstName == updateReq.FirstName &&
			p.LastName == updateReq.LastName &&
			p.Bio == updateReq.Bio
	})).Return(nil)

	// Call the method
	updatedProfile, err := svc.UpdateProfile(context.Background(), id, updateReq)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, updatedProfile)
	assert.Equal(t, updateReq.FirstName, updatedProfile.FirstName)
	assert.Equal(t, updateReq.LastName, updatedProfile.LastName)
	assert.Equal(t, updateReq.Bio, updatedProfile.Bio)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProfile_NotFound(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockRepository)
	// Create a noop logger
	logger := log.NewNopLogger()
	// Create a service with the mock repository
	svc := NewService(mockRepo, logger)

	// Setup expectations
	id := "non-existent-id"
	updateReq := model.UpdateProfileRequest{
		FirstName: "Updated",
		LastName:  "Name",
		Bio:       "Updated bio",
	}
	mockRepo.On("GetProfile", mock.Anything, id).Return(nil, errors.New("profile not found"))

	// Call the method
	profile, err := svc.UpdateProfile(context.Background(), id, updateReq)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, ErrProfileNotFound, err)
	assert.Nil(t, profile)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProfile(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockRepository)
	// Create a noop logger
	logger := log.NewNopLogger()
	// Create a service with the mock repository
	svc := NewService(mockRepo, logger)

	// Setup expectations
	id := uuid.New().String()
	mockRepo.On("DeleteProfile", mock.Anything, id).Return(nil)

	// Call the method
	err := svc.DeleteProfile(context.Background(), id)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProfile_NotFound(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockRepository)
	// Create a noop logger
	logger := log.NewNopLogger()
	// Create a service with the mock repository
	svc := NewService(mockRepo, logger)

	// Setup expectations
	id := "non-existent-id"
	mockRepo.On("DeleteProfile", mock.Anything, id).Return(errors.New("profile not found"))

	// Call the method
	err := svc.DeleteProfile(context.Background(), id)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, ErrProfileNotFound, err)
	mockRepo.AssertExpectations(t)
}
