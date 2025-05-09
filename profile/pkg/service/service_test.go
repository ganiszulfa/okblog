package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/ganis/okblog/profile/pkg/model"
	"github.com/go-kit/log"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
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

func (m *MockRepository) GetProfileByUsername(ctx context.Context, username string) (*model.Profile, error) {
	args := m.Called(ctx, username)
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

func (m *MockRepository) CountProfiles(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func TestRegisterProfile(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockRepository)
	// Create a noop logger
	logger := log.NewNopLogger()
	// Create a service with the mock repository
	svc := NewService(mockRepo, logger, false)

	// Setup test data
	req := model.RegisterProfileRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Bio:       "This is a test user",
	}

	// Setup expectations - we can't check exact password match since it's hashed
	mockRepo.On("CreateProfile", mock.Anything, mock.MatchedBy(func(p model.Profile) bool {
		return p.Username == req.Username &&
			p.Email == req.Email &&
			p.Password != req.Password && // Password should be hashed, not plaintext
			len(p.Password) > 0 && // Ensure password is not empty
			p.FirstName == req.FirstName &&
			p.LastName == req.LastName &&
			p.Bio == req.Bio
	})).Return(nil)

	// Call the method
	profile, err := svc.RegisterProfile(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, req.Username, profile.Username)
	assert.Equal(t, req.Email, profile.Email)
	assert.Equal(t, "", profile.Password) // Password should not be returned
	assert.Equal(t, req.FirstName, profile.FirstName)
	assert.Equal(t, req.LastName, profile.LastName)
	assert.Equal(t, req.Bio, profile.Bio)
	mockRepo.AssertExpectations(t)
}

func TestRegisterProfile_InvalidInput(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockRepository)
	// Create a noop logger
	logger := log.NewNopLogger()
	// Create a service with the mock repository
	svc := NewService(mockRepo, logger, false)

	// Test cases for invalid input
	testCases := []struct {
		name string
		req  model.RegisterProfileRequest
	}{
		{
			name: "Empty Username",
			req: model.RegisterProfileRequest{
				Username:  "",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
		},
		{
			name: "Empty Email",
			req: model.RegisterProfileRequest{
				Username:  "testuser",
				Email:     "",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
		},
		{
			name: "Empty Password",
			req: model.RegisterProfileRequest{
				Username:  "testuser",
				Email:     "test@example.com",
				Password:  "",
				FirstName: "Test",
				LastName:  "User",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the method
			profile, err := svc.RegisterProfile(context.Background(), tc.req)

			// Assertions
			assert.Error(t, err)
			assert.Equal(t, ErrInvalidInput, err)
			assert.Nil(t, profile)
		})
	}
}

func TestLogin(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockRepository)
	// Create a noop logger
	logger := log.NewNopLogger()
	// Create a service with the mock repository
	svc := NewService(mockRepo, logger, false)

	// Setup test data
	username := "testuser"
	password := "password123"
	id := uuid.New().String()

	// Generate a real bcrypt hash of the password for testing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err, "Password hashing should not error")

	profileData := &model.Profile{
		ID:        id,
		Username:  username,
		Email:     "test@example.com",
		Password:  string(hashedPassword), // Use the hashed password
		FirstName: "Test",
		LastName:  "User",
		Bio:       "This is a test user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	loginReq := model.LoginRequest{
		Username: username,
		Password: password, // Plain text password for login request
	}

	// Setup expectations
	mockRepo.On("GetProfileByUsername", mock.Anything, username).Return(profileData, nil)

	// Call the method
	loginResponse, err := svc.Login(context.Background(), loginReq)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, loginResponse)
	assert.NotNil(t, loginResponse.Profile)
	assert.Equal(t, profileData.ID, loginResponse.Profile.ID)
	assert.Equal(t, profileData.Username, loginResponse.Profile.Username)
	assert.Equal(t, "", loginResponse.Profile.Password) // Password should be cleared

	// Assert JWT token
	assert.NotEmpty(t, loginResponse.Token)

	// Validate token structure
	tokenParts := strings.Split(loginResponse.Token, ".")
	assert.Equal(t, 3, len(tokenParts), "JWT token should have 3 parts")

	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockRepository)
	// Create a noop logger
	logger := log.NewNopLogger()
	// Create a service with the mock repository
	svc := NewService(mockRepo, logger, false)

	// Setup test data
	username := "testuser"
	password := "password123"
	wrongPassword := "wrongpassword"
	id := uuid.New().String()

	// Generate a real bcrypt hash of the correct password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err, "Password hashing should not error")

	profileData := &model.Profile{
		ID:        id,
		Username:  username,
		Email:     "test@example.com",
		Password:  string(hashedPassword), // Use the hashed password
		FirstName: "Test",
		LastName:  "User",
		Bio:       "This is a test user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	loginReq := model.LoginRequest{
		Username: username,
		Password: wrongPassword, // Wrong password
	}

	// Setup expectations
	mockRepo.On("GetProfileByUsername", mock.Anything, username).Return(profileData, nil)

	// Call the method
	loginResponse, err := svc.Login(context.Background(), loginReq)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidCredentials, err)
	assert.Nil(t, loginResponse)
	mockRepo.AssertExpectations(t)
}

