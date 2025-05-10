package consumer

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"

	"okblog/tag/pkg/config"
	"okblog/tag/pkg/database"
	"okblog/tag/pkg/models"
)

func StartPostsConsumer() {
	kafkaBrokersEnv := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokersEnv == "" {
		log.Fatal("KAFKA_BROKERS environment variable must be set (comma-separated list)")
	}
	log.Printf("Kafka brokers: %s", kafkaBrokersEnv)
	kafkaBrokers := strings.Split(kafkaBrokersEnv, ",")

	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  kafkaBrokers,
		Topic:    config.TopicPostDB,
		GroupID:  config.KafkaGroupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		MaxWait:  1 * time.Second,
	})
	defer kafkaReader.Close()
	log.Printf("Kafka reader configured for topic %s, group %s", config.TopicPostDB, config.KafkaGroupID)

	log.Println("Starting Kafka message processing loop...")
	valkeyClient := database.GetClient()

	for {
		m, err := kafkaReader.FetchMessage(context.Background())
		if err != nil {
			log.Printf("Error fetching message: %v", err)
			time.Sleep(1 * time.Second) // Simple backoff
			continue
		}

		log.Printf("Message received: Topic %s, Partition %d, Offset %d, Key %s", m.Topic, m.Partition, m.Offset, string(m.Key))

		var debeziumMsg models.DebeziumMessage
		if err := json.Unmarshal(m.Value, &debeziumMsg); err != nil {
			log.Printf("Error unmarshalling Debezium message: %v. Message: %s", err, string(m.Value))
			if err := kafkaReader.CommitMessages(context.Background(), m); err != nil {
				log.Printf("Error committing message after Debezium unmarshal error: %v", err)
			}
			continue
		}

		// Check if this is a delete operation
		if debeziumMsg.Payload.Op == "d" && debeziumMsg.Payload.Before != nil {
			// For delete operations, remove from cache
			postID := debeziumMsg.Payload.Before.ID
			detailsKey := config.PostDetailsPrefix + postID
			if err := valkeyClient.Del(context.Background(), detailsKey).Err(); err != nil {
				log.Printf("Error deleting cached post details for deleted post ID %s: %v", postID, err)
			} else {
				log.Printf("Successfully deleted cached post details for deleted post ID %s", postID)
			}
			if err := kafkaReader.CommitMessages(context.Background(), m); err != nil {
				log.Printf("Error committing message for deleted post: %v", err)
			}
			continue
		}

		// Check if this is an update that unpublished a post
		if debeziumMsg.Payload.Op == "u" &&
			debeziumMsg.Payload.Before != nil &&
			debeziumMsg.Payload.After != nil &&
			debeziumMsg.Payload.Before.IsPublished &&
			!debeziumMsg.Payload.After.IsPublished {
			// Post was unpublished, remove from cache
			postID := debeziumMsg.Payload.After.ID
			detailsKey := config.PostDetailsPrefix + postID
			if err := valkeyClient.Del(context.Background(), detailsKey).Err(); err != nil {
				log.Printf("Error deleting cached post details for unpublished post ID %s: %v", postID, err)
			} else {
				log.Printf("Successfully deleted cached post details for unpublished post ID %s", postID)
			}
			if err := kafkaReader.CommitMessages(context.Background(), m); err != nil {
				log.Printf("Error committing message for unpublished post: %v", err)
			}
			continue
		}

		// We are interested in the 'after' state for creates and updates.
		payloadAfter := debeziumMsg.Payload.After
		if payloadAfter == nil {
			log.Printf("No 'after' data in Debezium message (Op: %s). Skipping post processing. Key: %s", debeziumMsg.Payload.Op, string(m.Key))
			if err := kafkaReader.CommitMessages(context.Background(), m); err != nil {
				log.Printf("Error committing message for op %s with no 'after' data: %v", debeziumMsg.Payload.Op, err)
			}
			continue
		}
		log.Printf("Payload after: %+v", payloadAfter)

		// Map data from payloadAfter to the application's Post struct
		var post models.Post
		post.ID = payloadAfter.ID
		post.Type = payloadAfter.Type
		post.IsPublished = payloadAfter.IsPublished
		post.Title = payloadAfter.Title
		post.Slug = payloadAfter.Slug
		post.ViewCount = int(payloadAfter.ViewCount)

		// Convert timestamps from microseconds epoch to time.Time
		if payloadAfter.PublishedAt != nil {
			post.PublishedAt = time.UnixMicro(*payloadAfter.PublishedAt)
		} else {
			post.PublishedAt = time.Time{} // Zero value for time.Time if null
		}

		if !post.IsPublished {
			log.Printf("Post ID %s is not published, skipping.", post.ID)
			if err := kafkaReader.CommitMessages(context.Background(), m); err != nil {
				log.Printf("Error committing message for unpublished post: %v", err)
			}
			continue
		}

		if strings.ToLower(post.Type) != "post" {
			log.Printf("Post ID %s is not a post but a %s, skipping.", post.ID, post.Type)
			if err := kafkaReader.CommitMessages(context.Background(), m); err != nil {
				log.Printf("Error committing message for non-post type: %v", err)
			}
			continue
		}

		// Process and store the post
		if err := processAndStorePost(context.Background(), valkeyClient, post); err != nil {
			log.Printf("Error processing post ID %s: %v. Skipping commit for now.", post.ID, err)
			continue // Skip commit if processing fails to allow reprocessing
		}

		if err := kafkaReader.CommitMessages(context.Background(), m); err != nil {
			log.Printf("Error committing message after successful processing: %v", err)
		}
	}
}

