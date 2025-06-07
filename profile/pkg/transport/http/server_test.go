package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ganis/okblog/profile/pkg/model"
	"github.com/go-kit/log"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of service.Service
type MockService struct {
	mock.Mock
}

func (m *MockService) RegisterProfile(ctx context.Context, req model.RegisterProfileRequest) (*model.Profile, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Profile), args.Error(1)
}

func (m *MockService) Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.LoginResponse), args.Error(1)
}

func (m *MockService) GetProfile(ctx context.Context, id string) (*model.Profile, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Profile), args.Error(1)
}

func (m *MockService) UpdateProfile(ctx context.Context, id string, req model.UpdateProfileRequest) (*model.Profile, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Profile), args.Error(1)
}

func (m *MockService) DeleteProfile(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockService) ValidateToken(ctx context.Context, token string) (*model.TokenClaims, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TokenClaims), args.Error(1)
}

func setupMockServer() (*MockService, *Server, *httptest.Server) {
	mockSvc := new(MockService)
	logger := log.NewNopLogger()
	server := NewServer(mockSvc, logger, nil)

	// Create a test HTTP server
	testServer := httptest.NewServer(server)

	return mockSvc, server, testServer
}

func TestRegisterProfileEndpoint(t *testing.T) {
	mockSvc, _, testServer := setupMockServer()
	defer testServer.Close()

	// Setup mock service
	profileReq := model.RegisterProfileRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Bio:       "This is a test user",
	}

	now := time.Now()
	expectedProfile := &model.Profile{
		ID:        uuid.New().String(),
		Username:  profileReq.Username,
		Email:     profileReq.Email,
		Password:  "", // Password should not be returned in response
		FirstName: profileReq.FirstName,
		LastName:  profileReq.LastName,
		Bio:       profileReq.Bio,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockSvc.On("RegisterProfile", mock.Anything, profileReq).Return(expectedProfile, nil)

	// Create request body
	reqBody, _ := json.Marshal(profileReq)

	// Send request
	resp, err := http.Post(testServer.URL+"/api/profiles/register", "application/json", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assertions
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseProfile model.Profile
	err = json.NewDecoder(resp.Body).Decode(&responseProfile)
	assert.NoError(t, err)

	assert.Equal(t, expectedProfile.ID, responseProfile.ID)
	assert.Equal(t, expectedProfile.Username, responseProfile.Username)
	assert.Equal(t, expectedProfile.Email, responseProfile.Email)
	assert.Empty(t, responseProfile.Password) // Password should not be in response

	mockSvc.AssertExpectations(t)
}

func TestLoginEndpoint(t *testing.T) {
	mockSvc, _, testServer := setupMockServer()
	defer testServer.Close()

	// Setup mock service
	loginReq := model.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	expectedProfile := &model.Profile{
		ID:        uuid.New().String(),
		Username:  loginReq.Username,
		Email:     "test@example.com",
		Password:  "", // Password should be empty in response
		FirstName: "Test",
		LastName:  "User",
		Bio:       "This is a test user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	expectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxMjM0NTY3ODkwIiwidXNlcm5hbWUiOiJ0ZXN0dXNlciJ9.4iN4aEJXDXY74C8uUe163X5PDF48FiRUUQJ-HbyX4WA"

	expectedResponse := &model.LoginResponse{
		Profile: expectedProfile,
		Token:   expectedToken,
	}

	mockSvc.On("Login", mock.Anything, loginReq).Return(expectedResponse, nil)

	// Create request body
	reqBody, _ := json.Marshal(loginReq)

	// Send request
	resp, err := http.Post(testServer.URL+"/api/profiles/login", "application/json", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assertions
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResponse model.LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	assert.NoError(t, err)

	assert.NotNil(t, loginResponse.Profile)
	assert.Equal(t, expectedProfile.ID, loginResponse.Profile.ID)
	assert.Equal(t, expectedProfile.Username, loginResponse.Profile.Username)
	assert.Equal(t, expectedProfile.Email, loginResponse.Profile.Email)
	assert.Empty(t, loginResponse.Profile.Password) // Password should not be in response

	// Verify token is returned
	assert.Equal(t, expectedToken, loginResponse.Token)

	mockSvc.AssertExpectations(t)
}

func TestGetProfileEndpoint(t *testing.T) {
	mockSvc, _, testServer := setupMockServer()
	defer testServer.Close()

	// Setup mock service
	id := uuid.New().String()
	now := time.Now()
	expectedProfile := &model.Profile{
		ID:        id,
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Bio:       "This is a test user",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockSvc.On("GetProfile", mock.Anything, id).Return(expectedProfile, nil)

	// Send request
	resp, err := http.Get(testServer.URL + "/api/profiles/" + id)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assertions
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseProfile model.Profile
	err = json.NewDecoder(resp.Body).Decode(&responseProfile)
	assert.NoError(t, err)

	assert.Equal(t, expectedProfile.ID, responseProfile.ID)
	assert.Equal(t, expectedProfile.Username, responseProfile.Username)
	assert.Equal(t, expectedProfile.Email, responseProfile.Email)

	mockSvc.AssertExpectations(t)
}

func TestUpdateProfileEndpoint(t *testing.T) {
	mockSvc, _, testServer := setupMockServer()
	defer testServer.Close()

	// Setup mock service
	id := uuid.New().String()
	updateReq := model.UpdateProfileRequest{
		FirstName: "Updated",
		LastName:  "Name",
		Bio:       "Updated bio",
	}

	now := time.Now()
	expectedProfile := &model.Profile{
		ID:        id,
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: updateReq.FirstName,
		LastName:  updateReq.LastName,
		Bio:       updateReq.Bio,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockSvc.On("UpdateProfile", mock.Anything, id, mock.AnythingOfType("model.UpdateProfileRequest")).Return(expectedProfile, nil)

	// Create request body
	reqBody, _ := json.Marshal(struct {
		ID   string                     `json:"id"`
		Data model.UpdateProfileRequest `json:"data"`
	}{
		ID:   id,
		Data: updateReq,
	})

	// Create request
	req, _ := http.NewRequest(http.MethodPut, testServer.URL+"/api/profiles/"+id, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assertions
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseProfile model.Profile
	err = json.NewDecoder(resp.Body).Decode(&responseProfile)
	assert.NoError(t, err)

	assert.Equal(t, expectedProfile.ID, responseProfile.ID)
	assert.Equal(t, updateReq.FirstName, responseProfile.FirstName)
	assert.Equal(t, updateReq.LastName, responseProfile.LastName)
	assert.Equal(t, updateReq.Bio, responseProfile.Bio)

	mockSvc.AssertExpectations(t)
}

func TestDeleteProfileEndpoint(t *testing.T) {
	mockSvc, _, testServer := setupMockServer()
	defer testServer.Close()

	// Setup mock service
	id := uuid.New().String()
	mockSvc.On("DeleteProfile", mock.Anything, id).Return(nil)

	// Create request
	req, _ := http.NewRequest(http.MethodDelete, testServer.URL+"/api/profiles/"+id, nil)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assertions
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	mockSvc.AssertExpectations(t)
}

func TestHandleNotFound(t *testing.T) {
	// Create a router
	router := mux.NewRouter()

	// Create a test HTTP server
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	// Send request to a non-existent endpoint
	resp, err := http.Get(testServer.URL + "/non-existent")
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assertions
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestValidateTokenEndpoint(t *testing.T) {
	mockSvc, _, testServer := setupMockServer()
	defer testServer.Close()

	// Setup mock service
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxMjM0NTY3ODkwIiwidXNlcm5hbWUiOiJ0ZXN0dXNlciJ9.4iN4aEJXDXY74C8uUe163X5PDF48FiRUUQJ-HbyX4WA"

	// Create expected claims
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	expectedClaims := &model.TokenClaims{
		UserID:    "1234567890",
		Username:  "testuser",
		IssuedAt:  now,
		ExpiresAt: expiresAt,
	}

	mockSvc.On("ValidateToken", mock.Anything, token).Return(expectedClaims, nil)

	// Create request with Authorization header
	req, _ := http.NewRequest(http.MethodPost, testServer.URL+"/api/profiles/validate-token", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assertions
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var validationResponse model.TokenValidationResponse
	err = json.NewDecoder(resp.Body).Decode(&validationResponse)
	assert.NoError(t, err)

	assert.True(t, validationResponse.Valid)
	assert.NotNil(t, validationResponse.Claims)
	assert.Equal(t, expectedClaims.UserID, validationResponse.Claims.UserID)
	assert.Equal(t, expectedClaims.Username, validationResponse.Claims.Username)

	mockSvc.AssertExpectations(t)
}

func TestValidateInvalidTokenEndpoint(t *testing.T) {
	mockSvc, _, testServer := setupMockServer()
	defer testServer.Close()

	// Setup mock service
	invalidToken := "invalid.token.format"

	mockSvc.On("ValidateToken", mock.Anything, invalidToken).Return(nil, errors.New("unauthorized: invalid token"))

	// Create request with Authorization header
	req, _ := http.NewRequest(http.MethodPost, testServer.URL+"/api/profiles/validate-token", nil)
	req.Header.Set("Authorization", "Bearer "+invalidToken)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Read error message from response
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "unauthorized: invalid token")

	mockSvc.AssertExpectations(t)
}

func TestMissingAuthorizationHeader(t *testing.T) {
	_, _, testServer := setupMockServer()
	defer testServer.Close()

	// Create request without Authorization header
	req, _ := http.NewRequest(http.MethodPost, testServer.URL+"/api/profiles/validate-token", nil)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestInvalidAuthorizationFormat(t *testing.T) {
	_, _, testServer := setupMockServer()
	defer testServer.Close()

	// Create request with invalid Authorization format
	req, _ := http.NewRequest(http.MethodPost, testServer.URL+"/api/profiles/validate-token", nil)
	req.Header.Set("Authorization", "InvalidFormat")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
