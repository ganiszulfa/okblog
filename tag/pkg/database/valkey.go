package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// clientInstance is the singleton Valkey client instance
	clientInstance *redis.Client
)

// InitValkeyClient initializes the Valkey client.
func InitValkeyClient() error {
	valkeyAddr := os.Getenv("VALKEY_ADDR")
	if valkeyAddr == "" {
		valkeyAddr = "localhost:6379" // Default Valkey address
		log.Printf("VALKEY_ADDR not set, using default: %s", valkeyAddr)
	}

	opts := &redis.Options{
		Addr: valkeyAddr,
	}
	// Potentially add more options like DB, Password, PoolSize etc. from env vars

	clientInstance = redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := clientInstance.Ping(ctx).Result(); err != nil {
		log.Printf("Could not connect to Valkey: %v", err)
		return err
	}
	log.Println("Successfully connected to Valkey")
	return nil
}

// GetClient returns the existing Valkey client instance.
// It's assumed InitValkeyClient has been called successfully before this.
func GetClient() *redis.Client {
	if clientInstance == nil {
		log.Fatal("Valkey client not initialized. Call InitValkeyClient first.")
	}
	return clientInstance
}
