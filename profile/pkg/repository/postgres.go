package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ganis/okblog/profile/pkg/model"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// Repository defines the interface for profile storage operations
type Repository interface {
	CreateProfile(ctx context.Context, profile model.Profile) error
	GetProfile(ctx context.Context, id string) (*model.Profile, error)
	GetProfileByUsername(ctx context.Context, username string) (*model.Profile, error)
	UpdateProfile(ctx context.Context, profile model.Profile) error
	DeleteProfile(ctx context.Context, id string) error
}

// PostgresRepository implements the Repository interface using PostgreSQL
type PostgresRepository struct {
	db     *sql.DB
	logger log.Logger
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *sql.DB, logger log.Logger) Repository {
	return &PostgresRepository{
		db:     db,
		logger: logger,
	}
}

// CreateProfile creates a new profile in the database
func (r *PostgresRepository) CreateProfile(ctx context.Context, profile model.Profile) error {
	query := `
		INSERT INTO profiles (id, username, email, password, first_name, last_name, bio, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		profile.ID,
		profile.Username,
		profile.Email,
		profile.Password,
		profile.FirstName,
		profile.LastName,
		profile.Bio,
		profile.CreatedAt,
		profile.UpdatedAt,
	)

	if err != nil {
		level.Error(r.logger).Log("msg", "Failed to create profile", "err", err)
		return err
	}

	return nil
}

// GetProfile retrieves a profile from the database by ID
func (r *PostgresRepository) GetProfile(ctx context.Context, id string) (*model.Profile, error) {
	query := `
		SELECT id, username, email, password, first_name, last_name, bio, created_at, updated_at
		FROM profiles
		WHERE id = $1
	`

	var profile model.Profile
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&profile.ID,
		&profile.Username,
		&profile.Email,
		&profile.Password,
		&profile.FirstName,
		&profile.LastName,
		&profile.Bio,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("profile not found")
		}
		level.Error(r.logger).Log("msg", "Failed to get profile", "err", err)
		return nil, err
	}

	return &profile, nil
}

// GetProfileByUsername retrieves a profile from the database by username
func (r *PostgresRepository) GetProfileByUsername(ctx context.Context, username string) (*model.Profile, error) {
	query := `
		SELECT id, username, email, password, first_name, last_name, bio, created_at, updated_at
		FROM profiles
		WHERE username = $1
	`

	var profile model.Profile
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&profile.ID,
		&profile.Username,
		&profile.Email,
		&profile.Password,
		&profile.FirstName,
		&profile.LastName,
		&profile.Bio,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("profile not found")
		}
		level.Error(r.logger).Log("msg", "Failed to get profile by username", "err", err)
		return nil, err
	}

	return &profile, nil
}

// UpdateProfile updates an existing profile in the database
func (r *PostgresRepository) UpdateProfile(ctx context.Context, profile model.Profile) error {
	query := `
		UPDATE profiles
		SET first_name = $1, last_name = $2, bio = $3, updated_at = $4
		WHERE id = $5
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		profile.FirstName,
		profile.LastName,
		profile.Bio,
		time.Now(),
		profile.ID,
	)

	if err != nil {
		level.Error(r.logger).Log("msg", "Failed to update profile", "err", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		level.Error(r.logger).Log("msg", "Failed to get rows affected", "err", err)
		return err
	}

	if rowsAffected == 0 {
		return errors.New("profile not found")
	}

	return nil
}

// DeleteProfile deletes a profile from the database by ID
func (r *PostgresRepository) DeleteProfile(ctx context.Context, id string) error {
	query := `DELETE FROM profiles WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		level.Error(r.logger).Log("msg", "Failed to delete profile", "err", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		level.Error(r.logger).Log("msg", "Failed to get rows affected", "err", err)
		return err
	}

	if rowsAffected == 0 {
		return errors.New("profile not found")
	}

	return nil
}
