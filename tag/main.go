package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"okblog/tag/pkg/consumer"
	"okblog/tag/pkg/database"
	"okblog/tag/pkg/handler"
)

// Post struct is now defined in models.go

func main() {
	// Initialize Valkey client
	if err := database.InitValkeyClient(); err != nil {
		log.Fatalf("Failed to initialize Valkey client: %v", err)
	}

	// Initialize MySQL database connection
	if err := database.InitMySQLDB(); err != nil {
		log.Fatalf("Failed to initialize MySQL database: %v", err)
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

	log.Println("Shutdown signal received, exiting...")
	// Any additional cleanup logic could be added here
}
