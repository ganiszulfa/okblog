package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ganis/okblog/profile/pkg/model"
	"github.com/ganis/okblog/profile/pkg/repository"
	"github.com/go-kit/log"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrProfileNotFound       = errors.New("profile not found")
	ErrInvalidInput          = errors.New("invalid input")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrHashingFailed         = errors.New("password hashing failed")
	ErrTokenGenerationFailed = errors.New("failed to generate token")
	ErrInvalidToken          = errors.New("invalid token")
)

// JWT signing key - loaded from environment variable or default value
var jwtSigningKey = getJWTSigningKey()

// getJWTSigningKey loads the JWT signing key from environment variable or uses default
func getJWTSigningKey() []byte {
	key := os.Getenv("JWT_SIGNING_KEY")
	if key == "" {
		// Default value if environment variable is not set
		return []byte("my_secret_key")
	}
	return []byte(key)
}

// JWT token expiration time
const jwtExpirationTime = 14 * 24 * time.Hour

// JWTClaims represents the data stored in the JWT token
type JWTClaims struct {
	UserID    string    `json:"userId"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// Service defines the interface for profile operations
type Service interface {
	RegisterProfile(ctx context.Context, req model.RegisterProfileRequest) (*model.Profile, error)
	Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error)
	ValidateToken(ctx context.Context, token string) (*model.TokenClaims, error)
	GetProfile(ctx context.Context, id string) (*model.Profile, error)
	UpdateProfile(ctx context.Context, id string, req model.UpdateProfileRequest) (*model.Profile, error)
	DeleteProfile(ctx context.Context, id string) error
}

// profileService implements the Service interface
type profileService struct {
	repo           repository.Repository
	logger         log.Logger
	onlyOneProfile bool
}

// NewService creates a new instance of the profile service
func NewService(repo repository.Repository, logger log.Logger, onlyOneProfile bool) Service {
	return &profileService{
		repo:           repo,
		logger:         logger,
		onlyOneProfile: onlyOneProfile,
	}
}

func (s *profileService) RegisterProfile(ctx context.Context, req model.RegisterProfileRequest) (*model.Profile, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, ErrInvalidInput
	}

	// Check if we only allow one profile
	if s.onlyOneProfile {
		// Count existing profiles
		count, err := s.repo.CountProfiles(ctx)
		if err != nil {
			s.logger.Log("err", err, "msg", "Failed to count profiles")
			return nil, err
		}

		// If profiles already exist, prevent registration
		if count > 0 {
			s.logger.Log("msg", "Registration blocked due to ONLY_ONE_PROFILE configuration")
			return nil, errors.New("registration is disabled")
		}
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

func (s *profileService) Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
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

	// Generate JWT token
	token, err := s.generateJWTToken(profile)
	if err != nil {
		s.logger.Log("err", err, "msg", "Failed to generate JWT token")
		return nil, ErrTokenGenerationFailed
	}

	// Don't return the password in the response
	profile.Password = ""

	// Create login response with profile and token
	response := &model.LoginResponse{
		Profile: profile,
		Token:   token,
	}

	return response, nil
}

// generateJWTToken creates a new JWT token for the user
func (s *profileService) generateJWTToken(profile *model.Profile) (string, error) {
	// Create JWT header (algorithm & token type)
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	// Convert header to JSON and encode to base64
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	headerBase64 := base64.RawURLEncoding.EncodeToString(headerJSON)

	// Create JWT payload (claims)
	now := time.Now()
	expiresAt := now.Add(jwtExpirationTime)

	claims := JWTClaims{
		UserID:    profile.ID,
		Username:  profile.Username,
		IssuedAt:  now,
		ExpiresAt: expiresAt,
	}

	// Convert payload to JSON and encode to base64
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payloadBase64 := base64.RawURLEncoding.EncodeToString(payloadJSON)

	// Create the signature
	signatureInput := fmt.Sprintf("%s.%s", headerBase64, payloadBase64)
	h := hmac.New(sha256.New, jwtSigningKey)
	h.Write([]byte(signatureInput))
	signature := h.Sum(nil)
	signatureBase64 := base64.RawURLEncoding.EncodeToString(signature)

	// Combine all parts to create the complete JWT token
	token := fmt.Sprintf("%s.%s.%s", headerBase64, payloadBase64, signatureBase64)

	return token, nil
}

// ValidateJWTToken validates a JWT token and returns the claims if valid
func (s *profileService) ValidateJWTToken(tokenString string) (*JWTClaims, error) {
	// Split the token into its parts
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	headerBase64, payloadBase64, signatureBase64 := parts[0], parts[1], parts[2]

	// Verify the signature
	signatureInput := fmt.Sprintf("%s.%s", headerBase64, payloadBase64)
	h := hmac.New(sha256.New, jwtSigningKey)
	h.Write([]byte(signatureInput))
	expectedSignature := h.Sum(nil)
	expectedSignatureBase64 := base64.RawURLEncoding.EncodeToString(expectedSignature)

	if signatureBase64 != expectedSignatureBase64 {
		return nil, errors.New("invalid token signature")
	}

	// Decode the payload
	payloadJSON, err := base64.RawURLEncoding.DecodeString(payloadBase64)
	if err != nil {
		return nil, errors.New("invalid token payload encoding")
	}

	var claims JWTClaims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return nil, errors.New("invalid token payload format")
	}

	// Check if the token is expired
	if time.Now().After(claims.ExpiresAt) {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}

// ValidateToken validates a JWT token and returns the claims if valid
func (s *profileService) ValidateToken(ctx context.Context, token string) (*model.TokenClaims, error) {
	if token == "" {
		return nil, ErrInvalidInput
	}

	// Validate the token
	claims, err := s.ValidateJWTToken(token)
	if err != nil {
		s.logger.Log("err", err, "msg", "Token validation failed")
		return nil, ErrInvalidToken
	}

	// Convert internal claims to the model claims
	tokenClaims := &model.TokenClaims{
		UserID:    claims.UserID,
		Username:  claims.Username,
		IssuedAt:  claims.IssuedAt,
		ExpiresAt: claims.ExpiresAt,
	}

	return tokenClaims, nil
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