func TestGetProfile(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockRepository)
	// Create a noop logger
	logger := log.NewNopLogger()
	// Create a service with the mock repository
	svc := NewService(mockRepo, logger, false)

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
	svc := NewService(mockRepo, logger, false)

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
	svc := NewService(mockRepo, logger, false)

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
	svc := NewService(mockRepo, logger, false)

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
	svc := NewService(mockRepo, logger, false)

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
	svc := NewService(mockRepo, logger, false)

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

func TestValidateToken(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockRepository)
	// Create a noop logger
	logger := log.NewNopLogger()
	// Create a service with the mock repository
	svc := NewService(mockRepo, logger, false)

	// Setup test data
	username := "testuser"
	password := "password123"
	id := uuid.New().String()

	// Generate a real bcrypt hash of the password for testing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err, "Password hashing should not error")

	profileData := &model.Profile{
		ID:        id,
		Username:  username,
		Email:     "test@example.com",
		Password:  string(hashedPassword), // Use the hashed password
		FirstName: "Test",
		LastName:  "User",
		Bio:       "This is a test user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	loginReq := model.LoginRequest{
		Username: username,
		Password: password, // Plain text password for login request
	}

	// Setup expectations
	mockRepo.On("GetProfileByUsername", mock.Anything, username).Return(profileData, nil)

	// Login to get a token
	loginResponse, err := svc.Login(context.Background(), loginReq)
	assert.NoError(t, err)
	assert.NotEmpty(t, loginResponse.Token)

	// Get the token from the login response
	token := loginResponse.Token

	// Validate the token
	claims, err := svc.ValidateToken(context.Background(), token)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, id, claims.UserID)
	assert.Equal(t, username, claims.Username)
	// Expiration time should be in the future
	assert.True(t, claims.ExpiresAt.After(time.Now()))

	mockRepo.AssertExpectations(t)
}

func TestValidateToken_Invalid(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockRepository)
	// Create a noop logger
	logger := log.NewNopLogger()
	// Create a service with the mock repository
	svc := NewService(mockRepo, logger, false)

	// Test with an invalid token
	invalidToken := "invalid.token.format"
	claims, err := svc.ValidateToken(context.Background(), invalidToken)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidToken, err)
	assert.Nil(t, claims)

	// Test with an empty token
	claims, err = svc.ValidateToken(context.Background(), "")

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidInput, err)
	assert.Nil(t, claims)
}

func TestRegisterProfile_OnlyOneProfile(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockRepository)
	// Create a noop logger
	logger := log.NewNopLogger()
	// Create a service with onlyOneProfile set to true
	svc := NewService(mockRepo, logger, true)

	// Setup test data
	req := model.RegisterProfileRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Bio:       "This is a test user",
	}

	// Test case 1: No existing profiles
	mockRepo.On("CountProfiles", mock.Anything).Return(0, nil).Once()
	mockRepo.On("CreateProfile", mock.Anything, mock.MatchedBy(func(p model.Profile) bool {
		return p.Username == req.Username &&
			p.Email == req.Email &&
			p.Password != req.Password && // Password should be hashed, not plaintext
			len(p.Password) > 0 && // Ensure password is not empty
			p.FirstName == req.FirstName &&
			p.LastName == req.LastName &&
			p.Bio == req.Bio
	})).Return(nil).Once()

	// Call the method when no profiles exist
	profile, err := svc.RegisterProfile(context.Background(), req)

	// Assertions for first call
	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, req.Username, profile.Username)

	// Test case 2: Existing profile should prevent registration
	mockRepo.On("CountProfiles", mock.Anything).Return(1, nil).Once()

	// Call the method when a profile already exists
	profile2, err2 := svc.RegisterProfile(context.Background(), req)

	// Assertions for second call
	assert.Error(t, err2)
	assert.Nil(t, profile2)
	assert.Equal(t, "registration is disabled", err2.Error())

	mockRepo.AssertExpectations(t)
}