// processAndStorePost handles saving the post to Valkey for tags and for direct lookup.
func processAndStorePost(ctx context.Context, valkeyOps *redis.Client, post models.Post) error {
	if post.ID == "" {
		log.Println("Post ID is empty, cannot process.")
		return nil
	}

	// Create CachedPostDetails for Valkey storage
	cachedData := models.CachedPostDetails{
		Title:       post.Title,
		PublishedAt: post.PublishedAt,
		Tags:        post.Tags,
		Slug:        post.Slug,
		ViewCount:   post.ViewCount,
	}

	postJSON, err := json.Marshal(cachedData)
	if err != nil {
		log.Printf("Error marshalling cached post details for ID %s to JSON: %v", post.ID, err)
		return err
	}
	detailsKey := config.PostDetailsPrefix + post.ID
	if err := valkeyOps.Set(ctx, detailsKey, postJSON, 0).Err(); err != nil { // 0 for no expiration
		log.Printf("Error storing cached post details for ID %s in Valkey: %v", post.ID, err)
		return err
	}
	log.Printf("Successfully stored cached post details for ID %s in Valkey at key %s", post.ID, detailsKey)

	return nil
}

// StartPostTagsConsumer initializes and starts the Kafka consumer for post_tags messages.
func StartPostTagsConsumer() {
	kafkaBrokersEnv := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokersEnv == "" {
		log.Fatal("KAFKA_BROKERS environment variable must be set for post_tags consumer (comma-separated list)")
	}
	kafkaBrokers := strings.Split(kafkaBrokersEnv, ",")

	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  kafkaBrokers,
		Topic:    config.TopicPostTags,
		GroupID:  config.KafkaGroupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
		MaxWait:  1 * time.Second,
	})
	defer kafkaReader.Close()
	log.Printf("Kafka reader configured for topic %s, group %s", config.TopicPostTags, config.KafkaGroupID)

	log.Println("Starting Kafka message processing loop for post_tags...")

	for {
		m, err := kafkaReader.FetchMessage(context.Background())
		if err != nil {
			log.Printf("Error fetching message from %s: %v", config.TopicPostTags, err)
			time.Sleep(1 * time.Second) // Simple backoff
			continue
		}

		log.Printf("Message received from %s: Partition %d, Offset %d, Key %s", m.Topic, m.Partition, m.Offset, string(m.Key))

		var msg models.PostTagDebeziumMessage
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			log.Printf("Error unmarshalling PostTagDebeziumMessage from %s: %v. Message: %s", config.TopicPostTags, err, string(m.Value))
			if err := kafkaReader.CommitMessages(context.Background(), m); err != nil {
				log.Printf("Error committing message on %s after unmarshal error: %v", config.TopicPostTags, err)
			}
			continue
		}

		processPostTag(msg, kafkaReader, m)

	}
}

