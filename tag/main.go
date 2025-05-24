package main

import (
	"os"
	"os/signal"
	"syscall"

	"okblog/tag/pkg/consumer"
	"okblog/tag/pkg/database"
	"okblog/tag/pkg/handler"
	"okblog/tag/pkg/logger"
)

// Post struct is now defined in models.go

func main() {
	// Initialize logger
	logger.Initialize()
	logger.Info("Starting tag service", nil)

	// Initialize Valkey client
	if err := database.InitValkeyClient(); err != nil {
		logger.Fatal("Failed to initialize Valkey client", err)
	}

	// Initialize MySQL database connection
	if err := database.InitMySQLDB(); err != nil {
		logger.Fatal("Failed to initialize MySQL database", err)
	}
	defer database.CloseMySQLDB()

	// Initialize Fiber app
	app := handler.InitFiberApp()

	// Start Kafka consumer for posts in a goroutine
	go consumer.StartPostsConsumer()

	// Start Kafka consumer for post tags in a goroutine
	go consumer.StartPostTagsConsumer()

	// Start Fiber server in a goroutine
	go handler.StartHTTPServer(app)

	// Graceful shutdown handling
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan // Block until a signal is received

	logger.Info("Shutdown signal received, exiting...", nil)
	// Any additional cleanup logic could be added here
}
