package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ganis/okblog/profile/pkg/model"
	"github.com/go-kit/log"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, Repository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	logger := log.NewNopLogger()
	repo := NewPostgresRepository(db, logger)

	return db, mock, repo
}

func TestCreateProfile(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	ctx := context.Background()
	now := time.Now()
	profile := model.Profile{
		ID:        uuid.New().String(),
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "$2a$10$hPkIwyYJBsmvKnXN9LBrNeoWsGnY6MQiEjgZXQvtdnVtPKQwvzBSG", // Bcrypt hash example
		FirstName: "Test",
		LastName:  "User",
		Bio:       "Test bio",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Set up expectations
	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO profiles (id, username, email, password, first_name, last_name, bio, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`)).WithArgs(
		profile.ID,
		profile.Username,
		profile.Email,
		profile.Password,
		profile.FirstName,
		profile.LastName,
		profile.Bio,
		profile.CreatedAt,
		profile.UpdatedAt,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the method
	err := repo.CreateProfile(ctx, profile)

	// Assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProfile(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	ctx := context.Background()
	id := uuid.New().String()
	now := time.Now()
	hashedPassword := "$2a$10$hPkIwyYJBsmvKnXN9LBrNeoWsGnY6MQiEjgZXQvtdnVtPKQwvzBSG" // Bcrypt hash example

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "first_name", "last_name", "bio", "created_at", "updated_at"}).
		AddRow(id, "testuser", "test@example.com", hashedPassword, "Test", "User", "Test bio", now, now)

	// Set up expectations
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, username, email, password, first_name, last_name, bio, created_at, updated_at
		FROM profiles
		WHERE id = $1
	`)).WithArgs(id).WillReturnRows(rows)

	// Call the method
	profile, err := repo.GetProfile(ctx, id)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, id, profile.ID)
	assert.Equal(t, "testuser", profile.Username)
	assert.Equal(t, "test@example.com", profile.Email)
	assert.Equal(t, hashedPassword, profile.Password) // Check password is retrieved correctly
	assert.Equal(t, "Test", profile.FirstName)
	assert.Equal(t, "User", profile.LastName)
	assert.Equal(t, "Test bio", profile.Bio)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProfile_NotFound(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	ctx := context.Background()
	id := uuid.New().String()

	// Set up expectations
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, username, email, password, first_name, last_name, bio, created_at, updated_at
		FROM profiles
		WHERE id = $1
	`)).WithArgs(id).WillReturnError(sql.ErrNoRows)

	// Call the method
	profile, err := repo.GetProfile(ctx, id)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, profile)
	assert.Equal(t, "profile not found", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProfileByUsername(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	ctx := context.Background()
	id := uuid.New().String()
	username := "testuser"
	now := time.Now()
	hashedPassword := "$2a$10$hPkIwyYJBsmvKnXN9LBrNeoWsGnY6MQiEjgZXQvtdnVtPKQwvzBSG" // Bcrypt hash example

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "first_name", "last_name", "bio", "created_at", "updated_at"}).
		AddRow(id, username, "test@example.com", hashedPassword, "Test", "User", "Test bio", now, now)

	// Set up expectations
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, username, email, password, first_name, last_name, bio, created_at, updated_at
		FROM profiles
		WHERE username = $1
	`)).WithArgs(username).WillReturnRows(rows)

	// Call the method
	profile, err := repo.GetProfileByUsername(ctx, username)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, id, profile.ID)
	assert.Equal(t, username, profile.Username)
	assert.Equal(t, "test@example.com", profile.Email)
	assert.Equal(t, hashedPassword, profile.Password) // Check password is retrieved correctly
	assert.Equal(t, "Test", profile.FirstName)
	assert.Equal(t, "User", profile.LastName)
	assert.Equal(t, "Test bio", profile.Bio)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProfileByUsername_NotFound(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	ctx := context.Background()
	username := "nonexistentuser"

	// Set up expectations
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, username, email, password, first_name, last_name, bio, created_at, updated_at
		FROM profiles
		WHERE username = $1
	`)).WithArgs(username).WillReturnError(sql.ErrNoRows)

	// Call the method
	profile, err := repo.GetProfileByUsername(ctx, username)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, profile)
	assert.Equal(t, "profile not found", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProfile(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	ctx := context.Background()
	now := time.Now()
	profile := model.Profile{
		ID:        uuid.New().String(),
		FirstName: "Updated",
		LastName:  "Name",
		Bio:       "Updated bio",
		UpdatedAt: now,
	}

	// Set up expectations for time.Now() in the UpdateProfile function
	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE profiles
		SET first_name = $1, last_name = $2, bio = $3, updated_at = $4
		WHERE id = $5
	`)).WithArgs(
		profile.FirstName,
		profile.LastName,
		profile.Bio,
		sqlmock.AnyArg(), // For updated_at which is set in the function
		profile.ID,
	).WillReturnResult(sqlmock.NewResult(0, 1))

	// Call the method
	err := repo.UpdateProfile(ctx, profile)

	// Assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProfile_NotFound(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	ctx := context.Background()
	profile := model.Profile{
		ID:        uuid.New().String(),
		FirstName: "Updated",
		LastName:  "Name",
		Bio:       "Updated bio",
	}

	// Set up expectations
	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE profiles
		SET first_name = $1, last_name = $2, bio = $3, updated_at = $4
		WHERE id = $5
	`)).WithArgs(
		profile.FirstName,
		profile.LastName,
		profile.Bio,
		sqlmock.AnyArg(), // For updated_at which is set in the function
		profile.ID,
	).WillReturnResult(sqlmock.NewResult(0, 0))

	// Call the method
	err := repo.UpdateProfile(ctx, profile)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, "profile not found", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteProfile(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	ctx := context.Background()
	id := uuid.New().String()

	// Set up expectations
	mock.ExpectExec("DELETE FROM profiles WHERE id = \\$1").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call the method
	err := repo.DeleteProfile(ctx, id)

	// Assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteProfile_NotFound(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	ctx := context.Background()
	id := uuid.New().String()

	// Set up expectations
	mock.ExpectExec("DELETE FROM profiles WHERE id = \\$1").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Call the method
	err := repo.DeleteProfile(ctx, id)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, "profile not found", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}