// Process the tag relationship by adding the post to the tag's sorted set
func processPostTag(msg models.PostTagDebeziumMessage, kafkaReader *kafka.Reader, m kafka.Message) {
	valkeyClient := database.GetClient()

	// Delete operation - remove post from tag's sorted set
	if msg.Payload.Op == "d" && msg.Payload.Before != nil {
		postID := msg.Payload.Before.PostID
		tag := msg.Payload.Before.Tag

		if postID == "" || tag == "" {
			log.Printf("Invalid post_tag relationship for deletion: PostID %s, Tag %s", postID, tag)
		} else {
			tagKey := config.TagPostsPrefix + tag
			if err := valkeyClient.ZRem(context.Background(), tagKey, postID).Err(); err != nil {
				log.Printf("Error removing post ID %s from tag %s sorted set: %v", postID, tag, err)
			} else {
				log.Printf("Successfully removed post ID %s from tag %s sorted set", postID, tag)
			}
		}

		if err := kafkaReader.CommitMessages(context.Background(), m); err != nil {
			log.Printf("Error committing message after tag deletion: %v", err)
		}
		return
	}

	// Extract post ID and tag from the message
	postID := msg.Payload.After.PostID
	tag := msg.Payload.After.Tag

	if postID == "" || tag == "" {
		log.Printf("Invalid post_tag relationship: PostID %s, Tag %s", postID, tag)
		if err := kafkaReader.CommitMessages(context.Background(), m); err != nil {
			log.Printf("Error committing message with invalid post_tag relationship: %v", err)
		}
		return
	}

	// Get post details to use timestamp as score for ordered set
	detailsKey := config.PostDetailsPrefix + postID
	postDetailsJSON, err := valkeyClient.Get(context.Background(), detailsKey).Result()
	if err != nil {
		log.Printf("Error retrieving post details for ID %s: %v", postID, err)
		// Continue processing even if we can't get post details
		// We'll use current time as score instead
	}

	// Default score to current time
	score := float64(time.Now().Unix())

	// If post details exist, use publication timestamp as score for proper time ordering
	if err == nil && postDetailsJSON != "" {
		var cachedPost models.CachedPostDetails
		if err := json.Unmarshal([]byte(postDetailsJSON), &cachedPost); err != nil {
			log.Printf("Error unmarshalling post details for ID %s: %v", postID, err)
		} else if !cachedPost.PublishedAt.IsZero() {
			// Use publication timestamp as score for the sorted set
			score = float64(cachedPost.PublishedAt.Unix())
		}
	}

	// Add post ID to the tag's sorted set with score based on publication time
	tagKey := config.TagPostsPrefix + tag
	if err := valkeyClient.ZAdd(context.Background(), tagKey, redis.Z{
		Score:  score,
		Member: postID,
	}).Err(); err != nil {
		log.Printf("Error adding post ID %s to tag %s tagkey %s sorted set: %v", postID, tag, tagKey, err)
	} else {
		log.Printf("Successfully added post ID %s to tag %s tagkey %s sorted set with score %f", postID, tag, tagKey, score)
	}

	if err := kafkaReader.CommitMessages(context.Background(), m); err != nil {
		log.Printf("Error committing processed message: %v", err)
	}
}
