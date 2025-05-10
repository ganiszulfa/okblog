package config

// Kafka configuration
const (
	// TopicPostDB is the Kafka topic for posts
	TopicPostDB = "post-db.okblog.posts"

	// TopicPostTags is the Kafka topic for post tags
	TopicPostTags = "post-db.okblog.post_tags"

	// KafkaGroupID is the consumer group ID
	KafkaGroupID = "tag-service-group"
)

// Valkey configuration
const (
	// PostDetailsPrefix is the prefix for post details in Valkey
	PostDetailsPrefix = "post:details:"

	// TagPostsPrefix is the prefix for tag posts sorted sets in Valkey
	TagPostsPrefix = "tag:posts:"
)

// HTTP configuration
const (
	// DefaultPage is the default page number for pagination
	DefaultPage = 1

	// DefaultPerPage is the default number of items per page
	DefaultPerPage = 10
)
