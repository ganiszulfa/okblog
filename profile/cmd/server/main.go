package main

import (
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/ganis/okblog/profile/pkg/database"
	"github.com/ganis/okblog/profile/pkg/logging"
	"github.com/ganis/okblog/profile/pkg/repository"
	"github.com/ganis/okblog/profile/pkg/service"
	httptransport "github.com/ganis/okblog/profile/pkg/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func main() {
	// Create a base logger
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	// Setup Kibana logger if enabled
	useKibana := getEnvBool("USE_KIBANA_LOGGING", false)
	if useKibana {
		// Get Kibana logger configuration
		kibanaConfig := logging.DefaultConfig()

		// Create Kibana logger
		kibanaLogger, err := logging.NewKibanaLogger(kibanaConfig, logger)
		if err != nil {
			level.Error(logger).Log("msg", "Failed to initialize Kibana logger", "err", err)
		} else {
			// Use Kibana logger as the primary logger
			logger = kibanaLogger
			level.Info(logger).Log("msg", "Kibana logging enabled")
		}
	}

	dbConfig := getDbConfig(logger)
	db, err := database.NewPostgresDB(dbConfig, logger)
	if err != nil {
		level.Error(logger).Log("msg", "Failed to connect to database", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	// Check if we should only allow one profile
	onlyOneProfile := getEnvBool("ONLY_ONE_PROFILE", true)

	// Initialize repository
	repo := repository.NewPostgresRepository(db, logger)

	// Create service with repository and logging middleware
	var svc service.Service
	svc = service.NewService(repo, logger, onlyOneProfile)
	svc = service.LoggingMiddleware(logger)(svc)

	// Create HTTP server with the service
	server := httptransport.NewServer(svc, logger)

	// Create a channel to listen for errors coming from the listener.
	errs := make(chan error, 2)

	// Start the server
	go func() {
		port := getEnv("PORT", "8080")
		logger.Log("transport", "HTTP", "addr", ":"+port, "msg", "Starting server")
		errs <- http.ListenAndServe(":"+port, server)
	}()

	// Listen for an interrupt signal
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		logger.Log("signal", "interrupt", "msg", "Shutting down")
		errs <- nil
	}()

	logger.Log("exit", <-errs)
}

func getDbConfig(logger log.Logger) database.Config {
	dbPortStr := getEnv("DB_PORT", "5432")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		level.Error(logger).Log("msg", "Invalid DB_PORT", "value", dbPortStr, "err", err)
		dbPort = 5432 // Default to 5432 if invalid
	}

	// Database configuration
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     dbPort,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "profile"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
	return dbConfig
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvBool gets a boolean environment variable or returns a default value
func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}
