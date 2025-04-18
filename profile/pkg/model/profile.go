package model

import "time"

// Profile represents a user profile
type Profile struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Password is not returned in JSON
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// RegisterProfileRequest represents the request to register a new profile
type RegisterProfileRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Bio       string `json:"bio"`
}

// LoginRequest represents the credentials needed for login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the response returned after successful login
type LoginResponse struct {
	Profile *Profile `json:"profile"`
	Token   string   `json:"token"`
}

// TokenClaims represents the data stored in the JWT token
type TokenClaims struct {
	UserID    string    `json:"userId"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// TokenValidationRequest represents the request to validate a token
type TokenValidationRequest struct {
	Token string `json:"token"`
}

// TokenValidationResponse represents the response to a token validation request
type TokenValidationResponse struct {
	Valid  bool         `json:"valid"`
	Claims *TokenClaims `json:"claims,omitempty"`
}

// UpdateProfileRequest represents the request to update an existing profile
type UpdateProfileRequest struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Bio       string `json:"bio,omitempty"`
}
