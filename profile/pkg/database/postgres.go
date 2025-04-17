package database

import (
	"database/sql"
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	_ "github.com/lib/pq"
)

// Config holds the database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(config Config, logger log.Logger) (*sql.DB, error) {
	// Create connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		level.Error(logger).Log("msg", "Failed to open database connection", "err", err)
		return nil, err
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		level.Error(logger).Log("msg", "Failed to ping database", "err", err)
		return nil, err
	}

	level.Info(logger).Log("msg", "Successfully connected to the database")
	return db, nil
}
