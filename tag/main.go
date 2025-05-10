package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"okblog/tag/pkg/consumer"
	"okblog/tag/pkg/database"
	"okblog/tag/pkg/handler"

	"github.com/redis/go-redis/v9"
)

// Post struct is now defined in models.go

var valkeyClient *redis.Client

// Global constants (could also be in a dedicated config.go)
const (
	kafkaTopic          = "post-db.okblog.posts"
	kafkaGroupID        = "tag-service-group"
	valkeyPostSetPrefix = "post_" // Used by kafka_consumer and http_handler
	// valkeyPostDetailsPrefix is now defined in valkey_client.go
)

func main() {
	// Initialize Valkey client
	if err := database.InitValkeyClient(); err != nil {
		log.Fatalf("Failed to initialize Valkey client: %v", err)
	}

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
