package http

import (
	"bytes"
	"context"
	"encoding/json"
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

func (m *MockService) CreateProfile(ctx context.Context, req model.CreateProfileRequest) (*model.Profile, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Profile), args.Error(1)
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

func setupMockServer() (*MockService, *Server, *httptest.Server) {
	mockSvc := new(MockService)
	logger := log.NewNopLogger()
	server := NewServer(mockSvc, logger)

	// Create a test HTTP server
	testServer := httptest.NewServer(server)

	return mockSvc, server, testServer
}

func TestCreateProfileEndpoint(t *testing.T) {
	mockSvc, _, testServer := setupMockServer()
	defer testServer.Close()

	// Setup mock service
	profileReq := model.CreateProfileRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Bio:       "This is a test user",
	}

	now := time.Now()
	expectedProfile := &model.Profile{
		ID:        uuid.New().String(),
		Username:  profileReq.Username,
		Email:     profileReq.Email,
		FirstName: profileReq.FirstName,
		LastName:  profileReq.LastName,
		Bio:       profileReq.Bio,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockSvc.On("CreateProfile", mock.Anything, profileReq).Return(expectedProfile, nil)

	// Create request body
	reqBody, _ := json.Marshal(profileReq)

	// Send request
	resp, err := http.Post(testServer.URL+"/profiles", "application/json", bytes.NewBuffer(reqBody))
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
	resp, err := http.Get(testServer.URL + "/profiles/" + id)
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
	req, _ := http.NewRequest(http.MethodPut, testServer.URL+"/profiles/"+id, bytes.NewBuffer(reqBody))
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
	req, _ := http.NewRequest(http.MethodDelete, testServer.URL+"/profiles/"+id, nil)

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
